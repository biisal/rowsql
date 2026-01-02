package database

var SqliteNumericDataTypes = []NumericDataType{
	{Type: "INT", HasSize: false, HasDigit: false, HasAutoIncrement: false},
	{Type: "INTEGER", HasSize: false, HasDigit: false, HasAutoIncrement: true},
	{Type: "TINYINT", HasSize: false, HasDigit: false, HasAutoIncrement: false},
	{Type: "SMALLINT", HasSize: false, HasDigit: false, HasAutoIncrement: false},
	{Type: "MEDIUMINT", HasSize: false, HasDigit: false, HasAutoIncrement: false},
	{Type: "BIGINT", HasSize: false, HasDigit: false, HasAutoIncrement: false},
	{Type: "UNSIGNED BIG INT", HasSize: false, HasDigit: false, HasAutoIncrement: false},
	{Type: "INT2", HasSize: false, HasDigit: false, HasAutoIncrement: false},
	{Type: "INT8", HasSize: false, HasDigit: false, HasAutoIncrement: false},
	{Type: "REAL", HasSize: false, HasDigit: false, HasAutoIncrement: false},
	{Type: "DOUBLE", HasSize: false, HasDigit: false, HasAutoIncrement: false},
	{Type: "DOUBLE PRECISION", HasSize: false, HasDigit: false, HasAutoIncrement: false},
	{Type: "FLOAT", HasSize: false, HasDigit: false, HasAutoIncrement: false},
	{Type: "NUMERIC", HasSize: true, HasDigit: true, HasAutoIncrement: false},
	{Type: "DECIMAL", HasSize: true, HasDigit: true, HasAutoIncrement: false},
	{Type: "BOOLEAN", HasSize: false, HasDigit: false, HasAutoIncrement: false},
	{Type: "DATE", HasSize: false, HasDigit: false, HasAutoIncrement: false},
	{Type: "DATETIME", HasSize: false, HasDigit: false, HasAutoIncrement: false},
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
