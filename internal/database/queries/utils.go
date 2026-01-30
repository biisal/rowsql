package queries

import (
	"fmt"
	"strings"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database"
	"github.com/biisal/rowsql/internal/logger"
)

func (b *Builder) getQuotedTableName(tableName string) (string, error) {
	if !strings.Contains(tableName, " ") {
		switch b.driver {
		case configs.DriverPostgres, configs.DriverMySQL, configs.DriverSQLite:
			return tableName, nil
		default:
			return "", apperr.ErrorInvalidDriver
		}
	}
	switch b.driver {
	case configs.DriverMySQL:
		return fmt.Sprintf("`%s`", tableName), nil
	case configs.DriverPostgres:
		return fmt.Sprintf("\"%s\"", tableName), nil
	case configs.DriverSQLite:
		return fmt.Sprintf("\"%s\"", tableName), nil
	default:
		return "", apperr.ErrorInvalidDriver
	}
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
