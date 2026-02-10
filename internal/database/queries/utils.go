package queries

import (
	"fmt"
	"strings"

	"github.com/biisal/rowsql/internal/database"
	"github.com/biisal/rowsql/internal/logger"
)

func (b *builder) formatColumnDefinition(input database.Input) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s %s", b.dialect.QuoteName(input.ColName), input.DataType.Type)
	if input.DataType.HasSize {
		fmt.Fprintf(&sb, "(%d)", input.DataType.Size)
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
	if input.DataType.AutoIncrement {
		if !input.IsPK {
			logger.Error("Auto-increment can only be set on primary key columns")
			return "", fmt.Errorf("auto-increment can only be set on primary key columns")
		}
		fmt.Fprintf(&sb, " %s", b.dialect.AutoIncrementKeyword())
	}

	return sb.String(), nil
}
