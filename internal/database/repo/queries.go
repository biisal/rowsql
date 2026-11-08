package repo

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database/models"
	"github.com/biisal/rowsql/internal/database/queries"
	"github.com/biisal/rowsql/internal/database/rowkey"
	"github.com/biisal/rowsql/internal/logger"
	"github.com/biisal/rowsql/internal/utils"
	"github.com/jmoiron/sqlx"
)

func ErrorInvalidTable(tableName string) error {
	return fmt.Errorf("table %s not found", tableName)
}

func closeRows(rows interface{ Close() error }) {
	if err := rows.Close(); err != nil {
		logger.Errorln(err)
	}
}

func (q *Queries) CheckTableExistsInDB(ctx context.Context, tableName string) error {
	query, args := q.queryBuilder.CheckTableExists(tableName)
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error("failed to query: %v", err)
		return err
	}
	defer closeRows(rows)
	if !rows.Next() {
		return ErrorInvalidTable(tableName)
	}
	return nil
}

func (q *Queries) ListColsMetaData(ctx context.Context, tableName string) ([]models.ColValue, error) {
	if cached, found := q.columnCache.get(tableName); found {
		return cached, nil
	}

	query, args := q.queryBuilder.ListColumns(tableName)
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error("failed to query: %v", err)
		return nil, err
	}
	defer closeRows(rows)

	var items []models.ColValue
	for rows.Next() {
		var column models.ColValue
		err := rows.Scan(
			&column.ColumnName,
			&column.ColumnType.DataType,
			&column.ColumnType.HasDefault,
			&column.ColumnType.IsUnique,
			&column.ColumnType.HasAutoIncrement,
			&column.ColumnType.IsNull,
			&column.ColumnType.PkPosition,
		)
		if err != nil {
			logger.Error("failed to scan rows in list cols: %v", err)
			return nil, err
		}
		column.ColumnType.IsPk = column.ColumnType.PkPosition > 0
		column.ColumnType.InputType = utils.GetInputType(column.ColumnType.DataType)
		items = append(items, column)
	}
	if err := rows.Err(); err != nil {
		logger.Error("failed to scan rows: %v", err)
		return nil, err
	}

	q.columnCache.set(tableName, items)
	return items, nil
}

func (q *Queries) ListTables(ctx context.Context) ([]models.ListTablesRow, error) {
	rows, err := q.db.QueryContext(ctx, q.queryBuilder.ListTables())
	if err != nil {
		return nil, err
	}
	defer closeRows(rows)

	var items []models.ListTablesRow
	for rows.Next() {
		var table models.ListTablesRow
		if err := rows.Scan(&table.TableSchema, &table.TableName); err != nil {
			logger.Error("failed to scan rows: %v", err)
			return nil, err
		}
		if table.TableName == historyTableName {
			continue
		}
		items = append(items, table)
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
	columns, err := q.ListColsMetaData(ctx, props.TableName)
	if err != nil {
		return nil, err
	}
	if props.Column != "" && !hasColumn(columns, props.Column) {
		return nil, apperr.ErrorInvalidColumn
	}

	plan := q.planRowKey(columns)
	query, args, err := q.queryBuilder.ListRows(props.TableName, queries.ListRowsOptions{
		OrderColumn:          props.Column,
		Descending:           props.Order == "DESC" || props.Order == "desc",
		Limit:                props.Limit,
		Offset:               props.Offset,
		DefaultOrderColumns:  plan.defaultOrderColumns(),
		IncludePhysicalRowID: plan.kind == rowkey.KindPhysical,
	})
	if err != nil {
		return nil, err
	}

	rows, err := q.db.QueryxContext(ctx, query, args...)
	if err != nil {
		logger.Errorln(err.Error())
		return nil, err
	}
	defer closeRows(rows)

	data := make([]models.RowSet, 0)
	for rows.Next() {
		rowColumns, physicalRowID, err := scanRow(rows, columns)
		if err != nil {
			return nil, err
		}
		key, err := plan.buildKey(valuesByName(rowColumns), physicalRowID).Encode()
		if err != nil {
			return nil, err
		}
		data = append(data, models.RowSet{Columns: rowColumns, Key: key})
	}
	if err := rows.Err(); err != nil {
		logger.Errorln(err.Error())
		return nil, err
	}
	return data, nil
}

func scanRow(rows *sqlx.Rows, schema []models.ColValue) ([]models.ColValue, any, error) {
	resultColumns, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	values, err := rows.SliceScan()
	if err != nil {
		logger.Errorln(err.Error())
		return nil, nil, err
	}

	schemaByName := make(map[string]models.ColValue, len(schema))
	for _, column := range schema {
		schemaByName[column.ColumnName] = column
	}

	var physicalRowID any
	rowColumns := make([]models.ColValue, 0, len(values))
	for index, value := range values {
		if bytes, isBytes := value.([]byte); isBytes {
			value = string(bytes)
		}
		if resultColumns[index] == queries.PhysicalRowIDAlias {
			physicalRowID = value
			continue
		}
		column := schemaByName[resultColumns[index]]
		column.ColumnName = resultColumns[index]
		column.Value = value
		rowColumns = append(rowColumns, column)
	}
	return rowColumns, physicalRowID, nil
}

func valuesByName(columns []models.ColValue) map[string]any {
	values := make(map[string]any, len(columns))
	for _, column := range columns {
		values[column.ColumnName] = column.Value
	}
	return values
}

func hasColumn(columns []models.ColValue, name string) bool {
	for _, column := range columns {
		if column.ColumnName == name {
			return true
		}
	}
	return false
}

func (q *Queries) GetRowCount(ctx context.Context, tableName string) (int, error) {
	countQuery, err := q.queryBuilder.CountRows(tableName)
	if err != nil {
		return 0, err
	}
	var count int
	if err := q.db.QueryRowxContext(ctx, countQuery).Scan(&count); err != nil {
		logger.Errorln(err.Error())
		return 0, err
	}
	return count, nil
}

func (q *Queries) InsertRow(ctx context.Context, tableName string, form []models.ColValue) error {
	schema, err := q.ListColsMetaData(ctx, tableName)
	if err != nil {
		return err
	}
	values, err := alignWithSchema(form, schema)
	if err != nil {
		return err
	}
	query, args, err := q.queryBuilder.InsertRow(tableName, values)
	if err != nil {
		return err
	}

	if _, err := q.db.ExecContext(ctx, query, args...); err != nil {
		logger.Errorln(err)
		return err
	}
	q.InsertHistory(ctx, fmt.Sprintf("Inserted row into table '%s'", tableName))
	return nil
}

func (q *Queries) GetRow(ctx context.Context, tableName, encodedKey string) ([]models.ColValue, error) {
	key, err := rowkey.Decode(encodedKey)
	if err != nil {
		return nil, err
	}
	schema, err := q.ListColsMetaData(ctx, tableName)
	if err != nil {
		return nil, err
	}
	query, args, err := q.queryBuilder.SelectRow(tableName, key)
	if err != nil {
		return nil, err
	}

	rows, err := q.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer closeRows(rows)

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return nil, apperr.ErrorRowNotFound
	}
	rowColumns, _, err := scanRow(rows, schema)
	return rowColumns, err
}

func (q *Queries) DeleteRow(ctx context.Context, tableName, encodedKey string) error {
	key, err := rowkey.Decode(encodedKey)
	if err != nil {
		return err
	}
	query, args, err := q.queryBuilder.DeleteRow(tableName, key)
	if err != nil {
		return err
	}
	if err := q.execExpectingOneRow(ctx, query, args); err != nil {
		return err
	}
	q.InsertHistory(ctx, fmt.Sprintf("Deleted row from table '%s'", tableName))
	return nil
}

func (q *Queries) UpdateRow(ctx context.Context, tableName, encodedKey string, form []models.ColValue) error {
	key, err := rowkey.Decode(encodedKey)
	if err != nil {
		return err
	}
	schema, err := q.ListColsMetaData(ctx, tableName)
	if err != nil {
		return err
	}
	changes, err := alignWithSchema(form, schema)
	if err != nil {
		return err
	}
	query, args, err := q.queryBuilder.UpdateRow(tableName, changes, key)
	if err != nil {
		return err
	}
	if err := q.execExpectingOneRow(ctx, query, args); err != nil {
		return err
	}
	q.InsertHistory(ctx, fmt.Sprintf("Updated row in table '%s'", tableName))
	return nil
}

func (q *Queries) execExpectingOneRow(ctx context.Context, query string, args []any) error {
	result, err := q.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return apperr.ErrorRowNotFound
	}
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
	if _, err := q.db.ExecContext(ctx, query); err != nil {
		logger.Errorln(err)
		return err
	}

	q.columnCache.clear()
	q.InsertHistory(ctx, fmt.Sprintf("Created table '%s'", props.TableName))
	return nil
}

func (q *Queries) DeleteTable(ctx context.Context, tableName string) error {
	if _, err := q.db.ExecContext(ctx, q.queryBuilder.DeleteTable(tableName)); err != nil {
		logger.Errorln(err)
		return err
	}

	q.columnCache.clear()
	q.InsertHistory(ctx, fmt.Sprintf("Dropped table '%s'", tableName))
	return nil
}

func (q *Queries) GetDriver() configs.Driver {
	return q.driver
}

const (
	maxSQLEditorRows       = 1000
	historyQueryPreviewLen = 100
)

var readOnlyStatementKeywords = map[string]bool{
	"select": true, "with": true, "explain": true, "show": true, "pragma": true, "describe": true,
}

func (q *Queries) RunSQLQuery(ctx context.Context, query string) (*models.RunSQLQueryOutput, error) {
	rows, err := q.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer closeRows(rows)

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	output := &models.RunSQLQueryOutput{Columns: cols}
	for rows.Next() {
		if len(output.Rows) == maxSQLEditorRows {
			output.Truncated = true
			break
		}
		values, err := rows.SliceScan()
		if err != nil {
			return nil, err
		}
		output.Rows = append(output.Rows, displayValues(values))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	q.afterUserSQL(ctx, query)
	return output, nil
}

func displayValues(values []any) []any {
	row := make([]any, len(values))
	for index, value := range values {
		switch typed := value.(type) {
		case []byte:
			row[index] = string(typed)
		case float64:
			if math.IsNaN(typed) || math.IsInf(typed, 0) {
				row[index] = nil
			} else {
				row[index] = typed
			}
		default:
			row[index] = value
		}
	}
	return row
}

func (q *Queries) afterUserSQL(ctx context.Context, query string) {
	q.columnCache.clear()

	fields := strings.Fields(query)
	if len(fields) == 0 || readOnlyStatementKeywords[strings.ToLower(fields[0])] {
		return
	}
	preview := strings.Join(fields, " ")
	if len(preview) > historyQueryPreviewLen {
		preview = preview[:historyQueryPreviewLen] + "..."
	}
	q.InsertHistory(ctx, "Ran SQL: "+preview)
}
