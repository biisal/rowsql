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
		values    []models.RowItem
		want      string
		args      []any
		err       error
	}{
		{
			name:      "Postgress",
			driver:    configs.DriverPostgres,
			tableName: "users",
			values: []models.RowItem{
				{
					ColumnName: "id",
					Value:      "1",
					Type:       "int",
				},
				{
					ColumnName: "name",
					Value:      "2",
					Type:       "string",
				},
				{
					ColumnName: "email",
					Value:      "3",
					Type:       "string",
				},
			},
			want: "INSERT INTO users (id, name, email) VALUES ($1, $2, $3)",
			args: []any{"1", "2", "3"},
			err:  nil,
		},
		{
			name:      "MySQL",
			driver:    configs.DriverMySQL,
			tableName: "users",
			values: []models.RowItem{
				{
					ColumnName: "id",
					Value:      "1",
					Type:       "int",
				},
				{
					ColumnName: "name",
					Value:      "2",
					Type:       "string",
				},
				{
					ColumnName: "email",
					Value:      "3",
					Type:       "string",
				},
			},
			want: "INSERT INTO users (id, name, email) VALUES (?, ?, ?)",
			args: []any{"1", "2", "3"},
			err:  nil,
		},
		{
			name:      "SQlite",
			driver:    configs.DriverSQLite,
			tableName: "users",
			values: []models.RowItem{
				{
					ColumnName: "id",
					Value:      "1",
					Type:       "int",
				},
				{
					ColumnName: "name",
					Value:      "2",
					Type:       "string",
				},
				{
					ColumnName: "email",
					Value:      "3",
					Type:       "string",
				},
			},
			want: "INSERT INTO users (id, name, email) VALUES ($1, $2, $3)",
			args: []any{"1", "2", "3"},
			err:  nil,
		},
		{
			name:      "Invalid driver",
			driver:    configs.Driver("invalid"),
			tableName: "users",
			values:    []models.RowItem{},
			err:       apperr.ErrorInvalidDriver,
		},
		{
			name:      "empty table name",
			driver:    configs.DriverPostgres,
			tableName: "",
			values:    []models.RowItem{},
			err:       apperr.ErrorEmptyTableName,
		},
		{
			name:      "no values postgres",
			driver:    configs.DriverPostgres,
			args:      []any{},
			tableName: "users",
			values:    []models.RowItem{},
			want:      "INSERT INTO users DEFAULT VALUES",
		},
		{
			name:      "no values mysql",
			driver:    configs.DriverMySQL,
			args:      []any{},
			tableName: "users",
			values:    []models.RowItem{},
			// MySQL requires the empty brackets syntax
			want: "INSERT INTO users () VALUES ()",
		},
		{
			name:      "Postgres Table with Spaces",
			driver:    configs.DriverPostgres,
			tableName: "order details",
			values: []models.RowItem{
				{ColumnName: "id", Value: "101", Type: "int"},
			},
			want: `INSERT INTO "order details" (id) VALUES ($1)`,
			args: []any{"101"},
		},
		{
			name:      "MySQL Table with Spaces",
			driver:    configs.DriverMySQL,
			tableName: "order details",
			values: []models.RowItem{
				{ColumnName: "id", Value: "101", Type: "int"},
			},
			want: "INSERT INTO `order details` (id) VALUES (?)",
			args: []any{"101"},
		},
		{
			name:      "SQL Injection Resistance",
			driver:    configs.DriverPostgres,
			tableName: "users",
			values: []models.RowItem{
				{
					ColumnName: "name",
					Value:      "'; DROP TABLE users; --",
					Type:       "string",
				},
			},
			want: "INSERT INTO users (name) VALUES ($1)",
			args: []any{"'; DROP TABLE users; --"},
		},
		{
			name:      "Postgres Incremental Placeholders",
			driver:    configs.DriverPostgres,
			tableName: "products",
			values: []models.RowItem{
				{ColumnName: "a", Value: "v1", Type: "string"},
				{ColumnName: "b", Value: "v2", Type: "string"},
				{ColumnName: "c", Value: "v3", Type: "string"},
			},
			want: "INSERT INTO products (a, b, c) VALUES ($1, $2, $3)",
			args: []any{"v1", "v2", "v3"},
		},
		{
			name:      "Single column Postgres",
			driver:    configs.DriverPostgres,
			tableName: "tags",
			values: []models.RowItem{
				{ColumnName: "name", Value: "golang", Type: "string"},
			},
			want: "INSERT INTO tags (name) VALUES ($1)",
			args: []any{"golang"},
		},
		{
			name:      "Many columns Postgres",
			driver:    configs.DriverPostgres,
			tableName: "wide_table",
			values: []models.RowItem{
				{ColumnName: "c1", Value: "v1"},
				{ColumnName: "c2", Value: "v2"},
				{ColumnName: "c3", Value: "v3"},
				{ColumnName: "c4", Value: "v4"},
				{ColumnName: "c5", Value: "v5"},
				{ColumnName: "c6", Value: "v6"},
				{ColumnName: "c7", Value: "v7"},
				{ColumnName: "c8", Value: "v8"},
				{ColumnName: "c9", Value: "v9"},
				{ColumnName: "c10", Value: "v10"},
				{ColumnName: "c11", Value: "v11"},
			},
			want: "INSERT INTO wide_table (c1, c2, c3, c4, c5, c6, c7, c8, c9, c10, c11) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
			args: []any{"v1", "v2", "v3", "v4", "v5", "v6", "v7", "v8", "v9", "v10", "v11"},
		},
		{
			name:      "Duplicate columns",
			driver:    configs.DriverPostgres,
			tableName: "users",
			values: []models.RowItem{
				{ColumnName: "email", Value: "a@b.com"},
				{ColumnName: "email", Value: "b@a.com"},
			},
			err: apperr.ErrorDuplicateColumn,
		},
		{
			name:      "Valid JSON Object Postgres",
			driver:    configs.DriverPostgres,
			tableName: "settings",
			values: []models.RowItem{
				{
					ColumnName: "metadata",
					Value:      `{"theme": "dark", "notifications": true}`,
					Type:       "json",
				},
			},
			want: "INSERT INTO settings (metadata) VALUES ($1)",
			args: []any{map[string]any{"theme": "dark", "notifications": true}},
			err:  nil,
		},
		{
			name:      "Valid JSON Array Postgres",
			driver:    configs.DriverPostgres,
			tableName: "posts",
			values: []models.RowItem{
				{
					ColumnName: "tags",
					Value:      `["golang", "sql", "backend"]`,
					Type:       "json",
				},
			},
			want: "INSERT INTO posts (tags) VALUES ($1)",
			args: []any{[]any{"golang", "sql", "backend"}},
			err:  nil,
		},
		{
			name:      "Malformed JSON error",
			driver:    configs.DriverPostgres,
			tableName: "users",
			values: []models.RowItem{
				{
					ColumnName: "extra_data",
					Value:      `{"missing_bracket": true`,
					Type:       "json",
				},
			},
			err: apperr.ErrorInvalidJSON,
		},
		{
			name:      "driver validation happens before values validation",
			driver:    configs.Driver("oracle"),
			tableName: "users",
			values: []models.RowItem{
				{
					ColumnName: "extra_data",
					Value:      ``,
					Type:       "json",
				},
			},
			err: apperr.ErrorInvalidDriver,
		},
		{
			name:      "Empty string as JSON error",
			driver:    configs.DriverPostgres,
			tableName: "users",
			values: []models.RowItem{
				{
					ColumnName: "extra_data",
					Value:      ``,
					Type:       "json",
				},
			},
			err: apperr.ErrorInvalidJSON,
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

func TestGetRows(t *testing.T) {
	tests := []struct {
		name      string
		driver    configs.Driver
		tableName string
		limit     int
		offset    int
		want      string
		arg       Arg
		err       error
	}{
		{
			name:      "Psql",
			driver:    configs.DriverPostgres,
			tableName: "users",
			limit:     10,
			offset:    10,
			want:      "SELECT * FROM users LIMIT $1 OFFSET $2",
			arg:       Arg{10, 10},
		},
		{
			name:      "Mysql",
			driver:    configs.DriverMySQL,
			tableName: "users",
			limit:     10,
			offset:    10,
			want:      "SELECT * FROM users LIMIT ? OFFSET ?",
			arg:       Arg{10, 10},
		},
		{
			name:      "Sqlite",
			driver:    configs.DriverSQLite,
			tableName: "users",
			limit:     10,
			offset:    10,
			want:      "SELECT * FROM users LIMIT $1 OFFSET $2",
			arg:       Arg{10, 10},
		},
		{
			name:      "zero offset",
			driver:    configs.DriverSQLite,
			tableName: "users",
			limit:     10,
			offset:    0,
			want:      "SELECT * FROM users LIMIT $1",
			arg:       Arg{10},
		},
		{
			name:      "zero limit",
			driver:    configs.DriverSQLite,
			tableName: "users",
			limit:     0,
			offset:    10,
			err:       apperr.ErrorInvalidPagination,
		},
		{
			name:      "Nagative offset",
			driver:    configs.DriverSQLite,
			tableName: "users",
			limit:     0,
			offset:    -1,
			err:       apperr.ErrorInvalidPagination,
		},
		{
			name:   "Unknown driver",
			driver: configs.Driver("unknown"),
			err:    apperr.ErrorInvalidDriver,
		},
		{
			name:      "whitespace table name",
			driver:    configs.DriverMySQL,
			tableName: "   ",
			limit:     10,
			offset:    0,
			err:       apperr.ErrorInvalidTableName,
		},
		{
			name:      "driver validation happens before pagination",
			driver:    configs.Driver("oracle"),
			tableName: "users",
			limit:     -1,
			offset:    -1,
			err:       apperr.ErrorInvalidDriver,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuilder(tt.driver, 10)
			query, args, err := builder.GetRows(tt.tableName, tt.limit, tt.offset)
			assertErr(t, err, tt.err)
			assertQuery(t, query, tt.want)
			assertArgs(t, args, tt.arg)
		})
	}
}

func TestWhereClause(t *testing.T) {
	tests := []struct {
		name    string
		driver  configs.Driver
		cols    []models.ListDataCol
		rows    []any
		argsIdx int
		want    string
		args    []any
		err     error
	}{
		{
			name:   "Postgress",
			driver: configs.DriverPostgres,
			cols: []models.ListDataCol{
				{
					ColumnName:       "id",
					DataType:         "int",
					IsUnique:         false,
					Value:            1,
					InputType:        "text",
					HasAutoIncrement: false,
				},
			},
			rows:    []any{1},
			want:    "id=$1",
			argsIdx: 1,
			args:    []any{1},
		},
		{
			name:   "MySQL",
			driver: configs.DriverMySQL,
			cols: []models.ListDataCol{
				{
					ColumnName:       "id",
					DataType:         "int",
					IsUnique:         false,
					Value:            1,
					InputType:        "text",
					HasAutoIncrement: false,
				},
			},
			rows:    []any{1},
			want:    "id=?",
			argsIdx: 1,
			args:    []any{1},
		},
		{
			name:   "SQLite",
			driver: configs.DriverSQLite,
			cols: []models.ListDataCol{
				{
					ColumnName:       "id",
					DataType:         "int",
					IsUnique:         false,
					Value:            1,
					InputType:        "text",
					HasAutoIncrement: false,
				},
			},
			rows:    []any{1},
			want:    "id=$1",
			argsIdx: 1,
			args:    []any{1},
		},
		{
			name:   "Sqlite zero index",
			driver: configs.DriverSQLite,
			cols: []models.ListDataCol{
				{
					ColumnName: "id",
					Value:      1,
				},
			},
			rows:    []any{1},
			argsIdx: 0,
			// err: apperr.Err,
			err: apperr.ErrorInvalidPlaceHolderIndex,
		},
		{
			name:   "Sqlite multiple rows and cols",
			driver: configs.DriverSQLite,
			cols: []models.ListDataCol{
				{
					ColumnName: "id",
					Value:      1,
				},
				{
					ColumnName: "name",
					Value:      "test",
				},
				{
					ColumnName: "email",
					Value:      "test",
				},
			},

			rows:    []any{1, "test", "test"},
			argsIdx: 1,
			want:    "id=$1 AND name=$2 AND email=$3",
			args:    []any{1, "test", "test"},
		},
		{
			name:   "Mysql multiple rows and cols",
			driver: configs.DriverMySQL,
			cols: []models.ListDataCol{
				{
					ColumnName: "id",
					Value:      1,
				},
				{
					ColumnName: "name",
					Value:      "test",
				},
				{
					ColumnName: "email",
					Value:      "test",
				},
			},

			rows:    []any{1, "test", "test"},
			argsIdx: 1,
			want:    "id=? AND name=? AND email=?",
			args:    []any{1, "test", "test"},
		},
		{
			name:   "Sqlite less cols then rows",
			driver: configs.DriverSQLite,
			cols: []models.ListDataCol{
				{
					ColumnName: "id",
					Value:      1,
				},
				{
					ColumnName: "name",
					Value:      "test",
				},
			},
			rows:    []any{1, "test", "test"},
			argsIdx: 1,
			err:     apperr.ErrorNotSameRowColsSize,
		},
		{
			name:   "Sqlite unique cols",
			driver: configs.DriverSQLite,
			cols: []models.ListDataCol{
				{
					ColumnName: "id",
					Value:      1,
					IsUnique:   true,
				},
				{
					ColumnName: "name",
					Value:      "test",
				},
			},
			rows:    []any{1, "test"},
			args:    []any{1},
			argsIdx: 1,
			want:    "id=$1",
		},
		{
			name:   "Sqlite multiple unique cols",
			driver: configs.DriverSQLite,
			cols: []models.ListDataCol{
				{
					ColumnName: "id",
					Value:      1,
					IsUnique:   true,
				},
				{
					ColumnName: "name",
					Value:      "test",
				},
				{
					ColumnName: "email",
					Value:      "test",
					IsUnique:   true,
				},
			},
			rows:    []any{1, "test", "test"},
			args:    []any{1},
			argsIdx: 1,
			want:    "id=$1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuilder(tt.driver, 10)
			query, args, err := builder.WhereCluse(tt.cols, tt.rows, tt.argsIdx)
			assertArgs(t, args, tt.args)
			assertErr(t, err, tt.err)
			assertQuery(t, query, tt.want)
		})
	}
}
