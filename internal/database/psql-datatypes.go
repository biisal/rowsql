package database

var PostgresNumericDataTypes = []NumericDataType{
	{Type: "SMALLINT", HasDigit: false},
	{Type: "INT2", HasDigit: false},
	{Type: "INTEGER", HasDigit: false},
	{Type: "INT", HasDigit: false},
	{Type: "INT4", HasDigit: false},
	{Type: "BIGINT", HasDigit: false},
	{Type: "INT8", HasDigit: false},
	{Type: "DECIMAL", HasDigit: true}, // DECIMAL(precision, scale)
	{Type: "NUMERIC", HasDigit: true}, // NUMERIC(precision, scale)
	{Type: "REAL", HasDigit: false},
	{Type: "FLOAT4", HasDigit: false},
	{Type: "DOUBLE PRECISION", HasDigit: false},
	{Type: "FLOAT8", HasDigit: false},
	{Type: "SMALLSERIAL", HasDigit: false},
	{Type: "SERIAL2", HasDigit: false},
	{Type: "SERIAL", HasDigit: false},
	{Type: "SERIAL4", HasDigit: false},
	{Type: "BIGSERIAL", HasDigit: false},
	{Type: "SERIAL8", HasDigit: false},
	{Type: "MONEY", HasDigit: false},
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
