package database

var SqliteNumericDataTypes = []NumericDataType{
	{Type: "INT", HasSize: false, HasDigit: false},
	{Type: "INTEGER", HasSize: false, HasDigit: false},
	{Type: "TINYINT", HasSize: false, HasDigit: false},
	{Type: "SMALLINT", HasSize: false, HasDigit: false},
	{Type: "MEDIUMINT", HasSize: false, HasDigit: false},
	{Type: "BIGINT", HasSize: false, HasDigit: false},
	{Type: "UNSIGNED BIG INT", HasSize: false, HasDigit: false},
	{Type: "INT2", HasSize: false, HasDigit: false},
	{Type: "INT8", HasSize: false, HasDigit: false},
	{Type: "REAL", HasSize: false, HasDigit: false},
	{Type: "DOUBLE", HasSize: false, HasDigit: false},
	{Type: "DOUBLE PRECISION", HasSize: false, HasDigit: false},
	{Type: "FLOAT", HasSize: false, HasDigit: false},
	{Type: "NUMERIC", HasSize: true, HasDigit: true},
	{Type: "DECIMAL", HasSize: true, HasDigit: true},
	{Type: "BOOLEAN", HasSize: false, HasDigit: false},
	{Type: "DATE", HasSize: false, HasDigit: false},
	{Type: "DATETIME", HasSize: false, HasDigit: false},
}

var SqliteStringDataTypes = []StringDataType{
	{Type: "TEXT", HasSize: false, HasValues: false},
	{Type: "CHARACTER", HasSize: true, HasValues: false},
	{Type: "VARCHAR", HasSize: true, HasValues: false},
	{Type: "VARYING CHARACTER", HasSize: true, HasValues: false},
	{Type: "NCHAR", HasSize: true, HasValues: false},
	{Type: "NATIVE CHARACTER", HasSize: true, HasValues: false},
	{Type: "CLOB", HasSize: false, HasValues: false},
	{Type: "BLOB", HasSize: false, HasValues: false},
}
