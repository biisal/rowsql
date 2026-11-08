package repo

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database/models"
	"github.com/biisal/rowsql/internal/database/queries"
	"github.com/jmoiron/sqlx"

	_ "modernc.org/sqlite"
)

type scenarioDatabase struct {
	driver   configs.Driver
	db       *sqlx.DB
	jsonType string
}

func (d scenarioDatabase) exec(t *testing.T, statements ...string) {
	t.Helper()
	for _, statement := range statements {
		if _, err := d.db.Exec(statement); err != nil {
			t.Fatalf("%s: %v", statement, err)
		}
	}
}

func (d scenarioDatabase) repository(t *testing.T) *Queries {
	t.Helper()
	builder, err := queries.NewBuilder(d.driver, 10)
	if err != nil {
		t.Fatal(err)
	}
	return New(d.db, d.driver, builder, 10)
}

func (d scenarioDatabase) freshTable(t *testing.T, name string, createStatement string) {
	t.Helper()
	d.exec(t, fmt.Sprintf("DROP TABLE IF EXISTS %s", name), createStatement)
}

func valueText(value any) string {
	if value == nil {
		return "<nil>"
	}
	return fmt.Sprint(value)
}

func rowValues(row models.RowSet) map[string]string {
	values := make(map[string]string, len(row.Columns))
	for _, column := range row.Columns {
		values[column.ColumnName] = valueText(column.Value)
	}
	return values
}

func allRows(t *testing.T, repository *Queries, table string) []map[string]string {
	t.Helper()
	rows, err := repository.ListRows(context.Background(), ListDataProps{TableName: table, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	result := make([]map[string]string, len(rows))
	for index, row := range rows {
		result[index] = rowValues(row)
	}
	return result
}

func findRow(t *testing.T, repository *Queries, table string, matches map[string]string) models.RowSet {
	t.Helper()
	rows, err := repository.ListRows(context.Background(), ListDataProps{TableName: table, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range rows {
		values := rowValues(row)
		matched := true
		for column, want := range matches {
			if values[column] != want {
				matched = false
			}
		}
		if matched {
			return row
		}
	}
	t.Fatalf("no row in %s matching %v", table, matches)
	return models.RowSet{}
}

func formFor(values map[string]any, order ...string) []models.ColValue {
	form := make([]models.ColValue, 0, len(order))
	for _, column := range order {
		form = append(form, models.ColValue{ColumnName: column, Value: values[column]})
	}
	return form
}

func scenarioDatabases(t *testing.T) map[string]func(t *testing.T) scenarioDatabase {
	return map[string]func(t *testing.T) scenarioDatabase{
		"sqlite": func(t *testing.T) scenarioDatabase {
			db, err := sqlx.Open("sqlite", filepath.Join(t.TempDir(), "scenario.db"))
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = db.Close() })
			return scenarioDatabase{driver: configs.DriverSQLite, db: db, jsonType: "JSON"}
		},
		"postgres": func(t *testing.T) scenarioDatabase {
			return scenarioDatabase{driver: configs.DriverPostgres, db: postgresDatabase(t), jsonType: "JSONB"}
		},
		"mysql": func(t *testing.T) scenarioDatabase {
			return scenarioDatabase{driver: configs.DriverMySQL, db: mysqlDatabase(t), jsonType: "JSON"}
		},
	}
}

func TestScenariosAcrossDatabases(t *testing.T) {
	scenarios := map[string]func(t *testing.T, database scenarioDatabase){
		"update row without primary key":         scenarioUpdateWithoutPrimaryKey,
		"delete row with null column":            scenarioDeleteRowWithNull,
		"delete from composite primary key":      scenarioCompositePrimaryKey,
		"delete one of two identical rows":       scenarioDuplicateRows,
		"delete twice reports row not found":     scenarioDeleteTwice,
		"change a primary key value":             scenarioChangePrimaryKey,
		"column metadata":                        scenarioColumnMetadata,
		"ordering and paging":                    scenarioOrderingAndPaging,
		"create table with auto increment":       scenarioCreateTable,
		"json column insert and update":          scenarioJSONColumn,
		"list tables stays in current database":  scenarioListTables,
		"history records writes":                 scenarioHistory,
		"hostile column name is rejected":        scenarioHostileColumn,
		"table named like a reserved word":       scenarioReservedWordTable,
		"key survives a new repository instance": scenarioKeyAcrossRepositories,
	}

	for databaseName, open := range scenarioDatabases(t) {
		t.Run(databaseName, func(t *testing.T) {
			database := open(t)
			for scenarioName, scenario := range scenarios {
				t.Run(scenarioName, func(t *testing.T) { scenario(t, database) })
			}
		})
	}
}

func scenarioUpdateWithoutPrimaryKey(t *testing.T, database scenarioDatabase) {
	database.freshTable(t, "sc_people", "CREATE TABLE sc_people (name VARCHAR(50) NOT NULL, age INTEGER)")
	database.exec(t, "INSERT INTO sc_people (name, age) VALUES ('old', 30), ('other', 40)")
	repository := database.repository(t)

	row := findRow(t, repository, "sc_people", map[string]string{"name": "old"})
	form := formFor(map[string]any{"name": "new", "age": 31}, "name", "age")
	if err := repository.UpdateRow(context.Background(), "sc_people", row.Key, form); err != nil {
		t.Fatal(err)
	}

	findRow(t, repository, "sc_people", map[string]string{"name": "new", "age": "31"})
	findRow(t, repository, "sc_people", map[string]string{"name": "other", "age": "40"})
}

func scenarioDeleteRowWithNull(t *testing.T, database scenarioDatabase) {
	database.freshTable(t, "sc_notes", "CREATE TABLE sc_notes (title VARCHAR(50), body VARCHAR(50))")
	database.exec(t, "INSERT INTO sc_notes (title, body) VALUES ('first', NULL), ('second', 'x')")
	repository := database.repository(t)

	row := findRow(t, repository, "sc_notes", map[string]string{"title": "first"})
	if err := repository.DeleteRow(context.Background(), "sc_notes", row.Key); err != nil {
		t.Fatal(err)
	}
	remaining := allRows(t, repository, "sc_notes")
	if len(remaining) != 1 || remaining[0]["title"] != "second" {
		t.Fatalf("remaining = %v", remaining)
	}
}

func scenarioCompositePrimaryKey(t *testing.T, database scenarioDatabase) {
	database.freshTable(
		t,
		"sc_memberships",
		"CREATE TABLE sc_memberships (user_id INTEGER, role_id INTEGER, PRIMARY KEY (user_id, role_id))",
	)
	database.exec(t, "INSERT INTO sc_memberships VALUES (1, 10), (1, 20), (1, 30)")
	repository := database.repository(t)

	row := findRow(t, repository, "sc_memberships", map[string]string{"role_id": "20"})
	if err := repository.DeleteRow(context.Background(), "sc_memberships", row.Key); err != nil {
		t.Fatal(err)
	}
	remaining := allRows(t, repository, "sc_memberships")
	if len(remaining) != 2 || remaining[0]["role_id"] != "10" || remaining[1]["role_id"] != "30" {
		t.Fatalf("remaining = %v", remaining)
	}
}

func scenarioDuplicateRows(t *testing.T, database scenarioDatabase) {
	database.freshTable(t, "sc_duplicates", "CREATE TABLE sc_duplicates (label VARCHAR(20), amount INTEGER)")
	database.exec(t, "INSERT INTO sc_duplicates VALUES ('same', 1), ('same', 1), ('other', 2)")
	repository := database.repository(t)

	row := findRow(t, repository, "sc_duplicates", map[string]string{"label": "same"})
	if err := repository.DeleteRow(context.Background(), "sc_duplicates", row.Key); err != nil {
		t.Fatal(err)
	}
	if got := len(allRows(t, repository, "sc_duplicates")); got != 2 {
		t.Fatalf("rows left = %d, want 2", got)
	}
}

func scenarioDeleteTwice(t *testing.T, database scenarioDatabase) {
	database.freshTable(t, "sc_once", "CREATE TABLE sc_once (id INTEGER PRIMARY KEY, label VARCHAR(20))")
	database.exec(t, "INSERT INTO sc_once VALUES (1, 'a')")
	repository := database.repository(t)

	key := findRow(t, repository, "sc_once", map[string]string{"id": "1"}).Key
	if err := repository.DeleteRow(context.Background(), "sc_once", key); err != nil {
		t.Fatal(err)
	}
	if err := repository.DeleteRow(context.Background(), "sc_once", key); !errors.Is(err, apperr.ErrorRowNotFound) {
		t.Fatalf("error = %v, want ErrorRowNotFound", err)
	}
}

func scenarioChangePrimaryKey(t *testing.T, database scenarioDatabase) {
	database.freshTable(t, "sc_keys", "CREATE TABLE sc_keys (id INTEGER PRIMARY KEY, label VARCHAR(20))")
	database.exec(t, "INSERT INTO sc_keys VALUES (1, 'a')")
	repository := database.repository(t)

	key := findRow(t, repository, "sc_keys", map[string]string{"id": "1"}).Key
	form := formFor(map[string]any{"id": 5, "label": "a"}, "id", "label")
	if err := repository.UpdateRow(context.Background(), "sc_keys", key, form); err != nil {
		t.Fatal(err)
	}
	findRow(t, repository, "sc_keys", map[string]string{"id": "5"})
}

func scenarioColumnMetadata(t *testing.T, database scenarioDatabase) {
	autoIncrement := map[configs.Driver]string{
		configs.DriverPostgres: "id SERIAL",
		configs.DriverMySQL:    "id INT AUTO_INCREMENT",
		configs.DriverSQLite:   "id INTEGER",
	}[database.driver]
	database.freshTable(t, "sc_meta", fmt.Sprintf(`CREATE TABLE sc_meta (
		%s PRIMARY KEY,
		email VARCHAR(50) UNIQUE,
		nickname VARCHAR(50),
		age INTEGER NOT NULL DEFAULT 0
	)`, autoIncrement))
	database.exec(
		t,
		"DROP TABLE IF EXISTS sc_meta_pair",
		"CREATE TABLE sc_meta_pair (first_id INTEGER, second_id INTEGER, PRIMARY KEY (first_id, second_id))",
	)
	repository := database.repository(t)

	columns, err := repository.ListColsMetaData(context.Background(), "sc_meta")
	if err != nil {
		t.Fatal(err)
	}
	byName := map[string]models.ColType{}
	var order []string
	for _, column := range columns {
		byName[column.ColumnName] = column.ColumnType
		order = append(order, column.ColumnName)
	}
	if fmt.Sprint(order) != "[id email nickname age]" {
		t.Fatalf("column order = %v", order)
	}
	if !byName["id"].IsPk || byName["id"].PkPosition != 1 || !byName["id"].HasAutoIncrement {
		t.Errorf("id = %+v", byName["id"])
	}
	if !byName["email"].IsUnique || byName["email"].IsPk || !byName["email"].IsNull {
		t.Errorf("email = %+v", byName["email"])
	}
	if byName["nickname"].IsUnique || !byName["nickname"].IsNull {
		t.Errorf("nickname = %+v", byName["nickname"])
	}
	if byName["age"].IsNull || !byName["age"].HasDefault {
		t.Errorf("age = %+v", byName["age"])
	}

	pairColumns, err := repository.ListColsMetaData(context.Background(), "sc_meta_pair")
	if err != nil {
		t.Fatal(err)
	}
	if len(pairColumns) != 2 || pairColumns[0].ColumnType.PkPosition != 1 || pairColumns[1].ColumnType.PkPosition != 2 {
		t.Errorf("composite positions = %+v", pairColumns)
	}
	if pairColumns[0].ColumnType.IsUnique || pairColumns[0].ColumnType.HasAutoIncrement {
		t.Errorf("composite key column flags = %+v", pairColumns[0].ColumnType)
	}
}

func scenarioOrderingAndPaging(t *testing.T, database scenarioDatabase) {
	database.freshTable(t, "sc_paging", "CREATE TABLE sc_paging (id INTEGER PRIMARY KEY, label VARCHAR(20))")
	for id := 15; id >= 1; id-- {
		database.exec(t, fmt.Sprintf("INSERT INTO sc_paging VALUES (%d, 'row%02d')", id, id))
	}
	repository := database.repository(t)
	ctx := context.Background()

	firstPage, err := repository.ListRows(ctx, ListDataProps{TableName: "sc_paging", Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	secondPage, err := repository.ListRows(ctx, ListDataProps{TableName: "sc_paging", Limit: 10, Offset: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(firstPage) != 10 || len(secondPage) != 5 {
		t.Fatalf("page sizes = %d, %d", len(firstPage), len(secondPage))
	}
	if rowValues(firstPage[0])["id"] != "1" || rowValues(secondPage[0])["id"] != "11" {
		t.Errorf("pages not ordered by primary key: %v then %v", rowValues(firstPage[0]), rowValues(secondPage[0]))
	}

	descending, err := repository.ListRows(
		ctx,
		ListDataProps{TableName: "sc_paging", Limit: 10, Column: "label", Order: "DESC"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if rowValues(descending[0])["label"] != "row15" {
		t.Errorf("descending first = %v", rowValues(descending[0]))
	}
	count, err := repository.GetRowCount(ctx, "sc_paging")
	if err != nil || count != 15 {
		t.Errorf("count = %d, %v", count, err)
	}
}

func scenarioCreateTable(t *testing.T, database scenarioDatabase) {
	database.exec(t, "DROP TABLE IF EXISTS sc_created")
	repository := database.repository(t)
	ctx := context.Background()

	err := repository.CreateTable(ctx, CreateTableProps{TableName: "sc_created", Inputs: []models.ColValue{
		{ColumnName: "id", ColumnType: models.ColType{DataType: "INTEGER", IsPk: true, HasAutoIncrement: true}},
		{
			ColumnName:   "status",
			DefaultValue: "it's new",
			Size:         50,
			ColumnType:   models.ColType{DataType: "VARCHAR", IsNull: true},
		},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.InsertRow(ctx, "sc_created", nil); err != nil {
		t.Fatal(err)
	}
	if err := repository.InsertRow(ctx, "sc_created", nil); err != nil {
		t.Fatal(err)
	}

	rows := allRows(t, repository, "sc_created")
	if len(rows) != 2 || rows[0]["id"] == rows[1]["id"] || rows[0]["status"] != "it's new" {
		t.Fatalf("rows = %v", rows)
	}
}

func scenarioJSONColumn(t *testing.T, database scenarioDatabase) {
	database.freshTable(
		t,
		"sc_documents",
		fmt.Sprintf("CREATE TABLE sc_documents (id INTEGER PRIMARY KEY, body %s)", database.jsonType),
	)
	repository := database.repository(t)
	ctx := context.Background()

	insert := formFor(map[string]any{"id": 1, "body": `{"a": 1}`}, "id", "body")
	if err := repository.InsertRow(ctx, "sc_documents", insert); err != nil {
		t.Fatal(err)
	}
	invalid := formFor(map[string]any{"id": 2, "body": `{broken`}, "id", "body")
	if err := repository.InsertRow(ctx, "sc_documents", invalid); !errors.Is(err, apperr.ErrorInvalidJSON) {
		t.Fatalf("invalid json error = %v", err)
	}
	nullBody := formFor(map[string]any{"id": 3, "body": nil}, "id", "body")
	if err := repository.InsertRow(ctx, "sc_documents", nullBody); err != nil {
		t.Fatal(err)
	}

	key := findRow(t, repository, "sc_documents", map[string]string{"id": "3"}).Key
	update := formFor(map[string]any{"id": 3, "body": `[1, 2, 3]`}, "id", "body")
	if err := repository.UpdateRow(ctx, "sc_documents", key, update); err != nil {
		t.Fatal(err)
	}
	if err := repository.DeleteRow(
		ctx,
		"sc_documents",
		findRow(t, repository, "sc_documents", map[string]string{"id": "1"}).Key,
	); err != nil {
		t.Fatal(err)
	}
	if got := len(allRows(t, repository, "sc_documents")); got != 1 {
		t.Fatalf("rows left = %d", got)
	}
}

func scenarioListTables(t *testing.T, database scenarioDatabase) {
	database.freshTable(t, "sc_listed", "CREATE TABLE sc_listed (id INTEGER)")
	switch database.driver {
	case configs.DriverPostgres:
		database.exec(
			t,
			"CREATE SCHEMA IF NOT EXISTS other_schema",
			"DROP TABLE IF EXISTS other_schema.sc_foreign",
			"CREATE TABLE other_schema.sc_foreign (id INTEGER)",
		)
	case configs.DriverMySQL:
		database.exec(
			t,
			"CREATE DATABASE IF NOT EXISTS other_db",
			"DROP TABLE IF EXISTS other_db.sc_foreign",
			"CREATE TABLE other_db.sc_foreign (id INTEGER)",
		)
	}
	repository := database.repository(t)
	repository.InsertHistory(context.Background(), "makes the history table exist")

	tables, err := repository.ListTables(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	names := map[string]bool{}
	var ordered []string
	for _, table := range tables {
		names[table.TableName] = true
		ordered = append(ordered, table.TableName)
	}
	if !names["sc_listed"] {
		t.Errorf("sc_listed missing from %v", ordered)
	}
	if names["sc_foreign"] || names[historyTableName] {
		t.Errorf("unexpected tables in %v", ordered)
	}
	for index := 1; index < len(ordered); index++ {
		if ordered[index-1] > ordered[index] {
			t.Errorf("tables not sorted: %v", ordered)
			break
		}
	}
	if err := repository.CheckTableExistsInDB(context.Background(), "sc_listed"); err != nil {
		t.Errorf("CheckTableExistsInDB: %v", err)
	}
	if err := repository.CheckTableExistsInDB(context.Background(), "sc_foreign"); err == nil {
		t.Error("table from another schema should not exist")
	}
}

func scenarioHistory(t *testing.T, database scenarioDatabase) {
	database.freshTable(t, "sc_history_target", "CREATE TABLE sc_history_target (label VARCHAR(20))")
	database.exec(t, fmt.Sprintf("DROP TABLE IF EXISTS %s", historyTableName))
	repository := database.repository(t)
	ctx := context.Background()

	if err := repository.InsertRow(
		ctx,
		"sc_history_target",
		formFor(map[string]any{"label": "a"}, "label"),
	); err != nil {
		t.Fatal(err)
	}
	if _, err := repository.RunSQLQuery(ctx, "UPDATE sc_history_target SET label = 'b'"); err != nil {
		t.Fatal(err)
	}

	history, err := repository.ListHistory(ctx, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 2 || history[0].Message != "Ran SQL: UPDATE sc_history_target SET label = 'b'" ||
		history[1].Message != "Inserted row into table 'sc_history_target'" {
		t.Fatalf("history = %+v", history)
	}
	if history[0].Time.IsZero() {
		t.Error("history time was not parsed")
	}
}

func scenarioHostileColumn(t *testing.T, database scenarioDatabase) {
	database.freshTable(t, "sc_safe", "CREATE TABLE sc_safe (id INTEGER PRIMARY KEY, label VARCHAR(20))")
	database.exec(t, "INSERT INTO sc_safe VALUES (1, 'a')")
	repository := database.repository(t)

	key := findRow(t, repository, "sc_safe", map[string]string{"id": "1"}).Key
	hostile := []models.ColValue{{ColumnName: `label" = 'x'; DROP TABLE sc_safe; --`, Value: "b"}}
	if err := repository.UpdateRow(
		context.Background(),
		"sc_safe",
		key,
		hostile,
	); !errors.Is(
		err,
		apperr.ErrorInvalidColumn,
	) {
		t.Fatalf("error = %v, want ErrorInvalidColumn", err)
	}
	if got := len(allRows(t, repository, "sc_safe")); got != 1 {
		t.Fatalf("rows = %d", got)
	}
}

func scenarioReservedWordTable(t *testing.T, database scenarioDatabase) {
	quote := database.repository(t).queryBuilder.QuoteName
	database.freshTable(
		t,
		quote("order"),
		fmt.Sprintf(
			"CREATE TABLE %s (%s INTEGER PRIMARY KEY, %s VARCHAR(20))",
			quote("order"),
			quote("group"),
			quote("select"),
		),
	)
	database.exec(t, fmt.Sprintf("INSERT INTO %s VALUES (1, 'a')", quote("order")))
	repository := database.repository(t)
	ctx := context.Background()

	row := findRow(t, repository, "order", map[string]string{"group": "1"})
	form := formFor(map[string]any{"group": 1, "select": "changed"}, "group", "select")
	if err := repository.UpdateRow(ctx, "order", row.Key, form); err != nil {
		t.Fatal(err)
	}
	findRow(t, repository, "order", map[string]string{"select": "changed"})
	if err := repository.DeleteRow(ctx, "order", row.Key); err != nil {
		t.Fatal(err)
	}
}

func scenarioKeyAcrossRepositories(t *testing.T, database scenarioDatabase) {
	database.freshTable(t, "sc_restart", "CREATE TABLE sc_restart (name VARCHAR(20), note VARCHAR(20))")
	database.exec(t, "INSERT INTO sc_restart VALUES ('a', NULL), ('b', 'x')")
	key := findRow(t, database.repository(t), "sc_restart", map[string]string{"name": "b"}).Key

	restarted := database.repository(t)
	columns, err := restarted.GetRow(context.Background(), "sc_restart", key)
	if err != nil {
		t.Fatal(err)
	}
	if valueText(columns[0].Value) != "b" {
		t.Fatalf("GetRow returned %v", columns)
	}
	if err := restarted.DeleteRow(context.Background(), "sc_restart", key); err != nil {
		t.Fatal(err)
	}
}
