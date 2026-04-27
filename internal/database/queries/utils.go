package queries

import (
	"fmt"
	"strings"

	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database/models"
	"github.com/biisal/rowsql/internal/logger"
)

func (b *builder) formatColumnDefinition(input models.ColValue) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s %s", b.dialect.QuoteName(input.ColumnName), input.ColumnType.DataType)
	if input.ColumnType.HasSize {
		fmt.Fprintf(&sb, "(%d)", input.Size)
	}
	if input.ColumnType.IsUnique {
		sb.WriteString(" UNIQUE")
	}
	if !input.ColumnType.IsNull {
		sb.WriteString(" NOT NULL")
	}
	if input.ColumnType.IsPk {
		sb.WriteString(" PRIMARY KEY")
	}
	if input.ColumnType.HasAutoIncrement {
		if !input.ColumnType.IsPk {
			logger.Errorln(apperr.ErrorInvalidAutoIncrement.Error())
			return "", apperr.ErrorInvalidAutoIncrement
		}
		fmt.Fprintf(&sb, " %s", b.dialect.AutoIncrementKeyword())
	}
	if input.DefaultValue != nil {
		var value = input.DefaultValue
		if str, ok := input.DefaultValue.(string); ok {
			value = b.dialect.QuoteName(str)
		}
		fmt.Fprintf(&sb, " DEFAULT %v", value)
	}

	return sb.String(), nil
}
