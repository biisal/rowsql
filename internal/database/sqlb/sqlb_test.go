package sqlb

import (
	"reflect"
	"testing"
)

type dollarDialect struct{}

func (dollarDialect) QuoteIdentifier(identifier string) string { return QuoteWith(identifier, `"`) }
func (dollarDialect) BindVariable(position int) string         { return PositionalDollar(position) }

type questionDialect struct{}

func (questionDialect) QuoteIdentifier(identifier string) string { return QuoteWith(identifier, "`") }
func (questionDialect) BindVariable(int) string                  { return "?" }

func TestQuoteWithEscapesEmbeddedQuotes(t *testing.T) {
	tests := []struct {
		identifier string
		quote      string
		want       string
	}{
		{"user", `"`, `"user"`},
		{"Order", `"`, `"Order"`},
		{`a"b`, `"`, `"a""b"`},
		{"a`b", "`", "`a``b`"},
		{"naïve", `"`, `"naïve"`},
		{`"x"; DROP TABLE y; --"`, `"`, `"""x""; DROP TABLE y; --"""`},
	}
	for _, tt := range tests {
		if got := QuoteWith(tt.identifier, tt.quote); got != tt.want {
			t.Errorf("QuoteWith(%q) = %s, want %s", tt.identifier, got, tt.want)
		}
	}
}

func TestEqualsWritesIsNullForNil(t *testing.T) {
	sql, arguments := New(dollarDialect{}).Raw("SELECT 1").
		WhereAllEqual([]string{"id", "note"}, []any{1, nil}).Build()

	if want := `SELECT 1 WHERE "id" = $1 AND "note" IS NULL`; sql != want {
		t.Fatalf("sql = %s, want %s", sql, want)
	}
	if !reflect.DeepEqual(arguments, []any{1}) {
		t.Fatalf("arguments = %v", arguments)
	}
}

func TestArgumentsAreNumberedInOrder(t *testing.T) {
	statement := New(dollarDialect{}).Raw("UPDATE t SET a = ").Argument("x").Raw(", b = ").Argument("y").
		WhereAllEqual([]string{"id"}, []any{7})
	sql, arguments := statement.Build()

	if want := `UPDATE t SET a = $1, b = $2 WHERE "id" = $3`; sql != want {
		t.Fatalf("sql = %s, want %s", sql, want)
	}
	if !reflect.DeepEqual(arguments, []any{"x", "y", 7}) {
		t.Fatalf("arguments = %v", arguments)
	}
}

func TestLimitOffset(t *testing.T) {
	tests := []struct {
		name          string
		limit, offset int
		wantSQL       string
		wantArguments []any
	}{
		{"both", 10, 20, " LIMIT $1 OFFSET $2", []any{10, 20}},
		{"limit only", 10, 0, " LIMIT $1", []any{10}},
		{"offset only", 0, 5, " OFFSET $1", []any{5}},
		{"neither", 0, 0, "", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, arguments := New(dollarDialect{}).LimitOffset(tt.limit, tt.offset).Build()
			if sql != tt.wantSQL || !reflect.DeepEqual(arguments, tt.wantArguments) {
				t.Fatalf("got (%q, %v), want (%q, %v)", sql, arguments, tt.wantSQL, tt.wantArguments)
			}
		})
	}
}

func TestQuestionMarkDialect(t *testing.T) {
	sql, _ := New(questionDialect{}).Identifier("a").Raw(" = ").Argument(1).Raw(" AND ").Equals("b", 2).Build()
	if want := "`a` = ? AND `b` = ?"; sql != want {
		t.Fatalf("sql = %s, want %s", sql, want)
	}
}

func TestQuoteStringLiteral(t *testing.T) {
	if got := QuoteStringLiteral("it's"); got != "'it''s'" {
		t.Fatalf("got %s", got)
	}
}
