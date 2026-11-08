package queries

import (
	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/database/sqlb"
)

type mysqlDialect struct{}

func (mysqlDialect) Name() configs.Driver { return configs.DriverMySQL }

func (mysqlDialect) QuoteIdentifier(identifier string) string {
	return sqlb.QuoteWith(identifier, "`")
}

func (mysqlDialect) BindVariable(int) string { return "?" }

func (mysqlDialect) ApplyAutoIncrement(typeName string) (string, string, error) {
	return typeName, " AUTO_INCREMENT", nil
}

func (mysqlDialect) InsertDefaultValuesClause() string { return "() VALUES ()" }

func (mysqlDialect) LimitsWritesToOneRow() bool { return true }

func (mysqlDialect) PhysicalRowID() (PhysicalRowID, bool) { return PhysicalRowID{}, false }

func (mysqlDialect) CheckTableExists(tableName string) (string, []any) {
	return `SELECT 1 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?`, []any{
		tableName,
	}
}

func (mysqlDialect) ListTables() string {
	return `
SELECT table_schema, table_name
FROM information_schema.tables
WHERE table_type = 'BASE TABLE'
  AND table_schema = DATABASE()
ORDER BY table_name`
}

func (mysqlDialect) ListColumns(tableName string) (string, []any) {
	return `
SELECT
    c.column_name,
    c.data_type,
    (c.column_default IS NOT NULL) AS has_default,
    (c.column_key = 'UNI') AS is_unique,
    (c.extra LIKE '%auto_increment%') AS is_auto_increment,
    (c.is_nullable = 'YES') AS is_nullable,
    COALESCE(k.ordinal_position, 0) AS pk_position
FROM information_schema.columns c
LEFT JOIN information_schema.key_column_usage k
    ON k.table_schema = c.table_schema
    AND k.table_name = c.table_name
    AND k.column_name = c.column_name
    AND k.constraint_name = 'PRIMARY'
WHERE c.table_schema = DATABASE() AND c.table_name = ?
ORDER BY c.ordinal_position`, []any{tableName}
}
