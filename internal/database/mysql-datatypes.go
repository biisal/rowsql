package database

var MySqlNumericDataTypes = []NumericDataType{
	{Type: "BIT", HasSize: true, HasDigit: false},
	{Type: "TINYINT", HasSize: true, HasDigit: false},
	{Type: "BOOL", HasSize: false, HasDigit: false},
	{Type: "BOOLEAN", HasSize: false, HasDigit: false},
	{Type: "SMALLINT", HasSize: true, HasDigit: false},
	{Type: "MEDIUMINT", HasSize: true, HasDigit: false},
	{Type: "INT", HasSize: true, HasDigit: false},
	{Type: "INTEGER", HasSize: true, HasDigit: false},
	{Type: "BIGINT", HasSize: true, HasDigit: false},
	{Type: "FLOAT", HasSize: true, HasDigit: true},            // FLOAT(size, d)
	{Type: "FLOAT_PRECISION", HasSize: true, HasDigit: false}, // FLOAT(p)
	{Type: "DOUBLE", HasSize: true, HasDigit: true},
	{Type: "DOUBLE PRECISION", HasSize: true, HasDigit: true},
	{Type: "DECIMAL", HasSize: true, HasDigit: true},
	{Type: "DEC", HasSize: true, HasDigit: true},
}

var MySqlStringDataTypes = []StringDataType{
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
