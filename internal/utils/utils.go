package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"

	"github.com/biisal/db-gui/configs"
)

func IsSafeIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if !(c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
			return false
		}
	}
	return true
}

func DetectDriver(connectionString string) (string, error) {
	if connectionString == "" {
		return "", fmt.Errorf("connection string is empty")
	}

	lower := strings.ToLower(connectionString)

	switch {
	case strings.HasPrefix(lower, "postgres://") ||
		strings.HasPrefix(lower, "postgresql://") ||
		strings.Contains(lower, "host=") && strings.Contains(lower, "dbname="):
		return configs.DRIVER_POSTGRES, nil

	case strings.HasPrefix(lower, "mysql://") ||
		strings.Contains(lower, "tcp(") ||
		strings.Contains(lower, "parseTime="):
		return configs.DRIVER_MYSQL, nil

	case strings.HasPrefix(lower, "sqlite://") ||
		strings.HasPrefix(lower, "file:") ||
		strings.HasSuffix(lower, ".db") ||
		strings.HasSuffix(lower, ".sqlite") ||
		strings.HasSuffix(lower, ".sqlite3") ||
		lower == ":memory:":
		return configs.DRIVER_SQLITE, nil

	default:
		return "", fmt.Errorf("unable to detect driver from connection string! Make sure your connection string is coorect and try again.")
	}
}

const (
	textInput     = "text"
	selectInput   = "select"
	checkboxInput = "checkbox"
	textAreaInput = "textarea"
	numberInput   = "number"
	jsonInput     = "json"
)

var dataTypeMap = map[string]string{
	"smallint":         numberInput,
	"integer":          numberInput,
	"bigint":           numberInput,
	"decimal":          numberInput,
	"numeric":          numberInput,
	"real":             numberInput,
	"double precision": numberInput,
	"smallserial":      numberInput,
	"serial":           numberInput,
	"bigserial":        numberInput,
	"money":            numberInput,

	"boolean": checkboxInput,
	"bool":    checkboxInput,

	"text":  textAreaInput,
	"json":  jsonInput,
	"jsonb": textAreaInput,
	"xml":   textAreaInput,
	"bytea": textAreaInput,

	"character":         textInput,
	"character varying": textInput,
	"varchar":           textInput,
	"uuid":              textInput,

	"date":                        textInput,
	"time":                        textInput,
	"time without time zone":      textInput,
	"time with time zone":         textInput,
	"timestamp":                   textInput,
	"timestamp without time zone": textInput,
	"timestamp with time zone":    textInput,
	"interval":                    textInput,

	"inet":          textInput,
	"cidr":          textInput,
	"macaddr":       textInput,
	"macaddr8":      textInput,
	"bit":           textInput,
	"bit varying":   textInput,
	"point":         textInput,
	"line":          textInput,
	"lseg":          textInput,
	"box":           textInput,
	"path":          textInput,
	"polygon":       textInput,
	"circle":        textInput,
	"tsvector":      textInput,
	"tsquery":       textInput,
	"pg_lsn":        textInput,
	"pg_snapshot":   textInput,
	"txid_snapshot": textInput,

	"tinyint":   numberInput,
	"mediumint": numberInput,
	"float":     numberInput,
	"double":    numberInput,
	"year":      numberInput,

	"tinyblob":   textAreaInput,
	"mediumblob": textAreaInput,
	"blob":       textAreaInput,
	"longblob":   textAreaInput,
	"binary":     textAreaInput,
	"varbinary":  textAreaInput,

	"tinytext":   textAreaInput,
	"mediumtext": textAreaInput,
	"longtext":   textAreaInput,

	"enum": selectInput,
	"set":  selectInput,

	"int":              numberInput,
	"int2":             numberInput,
	"int8":             numberInput,
	"unsigned big int": numberInput,

	"clob":     textAreaInput,
	"datetime": textInput,

	"native character":  textInput,
	"nchar":             textInput,
	"nvarchar":          textInput,
	"varying character": textInput,
}

func GetInputType(dbType string) string {
	lowerType := strings.ToLower(dbType)

	if idx := strings.Index(lowerType, "("); idx != -1 {
		lowerType = lowerType[:idx]
	}

	lowerType = strings.TrimSuffix(lowerType, " unsigned")

	if val, ok := dataTypeMap[lowerType]; ok {
		return val
	}
	slog.Warn("unknown data type", "type", dbType)
	return textInput
}

func MakeRowHash(data []any) string {
	h := sha256.New()
	fmt.Fprint(h, data)
	hashBytes := h.Sum(nil)
	fullHash := hex.EncodeToString(hashBytes)

	shortHash := fullHash[:8]
	return shortHash
}
