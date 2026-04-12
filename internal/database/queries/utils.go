package queries

import (
	"fmt"
	"strings"

	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database/models"
	"github.com/biisal/rowsql/internal/logger"
)

func (b *builder) formatColumnDefinition(input models.ColValues) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s %s", b.dialect.QuoteName(input.Name), input.Type)
	if input.Size > 0 {
		fmt.Fprintf(&sb, "(%d)", input.Size)
	}
	if input.IsUnique {
		sb.WriteString(" UNIQUE")
	}
	if !input.IsNull {
		sb.WriteString(" NOT NULL")
	}
	if input.IsPK {
		sb.WriteString(" PRIMARY KEY")
	}
	if input.AutoIncrement {
		if !input.IsPK {
			logger.Errorln(apperr.ErrorInvalidAutoIncrement.Error())
			return "", apperr.ErrorInvalidAutoIncrement
		}
		fmt.Fprintf(&sb, " %s", b.dialect.AutoIncrementKeyword())
	}
	if input.Default != nil {
		var value = input.Default
		if str, ok := input.Default.(string); ok {
			value = b.dialect.QuoteName(str)
		}
		fmt.Fprintf(&sb, " DEFAULT %v", value)
	}

	return sb.String(), nil
}
