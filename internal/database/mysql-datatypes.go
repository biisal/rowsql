package database

var MySQLNumericDataTypes = []VarianDataType{
	{Type: "BIT", HasSize: true},
	{Type: "TINYINT", HasSize: true, HasAutoIncrement: true},
	{Type: "BOOL"},
	{Type: "BOOLEAN"},
	{Type: "SMALLINT", HasSize: true, HasAutoIncrement: true},
	{Type: "MEDIUMINT", HasSize: true, HasAutoIncrement: true},
	{Type: "INT", HasSize: true, HasAutoIncrement: true},
	{Type: "INTEGER", HasSize: true, HasAutoIncrement: true},
	{Type: "BIGINT", HasSize: true, HasAutoIncrement: true},
	{Type: "FLOAT", HasSize: true, HasDigit: true},
	{Type: "FLOAT_PRECISION", HasSize: true},
	{Type: "DOUBLE", HasSize: true, HasDigit: true},
	{Type: "DOUBLE PRECISION", HasSize: true, HasDigit: true},
	{Type: "DECIMAL", HasSize: true, HasDigit: true},
	{Type: "DEC", HasSize: true, HasDigit: true},
}

func init() {
	for i := range MySQLNumericDataTypes {
		MySQLNumericDataTypes[i].IsNumeric = true
	}
}

var MySQLStringDataTypes = []VarianDataType{
	{Type: "CHAR", HasSize: true},
	{Type: "VARCHAR", HasSize: true},
	{Type: "BINARY", HasSize: true},
	{Type: "VARBINARY", HasSize: true},

	{Type: "TINYBLOB"},
	{Type: "TINYTEXT"},

	{Type: "TEXT", HasSize: true}, // TEXT(size)
	{Type: "BLOB", HasSize: true}, // BLOB(size)

	{Type: "MEDIUMTEXT"},
	{Type: "MEDIUMBLOB"},

	{Type: "LONGTEXT"},
	{Type: "LONGBLOB"},

	{Type: "ENUM", HasValues: true}, // ENUM(val1, val2, ...)
	{Type: "SET", HasValues: true},  // SET(val1, val2, ...)
}
