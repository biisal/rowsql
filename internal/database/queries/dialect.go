package queries

import (
	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database/models"
)

func GetDialect(name configs.Driver) (Dialect, error) {
	switch name {
	case configs.DriverPostgres:
		return &PostgresDialect{}, nil
	case configs.DriverMySQL:
		return &MySQLDialect{}, nil
	case configs.DriverSQLite:
		return &SQLiteDialect{}, nil
	default:
		return nil, apperr.ErrorInvalidDriver
	}
}

type Dialect interface {
	Name() configs.Driver
	CheckTableExists(tableName string) (string, []any)
	ListTables() string
	ListColumns(tableName string) (string, []any)
	PlaceHolder(n int) (string, error)
	QuoteName(tableName string) string
	AutoIncrementKeyword() string
	InsertDefaultValues(tableName string) string
	DeleteRow(tableName string, whereClause string) string
	FilterOneRowClause(tableName, whereClause string) string
	WhereCluse(cols []models.ListDataCol, rows []any, argsIdx int) (string, []any, error)
}
