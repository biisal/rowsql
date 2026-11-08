package queries

import (
	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database/sqlb"
	"strings"
)

type sqliteDialect struct{}

func (sqliteDialect) Name() configs.Driver { return configs.DriverSQLite }

func (sqliteDialect) QuoteIdentifier(identifier string) string {
	return sqlb.QuoteWith(identifier, `"`)
}

func (sqliteDialect) BindVariable(position int) string { return sqlb.PositionalDollar(position) }

func (sqliteDialect) ApplyAutoIncrement(typeName string) (string, string, error) {
	if !strings.EqualFold(typeName, "INTEGER") {
		return "", "", apperr.ErrorInvalidAutoIncrement
	}
	return typeName, " AUTOINCREMENT", nil
}

func (sqliteDialect) InsertDefaultValuesClause() string { return "DEFAULT VALUES" }

func (sqliteDialect) LimitsWritesToOneRow() bool { return false }

func (sqliteDialect) PhysicalRowID() (PhysicalRowID, bool) {
	return PhysicalRowID{Column: "rowid", SelectExpr: "rowid"}, true
}

func (sqliteDialect) CheckTableExists(tableName string) (string, []any) {
	return `SELECT 1 FROM sqlite_master WHERE type = 'table' AND name = ?`, []any{tableName}
}

func (sqliteDialect) ListTables() string {
	return `
SELECT '' AS table_schema, name AS table_name
FROM sqlite_master
WHERE type = 'table' AND name NOT LIKE 'sqlite_%'
ORDER BY name`
}

func (sqliteDialect) ListColumns(tableName string) (string, []any) {
	return `
SELECT
    p.name AS column_name,
    p.type AS data_type,
    (p.dflt_value IS NOT NULL) AS has_default,
    EXISTS (
        SELECT 1
        FROM pragma_index_list(?) il
        JOIN pragma_index_info(il.name) ii ON ii.name = p.name
        WHERE il."unique" = 1 AND il.origin <> 'pk'
          AND (SELECT COUNT(*) FROM pragma_index_info(il.name)) = 1
    ) AS is_unique,
    (p.pk = 1 AND lower(p.type) = 'integer'
        AND (SELECT COUNT(*) FROM pragma_table_info(?) WHERE pk > 0) = 1) AS is_auto_increment,
    (p."notnull" = 0) AS is_nullable,
    p.pk AS pk_position
FROM pragma_table_info(?) AS p
ORDER BY p.cid`, []any{tableName, tableName, tableName}
}
