package queries

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database"
	"github.com/biisal/rowsql/internal/database/models"
	"github.com/biisal/rowsql/internal/database/sqlb"
)

var validTableName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

var defaultExpressions = map[string]bool{
	"CURRENT_TIMESTAMP": true,
	"CURRENT_DATE":      true,
	"CURRENT_TIME":      true,
	"NULL":              true,
}

func (b *builder) CreateTable(tableName string, columns []models.ColValue) (string, error) {
	if !validTableName.MatchString(tableName) {
		return "", apperr.ErrorInvalidCreateTableName
	}

	definitions := make([]string, 0, len(columns))
	for _, column := range columns {
		if column.ColumnName == "" {
			continue
		}
		definition, err := b.columnDefinition(column)
		if err != nil {
			return "", err
		}
		definitions = append(definitions, definition)
	}
	if len(definitions) == 0 {
		return "", apperr.ErrorNoColumnsInArgs
	}

	query, _ := b.statement().Raw("CREATE TABLE ").Identifier(tableName).
		Raw(" (" + strings.Join(definitions, ", ") + ")").Build()
	return query, nil
}

func (b *builder) columnDefinition(column models.ColValue) (string, error) {
	dataType, found := b.findDataType(column.ColumnType.DataType)
	if !found {
		return "", apperr.ErrorUnsupportedDataType
	}

	typeName := dataType.Type
	autoIncrementSuffix := ""
	if column.ColumnType.HasAutoIncrement {
		if !column.ColumnType.IsPk {
			return "", apperr.ErrorInvalidAutoIncrement
		}
		var err error
		typeName, autoIncrementSuffix, err = b.dialect.ApplyAutoIncrement(typeName)
		if err != nil {
			return "", err
		}
	}

	var definition strings.Builder
	fmt.Fprintf(&definition, "%s %s", b.dialect.QuoteIdentifier(column.ColumnName), typeName)
	if dataType.HasSize && column.Size > 0 && !column.ColumnType.HasAutoIncrement {
		fmt.Fprintf(&definition, "(%d)", column.Size)
	}
	if column.ColumnType.IsUnique {
		definition.WriteString(" UNIQUE")
	}
	if !column.ColumnType.IsNull {
		definition.WriteString(" NOT NULL")
	}
	if column.ColumnType.IsPk {
		definition.WriteString(" PRIMARY KEY")
	}
	definition.WriteString(autoIncrementSuffix)

	if column.DefaultValue != nil {
		literal, err := defaultLiteral(column.DefaultValue)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&definition, " DEFAULT %s", literal)
	}
	return definition.String(), nil
}

func (b *builder) findDataType(requested string) (database.VarianDataType, bool) {
	numeric, text := database.DataTypesFor(b.dialect.Name())
	for _, dataType := range append(numeric, text...) {
		if strings.EqualFold(dataType.Type, strings.TrimSpace(requested)) {
			return dataType, true
		}
	}
	return database.VarianDataType{}, false
}

func defaultLiteral(value any) (string, error) {
	switch typed := value.(type) {
	case string:
		if defaultExpressions[strings.ToUpper(typed)] {
			return strings.ToUpper(typed), nil
		}
		return sqlb.QuoteStringLiteral(typed), nil
	case bool:
		return strconv.FormatBool(typed), nil
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64), nil
	case int:
		return strconv.Itoa(typed), nil
	case int64:
		return strconv.FormatInt(typed, 10), nil
	default:
		return "", apperr.ErrorUnsupportedDefault
	}
}
