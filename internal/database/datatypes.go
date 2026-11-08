package database

import "github.com/biisal/rowsql/configs"

type VarianDataType struct {
	Type             string `json:"type"`
	HasSize          bool   `json:"hasSize"`
	HasDigit         bool   `json:"hasDigit"`
	HasAutoIncrement bool   `json:"hasAutoIncrement"`
	HasValues        bool   `json:"hasValues"`
	IsNumeric        bool   `json:"isNumeric"`
}

func DataTypesFor(driver configs.Driver) (numeric []VarianDataType, text []VarianDataType) {
	switch driver {
	case configs.DriverMySQL:
		return MySQLNumericDataTypes, MySQLStringDataTypes
	case configs.DriverPostgres:
		return PostgresNumericDataTypes, PostgresStringDataTypes
	case configs.DriverSQLite:
		return SqliteNumericDataTypes, SqliteStringDataTypes
	default:
		return nil, nil
	}
}
