package database

var PostgresNumericDataTypes = []VarianDataType{
	{Type: "SMALLINT", HasAutoIncrement: true},
	{Type: "INT2", HasAutoIncrement: true},
	{Type: "INTEGER", HasAutoIncrement: true},
	{Type: "INT", HasAutoIncrement: true},
	{Type: "INT4", HasAutoIncrement: true},
	{Type: "BIGINT", HasAutoIncrement: true},
	{Type: "INT8", HasAutoIncrement: true},
	{Type: "DECIMAL", HasDigit: true}, // DECIMAL(p, s)
	{Type: "NUMERIC", HasDigit: true}, // NUMERIC(p, s)
	{Type: "REAL"},
	{Type: "FLOAT4"},
	{Type: "DOUBLE PRECISION"},
	{Type: "FLOAT8"},
	{Type: "SMALLSERIAL", HasAutoIncrement: true},
	{Type: "SERIAL2", HasAutoIncrement: true},
	{Type: "SERIAL", HasAutoIncrement: true},
	{Type: "SERIAL4", HasAutoIncrement: true},
	{Type: "BIGSERIAL", HasAutoIncrement: true},
	{Type: "SERIAL8", HasAutoIncrement: true},
	{Type: "MONEY"},
}

func init() {
	for i := range PostgresNumericDataTypes {
		PostgresNumericDataTypes[i].IsNumeric = true
	}
}

var PostgresStringDataTypes = []VarianDataType{
	{Type: "CHAR", HasSize: true},              // CHAR(n)
	{Type: "CHARACTER", HasSize: true},         // CHARACTER(n)
	{Type: "VARCHAR", HasSize: true},           // VARCHAR(n)
	{Type: "CHARACTER VARYING", HasSize: true}, // CHARACTER VARYING(n)
	{Type: "TEXT"},                             // No size limit
	{Type: "BPCHAR", HasSize: true},            // Internal name for CHAR
	{Type: "BYTEA"},                            // Binary data
	{Type: "UUID"},
	{Type: "JSON"},
	{Type: "JSONB"},
	{Type: "XML"},
	{Type: "CITEXT"}, // Case-insensitive text (requires extension)
}
