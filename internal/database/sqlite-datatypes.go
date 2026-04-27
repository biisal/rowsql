package database

var SqliteNumericDataTypes = []VarianDataType{
	{Type: "INT"},
	{Type: "INTEGER", HasAutoIncrement: true},
	{Type: "TINYINT"},
	{Type: "SMALLINT"},
	{Type: "MEDIUMINT"},
	{Type: "BIGINT"},
	{Type: "UNSIGNED BIG INT"},
	{Type: "INT2"},
	{Type: "INT8"},
	{Type: "REAL"},
	{Type: "DOUBLE"},
	{Type: "DOUBLE PRECISION"},
	{Type: "FLOAT"},
	{Type: "NUMERIC", HasSize: true, HasDigit: true},
	{Type: "DECIMAL", HasSize: true, HasDigit: true},
	{Type: "BOOLEAN"},
	{Type: "DATE"},
	{Type: "DATETIME"},
}

func init() {
	for i := range SqliteNumericDataTypes {
		SqliteNumericDataTypes[i].IsNumeric = true
	}
}

var SqliteStringDataTypes = []VarianDataType{
	{Type: "TEXT"},
	{Type: "CHARACTER", HasSize: true},
	{Type: "VARCHAR", HasSize: true},
	{Type: "VARYING CHARACTER", HasSize: true},
	{Type: "NCHAR", HasSize: true},
	{Type: "NATIVE CHARACTER", HasSize: true},
	{Type: "CLOB"},
	{Type: "BLOB"},
}
