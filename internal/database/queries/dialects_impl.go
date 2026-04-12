package queries

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database/models"
	"github.com/biisal/rowsql/internal/logger"
)

// PostgresDialect implementation
type PostgresDialect struct{}

func (d *PostgresDialect) Name() configs.Driver {
	return configs.DriverPostgres
}

func quoteName(name string, format ...string) string {
	if name == "" {
		return "''"
	}

	formator := "%q"
	if len(format) > 0 {
		formator = format[0]
	}

	quoteChar := "\""
	if strings.Contains(formator, "`") {
		quoteChar = "`"
	}

	if len(name) >= 2 && name[0] == quoteChar[0] && name[len(name)-1] == quoteChar[0] {
		return name
	}

	needsQuoting := false
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_') {
			needsQuoting = true
			break
		}
	}

	if !needsQuoting {
		return name
	}

	if formator == "%q" {
		return fmt.Sprintf(formator, name)
	}

	escaped := strings.ReplaceAll(name, quoteChar, quoteChar+quoteChar)
	return fmt.Sprintf(formator, escaped)
}

func whereClause(d Dialect, cols []models.ColValues, rows []any, argsIdx int) (string, []any, error) {
	if len(cols) != len(rows) {
		return "", nil, apperr.ErrorNotSameRowColsSize
	}

	var mixed []string
	var args []any
	for i, col := range cols {
		ph, err := d.PlaceHolder(argsIdx + i)
		if err != nil {
			return "", nil, err
		}
		colName := d.QuoteName(col.Name)

		if col.IsUnique {
			ph, err = d.PlaceHolder(argsIdx)
			if err != nil {
				return "", nil, err
			}
			return fmt.Sprintf("%s=%s", colName, ph), []any{col.Value}, nil
		}
		colVal := rows[i]

		if col.Type == "json" {
			var jsonVal map[string]any
			if err := json.Unmarshal([]byte(colVal.(string)), &jsonVal); err != nil {
				logger.Errorln(err)
				return "", nil, err
			}
			colVal = jsonVal
			mixed = append(mixed, fmt.Sprintf("%s::jsonb @> %s::jsonb", colName, ph))
		} else {
			mixed = append(mixed, fmt.Sprintf("%s=%s", colName, ph))
		}
		args = append(args, colVal)
	}
	return strings.Join(mixed, " AND "), args, nil
}

func (d *PostgresDialect) CheckTableExists(tableName string) (string, []any) {
	return `SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = $1`, []any{tableName}
}

func (d *PostgresDialect) ListTables() string {
	return `
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
}

func (d *PostgresDialect) ListColumns(tableName string) (string, []any) {
	return `
SELECT
    c.column_name,
    c.data_type,
    (c.column_default IS NOT NULL) AS default_value,
    COALESCE(
        bool_or(tc.constraint_type IN ('UNIQUE', 'PRIMARY KEY')),
        false
    ) AS is_unique,
    COALESCE(
	bool_or(
        c.is_identity = 'YES'
        OR c.column_default LIKE 'nextval(%'
    ) , false) AS is_auto_increment
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
`, []any{tableName}
}

func (d *PostgresDialect) PlaceHolder(n int) (string, error) {
	if n <= 0 {
		return "", apperr.ErrorInvalidPlaceHolderIndex
	}
	return fmt.Sprintf("$%d", n), nil
}

func (d *PostgresDialect) QuoteName(tableName string) string {
	return quoteName(tableName)
}

func (d *PostgresDialect) AutoIncrementKeyword() string {
	return "SERIAL"
}

func (d *PostgresDialect) InsertDefaultValues(tableName string) string {
	return fmt.Sprintf("INSERT INTO %s DEFAULT VALUES", d.QuoteName(tableName))
}

func (d *PostgresDialect) DeleteRow(tableName string, whereClause string) string {
	quoted := d.QuoteName(tableName)
	return fmt.Sprintf("DELETE FROM %s WHERE ctid IN (SELECT ctid FROM %s WHERE %s LIMIT 1)", quoted, quoted, whereClause)
}

func (d *PostgresDialect) FilterOneRowClause(tableName, whereClause string) string {
	return fmt.Sprintf("ctid IN (SELECT ctid FROM %s WHERE %s LIMIT 1)", d.QuoteName(tableName), whereClause)
}

func (d *PostgresDialect) WhereCluse(cols []models.ColValues, rows []any, argsIdx int) (string, []any, error) {
	return whereClause(d, cols, rows, argsIdx)
}

// MySQLDialect implementation
type MySQLDialect struct{}

func (d *MySQLDialect) Name() configs.Driver {
	return configs.DriverMySQL
}

func (d *MySQLDialect) CheckTableExists(tableName string) (string, []any) {
	return `SELECT 1 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?`, []any{tableName}
}

func (d *MySQLDialect) ListTables() string {
	// Reusing Postgres query as it was in original code (postgresMySQLTablesListQuery)
	return `
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
}

func (d *MySQLDialect) ListColumns(tableName string) (string, []any) {
	return `
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
`, []any{tableName}
}

func (d *MySQLDialect) PlaceHolder(n int) (string, error) {
	if n <= 0 {
		return "", apperr.ErrorInvalidPlaceHolderIndex
	}
	return "?", nil
}

func (d *MySQLDialect) QuoteName(tableName string) string {
	return quoteName(tableName, "`%s`")
}

func (d *MySQLDialect) AutoIncrementKeyword() string {
	return "AUTO_INCREMENT"
}

func (d *MySQLDialect) InsertDefaultValues(tableName string) string {
	return fmt.Sprintf("INSERT INTO %s () VALUES ()", d.QuoteName(tableName))
}

func (d *MySQLDialect) DeleteRow(tableName string, whereClause string) string {
	return fmt.Sprintf("DELETE FROM %s WHERE %s LIMIT 1", d.QuoteName(tableName), whereClause)
}

func (d *MySQLDialect) FilterOneRowClause(tableName, whereClause string) string {
	return fmt.Sprintf("%s LIMIT 1", whereClause)
}

func (d *MySQLDialect) WhereCluse(cols []models.ColValues, rows []any, argsIdx int) (string, []any, error) {
	return whereClause(d, cols, rows, argsIdx)
}

// SQLiteDialect implementation
type SQLiteDialect struct{}

func (d *SQLiteDialect) Name() configs.Driver {
	return configs.DriverSQLite
}

func (d *SQLiteDialect) CheckTableExists(tableName string) (string, []any) {
	return `SELECT 1 FROM sqlite_master WHERE type = 'table' AND name = ?`, []any{tableName}
}

func (d *SQLiteDialect) ListTables() string {
	return `
SELECT
  '' AS table_schema,
  name AS table_name
FROM sqlite_master
WHERE type = 'table'
  AND name NOT LIKE 'sqlite_%'
ORDER BY name;
`
}

func (d *SQLiteDialect) ListColumns(tableName string) (string, []any) {
	return `
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
`, []any{tableName, tableName}
}

func (d *SQLiteDialect) PlaceHolder(n int) (string, error) {
	if n <= 0 {
		return "", apperr.ErrorInvalidPlaceHolderIndex
	}
	return fmt.Sprintf("$%d", n), nil
}

func (d *SQLiteDialect) QuoteName(tableName string) string {
	return quoteName(tableName, "\"%s\"")
	// return quoteName(tableName)
}

func (d *SQLiteDialect) AutoIncrementKeyword() string {
	return "AUTOINCREMENT"
}

func (d *SQLiteDialect) InsertDefaultValues(tableName string) string {
	return fmt.Sprintf("INSERT INTO %s DEFAULT VALUES", d.QuoteName(tableName))
}

func (d *SQLiteDialect) DeleteRow(tableName string, whereClause string) string {
	quoted := d.QuoteName(tableName)
	return fmt.Sprintf("DELETE FROM %s WHERE rowid IN (SELECT rowid FROM %s WHERE %s LIMIT 1)", quoted, quoted, whereClause)
}

func (d *SQLiteDialect) FilterOneRowClause(tableName, whereClause string) string {
	return fmt.Sprintf("rowid IN (SELECT rowid FROM %s WHERE %s LIMIT 1)", d.QuoteName(tableName), whereClause)
}

func (d *SQLiteDialect) WhereCluse(cols []models.ColValues, rows []any, argsIdx int) (string, []any, error) {
	return whereClause(d, cols, rows, argsIdx)
}
