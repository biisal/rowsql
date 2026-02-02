package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/database"
	"github.com/biisal/rowsql/internal/database/models"
	"github.com/biisal/rowsql/internal/logger"
	"github.com/biisal/rowsql/internal/utils"
)

func ErrorInvalidTable(tableName string) error {
	return fmt.Errorf("table %s not found", tableName)
}

var ErrorNotFound = errors.New("not found")

func (q *Queries) GetQuotedTableName(tableName string) string {
	switch q.GetDriver() {
	case configs.DriverMySQL:
		tableName = fmt.Sprintf("`%s`", tableName)
	case configs.DriverPostgres:
		tableName = fmt.Sprintf("\"%s\"", tableName)
	case configs.DriverSQLite:
		tableName = fmt.Sprintf("\"%s\"", tableName)
	}
	return tableName
}

func (q *Queries) CheckTableExitsInDB(ctx context.Context, tableName string) error {
	query, args, err := q.queryBuilder.CheckTableExitsQuery(tableName)
	if err != nil {
		return err
	}
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

func (q *Queries) ListCols(ctx context.Context, tableName string) ([]models.ListDataCol, error) {
	query, args, err := q.queryBuilder.ColumnsList(tableName)
	if err != nil {
		logger.Error("failed to build query : %v", err)
		return nil, err
	}
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
	var items []models.ListDataCol
	for rows.Next() {
		var i models.ListDataCol
		if err := rows.Scan(&i.ColumnName, &i.DataType, &i.HasDefault, &i.IsUnique, &i.HasAutoIncrement); err != nil {
			logger.Error("failed to scan rows in list cols: %v", err)
			return nil, err
		}
		i.InputType = utils.GetInputType(i.DataType)
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		logger.Error("failed to scan rows: %v", err)
		return nil, err
	}
	return items, nil
}

func (q *Queries) ListTables(ctx context.Context) ([]models.ListTablesRow, error) {
	query, err := q.queryBuilder.ListTables()
	if err != nil {
		logger.Error("failed to build query : %v", err)
		return nil, err
	}
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

func (q *Queries) ListRows(ctx context.Context, props models.ListDataProps) (models.ListDataRow, error) {
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
		if err := rows.Close(); err != nil {
			logger.Errorln(err)
		}
	}()
	data := make(models.ListDataRow, 0)
	for rows.Next() {
		row, err := rows.SliceScan()
		if err != nil {
			logger.Errorln(err.Error())
			return nil, err
		}

		for i, v := range row {
			if b, ok := v.([]byte); ok {
				row[i] = string(b)
			}
		}
		rowHash, err := utils.MakeRowHash(row)
		if err != nil {
			logger.Error("failed to hash row: %v", err)
			continue
		}
		q.cache.Set(rowHash, row)
		row = append([]any{rowHash}, row...)
		data = append(data, row)
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

func (q *Queries) InsertRow(ctx context.Context, props models.InsertDataProps) error {
	query, args, err := q.queryBuilder.InsertRow(props.TableName, props.Values)
	if err != nil {
		return err
	}

	logger.Info("Query: %s", query)
	_, err = q.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	historyMsg := fmt.Sprintf("Inserted row into table '%s'", props.TableName)
	q.InsertHistory(ctx, historyMsg)

	return nil
}

func (q *Queries) GetRow(ctx context.Context, tableName, hash string, offest, limit int) ([]any, error) {
	if row := q.cache.Get(hash); row != nil {
		logger.Info("found data in cache: %v", row)
		return row, nil
	}
	logger.Info("not found in cache! Fetching from db limit=%d offset=%d", limit, offest)
	for offest <= limit {
		query, args, err := q.queryBuilder.GetRows(tableName, offest+1, offest)
		if err != nil {
			return nil, err
		}
		logger.Info("Query: %s offset=%d tableName=%s", query, offest, tableName)
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
		}
		rowHash, err := utils.MakeRowHash(data)
		if err != nil {
			logger.Error("failed to hash row: %v", err)
			continue
		}
		if rowHash == hash {
			q.cache.Set(rowHash, data)
			logger.Info("found data in db: %v", data)
			return data, nil
		}
		offest++
	}
	return nil, ErrorNotFound
}

func (q *Queries) DeleteRow(ctx context.Context, props UpdateOrDeleteRowProps) error {
	row := q.cache.Get(props.Hash)
	if row == nil {
		var err error
		row, err = q.GetRow(ctx, props.TableName, props.Hash, props.Offset, props.Limit)
		if err != nil {
			return err
		}
	}
	cols, err := q.ListCols(ctx, props.TableName)
	if err != nil {
		return err
	}
	query, args, err := q.queryBuilder.DeleteRow(props.TableName, cols, row, 1)
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
	Values    []models.RowItem
	Hash      string
	Limit     int
	Offset    int
}

func (q *Queries) UpdateRow(ctx context.Context, props UpdateOrDeleteRowProps) error {
	row := q.cache.Get(props.Hash)
	if row == nil {
		var err error
		row, err = q.GetRow(ctx, props.TableName, props.Hash, props.Offset, props.Limit)
		if err != nil {
			return err
		}
	}

	cols, err := q.ListCols(ctx, props.TableName)
	if err != nil {
		return err
	}

	query, args, err := q.queryBuilder.UpdateRow(props.TableName, props.Values, cols, row)
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
	TableName string           `json:"tableName"`
	Inputs    []database.Input `json:"inputs"`
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
