package database

var MySQLNumericDataTypes = []NumericDataType{
	{Type: "BIT", HasSize: true, HasDigit: false, HasAutoIncrement: false},
	{Type: "TINYINT", HasSize: true, HasDigit: false, HasAutoIncrement: true},
	{Type: "BOOL", HasSize: false, HasDigit: false, HasAutoIncrement: false},
	{Type: "BOOLEAN", HasSize: false, HasDigit: false, HasAutoIncrement: false},
	{Type: "SMALLINT", HasSize: true, HasDigit: false, HasAutoIncrement: true},
	{Type: "MEDIUMINT", HasSize: true, HasDigit: false, HasAutoIncrement: true},
	{Type: "INT", HasSize: true, HasDigit: false, HasAutoIncrement: true},
	{Type: "INTEGER", HasSize: true, HasDigit: false, HasAutoIncrement: true},
	{Type: "BIGINT", HasSize: true, HasDigit: false, HasAutoIncrement: true},
	{Type: "FLOAT", HasSize: true, HasDigit: true, HasAutoIncrement: false},
	{Type: "FLOAT_PRECISION", HasSize: true, HasDigit: false, HasAutoIncrement: false},
	{Type: "DOUBLE", HasSize: true, HasDigit: true, HasAutoIncrement: false},
	{Type: "DOUBLE PRECISION", HasSize: true, HasDigit: true, HasAutoIncrement: false},
	{Type: "DECIMAL", HasSize: true, HasDigit: true, HasAutoIncrement: false},
	{Type: "DEC", HasSize: true, HasDigit: true, HasAutoIncrement: false},
}

var MySQLStringDataTypes = []StringDataType{
	{Type: "CHAR", HasSize: true, HasValues: false},
	{Type: "VARCHAR", HasSize: true, HasValues: false},
	{Type: "BINARY", HasSize: true, HasValues: false},
	{Type: "VARBINARY", HasSize: true, HasValues: false},

	{Type: "TINYBLOB", HasSize: false, HasValues: false},
	{Type: "TINYTEXT", HasSize: false, HasValues: false},

	{Type: "TEXT", HasSize: true, HasValues: false}, // TEXT(size)
	{Type: "BLOB", HasSize: true, HasValues: false}, // BLOB(size)

	{Type: "MEDIUMTEXT", HasSize: false, HasValues: false},
	{Type: "MEDIUMBLOB", HasSize: false, HasValues: false},

	{Type: "LONGTEXT", HasSize: false, HasValues: false},
	{Type: "LONGBLOB", HasSize: false, HasValues: false},

	{Type: "ENUM", HasSize: false, HasValues: true}, // ENUM(val1, val2, ...)
	{Type: "SET", HasSize: false, HasValues: true},  // SET(val1, val2, ...)
}
