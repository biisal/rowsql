package repo

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/biisal/db-gui/configs"
	"github.com/biisal/db-gui/internal/database"
	"github.com/biisal/db-gui/internal/logger"
)

const fallbackMsg = "Unknown driver – falling back to SQLite"

const pgColsQuery = `
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

const mysqlColsQuery = `
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

const sqliteColsQuery = `
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

func colsQuery(driver string, tableName string) (string, []any) {
	switch driver {
	case configs.DRIVER_POSTGRES:
		return pgColsQuery, []any{tableName}
	case configs.DRIVER_MYSQL:
		return mysqlColsQuery, []any{tableName}
	case configs.DRIVER_SQLITE:
		return sqliteColsQuery, []any{tableName, tableName}
	default:
		logger.Info("%s driver=%s", fallbackMsg, driver)
		return sqliteColsQuery, []any{tableName, tableName}
	}
}

const pgMysqlTablesQuery = `
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

const sqliteTablesQuery = `
SELECT
  '' AS table_schema,
  name AS table_name
FROM sqlite_master
WHERE type = 'table'
  AND name NOT LIKE 'sqlite_%'
ORDER BY name;
`

func tablesQuery(driver string) string {
	switch driver {
	case configs.DRIVER_POSTGRES, configs.DRIVER_MYSQL:
		return pgMysqlTablesQuery
	case configs.DRIVER_SQLITE:
		return sqliteTablesQuery
	default:
		logger.Info("%s driver=%s", fallbackMsg, driver)
		return sqliteTablesQuery
	}
}

func buildQueryParts(form map[string]FormValue, driver string) (*QueryParts, error) {
	columns := make([]string, 0, len(form))
	placeholders := make([]string, 0, len(form))
	values := make([]any, 0, len(form))

	isPostgres := driver == configs.DRIVER_POSTGRES

	paramIndex := 1

	var ph = "?"
	for col, field := range form {
		if field.Value == "" {
			continue
		}
		columns = append(columns, col)
		if isPostgres {
			ph = "$" + strconv.Itoa(paramIndex)
		}
		placeholders = append(placeholders, ph)

		if field.Type == "json" {
			var jsonVal map[string]any
			if err := json.Unmarshal([]byte(field.Value), &jsonVal); err != nil {
				logger.Errorln(err)
				return nil, err
			}
			values = append(values, jsonVal)
		} else {
			values = append(values, field.Value)
		}

		paramIndex++
	}

	if len(columns) == 0 {
		return nil, fmt.Errorf("no data provided")
	}

	return &QueryParts{
		Columns:      strings.Join(columns, ","),
		Placeholders: strings.Join(placeholders, ","),
		Args:         values,
	}, nil
}

func buildUpdateParts(form map[string]FormValue, driver string) (string, []any) {
	parts := make([]string, 0, len(form))
	ph := "?"
	index := 1
	args := make([]any, 0, len(form))
	for k, v := range form {
		if v.Value == "" {
			continue
		}
		if driver == configs.DRIVER_POSTGRES {
			ph = "$" + strconv.Itoa(index)
		}
		parts = append(parts, fmt.Sprintf("%s=%s", k, ph))
		args = append(args, v.Value)
		index++
	}
	return strings.Join(parts, ","), args
}

func buildQueryWhereClause(cols []ListDataCol, data []any, driver string, argsIdx int) (string, []any, error) {
	if len(cols) != len(data) {
		return "", nil, fmt.Errorf("cols and rows aren't same in length")
	}

	var mixed []string
	var ph = "?"
	var args []any
	for i, val := range cols {
		if data[i] == "" {
			continue
		}
		if driver == configs.DRIVER_POSTGRES {
			ph = "$" + strconv.Itoa(argsIdx+i)
		}
		if val.IsUnique {
			return fmt.Sprintf("%s=%s", val.ColumnName, ph), []any{data[i]}, nil
		}
		var colVal any = data[i]
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

func buildCreateTableQuiry(driver string, parts []database.Input) (string, error) {
	logger.Info("Building create table query")
	var s = strings.Builder{}
	partsLen := len(parts)
	for i, input := range parts {
		if input.ColName == "" {
			continue
		}
		fmt.Fprintf(&s, "%s %s", input.ColName, input.DataType.Type)
		if input.DataType.HasSize {
			fmt.Fprintf(&s, "(%d)", input.DataType.Size)
		}
		if input.IsUnique {
			s.WriteString(" UNIQUE")
		}
		if !input.IsNull {
			s.WriteString(" NOT NULL")
		}
		if input.IsPK {
			s.WriteString(" PRIMARY KEY")
		}
		if input.DataType.AutoIncrement {
			if !input.IsPK {
				logger.Error("Auto-increment can only be set on primary key columns")
				return "", fmt.Errorf("auto-increment can only be set on primary key columns")
			}
			text := "AUTO_INCREMENT"
			switch driver {
			case configs.DRIVER_POSTGRES:
				text = "SERIAL"
			case configs.DRIVER_SQLITE:
				text = "AUTOINCREMENT"
			}
			fmt.Fprintf(&s, " %s", text)
		}
		if i < partsLen-1 {
			s.WriteString(",")
		}
	}
	logger.Info("Query: %s", s.String())
	return s.String(), nil

}
