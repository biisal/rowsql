package queries

import (
	"errors"
	"reflect"
	"testing"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database/models"
)

type Arg []any

func TestListTables(t *testing.T) {
	assertQuery := func(t testing.TB, driver configs.Driver, want string, err error) {
		builder := NewBuilder(driver, 10)
		t.Helper()
		query, qErr := builder.ListTables()
		if !errors.Is(qErr, err) {
			t.Errorf("expected error %#v but got %#v", err, qErr)
		}
		if query != want {
			t.Errorf("expected %s but got %s", want, query)
		}
	}

	tablesListTest := []struct {
		name   string
		driver configs.Driver
		want   string
		err    error
	}{
		{
			name:   "Postgress",
			driver: configs.DriverPostgres,
			want:   postgresMySQLTablesListQuery,
			err:    nil,
		},
		{
			name:   "MySQL",
			driver: configs.DriverMySQL,
			want:   postgresMySQLTablesListQuery,
			err:    nil,
		},
		{
			name:   "SQLite",
			driver: configs.DriverSQLite,
			want:   sqliteTablesListQuery,
			err:    nil,
		},
		{
			name:   "Empty driver",
			driver: "",
			want:   "",
			err:    ErrUnknownDriver,
		},
	}

	for _, tt := range tablesListTest {
		t.Run(tt.name, func(t *testing.T) {
			assertQuery(t, tt.driver, tt.want, tt.err)
		})
	}
}

func TestColumnsList(t *testing.T) {
	assertQuery := func(t testing.TB, driver configs.Driver, want string, wantArgs []any, err error) {
		builder := NewBuilder(driver, 10)
		t.Helper()
		query, args, qErr := builder.ColumnsList("users")
		if !errors.Is(qErr, err) {
			t.Errorf("expected error %#v but got %#v", err, qErr)
		}
		if query != want {
			t.Errorf("expected %s but got %s", want, query)
		}
		if len(args) != len(wantArgs) {
			t.Errorf("expected %d arguments but got %d", len(wantArgs), len(args))
		}
		for i, arg := range args {
			if arg != wantArgs[i] {
				t.Errorf("expected %#v but got %#v", wantArgs[i], arg)
			}
		}
	}

	columnsListTest := []struct {
		name   string
		driver configs.Driver
		want   string
		args   []any
		err    error
	}{
		{
			name:   "Postgress",
			driver: configs.DriverPostgres,
			want:   postgresColumnsListsQuery,
			args:   []any{"users"},
			err:    nil,
		},
		{
			name:   "MySQL",
			driver: configs.DriverMySQL,
			want:   mysqlColumnsListsQuery,
			args:   []any{"users"},
			err:    nil,
		},
		{
			name:   "SQLite",
			driver: configs.DriverSQLite,
			want:   sqliteColumnsListQuery,
			args:   []any{"users", "users"},
			err:    nil,
		},
		{
			name:   "Empty driver",
			driver: "",
			want:   "",
			args:   []any{},
			err:    ErrUnknownDriver,
		},
	}

	for _, tt := range columnsListTest {
		t.Run(tt.name, func(t *testing.T) {
			assertQuery(t, tt.driver, tt.want, tt.args, tt.err)
		})
	}
}

func assertQuery(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertArgs(t testing.TB, got, want Arg) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v , want %#v", got, want)
	}
}

func assertErr(t testing.TB, got, want error) {
	t.Helper()
	if want == nil {
		return
	}
	if got == nil {
		t.Fatal("got no error but expected one")
	}
	if got.Error() != want.Error() {
		t.Errorf("got %q , want %q", got.Error(), want.Error())
	}
}
func TestListRows(t *testing.T) {
	tests := []struct {
		name      string
		driver    configs.Driver
		tableName string
		orderCol  string
		orderBy   string
		limit     int
		offset    int
		want      string
		arg       Arg
		err       error
		maxLimit  int
	}{
		{
			name:      "Psql",
			driver:    configs.DriverPostgres,
			tableName: "users",
			orderCol:  "id",
			orderBy:   "DESC",
			limit:     10,
			offset:    10,
			want:      "SELECT * FROM users ORDER BY id DESC LIMIT $1 OFFSET $2",
			arg:       Arg{10, 10},
			maxLimit:  20,
		},
		{
			name:      "MySQL test",
			driver:    configs.DriverMySQL,
			tableName: "users",
			orderCol:  "id",
			orderBy:   "DESC",
			limit:     10,
			offset:    10,
			want:      "SELECT * FROM users ORDER BY id DESC LIMIT ? OFFSET ?",
			arg:       Arg{10, 10},
			maxLimit:  20,
		},
		{
			name:      "SQLite test",
			driver:    configs.DriverSQLite,
			tableName: "users",
			orderCol:  "id",
			orderBy:   "DESC",
			limit:     10,
			offset:    10,
			want:      "SELECT * FROM users ORDER BY id DESC LIMIT $1 OFFSET $2",
			arg:       Arg{10, 10},
			maxLimit:  20,
		},
		{
			name:      "SQLite test",
			driver:    configs.DriverSQLite,
			tableName: "users",
			orderCol:  "id",
			orderBy:   "DESC",
			limit:     10,
			offset:    10,
			want:      "SELECT * FROM users ORDER BY id DESC LIMIT $1 OFFSET $2",
			arg:       Arg{10, 10},
			maxLimit:  20,
		},
		{
			name:      "Invalid driver name test",
			driver:    configs.Driver("invald driver"),
			tableName: "users",
			orderCol:  "id",
			orderBy:   "DESC",
			limit:     10,
			offset:    10,

			err:      apperr.ErrorInvalidDriver,
			maxLimit: 20,
		},
		{
			name:      "Invalid orderby fallback to deafault ASC",
			driver:    configs.DriverPostgres,
			tableName: "users",
			orderCol:  "id",
			orderBy:   "invalid",
			limit:     10,
			offset:    10,
			want:      "SELECT * FROM users ORDER BY id ASC LIMIT $1 OFFSET $2",
			arg:       Arg{10, 10},
			maxLimit:  20,
		},
		{
			name:      "Zero limit and offset",
			driver:    configs.DriverPostgres,
			tableName: "users",
			orderCol:  "id",
			orderBy:   "ASC",
			limit:     0,
			offset:    0,
			want:      "SELECT * FROM users ORDER BY id ASC",
			arg:       Arg{},
			maxLimit:  20,
		},
		{
			name:      "Zero limit",
			driver:    configs.DriverPostgres,
			tableName: "users",
			orderCol:  "id",
			orderBy:   "ASC",
			limit:     0,
			offset:    10,
			want:      "SELECT * FROM users ORDER BY id ASC OFFSET $1",
			arg:       Arg{10},
			maxLimit:  20,
		},
		{
			name:      "Zero offset",
			driver:    configs.DriverPostgres,
			tableName: "users",
			orderCol:  "id",
			orderBy:   "ASC",
			limit:     10,
			offset:    0,
			want:      "SELECT * FROM users ORDER BY id ASC LIMIT $1",
			arg:       Arg{10},
			maxLimit:  20,
		},
		{
			name:      "Lowercase orderby DESC",
			driver:    configs.DriverPostgres,
			tableName: "users",
			orderCol:  "id",
			orderBy:   "desc",
			limit:     10,
			offset:    10,
			want:      "SELECT * FROM users ORDER BY id DESC LIMIT $1 OFFSET $2",
			arg:       Arg{10, 10},
			maxLimit:  20,
		},
		{
			name:      "Lowercase orderby ASC",
			driver:    configs.DriverPostgres,
			tableName: "users",
			orderCol:  "id",
			orderBy:   "asc",
			limit:     10,
			offset:    10,
			want:      "SELECT * FROM users ORDER BY id ASC LIMIT $1 OFFSET $2",
			arg:       Arg{10, 10},
			maxLimit:  20,
		},
		{
			name:      "Empty order column fallback",
			driver:    configs.DriverPostgres,
			tableName: "users",
			orderCol:  "",
			orderBy:   "DESC",
			limit:     10,
			offset:    10,
			want:      "SELECT * FROM users LIMIT $1 OFFSET $2",
			arg:       Arg{10, 10},
			maxLimit:  20,
		},
		{
			name:      "Empty order by fallback",
			driver:    configs.DriverPostgres,
			tableName: "users",
			orderCol:  "id",
			orderBy:   "",
			limit:     10,
			offset:    10,
			want:      "SELECT * FROM users ORDER BY id ASC LIMIT $1 OFFSET $2",
			arg:       Arg{10, 10},
			maxLimit:  20,
		},
		{
			name:      "space in table name",
			driver:    configs.DriverPostgres,
			tableName: "users table",
			orderCol:  "id",
			orderBy:   "",
			limit:     10,
			offset:    10,
			want:      "SELECT * FROM \"users table\" ORDER BY id ASC LIMIT $1 OFFSET $2",
			arg:       Arg{10, 10},
			maxLimit:  20,
		},
		{
			name:      "empty table name",
			driver:    configs.DriverPostgres,
			tableName: "",
			orderCol:  "id",
			orderBy:   "",
			limit:     10,
			offset:    10,
			err:       apperr.ErrorEmptyTableName,
			maxLimit:  20,
		},
		{
			name:      "Negative limit",
			driver:    configs.DriverPostgres,
			tableName: "users",
			orderCol:  "id",
			orderBy:   "ASC",
			limit:     -5,
			offset:    0,
			want:      "",
			err:       apperr.ErrorInvalidPagination,
			maxLimit:  20,
		},
		{
			name:      "Negative offset",
			driver:    configs.DriverPostgres,
			tableName: "users",
			orderCol:  "id",
			orderBy:   "ASC",
			limit:     10,
			offset:    -1,
			err:       apperr.ErrorInvalidPagination,
			maxLimit:  20,
		},
		{
			maxLimit:  10,
			name:      "Very large limit",
			driver:    configs.DriverPostgres,
			tableName: "users",
			orderCol:  "id",
			orderBy:   "ASC",
			limit:     1000000,
			offset:    0,
			err:       apperr.ErrorLimitTooLarge(10),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuilder(tt.driver, tt.maxLimit)
			query, args, err := builder.ListRows(tt.tableName, tt.orderCol, tt.orderBy, tt.limit, tt.offset)
			assertErr(t, err, tt.err)
			assertQuery(t, query, tt.want)
			assertArgs(t, args, tt.arg)
		})
	}
}

func TestInsertRow(t *testing.T) {
	tests := []struct {
		name      string
		driver    configs.Driver
		tableName string
		values    map[string]models.FormValue
		want      string
		args      []any
		err       error
	}{
		{
			name:      "Postgress",
			driver:    configs.DriverPostgres,
			tableName: "users",
			values: map[string]models.FormValue{
				"id":    {Value: "1"},
				"name":  {Value: "2"},
				"email": {Value: "3"},
			},
			want: "INSERT INTO users (id, name, email) VALUES ($1, $2, $3)",
			args: []any{"1", "2", "3"},
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuilder(tt.driver, 10)
			query, args, err := builder.InsertRow(tt.tableName, tt.values)
			assertErr(t, err, tt.err)
			assertQuery(t, query, tt.want)
			assertArgs(t, args, tt.args)
		})
	}

}
