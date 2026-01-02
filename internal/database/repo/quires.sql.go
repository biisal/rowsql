package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"regexp"

	"github.com/biisal/db-gui/internal/database"
	"github.com/biisal/db-gui/internal/utils"
)

const ErrorInvalidTable = "invalid table name"
const ErrorNotFound = "not found"

func (q *Queries) isAllowedTable(table string) bool {
	for _, t := range q.Tables {
		if t.TableName == table {
			return true
		}
	}
	return false
}

func (q *Queries) ListCols(ctx context.Context, tableName string) ([]ListDataCol, error) {
	query, args := colsQuery(q.db.DriverName(), tableName)
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		slog.Error("failed to query", "err", err)
		return nil, err
	}
	defer rows.Close()
	var items []ListDataCol
	for rows.Next() {
		var i ListDataCol
		if err := rows.Scan(&i.ColumnName, &i.DataType, &i.IsUnique, &i.HasAutoIncrement); err != nil {
			slog.Error("failed to scan rows in list cols", "err", err)
			return nil, err
		}
		i.InputType = utils.GetInputType(i.DataType)
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		slog.Error("failed to scan rows", "err", err)
		return nil, err
	}
	return items, nil
}

type ListTablesRow struct {
	TableSchema string `json:"tableSchema"`
	TableName   string `json:"tableName"`
}

func (q *Queries) ListTables(ctx context.Context) ([]ListTablesRow, error) {
	rows, err := q.db.QueryContext(ctx, tablesQuery(q.db.DriverName()))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListTablesRow
	for rows.Next() {
		var i ListTablesRow
		if err := rows.Scan(&i.TableSchema, &i.TableName); err != nil {
			slog.Error("failed to scan rows", "err", err)
			return nil, err
		}
		if i.TableName == historyTableName {
			continue
		}
		items = append([]ListTablesRow{i}, items...)
	}
	if err := rows.Err(); err != nil {
		slog.Error("failed to scan rows", "err", err)
		return nil, err
	}
	slog.Info("Tables", "tables", items)
	q.Tables = items
	return items, nil
}

func (q *Queries) ListRows(ctx context.Context, props ListDataProps) (ListDataRow, error) {
	if !q.isAllowedTable(props.TableName) {
		return nil, fmt.Errorf(ErrorInvalidTable)
	}
	query := fmt.Sprintf("SELECT * FROM %s LIMIT $1 OFFSET $2", props.TableName)
	rows, err := q.db.QueryxContext(ctx, query, props.Limit, props.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	data := make(ListDataRow, 0)
	for rows.Next() {
		row, err := rows.SliceScan()
		if err != nil {
			return nil, err
		}

		for i, v := range row {
			if b, ok := v.([]byte); ok {
				row[i] = string(b)
			}
		}
		row_hash := utils.MakeRowHash(row)
		q.cache.Set(row_hash, row)
		row = append([]any{row_hash}, row...)
		data = append(data, row)
	}

	return data, nil
}

func (q *Queries) InsertRow(ctx context.Context, props InsertDataProps) error {
	if !q.isAllowedTable(props.TableName) {
		return fmt.Errorf(ErrorInvalidTable)
	}

	qParts, err := buildQueryParts(props.Values, q.db.DriverName())
	if err != nil {
		return err
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", props.TableName, qParts.Columns, qParts.Placeholders)
	slog.Info("Query", "query", query)
	slog.Info("Args", "args", qParts.Args)
	_, err = q.db.ExecContext(ctx, query, qParts.Args...)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	historyMsg := fmt.Sprintf("Inserted row into table '%s'", props.TableName)
	q.InsertHistory(ctx, historyMsg)

	return nil
}

func (q *Queries) GetRow(ctx context.Context, tableName, hash string, offest, limit int) ([]any, error) {
	if !q.isAllowedTable(tableName) {
		return nil, fmt.Errorf(ErrorInvalidTable)
	}

	if row := q.cache.Get(hash); row != nil {
		slog.Info("found data in cache", "data", row)
		return row, nil
	}
	slog.Info("not found in cache! Fetching from db", "limit", limit, "offset", offest)
	for offest <= limit {
		query := fmt.Sprintf("SELECT * FROM %s LIMIT $1 OFFSET $2", tableName)
		slog.Info("Query", "query", query, "offset", offest, "tableName", tableName)
		data, err := q.db.QueryRowxContext(ctx, query, offest+1, offest).SliceScan()
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
			slog.Info("row not found! Continuing to next row")
		}
		for i, v := range data {
			if b, ok := v.([]byte); ok {
				data[i] = string(b)
			}
		}
		row_hash := utils.MakeRowHash(data)
		if row_hash == hash {
			q.cache.Set(row_hash, data)
			slog.Info("found data in db", "data", data)
			return data, nil
		}
		offest++
	}
	return nil, fmt.Errorf(ErrorNotFound)
}

func (q *Queries) DeleteRow(ctx context.Context, props UpdateOrDeleteRowProps) error {
	if !q.isAllowedTable(props.TableName) {
		return fmt.Errorf(ErrorInvalidTable)
	}
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
	clause, args, err := buildQueryWhereClause(cols, row, q.db.DriverName(), 1)
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", props.TableName, clause)
	slog.Info("Query", "query", query)
	_, err = q.db.ExecContext(ctx, query, args...)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	historyMsg := fmt.Sprintf("Deleted row from table '%s'", props.TableName)
	q.InsertHistory(ctx, historyMsg)

	return nil
}

type UpdateOrDeleteRowProps struct {
	TableName string
	Values    map[string]FormValue
	Hash      string
	Limit     int
	Offset    int
}

func (q *Queries) UpdateRow(ctx context.Context, props UpdateOrDeleteRowProps) error {
	if !q.isAllowedTable(props.TableName) {
		return fmt.Errorf(ErrorInvalidTable)
	}
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

	data, args := buildUpdateParts(props.Values, q.db.DriverName())

	whereClause, wcArgs, err := buildQueryWhereClause(cols, row, q.db.DriverName(), len(args)+1)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", props.TableName, data, whereClause)
	slog.Info("Query", "query", query)
	fullargs := append(args, wcArgs...)
	_, err = q.db.ExecContext(ctx, query, fullargs...)
	if err != nil {
		slog.Error(err.Error())
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
	parts, err := buildCreateTableQuiry(q.driver, props.Inputs)
	if err != nil {
		return err
	}
	validName := regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	if !validName.MatchString(props.TableName) {
		slog.Error("invalid table name! Only alphanumeric and _ are allowed", "table_name", props.TableName)
		return errors.New("invalid table name! Only alphanumeric and _ are allowed")
	}
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s) ;", props.TableName, parts)
	slog.Info("CREATE Query", "query", query)
	result, err := q.db.ExecContext(ctx, query)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	historyMsg := fmt.Sprintf("Created table '%s'", props.TableName)
	q.InsertHistory(ctx, historyMsg)

	return nil

}

func (q *Queries) DeleteTable(ctx context.Context, tableName string) error {
	tables, err := q.ListTables(ctx)
	found := false
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	for _, t := range tables {
		if t.TableName == tableName {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("table %s not found", tableName)
	}
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	slog.Info("Query", "query", query)
	_, err = q.db.ExecContext(ctx, query)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	historyMsg := fmt.Sprintf("Dropped table '%s'", tableName)
	q.InsertHistory(ctx, historyMsg)

	return nil
}

func (q *Queries) GetDriver() string {
	return q.db.DriverName()
}
