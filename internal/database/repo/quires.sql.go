package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/biisal/db-gui/configs"
	"github.com/biisal/db-gui/internal/database"
	"github.com/biisal/db-gui/internal/logger"
	"github.com/biisal/db-gui/internal/utils"
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
	query := ""
	switch q.db.DriverName() {
	case configs.DriverMySQL:
		query = `SELECT 1 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?`
	case configs.DriverPostgres:
		query = `SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1`
	case configs.DriverSQLite:
		query = `SELECT 1 FROM sqlite_master WHERE type = 'table' AND name = ?`
	default:
		logger.Info("%s driver=%s", fallbackMsg, q.db.DriverName())
		return errors.New("unsupported driver")

	}
	rows, err := q.db.QueryContext(ctx, query, tableName)
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

func (q *Queries) ListCols(ctx context.Context, tableName string) ([]ListDataCol, error) {
	query, args := colsQuery(q.db.DriverName(), tableName)
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
	var items []ListDataCol
	for rows.Next() {
		var i ListDataCol
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

type ListTablesRow struct {
	TableSchema string `json:"tableSchema"`
	TableName   string `json:"tableName"`
}

func (q *Queries) ListTables(ctx context.Context) ([]ListTablesRow, error) {
	rows, err := q.db.QueryContext(ctx, tablesQuery(q.db.DriverName()))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Errorln(err)
		}
	}()
	var items []ListTablesRow
	for rows.Next() {
		var i ListTablesRow
		if err := rows.Scan(&i.TableSchema, &i.TableName); err != nil {
			logger.Error("failed to scan rows: %v", err)
			return nil, err
		}
		if i.TableName == historyTableName {
			continue
		}
		items = append([]ListTablesRow{i}, items...)
	}
	if err := rows.Err(); err != nil {
		logger.Error("failed to scan rows: %v", err)
		return nil, err
	}
	return items, nil
}

func placeHolder(driver string, n int) string {
	if driver == configs.DriverMySQL {
		return "?"
	}
	return fmt.Sprintf("$%d", n)
}

func (q *Queries) ListRows(ctx context.Context, props ListDataProps) (ListDataRow, error) {
	props.Column = strings.TrimSpace(props.Column)

	orderByClause := ""
	if props.Column != "" {
		order := "ASC"
		if strings.ToLower(props.Order) == "desc" {
			order = "DESC"
		}
		orderByClause = fmt.Sprintf("ORDER BY %s %s", props.Column, order)
	}

	driver := q.db.DriverName()
	props.TableName = q.GetQuotedTableName(props.TableName)
	query := fmt.Sprintf("SELECT * FROM %s %s LIMIT %s OFFSET %s", props.TableName, orderByClause, placeHolder(driver, 1), placeHolder(driver, 2))

	logger.Info("Query : %s", query)
	rows, err := q.db.QueryxContext(ctx, query, props.Limit, props.Offset)
	if err != nil {
		logger.Errorln(err.Error())
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Errorln(err)
		}
	}()
	data := make(ListDataRow, 0)
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

func (q *Queries) InsertRow(ctx context.Context, props InsertDataProps) error {
	qParts, err := buildQueryParts(props.Values, q.db.DriverName())
	if err != nil {
		return err
	}
	var query string
	tableName := q.GetQuotedTableName(props.TableName)
	if qParts.Columns == "" {
		if q.db.DriverName() == configs.DriverMySQL {
			query = fmt.Sprintf("INSERT INTO %s VALUES (%s)", tableName, qParts.Placeholders)
		} else {
			query = fmt.Sprintf("INSERT INTO %s DEFAULT VALUES", tableName)
		}
	} else {
		query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, qParts.Columns, qParts.Placeholders)
	}
	logger.Info("Query: %s", query)
	logger.Info("Args: %v", qParts.Args)
	_, err = q.db.ExecContext(ctx, query, qParts.Args...)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	historyMsg := fmt.Sprintf("Inserted row into table '%s'", tableName)
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
		query := fmt.Sprintf("SELECT * FROM %s LIMIT $1 OFFSET $2", tableName)
		logger.Info("Query: %s offset=%d tableName=%s", query, offest, tableName)
		data, err := q.db.QueryRowxContext(ctx, query, offest+1, offest).SliceScan()
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
			logger.Info("row not found! Continuing to next row")
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
	clause, args, err := buildQueryWhereClause(cols, row, q.db.DriverName(), 1)
	if err != nil {
		logger.Errorln(err.Error())
		return err
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", props.TableName, clause)
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
	Values    map[string]FormValue
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

	data, args := buildUpdateParts(props.Values, q.db.DriverName())

	whereClause, wcArgs, err := buildQueryWhereClause(cols, row, q.db.DriverName(), len(args)+1)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", props.TableName, data, whereClause)
	logger.Info("Query: %s", query)
	fullargs := append(args, wcArgs...)
	_, err = q.db.ExecContext(ctx, query, fullargs...)
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
	parts, err := buildCreateTableQuiry(q.driver, props.Inputs)
	if err != nil {
		return err
	}
	validName := regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	if !validName.MatchString(props.TableName) {
		logger.Error("invalid table name! Only alphanumeric and _ are allowed table_name=%s", props.TableName)
		return errors.New("invalid table name! Only alphanumeric and _ are allowed")
	}
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s) ;", props.TableName, parts)
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
	tables, err := q.ListTables(ctx)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	found := false
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
	logger.Info("Query: %s", query)
	_, err = q.db.ExecContext(ctx, query)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	historyMsg := fmt.Sprintf("Dropped table '%s'", tableName)
	q.InsertHistory(ctx, historyMsg)

	return nil
}

func (q *Queries) GetDriver() string {
	return q.db.DriverName()
}
