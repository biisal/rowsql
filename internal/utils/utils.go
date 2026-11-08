// Package utils provides utility functions for the rowsql application.
package utils

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/logger"
	"github.com/go-sql-driver/mysql"
)

func ReplaceTildeWithHomeDir(text string) (string, error) {
	if text == "~" || strings.HasPrefix(text, "~/") {
		homeDir, homeErr := os.UserHomeDir()
		if homeErr != nil {
			logger.Error("failed to get user home directory: %s", homeErr)
			return "", homeErr
		}
		text = strings.Replace(text, "~", homeDir, 1)
	}
	return text, nil
}

func isPostgresConnString(connString string) bool {
	return strings.HasPrefix(connString, "postgres://") ||
		strings.HasPrefix(connString, "postgresql://") ||
		strings.Contains(connString, "host=") && strings.Contains(connString, "dbname=")
}

func isSQLiteConnString(connString string) bool {
	return strings.HasPrefix(connString, "sqlite://") ||
		strings.HasPrefix(connString, "file:") ||
		strings.HasSuffix(connString, ".db") ||
		strings.HasSuffix(connString, ".sqlite") ||
		strings.HasSuffix(connString, ".sqlite3") ||
		connString == ":memory:"
}

func isMySQLConnString(connString string) bool {
	return strings.HasPrefix(connString, "mysql://") ||
		strings.Contains(connString, "tcp(") ||
		strings.Contains(connString, "parseTime=")
}

func DetectDriver(connectionString *string) (configs.Driver, error) {
	original := *connectionString
	if original == "" {
		return "", fmt.Errorf("connection string is empty")
	}

	lower := strings.ToLower(original)

	switch {
	case isPostgresConnString(lower):
		return configs.DriverPostgres, nil

	case isMySQLConnString(lower):
		dsn, err := normalizeMySQLDSN(original)
		if err != nil {
			return "", err
		}
		*connectionString = dsn
		return configs.DriverMySQL, nil

	case isSQLiteConnString(lower):
		path, err := ReplaceTildeWithHomeDir(strings.TrimPrefix(original, sqliteURLPrefix))
		if err != nil {
			return "", err
		}
		*connectionString = path
		return configs.DriverSQLite, nil

	default:
		err := fmt.Errorf("unable to detect driver from connection string! Make sure your connection string is correct")
		logger.Errorln(err)
		return "", err
	}
}

const (
	sqliteURLPrefix = "sqlite://"
	mysqlURLPrefix  = "mysql://"
)

func normalizeMySQLDSN(connectionString string) (string, error) {
	if !strings.HasPrefix(strings.ToLower(connectionString), mysqlURLPrefix) {
		return withParseTime(connectionString)
	}

	parsedURL, err := url.Parse(connectionString)
	if err != nil {
		return "", fmt.Errorf("invalid mysql connection url: %w", err)
	}

	config := mysql.NewConfig()
	config.Net = "tcp"
	config.Addr = parsedURL.Host
	if !strings.Contains(config.Addr, ":") {
		config.Addr += ":3306"
	}
	config.DBName = strings.TrimPrefix(parsedURL.Path, "/")
	if parsedURL.User != nil {
		config.User = parsedURL.User.Username()
		config.Passwd, _ = parsedURL.User.Password()
	}
	config.Params = map[string]string{}
	for key, values := range parsedURL.Query() {
		config.Params[key] = values[0]
	}
	config.ParseTime = true
	return config.FormatDSN(), nil
}

func withParseTime(dsn string) (string, error) {
	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		return "", fmt.Errorf("invalid mysql connection string: %w", err)
	}
	config.ParseTime = true
	return config.FormatDSN(), nil
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
	logger.Warning("unknown data type: %s", dbType)
	return textInput
}
