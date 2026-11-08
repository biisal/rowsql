package queries

import (
	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database/sqlb"
	"strings"
)

type postgresDialect struct{}

func (postgresDialect) Name() configs.Driver { return configs.DriverPostgres }

func (postgresDialect) QuoteIdentifier(identifier string) string {
	return sqlb.QuoteWith(identifier, `"`)
}

func (postgresDialect) BindVariable(position int) string { return sqlb.PositionalDollar(position) }

var postgresSerialTypes = map[string]string{
	"SMALLINT": "SMALLSERIAL", "INT2": "SMALLSERIAL", "SMALLSERIAL": "SMALLSERIAL", "SERIAL2": "SERIAL2",
	"INTEGER": "SERIAL", "INT": "SERIAL", "INT4": "SERIAL", "SERIAL": "SERIAL", "SERIAL4": "SERIAL4",
	"BIGINT": "BIGSERIAL", "INT8": "BIGSERIAL", "BIGSERIAL": "BIGSERIAL", "SERIAL8": "SERIAL8",
}

func (postgresDialect) ApplyAutoIncrement(typeName string) (string, string, error) {
	serialType, supported := postgresSerialTypes[strings.ToUpper(typeName)]
	if !supported {
		return "", "", apperr.ErrorInvalidAutoIncrement
	}
	return serialType, "", nil
}

func (postgresDialect) InsertDefaultValuesClause() string { return "DEFAULT VALUES" }

func (postgresDialect) LimitsWritesToOneRow() bool { return false }

func (postgresDialect) PhysicalRowID() (PhysicalRowID, bool) {
	return PhysicalRowID{Column: "ctid", SelectExpr: "ctid::text", BindSuffix: "::text::tid"}, true
}

func (postgresDialect) CheckTableExists(tableName string) (string, []any) {
	return `SELECT 1 FROM information_schema.tables WHERE table_schema = current_schema() AND table_name = $1`, []any{
		tableName,
	}
}

func (postgresDialect) ListTables() string {
	return `
SELECT table_schema, table_name
FROM information_schema.tables
WHERE table_type = 'BASE TABLE'
  AND table_schema = current_schema()
ORDER BY table_name`
}

func (postgresDialect) ListColumns(tableName string) (string, []any) {
	return `
SELECT
    c.column_name,
    c.data_type,
    (c.column_default IS NOT NULL) AS has_default,
    COALESCE(single_unique.is_unique, false) AS is_unique,
    (c.is_identity = 'YES' OR COALESCE(c.column_default, '') LIKE 'nextval(%') AS is_auto_increment,
    (c.is_nullable = 'YES') AS is_nullable,
    COALESCE(primary_key.position, 0) AS pk_position
FROM information_schema.columns c
LEFT JOIN (
    SELECT a.attname AS column_name, key_column.position::int AS position
    FROM pg_index i
    JOIN pg_class t ON t.oid = i.indrelid
    JOIN pg_namespace n ON n.oid = t.relnamespace
    CROSS JOIN LATERAL unnest(i.indkey) WITH ORDINALITY AS key_column(attnum, position)
    JOIN pg_attribute a ON a.attrelid = t.oid AND a.attnum = key_column.attnum
    WHERE i.indisprimary AND t.relname = $1 AND n.nspname = current_schema()
) primary_key ON primary_key.column_name = c.column_name
LEFT JOIN (
    SELECT a.attname AS column_name, true AS is_unique
    FROM pg_index i
    JOIN pg_class t ON t.oid = i.indrelid
    JOIN pg_namespace n ON n.oid = t.relnamespace
    JOIN pg_attribute a ON a.attrelid = t.oid AND a.attnum = i.indkey[0]
    WHERE i.indisunique AND NOT i.indisprimary AND i.indnatts = 1
      AND t.relname = $1 AND n.nspname = current_schema()
) single_unique ON single_unique.column_name = c.column_name
WHERE c.table_name = $1 AND c.table_schema = current_schema()
ORDER BY c.ordinal_position`, []any{tableName}
}
