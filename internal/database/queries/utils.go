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
	fmt.Fprintf(&sb, "%s %s", b.dialect.QuoteName(input.ColumnName), input.Type)
	if input.Size > 0 {
		fmt.Fprintf(&sb, "(%d)", input.Size)
	}
	if input.IsUnique {
		sb.WriteString(" UNIQUE")
	}
	if !input.IsNull {
		sb.WriteString(" NOT NULL")
	}
	if input.IsPk {
		sb.WriteString(" PRIMARY KEY")
	}
	if input.HasAutoIncrement {
		if !input.IsPk {
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
