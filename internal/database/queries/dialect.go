package queries

import (
	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database/sqlb"
)

func GetDialect(name configs.Driver) (Dialect, error) {
	switch name {
	case configs.DriverPostgres:
		return postgresDialect{}, nil
	case configs.DriverMySQL:
		return mysqlDialect{}, nil
	case configs.DriverSQLite:
		return sqliteDialect{}, nil
	default:
		return nil, apperr.ErrorInvalidDriver
	}
}

type Dialect interface {
	sqlb.Dialect
	Name() configs.Driver
	CheckTableExists(tableName string) (string, []any)
	ListTables() string
	ListColumns(tableName string) (string, []any)
	ApplyAutoIncrement(typeName string) (resultType string, suffix string, err error)
	InsertDefaultValuesClause() string
	LimitsWritesToOneRow() bool
	PhysicalRowID() (physicalRowID PhysicalRowID, supported bool)
}

type PhysicalRowID struct {
	Column     string
	SelectExpr string
	BindSuffix string
}
