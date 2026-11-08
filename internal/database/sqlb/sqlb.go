// Package sqlb builds SQL statements with correctly quoted identifiers and numbered bind variables.
package sqlb

import (
	"strconv"
	"strings"
)

type Dialect interface {
	QuoteIdentifier(identifier string) string
	BindVariable(position int) string
}

type Statement struct {
	dialect   Dialect
	text      strings.Builder
	Arguments []any
}

func New(dialect Dialect) *Statement {
	return &Statement{dialect: dialect}
}

func (s *Statement) Raw(sql string) *Statement {
	s.text.WriteString(sql)
	return s
}

func (s *Statement) Identifier(name string) *Statement {
	s.text.WriteString(s.dialect.QuoteIdentifier(name))
	return s
}

func (s *Statement) Argument(value any) *Statement {
	s.Arguments = append(s.Arguments, value)
	s.text.WriteString(s.dialect.BindVariable(len(s.Arguments)))
	return s
}

func (s *Statement) Equals(column string, value any) *Statement {
	s.Identifier(column)
	if value == nil {
		return s.Raw(" IS NULL")
	}
	return s.Raw(" = ").Argument(value)
}

func (s *Statement) WhereAllEqual(columns []string, values []any) *Statement {
	s.Raw(" WHERE ")
	for index, column := range columns {
		if index > 0 {
			s.Raw(" AND ")
		}
		s.Equals(column, values[index])
	}
	return s
}

func (s *Statement) LimitOffset(limit, offset int) *Statement {
	if limit > 0 {
		s.Raw(" LIMIT ").Argument(limit)
	}
	if offset > 0 {
		s.Raw(" OFFSET ").Argument(offset)
	}
	return s
}

func (s *Statement) Build() (string, []any) {
	return s.text.String(), s.Arguments
}

func QuoteWith(identifier string, quoteCharacter string) string {
	escaped := strings.ReplaceAll(identifier, quoteCharacter, quoteCharacter+quoteCharacter)
	return quoteCharacter + escaped + quoteCharacter
}

func QuoteStringLiteral(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}

func PositionalDollar(position int) string {
	return "$" + strconv.Itoa(position)
}
