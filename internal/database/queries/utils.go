package queries

import (
	"fmt"
	"strings"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/database"
	"github.com/biisal/rowsql/internal/logger"
)

func (b *Builder) getQuotedTableName(tableName string) string {
	switch b.driver {
	case configs.DriverMySQL:
		tableName = fmt.Sprintf("`%s`", tableName)
	case configs.DriverPostgres:
		tableName = fmt.Sprintf("\"%s\"", tableName)
	case configs.DriverSQLite:
		tableName = fmt.Sprintf("\"%s\"", tableName)
	}
	return tableName
}

func (b *Builder) placeHolder(n int) string {
	if b.driver == configs.DriverMySQL {
		return "?"
	}
	return fmt.Sprintf("$%d", n)
}

func (b *Builder) getAutoIncrementKeyword() string {
	switch b.driver {
	case configs.DriverPostgres:
		return "SERIAL"
	case configs.DriverSQLite:
		return "AUTOINCREMENT"
	default:
		return "AUTO_INCREMENT"
	}
}

func (b *Builder) formatColumnDefinition(input database.Input) (string, error) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s %s", input.ColName, input.DataType.Type)
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
		fmt.Fprintf(&sb, " %s", b.getAutoIncrementKeyword())
	}

	return sb.String(), nil
}
