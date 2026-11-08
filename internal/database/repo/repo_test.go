package repo

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database/models"
	"github.com/biisal/rowsql/internal/database/queries"
	"github.com/jmoiron/sqlx"

	_ "modernc.org/sqlite"
)

func newSQLiteRepo(t *testing.T, schemaStatements ...string) *Queries {
	t.Helper()
	db, err := sqlx.Open("sqlite", filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	for _, statement := range schemaStatements {
		if _, err := db.Exec(statement); err != nil {
			t.Fatalf("%s: %v", statement, err)
		}
	}

	builder, err := queries.NewBuilder(configs.DriverSQLite, 10)
	if err != nil {
		t.Fatal(err)
	}
	return New(db, configs.DriverSQLite, builder, 10)
}

func listTableRows(t *testing.T, repository *Queries, tableName string) []models.RowSet {
	t.Helper()
	rows, err := repository.ListRows(context.Background(), ListDataProps{TableName: tableName, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	return rows
}

func countRows(t *testing.T, repository *Queries, tableName string) int {
	t.Helper()
	count, err := repository.GetRowCount(context.Background(), tableName)
	if err != nil {
		t.Fatal(err)
	}
	return count
}

func columnValue(row models.RowSet, columnName string) any {
	for _, column := range row.Columns {
		if column.ColumnName == columnName {
			return column.Value
		}
	}
	return nil
}

func TestUpdateRowChangesNonKeyColumn(t *testing.T) {
	repository := newSQLiteRepo(t,
		`CREATE TABLE people (name TEXT NOT NULL, age INTEGER)`,
		`INSERT INTO people (name, age) VALUES ('old', 30)`,
	)
	original := listTableRows(t, repository, "people")[0]

	updatedColumns := make([]models.ColValue, len(original.Columns))
	copy(updatedColumns, original.Columns)
	for index := range updatedColumns {
		if updatedColumns[index].ColumnName == "name" {
			updatedColumns[index].Value = "new"
		}
	}

	err := repository.UpdateRow(context.Background(), "people", original.Key, updatedColumns)
	if err != nil {
		t.Fatal(err)
	}

	if got := columnValue(listTableRows(t, repository, "people")[0], "name"); got != "new" {
		t.Fatalf("name = %v, want new", got)
	}
}

func TestDeleteRowWithNullColumn(t *testing.T) {
	repository := newSQLiteRepo(t,
		`CREATE TABLE notes (title TEXT, body TEXT)`,
		`INSERT INTO notes (title, body) VALUES ('first', NULL)`,
	)
	row := listTableRows(t, repository, "notes")[0]

	err := repository.DeleteRow(context.Background(), "notes", row.Key)
	if err != nil {
		t.Fatal(err)
	}
	if got := countRows(t, repository, "notes"); got != 0 {
		t.Fatalf("row count = %d, want 0", got)
	}
}

func TestDeleteRowInCompositeKeyTableRemovesOnlyThatRow(t *testing.T) {
	repository := newSQLiteRepo(t,
		`CREATE TABLE memberships (user_id INTEGER, role_id INTEGER, PRIMARY KEY (user_id, role_id))`,
		`INSERT INTO memberships VALUES (1, 10), (1, 20), (1, 30)`,
	)
	var target models.RowSet
	for _, row := range listTableRows(t, repository, "memberships") {
		if columnValue(row, "role_id") == int64(20) {
			target = row
		}
	}

	err := repository.DeleteRow(context.Background(), "memberships", target.Key)
	if err != nil {
		t.Fatal(err)
	}

	remaining := listTableRows(t, repository, "memberships")
	if len(remaining) != 2 {
		t.Fatalf("remaining rows = %d, want 2", len(remaining))
	}
	for _, row := range remaining {
		if columnValue(row, "role_id") == int64(20) {
			t.Fatal("row with role_id 20 should have been deleted")
		}
	}
}

func TestDeleteRowWithKeyFromEarlierSession(t *testing.T) {
	repository := newSQLiteRepo(t,
		`CREATE TABLE tags (id INTEGER PRIMARY KEY, label TEXT NOT NULL)`,
		`INSERT INTO tags (id, label) VALUES (1, 'a'), (2, 'b')`,
	)
	keyOfSecondRow := listTableRows(t, repository, "tags")[1].Key

	restartedRepository := New(repository.db, configs.DriverSQLite, repository.queryBuilder, 10)
	if err := restartedRepository.DeleteRow(context.Background(), "tags", keyOfSecondRow); err != nil {
		t.Fatal(err)
	}
	if got := countRows(t, restartedRepository, "tags"); got != 1 {
		t.Fatalf("row count = %d, want 1", got)
	}
}

func TestDeleteRowTwiceReportsRowNotFound(t *testing.T) {
	repository := newSQLiteRepo(t,
		`CREATE TABLE tags (id INTEGER PRIMARY KEY, label TEXT)`,
		`INSERT INTO tags (id, label) VALUES (1, 'a')`,
	)
	key := listTableRows(t, repository, "tags")[0].Key

	if err := repository.DeleteRow(context.Background(), "tags", key); err != nil {
		t.Fatal(err)
	}
	err := repository.DeleteRow(context.Background(), "tags", key)
	if !errors.Is(err, apperr.ErrorRowNotFound) {
		t.Fatalf("error = %v, want ErrorRowNotFound", err)
	}
}

func TestUpdateRowCanChangePrimaryKey(t *testing.T) {
	repository := newSQLiteRepo(t,
		`CREATE TABLE tags (id INTEGER PRIMARY KEY, label TEXT)`,
		`INSERT INTO tags (id, label) VALUES (1, 'a')`,
	)
	row := listTableRows(t, repository, "tags")[0]
	changes := []models.ColValue{{ColumnName: "id", Value: int64(5)}, {ColumnName: "label", Value: "a"}}

	if err := repository.UpdateRow(context.Background(), "tags", row.Key, changes); err != nil {
		t.Fatal(err)
	}
	if got := columnValue(listTableRows(t, repository, "tags")[0], "id"); got != int64(5) {
		t.Fatalf("id = %v, want 5", got)
	}
}

func TestUpdateRowRejectsUnknownColumn(t *testing.T) {
	repository := newSQLiteRepo(t,
		`CREATE TABLE tags (id INTEGER PRIMARY KEY, label TEXT)`,
		`INSERT INTO tags (id, label) VALUES (1, 'a')`,
	)
	row := listTableRows(t, repository, "tags")[0]
	changes := []models.ColValue{{ColumnName: `label" = 'x'; --`, Value: "b"}}

	err := repository.UpdateRow(context.Background(), "tags", row.Key, changes)
	if !errors.Is(err, apperr.ErrorInvalidColumn) {
		t.Fatalf("error = %v, want ErrorInvalidColumn", err)
	}
}

func TestListRowsOrdersByPrimaryKeyAndSupportsDescending(t *testing.T) {
	repository := newSQLiteRepo(t,
		`CREATE TABLE tags (id INTEGER PRIMARY KEY, label TEXT)`,
		`INSERT INTO tags (id, label) VALUES (3, 'c'), (1, 'a'), (2, 'b')`,
	)

	rows, err := repository.ListRows(context.Background(), ListDataProps{TableName: "tags", Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	var ids []any
	for _, row := range rows {
		ids = append(ids, columnValue(row, "id"))
	}
	if !reflect.DeepEqual(ids, []any{int64(1), int64(2), int64(3)}) {
		t.Fatalf("ids = %v", ids)
	}

	descending, err := repository.ListRows(
		context.Background(),
		ListDataProps{TableName: "tags", Limit: 10, Column: "label", Order: "DESC"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if got := columnValue(descending[0], "label"); got != "c" {
		t.Fatalf("first label = %v, want c", got)
	}
}

func TestListTablesHidesHistoryAndKeepsOrder(t *testing.T) {
	repository := newSQLiteRepo(t,
		`CREATE TABLE b_table (id INTEGER)`,
		`CREATE TABLE a_table (id INTEGER)`,
	)
	repository.InsertHistory(context.Background(), "created")

	tables, err := repository.ListTables(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, table := range tables {
		names = append(names, table.TableName)
	}
	if !reflect.DeepEqual(names, []string{"a_table", "b_table"}) {
		t.Fatalf("names = %v", names)
	}
}

func TestHistoryRecordsInsertedRows(t *testing.T) {
	repository := newSQLiteRepo(t, `CREATE TABLE tags (label TEXT)`)
	ctx := context.Background()

	err := repository.InsertRow(ctx, "tags", []models.ColValue{{ColumnName: "label", Value: "a"}})
	if err != nil {
		t.Fatal(err)
	}

	history, err := repository.ListHistory(ctx, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 1 || history[0].Message != "Inserted row into table 'tags'" {
		t.Fatalf("unexpected history: %+v", history)
	}
}

func TestRunSQLQueryCapsRowsAndRecordsWrites(t *testing.T) {
	repository := newSQLiteRepo(t, `CREATE TABLE numbers (value INTEGER)`)
	ctx := context.Background()

	var insertStatements []string
	for value := 0; value < maxSQLEditorRows+5; value++ {
		insertStatements = append(insertStatements, fmt.Sprintf("(%d)", value))
	}
	if _, err := repository.RunSQLQuery(
		ctx,
		"INSERT INTO numbers VALUES "+strings.Join(insertStatements, ","),
	); err != nil {
		t.Fatal(err)
	}

	output, err := repository.RunSQLQuery(ctx, "SELECT value FROM numbers")
	if err != nil {
		t.Fatal(err)
	}
	if len(output.Rows) != maxSQLEditorRows || !output.Truncated {
		t.Fatalf("rows = %d, truncated = %v", len(output.Rows), output.Truncated)
	}

	history, err := repository.ListHistory(ctx, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 1 || !strings.HasPrefix(history[0].Message, "Ran SQL: INSERT INTO numbers") {
		t.Fatalf("history = %+v", history)
	}
}

func TestRunSQLQueryRefreshesColumnMetadataAfterAlter(t *testing.T) {
	repository := newSQLiteRepo(t, `CREATE TABLE items (name TEXT)`)
	ctx := context.Background()

	before, err := repository.ListColsMetaData(ctx, "items")
	if err != nil || len(before) != 1 {
		t.Fatalf("before = %v, %v", before, err)
	}
	if _, err := repository.RunSQLQuery(ctx, "ALTER TABLE items ADD COLUMN size INTEGER"); err != nil {
		t.Fatal(err)
	}
	after, err := repository.ListColsMetaData(ctx, "items")
	if err != nil || len(after) != 2 {
		t.Fatalf("after = %v, %v", after, err)
	}
}
