// Package queries provides a set of SQL queries based on the driver
package queries

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database"
	"github.com/biisal/rowsql/internal/database/models"
	"github.com/biisal/rowsql/internal/logger"
)

var ErrUnknownDriver = errors.New("unknown driver")

type Builder struct {
	driver   configs.Driver
	maxLimit int
}

func NewBuilder(driver configs.Driver, maxLimit int) *Builder {
	return &Builder{
		driver:   driver,
		maxLimit: maxLimit,
	}
}

func (b *Builder) CheckTableExitsQuery(tableName string) (string, []any, error) {
	args := []any{tableName}
	switch b.driver {
	case configs.DriverMySQL:
		return `SELECT 1 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?`, args, nil
	case configs.DriverPostgres:
		return `SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1`, args, nil
	case configs.DriverSQLite:
		return `SELECT 1 FROM sqlite_master WHERE type = 'table' AND name = ?`, args, nil
	}
	logger.Error("error : %s , %s", ErrUnknownDriver.Error(), b.driver)
	return "", nil, ErrUnknownDriver
}

const postgresColumnsListsQuery = `
SELECT
    c.column_name,
    c.data_type,
    (c.column_default IS NOT NULL) AS default_value,
    COALESCE(
        bool_or(tc.constraint_type IN ('UNIQUE', 'PRIMARY KEY')),
        false
    ) AS is_unique,
    (
        c.is_identity = 'YES'
        OR c.column_default LIKE 'nextval(%'
    ) AS is_auto_increment
FROM information_schema.columns c
LEFT JOIN information_schema.key_column_usage kcu
    ON c.table_name = kcu.table_name
    AND c.column_name = kcu.column_name
    AND c.table_schema = kcu.table_schema
LEFT JOIN information_schema.table_constraints tc
    ON kcu.constraint_name = tc.constraint_name
    AND kcu.table_schema = tc.table_schema
WHERE c.table_name = $1
GROUP BY
    c.column_name,
    c.data_type,
    c.ordinal_position,
    c.is_identity,
    c.column_default
ORDER BY c.ordinal_position;
`

const mysqlColumnsListsQuery = `
SELECT
    c.column_name,
    c.data_type,
    (c.column_default IS NOT NULL) AS default_value,
    COALESCE(
        MAX(CASE
            WHEN tc.constraint_type IN ('UNIQUE', 'PRIMARY KEY') THEN 1
            ELSE 0
        END) = 1,
        false
    ) AS is_unique,
    (c.extra LIKE '%auto_increment%') AS is_auto_increment
FROM information_schema.columns c
LEFT JOIN information_schema.key_column_usage kcu
    ON c.table_name = kcu.table_name
    AND c.column_name = kcu.column_name
    AND c.table_schema = kcu.table_schema
LEFT JOIN information_schema.table_constraints tc
    ON kcu.constraint_name = tc.constraint_name
    AND kcu.table_schema = tc.table_schema
WHERE c.table_name = ?
  AND c.table_schema = DATABASE()
GROUP BY
    c.column_name,
    c.data_type,
    c.ordinal_position,
    c.extra,
    c.column_default
ORDER BY c.ordinal_position;
`

const sqliteColumnsListQuery = `
SELECT
    p.name AS column_name,
    p.type AS data_type,
    (p.dflt_value IS NOT NULL) AS default_value,
    CASE
        WHEN p.pk = 1 THEN 1
        WHEN EXISTS (
            SELECT 1
            FROM pragma_index_list(?) il
            JOIN pragma_index_info(il.name) ii
                ON ii.name = p.name
            WHERE il."unique" = 1
        ) THEN 1
        ELSE 0
    END AS is_unique,
    CASE
        WHEN p.pk = 1
             AND lower(p.type) = 'integer'
        THEN 1
        ELSE 0
    END AS is_auto_increment
FROM pragma_table_info(?) AS p;
`

func (b *Builder) ColumnsList(tableName string) (string, []any, error) {
	switch b.driver {
	case configs.DriverPostgres:
		return postgresColumnsListsQuery, []any{tableName}, nil
	case configs.DriverMySQL:
		return mysqlColumnsListsQuery, []any{tableName}, nil
	case configs.DriverSQLite:
		return sqliteColumnsListQuery, []any{tableName, tableName}, nil
	}

	logger.Error("error : %s , %s", ErrUnknownDriver.Error(), b.driver)
	return "", nil, ErrUnknownDriver
}

const postgresMySQLTablesListQuery = `
SELECT
  table_schema,
  table_name
FROM information_schema.tables
WHERE table_type = 'BASE TABLE'
  AND table_schema NOT IN (
    'pg_catalog',
    'information_schema',
    'mysql',
    'performance_schema',
    'sys'
  )
ORDER BY table_schema, table_name;
`

const sqliteTablesListQuery = `
SELECT
  '' AS table_schema,
  name AS table_name
FROM sqlite_master
WHERE type = 'table'
  AND name NOT LIKE 'sqlite_%'
ORDER BY name;
`

func (b *Builder) ListTables() (string, error) {
	switch b.driver {
	case configs.DriverPostgres, configs.DriverMySQL:
		return postgresMySQLTablesListQuery, nil
	case configs.DriverSQLite:
		return sqliteTablesListQuery, nil
	}

	logger.Error("error : %s , %s", ErrUnknownDriver.Error(), b.driver)
	return "", ErrUnknownDriver
}

func (b *Builder) ListRows(tableName, orderCol, orderBy string, limit, offset int) (string, []any, error) {
	if tableName == "" {
		return "", nil, apperr.ErrorEmptyTableName
	}
	if limit < 0 || offset < 0 {
		return "", nil, apperr.ErrorInvalidPagination
	}
	if limit > b.maxLimit {
		return "", nil, apperr.ErrorLimitTooLarge(b.maxLimit)
	}
	tableName, err := b.getQuotedTableName(tableName)
	if err != nil {
		logger.Errorln(err.Error())
		return "", nil, err
	}
	parts := []string{fmt.Sprintf("SELECT * FROM %s", tableName)}
	if orderCol != "" {
		order := "ASC"
		if strings.ToLower(orderBy) == "desc" {
			order = "DESC"
		}
		parts = append(parts, fmt.Sprintf("ORDER BY %s %s", orderCol, order))
	}

	args := []any{}
	if limit > 0 {
		ph, err := b.placeHolder(1)
		if err != nil {
			return "", nil, err
		}
		parts = append(parts, fmt.Sprintf("LIMIT %s", ph))
		args = append(args, limit)
	}
	if offset > 0 {
		placeholder, err := b.placeHolder(1)
		if err != nil {
			return "", nil, err
		}
		if limit > 0 {
			placeholder, err = b.placeHolder(2)
			if err != nil {
				return "", nil, err
			}
		}
		parts = append(parts, fmt.Sprintf("OFFSET %s", placeholder))
		args = append(args, offset)
	}

	return strings.Join(parts, " "), args, nil
}

func (b *Builder) InsertRow(tableName string, form []models.RowItem) (string, []any, error) {
	if err := b.checkValidDriver(); err != nil {
		return "", nil, err
	}
	if tableName == "" {
		return "", nil, apperr.ErrorEmptyTableName
	}
	columns := make([]string, 0, len(form))
	placeholders := make([]string, 0, len(form))
	args := make([]any, 0, len(form))

	isPostgresOrSqLite := b.driver == configs.DriverPostgres || b.driver == configs.DriverSQLite

	paramIndex := 1

	ph := "?"

	seen := make(map[string]bool)

	for _, field := range form {
		if _, ok := seen[field.ColumnName]; ok {
			return "", nil, apperr.ErrorDuplicateColumn
		}
		seen[field.ColumnName] = true
		columns = append(columns, field.ColumnName)
		if isPostgresOrSqLite {
			ph = "$" + strconv.Itoa(paramIndex)
		}
		placeholders = append(placeholders, ph)

		if field.Type == "json" {
			var jsonVal any
			if err := json.Unmarshal([]byte(field.Value), &jsonVal); err != nil {
				logger.Errorln(err)
				var syntaxErr *json.SyntaxError
				if errors.As(err, &syntaxErr) || errors.Is(err, io.ErrUnexpectedEOF) {
					return "", nil, apperr.ErrorInvalidJSON
				}
				return "", nil, err

			}
			logger.Info("jsonVal: %v", jsonVal)
			args = append(args, jsonVal)
		} else {
			args = append(args, field.Value)
		}

		paramIndex++
	}

	qColumns := ""
	qPlaceholders := ""

	if len(columns) > 0 {
		qColumns = strings.Join(columns, ", ")
		qPlaceholders = strings.Join(placeholders, ", ")
	}

	var query string
	tableName, err := b.getQuotedTableName(tableName)
	if err != nil {
		return "", nil, err
	}
	if qColumns == "" {
		if b.driver == configs.DriverMySQL {
			query = fmt.Sprintf("INSERT INTO %s () VALUES ()", tableName)
		} else {
			query = fmt.Sprintf("INSERT INTO %s DEFAULT VALUES", tableName)
		}
	} else {
		query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, qColumns, qPlaceholders)
	}
	logger.Info("Query: %s", query)
	return query, args, nil
}

func (b *Builder) GetRows(tableName string, limit, offset int) (string, []any, error) {
	if err := b.checkValidDriver(); err != nil {
		return "", nil, err
	}
	if limit <= 0 || offset < 0 {
		return "", nil, apperr.ErrorInvalidPagination
	}

	tableName = strings.TrimSpace(tableName)
	if tableName == "" {
		return "", nil, apperr.ErrorInvalidTableName
	}

	args := []any{limit}
	ph, err := b.placeHolder(1)
	if err != nil {
		return "", nil, err
	}
	parts := []string{fmt.Sprintf("SELECT * FROM %s LIMIT %s", tableName, ph)}
	if offset > 0 {
		ph, err = b.placeHolder(2)
		if err != nil {
			return "", nil, err
		}
		parts = append(parts, fmt.Sprintf("OFFSET %s", ph))
		args = append(args, offset)
	}

	return strings.Join(parts, " "), args, nil
}

func (b *Builder) WhereCluse(cols []models.ListDataCol, rows []any, argsIdx int) (string, []any, error) {
	if len(cols) != len(rows) {
		return "", nil, apperr.ErrorNotSameRowColsSize
	}

	var mixed []string
	var args []any
	for i, val := range cols {
		ph, err := b.placeHolder(argsIdx + i)
		if err != nil {
			return "", nil, err
		}

		if val.IsUnique {
			return fmt.Sprintf("%s=%s", val.ColumnName, ph), []any{rows[i]}, nil
		}
		colVal := rows[i]
		if val.DataType == "json" {
			var jsonVal map[string]any
			if err := json.Unmarshal([]byte(colVal.(string)), &jsonVal); err != nil {
				logger.Errorln(err)
				return "", nil, err
			}
			colVal = jsonVal
			mixed = append(mixed, fmt.Sprintf("%s::jsonb @> %s::jsonb", val.ColumnName, ph))
		} else {
			mixed = append(mixed, fmt.Sprintf("%s=%s", val.ColumnName, ph))
		}
		args = append(args, colVal)
	}
	return strings.Join(mixed, " AND "), args, nil
}

func (b *Builder) DeleteRow(tableName string, columns []models.ListDataCol, rows []any, argIdx int) (string, []any, error) {
	clause, args, err := b.WhereCluse(columns, rows, argIdx)
	logger.Info("Clause: %s", clause)
	if err != nil {
		return "", nil, err
	}
	query := ""
	switch b.driver {
	case configs.DriverMySQL:
		query = fmt.Sprintf("DELETE FROM %s WHERE %s LIMIT 1", tableName, clause)
	case configs.DriverPostgres:
		query = fmt.Sprintf("DELETE FROM %s WHERE ctid IN (SELECT ctid FROM %s WHERE %s LIMIT 1)", tableName, tableName, clause)
	case configs.DriverSQLite:
		query = fmt.Sprintf("DELETE FROM %s WHERE rowid IN (SELECT rowid FROM %s WHERE %s LIMIT 1)", tableName, tableName, clause)
	default:
		return "", nil, apperr.ErrorInvalidDriver

	}
	return query, args, nil
}

func (b *Builder) UpdateRow(tableName string, form []models.RowItem, columns []models.ListDataCol, row []any) (string, []any, error) {
	parts := make([]string, 0, len(form))
	ph := "?"
	index := 1
	args := make([]any, 0, len(form))
	for _, v := range form {
		if b.driver == configs.DriverPostgres {
			ph = "$" + strconv.Itoa(index)
		}
		parts = append(parts, fmt.Sprintf("%s=%s", v.ColumnName, ph))
		args = append(args, v.Value)
		index++
	}
	updateQuery := strings.Join(parts, ",")
	whereClause, whereClauseArgs, err := b.WhereCluse(columns, row, len(args)+1)
	if err != nil {
		return "", nil, err
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tableName, updateQuery, whereClause)
	return query, append(args, whereClauseArgs...), nil
}

func (b *Builder) CreateTable(tableName string, inputs []database.Input) (string, error) {
	logger.Info("Building create table query")
	columnDefs := make([]string, 0, len(inputs))
	for _, input := range inputs {
		if input.ColName == "" {
			continue
		}
		formattedColDef, err := b.formatColumnDefinition(input)
		if err != nil {
			return "", err
		}
		columnDefs = append(columnDefs, formattedColDef)
	}
	logger.Info("Query: %s", strings.Join(columnDefs, ", "))

	validName := regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	if !validName.MatchString(tableName) {
		logger.Error("invalid table name! Only alphanumeric and _ are allowed table_name=%s", tableName)
		return "", errors.New("invalid table name! Only alphanumeric and _ are allowed")
	}
	parts := strings.Join(columnDefs, ", ")

	query := fmt.Sprintf("CREATE TABLE %s (%s) ;", tableName, parts)
	return query, nil
}

func (b *Builder) DeleteTable(tableName string) string {
	return fmt.Sprintf("DROP TABLE %s", tableName)
}
