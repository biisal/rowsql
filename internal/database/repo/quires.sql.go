package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/database/models"
	"github.com/biisal/rowsql/internal/logger"
	"github.com/biisal/rowsql/internal/utils"
)

func ErrorInvalidTable(tableName string) error {
	return fmt.Errorf("table %s not found", tableName)
}

var ErrorNotFound = errors.New("not found")

// func getColValues(rows []any, cols []models.ColValue) ([]models.ColValue, error) {
// 	if len(cols) != len(rows) {
// 		return nil, apperr.ErrorNotSameRowColsSize
// 	}

// 	colWithValues := make([]models.ColValue, len(cols))
// 	for i, col := range cols {
// 		colWithValues[i] = models.ColValue{
// 			ColumnName: col.ColumnName,
// 			Value:      rows[i],
// 			ColumnType: models.ColType{
// 				DataType:         col.ColumnType.DataType,
// 				IsUnique:         col.ColumnType.IsUnique,
// 				HasAutoIncrement: col.ColumnType.HasAutoIncrement,
// 				HasDefault:       col.ColumnType.HasDefault,
// 			},
// 		}
// 	}
// 	return colWithValues, nil
// }

func (q *Queries) GetQuotedTableName(tableName string) string {
	return q.queryBuilder.QuoteName(tableName)
}

func (q *Queries) CheckTableExitsInDB(ctx context.Context, tableName string) error {
	query, args := q.queryBuilder.CheckTableExitsQuery(tableName)
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error("failed to query: %v", err)
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Errorln(err)
		}
	}()
	if !rows.Next() {
		logger.Info("Query : %s", query)
		return ErrorInvalidTable(tableName)
	}
	return nil
}

func (q *Queries) ListColsMetaData(ctx context.Context, tableName string) ([]models.ColValue, error) {
	query, args := q.queryBuilder.ColumnsList(tableName)
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error("failed to query: %v", err)
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Errorln(err)
		}
	}()
	var items []models.ColValue
	for rows.Next() {
		var i models.ColValue
		if err := rows.Scan(&i.ColumnName, &i.ColumnType.DataType, &i.ColumnType.HasDefault, &i.ColumnType.IsUnique, &i.ColumnType.HasAutoIncrement); err != nil {
			logger.Error("failed to scan rows in list cols: %v", err)
			return nil, err
		}
		i.ColumnType.InputType = utils.GetInputType(i.ColumnType.DataType)
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		logger.Error("failed to scan rows: %v", err)
		return nil, err
	}
	return items, nil
}

func (q *Queries) ListTables(ctx context.Context) ([]models.ListTablesRow, error) {
	query := q.queryBuilder.ListTables()
	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Errorln(err)
		}
	}()
	var items []models.ListTablesRow
	for rows.Next() {
		var i models.ListTablesRow
		if err := rows.Scan(&i.TableSchema, &i.TableName); err != nil {
			logger.Error("failed to scan rows: %v", err)
			return nil, err
		}
		if i.TableName == historyTableName {
			continue
		}
		items = append([]models.ListTablesRow{i}, items...)
	}
	if err := rows.Err(); err != nil {
		logger.Error("failed to scan rows: %v", err)
		return nil, err
	}
	return items, nil
}

type ListDataProps struct {
	TableName string `json:"tableName"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
	Column    string `json:"column"`
	Order     string `json:"order"`
}

func (q *Queries) ListRows(ctx context.Context, props ListDataProps) ([]models.RowSet, error) {
	query, args, err := q.queryBuilder.ListRows(props.TableName, props.Column, props.Order, props.Limit, props.Offset)
	if err != nil {
		return nil, err
	}

	logger.Info("Query : %s", query)
	rows, err := q.db.QueryxContext(ctx, query, args...)
	if err != nil {
		logger.Errorln(err.Error())
		return nil, err
	}
	defer func() {
		if err = rows.Close(); err != nil {
			logger.Errorln(err)
		}
	}()
	data := make([]models.RowSet, 0)
	colValues, err := q.ListColsMetaData(ctx, props.TableName)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		cvs := make([]models.ColValue, len(colValues))
		copy(cvs, colValues)

		row, err := rows.SliceScan()
		if err != nil {
			logger.Errorln(err.Error())
			return nil, err
		}

		for i, v := range row {
			if b, ok := v.([]byte); ok {
				row[i] = string(b)
			}
			cvs[i].Value = v
		}
		rowHash, err := utils.MakeRowHash(row)
		if err != nil {
			logger.Error("failed to hash row: %v", err)
			continue
		}
		q.cache.Set(rowHash, cvs)
		data = append(data, models.RowSet{
			Columns: cvs,
			Hash:    rowHash,
		})
	}

	if err := rows.Err(); err != nil {
		logger.Errorln(err.Error())
		return nil, err
	}

	return data, nil
}

func (q *Queries) GetRowCount(ctx context.Context, tableName string) (int, error) {
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", q.GetQuotedTableName(tableName))
	var count int
	err := q.db.QueryRowxContext(ctx, countQuery).Scan(&count)
	if err != nil {
		logger.Errorln(err.Error())
		return 0, err
	}

	return count, nil
}

func (q *Queries) InsertRow(ctx context.Context, tableName string, form []models.ColValue) error {
	query, args, err := q.queryBuilder.InsertRow(tableName, form)
	if err != nil {
		return err
	}

	logger.Info("Query: %s", query)
	_, err = q.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	historyMsg := fmt.Sprintf("Inserted row into table '%s'", tableName)
	q.InsertHistory(ctx, historyMsg)

	return nil
}

func (q *Queries) GetRow(ctx context.Context, tableName, hash string, offest, limit int) ([]models.ColValue, error) {
	if row := q.cache.Get(hash); row != nil {
		if cv, ok := row.([]models.ColValue); ok {
			return cv, nil
		}
	}
	colValues, err := q.ListColsMetaData(ctx, tableName)
	if err != nil {
		return nil, err
	}
	logger.Info("not found in cache! Fetching from db limit=%d offset=%d", limit, offest)
	for offest <= limit {
		colValue := make([]models.ColValue, len(colValues))
		copy(colValue, colValues)
		query, args, err := q.queryBuilder.GetRows(tableName, offest+1, offest)
		if err != nil {
			return nil, err
		}
		data, err := q.db.QueryRowxContext(ctx, query, args...).SliceScan()
		if err != nil {
			logger.Error("failed to query: %v", err)
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
		}
		for i, v := range data {
			if b, ok := v.([]byte); ok {
				data[i] = string(b)
			}
			colValue[i].Value = v
		}
		rowHash, err := utils.MakeRowHash(colValue)
		if err != nil {
			logger.Error("failed to hash row: %v", err)
			continue
		}
		if rowHash == hash {
			q.cache.Set(rowHash, colValue)
			logger.Info("found data in db: %v", data)
			return colValue, nil
		}
		offest++
	}
	return nil, ErrorNotFound
}

func (q *Queries) DeleteRow(ctx context.Context, props UpdateOrDeleteRowProps) error {
	var rowValues []models.ColValue
	if cached := q.cache.Get(props.Hash); cached != nil {
		rowValues = cached.([]models.ColValue)
	} else {
		row, err := q.GetRow(ctx, props.TableName, props.Hash, props.Offset, props.Limit)
		if err != nil {
			return err
		}
		rowValues = row
	}

	query, args, err := q.queryBuilder.DeleteRow(props.TableName, rowValues, 1)
	if err != nil {
		return err
	}
	logger.Info("Query: %s", query)
	_, err = q.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	historyMsg := fmt.Sprintf("Deleted row from table '%s'", props.TableName)
	q.InsertHistory(ctx, historyMsg)

	return nil
}

type UpdateOrDeleteRowProps struct {
	TableName string
	Values    []models.ColValue
	Hash      string
	Limit     int
	Offset    int
}

func (q *Queries) UpdateRow(ctx context.Context, props UpdateOrDeleteRowProps) error {
	query, args, err := q.queryBuilder.UpdateRow(props.TableName, props.Values)
	logger.Info("Query to Update : %s", query)
	if err != nil {
		return err
	}
	_, err = q.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	historyMsg := fmt.Sprintf("Updated row in table '%s'", props.TableName)
	q.InsertHistory(ctx, historyMsg)

	return nil
}

type CreateTableProps struct {
	TableName string            `json:"tableName"`
	Inputs    []models.ColValue `json:"inputs"`
}

func (q *Queries) CreateTable(ctx context.Context, props CreateTableProps) error {
	query, err := q.queryBuilder.CreateTable(props.TableName, props.Inputs)
	if err != nil {
		return err
	}
	logger.Info("CREATE Query: %s", query)
	result, err := q.db.ExecContext(ctx, query)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		logger.Errorln(err)
		return err
	}

	historyMsg := fmt.Sprintf("Created table '%s'", props.TableName)
	q.InsertHistory(ctx, historyMsg)

	// TODO: get table info and add to q.Tables
	// temp refresh
	if _, err := q.ListTables(ctx); err != nil {
		logger.Errorln(err)
	}

	return nil
}

func (q *Queries) DeleteTable(ctx context.Context, tableName string) error {
	query := q.queryBuilder.DeleteTable(tableName)
	logger.Info("Query: %s", query)
	_, err := q.db.ExecContext(ctx, query)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	historyMsg := fmt.Sprintf("Dropped table '%s'", tableName)
	q.InsertHistory(ctx, historyMsg)

	return nil
}

func (q *Queries) GetDriver() configs.Driver {
	return q.driver
}

func (q *Queries) RunSQLQuery(ctx context.Context, query string) (*models.RunSQLQueryOutput, error) {
	rows, err := q.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(e error) {
		if e := rows.Close(); e != nil {
			logger.Errorln(e)
		}
	}(err)

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result [][]any
	for rows.Next() {
		values, err := rows.SliceScan()
		if err != nil {
			return nil, err
		}
		row := make([]any, len(values))
		for i, v := range values {
			if b, ok := v.([]byte); ok {
				row[i] = string(b)
			} else {
				row[i] = v
			}
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &models.RunSQLQueryOutput{
		Columns: cols,
		Rows:    result,
	}, nil
}
