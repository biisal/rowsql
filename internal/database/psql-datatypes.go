package database

var PostgresNumericDataTypes = []NumericDataType{
	{Type: "SMALLINT", HasDigit: false, HasAutoIncrement: true},
	{Type: "INT2", HasDigit: false, HasAutoIncrement: true},
	{Type: "INTEGER", HasDigit: false, HasAutoIncrement: true},
	{Type: "INT", HasDigit: false, HasAutoIncrement: true},
	{Type: "INT4", HasDigit: false, HasAutoIncrement: true},
	{Type: "BIGINT", HasDigit: false, HasAutoIncrement: true},
	{Type: "INT8", HasDigit: false, HasAutoIncrement: true},

	{Type: "DECIMAL", HasDigit: true, HasAutoIncrement: false}, // DECIMAL(p, s)
	{Type: "NUMERIC", HasDigit: true, HasAutoIncrement: false}, // NUMERIC(p, s)

	{Type: "REAL", HasDigit: false, HasAutoIncrement: false},
	{Type: "FLOAT4", HasDigit: false, HasAutoIncrement: false},
	{Type: "DOUBLE PRECISION", HasDigit: false, HasAutoIncrement: false},
	{Type: "FLOAT8", HasDigit: false, HasAutoIncrement: false},

	{Type: "SMALLSERIAL", HasDigit: false, HasAutoIncrement: true},
	{Type: "SERIAL2", HasDigit: false, HasAutoIncrement: true},
	{Type: "SERIAL", HasDigit: false, HasAutoIncrement: true},
	{Type: "SERIAL4", HasDigit: false, HasAutoIncrement: true},
	{Type: "BIGSERIAL", HasDigit: false, HasAutoIncrement: true},
	{Type: "SERIAL8", HasDigit: false, HasAutoIncrement: true},

	{Type: "MONEY", HasDigit: false, HasAutoIncrement: false},
}

var PostgresStringDataTypes = []StringDataType{
	{Type: "CHAR", HasSize: true, HasValues: false},              // CHAR(n)
	{Type: "CHARACTER", HasSize: true, HasValues: false},         // CHARACTER(n)
	{Type: "VARCHAR", HasSize: true, HasValues: false},           // VARCHAR(n)
	{Type: "CHARACTER VARYING", HasSize: true, HasValues: false}, // CHARACTER VARYING(n)
	{Type: "TEXT", HasValues: false},                             // No size limit
	{Type: "BPCHAR", HasSize: true, HasValues: false},            // Internal name for CHAR
	{Type: "BYTEA", HasValues: false},                            // Binary data
	{Type: "UUID", HasValues: false},
	{Type: "JSON", HasValues: false},
	{Type: "JSONB", HasValues: false},
	{Type: "XML", HasValues: false},
	{Type: "CITEXT", HasValues: false}, // Case-insensitive text (requires extension)
}
