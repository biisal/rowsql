package queries

import (
	"errors"
	"reflect"
	"testing"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database/models"
	"github.com/biisal/rowsql/internal/database/rowkey"
)

func newTestBuilder(t *testing.T, driver configs.Driver) Builder {
	t.Helper()
	builder, err := NewBuilder(driver, 10)
	if err != nil {
		t.Fatal(err)
	}
	return builder
}

func primaryKey(columns []string, values ...any) rowkey.Key {
	return rowkey.Key{Kind: rowkey.KindPrimary, Columns: columns, Values: values}
}

func TestListRows(t *testing.T) {
	tests := []struct {
		name     string
		driver   configs.Driver
		table    string
		options  ListRowsOptions
		wantSQL  string
		wantArgs []any
		wantErr  error
	}{
		{
			name: "postgres with pagination", driver: configs.DriverPostgres, table: "users",
			options: ListRowsOptions{Limit: 10, Offset: 20},
			wantSQL: `SELECT * FROM "users" LIMIT $1 OFFSET $2`, wantArgs: []any{10, 20},
		},
		{
			name: "mysql uses question marks and backticks", driver: configs.DriverMySQL, table: "users",
			options: ListRowsOptions{Limit: 10, Offset: 20},
			wantSQL: "SELECT * FROM `users` LIMIT ? OFFSET ?", wantArgs: []any{10, 20},
		},
		{
			name: "reserved word table", driver: configs.DriverSQLite, table: "order",
			options: ListRowsOptions{Limit: 5},
			wantSQL: `SELECT * FROM "order" LIMIT $1`, wantArgs: []any{5},
		},
		{
			name:   "order descending with key tie breaker",
			driver: configs.DriverPostgres,
			table:  "users",
			options: ListRowsOptions{
				OrderColumn:         "name",
				Descending:          true,
				Limit:               10,
				DefaultOrderColumns: []string{"id"},
			},
			wantSQL:  `SELECT * FROM "users" ORDER BY "name" DESC, "id" LIMIT $1`,
			wantArgs: []any{10},
		},
		{
			name: "default order by key", driver: configs.DriverPostgres, table: "memberships",
			options: ListRowsOptions{Limit: 10, DefaultOrderColumns: []string{"user_id", "role_id"}},
			wantSQL: `SELECT * FROM "memberships" ORDER BY "user_id", "role_id" LIMIT $1`, wantArgs: []any{10},
		},
		{
			name:     "order column that is also a key column is not repeated",
			driver:   configs.DriverPostgres,
			table:    "users",
			options:  ListRowsOptions{OrderColumn: "id", Limit: 10, DefaultOrderColumns: []string{"id"}},
			wantSQL:  `SELECT * FROM "users" ORDER BY "id" LIMIT $1`,
			wantArgs: []any{10},
		},
		{
			name: "sqlite physical row id", driver: configs.DriverSQLite, table: "notes",
			options: ListRowsOptions{Limit: 10, IncludePhysicalRowID: true, DefaultOrderColumns: []string{"rowid"}},
			wantSQL: `SELECT rowid AS "_rowsql_row_id", * FROM "notes" ORDER BY "rowid" LIMIT $1`, wantArgs: []any{10},
		},
		{
			name: "postgres physical row id", driver: configs.DriverPostgres, table: "notes",
			options: ListRowsOptions{Limit: 10, IncludePhysicalRowID: true},
			wantSQL: `SELECT ctid::text AS "_rowsql_row_id", * FROM "notes" LIMIT $1`, wantArgs: []any{10},
		},
		{
			name: "mysql has no physical row id", driver: configs.DriverMySQL, table: "notes",
			options: ListRowsOptions{Limit: 10, IncludePhysicalRowID: true},
			wantErr: apperr.ErrorPhysicalRowIDUnsupported,
		},
		{
			name:    "empty table",
			driver:  configs.DriverPostgres,
			table:   "",
			options: ListRowsOptions{Limit: 1},
			wantErr: apperr.ErrorEmptyTableName,
		},
		{
			name:    "negative offset",
			driver:  configs.DriverPostgres,
			table:   "t",
			options: ListRowsOptions{Limit: 1, Offset: -1},
			wantErr: apperr.ErrorInvalidPagination,
		},
		{
			name:    "limit above maximum",
			driver:  configs.DriverPostgres,
			table:   "t",
			options: ListRowsOptions{Limit: 11},
			wantErr: apperr.ErrorLimitTooLarge(10),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, args, err := newTestBuilder(t, tt.driver).ListRows(tt.table, tt.options)
			if tt.wantErr != nil {
				if err == nil || err.Error() != tt.wantErr.Error() {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if sql != tt.wantSQL || !reflect.DeepEqual(args, tt.wantArgs) {
				t.Fatalf("got (%s, %v), want (%s, %v)", sql, args, tt.wantSQL, tt.wantArgs)
			}
		})
	}
}

func TestSelectRow(t *testing.T) {
	sql, args, err := newTestBuilder(t, configs.DriverPostgres).
		SelectRow("memberships", primaryKey([]string{"user_id", "role_id"}, 1, nil))
	if err != nil {
		t.Fatal(err)
	}
	if want := `SELECT * FROM "memberships" WHERE "user_id" = $1 AND "role_id" IS NULL LIMIT $2`; sql != want {
		t.Fatalf("sql = %s, want %s", sql, want)
	}
	if !reflect.DeepEqual(args, []any{1, 1}) {
		t.Fatalf("args = %v", args)
	}
}

func TestDeleteRow(t *testing.T) {
	tests := []struct {
		name     string
		driver   configs.Driver
		key      rowkey.Key
		wantSQL  string
		wantArgs []any
	}{
		{
			name: "primary key", driver: configs.DriverPostgres,
			key:     primaryKey([]string{"id"}, 7),
			wantSQL: `DELETE FROM "users" WHERE "id" = $1`, wantArgs: []any{7},
		},
		{
			name: "composite primary key", driver: configs.DriverSQLite,
			key:     primaryKey([]string{"user_id", "role_id"}, 1, 20),
			wantSQL: `DELETE FROM "users" WHERE "user_id" = $1 AND "role_id" = $2`, wantArgs: []any{1, 20},
		},
		{
			name: "postgres physical row id", driver: configs.DriverPostgres,
			key:     rowkey.Key{Kind: rowkey.KindPhysical, Columns: []string{"ctid"}, Values: []any{"(0,1)"}},
			wantSQL: `DELETE FROM "users" WHERE "ctid" = $1::text::tid`, wantArgs: []any{"(0,1)"},
		},
		{
			name: "mysql all columns limits to one row and matches null", driver: configs.DriverMySQL,
			key:     rowkey.Key{Kind: rowkey.KindAllColumns, Columns: []string{"a", "b"}, Values: []any{1, nil}},
			wantSQL: "DELETE FROM `users` WHERE `a` = ? AND `b` IS NULL LIMIT 1", wantArgs: []any{1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, args, err := newTestBuilder(t, tt.driver).DeleteRow("users", tt.key)
			if err != nil {
				t.Fatal(err)
			}
			if sql != tt.wantSQL || !reflect.DeepEqual(args, tt.wantArgs) {
				t.Fatalf("got (%s, %v), want (%s, %v)", sql, args, tt.wantSQL, tt.wantArgs)
			}
		})
	}
}

func TestDeleteRowRejectsInvalidKeys(t *testing.T) {
	builder := newTestBuilder(t, configs.DriverMySQL)
	invalidKeys := []rowkey.Key{
		{},
		{Kind: rowkey.KindPrimary, Columns: []string{"id"}},
		{Kind: rowkey.KindPhysical, Columns: []string{"rowid"}, Values: []any{1}},
	}
	for _, key := range invalidKeys {
		if _, _, err := builder.DeleteRow("users", key); !errors.Is(err, rowkey.ErrInvalidKey) {
			t.Errorf("key %+v: error = %v", key, err)
		}
	}
}

func TestUpdateRowNumbersArgumentsAcrossSetAndWhere(t *testing.T) {
	changes := []models.ColValue{{ColumnName: "name", Value: "new"}, {ColumnName: "age", Value: 5}}
	sql, args, err := newTestBuilder(t, configs.DriverPostgres).
		UpdateRow("people", changes, primaryKey([]string{"id"}, 9))
	if err != nil {
		t.Fatal(err)
	}
	if want := `UPDATE "people" SET "name" = $1, "age" = $2 WHERE "id" = $3`; sql != want {
		t.Fatalf("sql = %s, want %s", sql, want)
	}
	if !reflect.DeepEqual(args, []any{"new", 5, 9}) {
		t.Fatalf("args = %v", args)
	}
}

func TestUpdateRowErrors(t *testing.T) {
	builder := newTestBuilder(t, configs.DriverPostgres)
	key := primaryKey([]string{"id"}, 1)

	if _, _, err := builder.UpdateRow("t", nil, key); !errors.Is(err, apperr.ErrorNoValuesProvided) {
		t.Errorf("no changes: %v", err)
	}
	duplicates := []models.ColValue{{ColumnName: "a"}, {ColumnName: "a"}}
	if _, _, err := builder.UpdateRow("t", duplicates, key); !errors.Is(err, apperr.ErrorDuplicateColumn) {
		t.Errorf("duplicates: %v", err)
	}
	if _, _, err := builder.UpdateRow(
		"",
		[]models.ColValue{{ColumnName: "a"}},
		key,
	); !errors.Is(
		err,
		apperr.ErrorInvalidTableName,
	) {
		t.Errorf("empty table: %v", err)
	}
}

func TestInsertRow(t *testing.T) {
	values := []models.ColValue{{ColumnName: "name", Value: "Ann"}, {ColumnName: "age", Value: 3}}
	tests := []struct {
		driver   configs.Driver
		wantSQL  string
		wantArgs []any
	}{
		{configs.DriverPostgres, `INSERT INTO "people" ("name", "age") VALUES ($1, $2)`, []any{"Ann", 3}},
		{configs.DriverMySQL, "INSERT INTO `people` (`name`, `age`) VALUES (?, ?)", []any{"Ann", 3}},
	}
	for _, tt := range tests {
		sql, args, err := newTestBuilder(t, tt.driver).InsertRow("people", values)
		if err != nil {
			t.Fatal(err)
		}
		if sql != tt.wantSQL || !reflect.DeepEqual(args, tt.wantArgs) {
			t.Errorf("%s: got (%s, %v)", tt.driver, sql, args)
		}
	}
}

func TestInsertRowWithoutValuesUsesDefaults(t *testing.T) {
	tests := map[configs.Driver]string{
		configs.DriverPostgres: `INSERT INTO "t" DEFAULT VALUES`,
		configs.DriverSQLite:   `INSERT INTO "t" DEFAULT VALUES`,
		configs.DriverMySQL:    "INSERT INTO `t` () VALUES ()",
	}
	for driver, want := range tests {
		sql, _, err := newTestBuilder(t, driver).InsertRow("t", nil)
		if err != nil || sql != want {
			t.Errorf("%s: got (%s, %v), want %s", driver, sql, err, want)
		}
	}
}

func TestInsertRowJSONColumns(t *testing.T) {
	builder := newTestBuilder(t, configs.DriverPostgres)
	jsonColumn := func(value any) []models.ColValue {
		return []models.ColValue{{ColumnName: "data", Value: value, ColumnType: models.ColType{DataType: "jsonb"}}}
	}

	for _, valid := range []any{`{"a":1}`, `[1,2]`, `"text"`, `5`, nil} {
		if _, args, err := builder.InsertRow(
			"t",
			jsonColumn(valid),
		); err != nil ||
			!reflect.DeepEqual(args, []any{valid}) {
			t.Errorf("valid %v: args=%v err=%v", valid, args, err)
		}
	}
	for _, invalid := range []any{`{broken`, 12, map[string]any{"a": 1}} {
		if _, _, err := builder.InsertRow("t", jsonColumn(invalid)); !errors.Is(err, apperr.ErrorInvalidJSON) {
			t.Errorf("invalid %v: err=%v", invalid, err)
		}
	}
}

func TestInsertRowRejectsDuplicateColumns(t *testing.T) {
	values := []models.ColValue{{ColumnName: "a"}, {ColumnName: "a"}}
	if _, _, err := newTestBuilder(
		t,
		configs.DriverPostgres,
	).InsertRow("t", values); !errors.Is(
		err,
		apperr.ErrorDuplicateColumn,
	) {
		t.Fatalf("error = %v", err)
	}
}

func TestQuotingNeverLetsIdentifiersEscape(t *testing.T) {
	hostile := `x"; DROP TABLE users; --`
	sql, _, err := newTestBuilder(t, configs.DriverPostgres).DeleteRow(hostile, primaryKey([]string{hostile}, 1))
	if err != nil {
		t.Fatal(err)
	}
	want := `DELETE FROM "x""; DROP TABLE users; --" WHERE "x""; DROP TABLE users; --" = $1`
	if sql != want {
		t.Fatalf("sql = %s, want %s", sql, want)
	}
}

func TestDeleteTableAndCountRows(t *testing.T) {
	builder := newTestBuilder(t, configs.DriverMySQL)
	if got := builder.DeleteTable("user"); got != "DROP TABLE `user`" {
		t.Errorf("DeleteTable = %s", got)
	}
	count, err := builder.CountRows("user")
	if err != nil || count != "SELECT COUNT(*) FROM `user`" {
		t.Errorf("CountRows = %s, %v", count, err)
	}
	if _, err := builder.CountRows(""); !errors.Is(err, apperr.ErrorEmptyTableName) {
		t.Errorf("CountRows empty: %v", err)
	}
}
