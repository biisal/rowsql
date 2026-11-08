package router

import (
	"context"
	"reflect"
	"slices"
	"testing"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/database/queries"
	"github.com/biisal/rowsql/internal/database/repo"
	"github.com/biisal/rowsql/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// func TestNewHandler(t *testing.T) {
// 	s := service.DBService(nil)
// 	itemsLimtit := 10
//
// 	handler := NewHandler(s, itemsLimtit)
//
// 	if handler.itemsLimit != itemsLimtit {
// 		t.Errorf("Expected itemsLimit %d, got %d", itemsLimtit, handler.itemsLimit)
// 	}
//
// 	if handler.service != s {
// 		t.Errorf("Expected service %v, got %v", s, handler.service)
// 	}
// }

func setUpDB(ctx context.Context, t *testing.T, driver string) (*sqlx.DB, error) {
	t.Helper()
	pg, err := postgres.Run(
		ctx, "postgres:16-alpine", postgres.WithDatabase("test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatal(err)
	}

	connString, err := pg.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	return sqlx.ConnectContext(ctx, driver, connString)
}

func newTX(t *testing.T, db *sqlx.DB) *sqlx.Tx {
	t.Helper()
	tx, err := db.Beginx()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := tx.Rollback(); err != nil {
			t.Errorf("Failed to rollback transaction: %v", err)
		}
	})
	return tx
}

func getDBHander(t *testing.T, db sqlx.ExtContext) DBHandler {
	t.Helper()
	queryBuilder, err := queries.NewBuilder(configs.DriverPostgres, 10)
	if err != nil {
		t.Fatal(err)
	}
	dbRepo := repo.New(db, "pgx", queryBuilder, 10)
	dbService := service.NewService(dbRepo, queryBuilder, 10)
	return NewHandler(dbService, 10)
}

func dropAllTables(ctx context.Context, db *sqlx.DB) error {
	_, err := db.ExecContext(ctx, `
		DO $$ DECLARE
			r RECORD;
		BEGIN
			FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
				EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
			END LOOP;
		END $$;
	`)
	return err
}

func TestListTablesWithPostgerss(t *testing.T) {
	ctx := context.Background()

	db, err := setUpDB(ctx, t, "pgx")
	if err != nil {
		t.Fatal(err)
	}

	if err := dropAllTables(ctx, db); err != nil {
		t.Fatal(err)
	}

	type listTablesExpectation struct {
		tableNames []string
		schemas    []string
	}

	tests := []struct {
		name                    string
		queryToExcuteBeforeTest string
		expected                listTablesExpectation
	}{
		{
			name:                    "Single table",
			queryToExcuteBeforeTest: "CREATE TABLE IF NOT EXISTS table1 (id SERIAL PRIMARY KEY, name VARCHAR(255))",
			expected: listTablesExpectation{
				tableNames: []string{"table1"},
				schemas:    []string{"public"},
			},
		},
		{
			name: "Multiple tables",
			queryToExcuteBeforeTest: `
CREATE TABLE IF NOT EXISTS table1 (id SERIAL PRIMARY KEY, name VARCHAR(255));
CREATE TABLE IF NOT EXISTS table2 (id SERIAL PRIMARY KEY, name VARCHAR(255));
`,
			expected: listTablesExpectation{
				tableNames: []string{"table1", "table2"},
				schemas:    []string{"public", "public"},
			},
		},
		{
			name:                    "No tables - empty database",
			queryToExcuteBeforeTest: "",
			expected: listTablesExpectation{
				tableNames: []string{},
				schemas:    []string{},
			},
		},
		{
			name:                    "Only rowsql_history table is filtered",
			queryToExcuteBeforeTest: "CREATE TABLE IF NOT EXISTS rowsql_history (id SERIAL PRIMARY KEY, message TEXT, time TIMESTAMP)",
			expected: listTablesExpectation{
				tableNames: []string{},
				schemas:    []string{},
			},
		},
		{
			name: "rowsql_history filtered alongside regular tables",
			queryToExcuteBeforeTest: `
CREATE TABLE IF NOT EXISTS rowsql_history (id SERIAL PRIMARY KEY, message TEXT, time TIMESTAMP);
CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name VARCHAR(255));
CREATE TABLE IF NOT EXISTS orders (id SERIAL PRIMARY KEY, amount INT);
`,
			expected: listTablesExpectation{
				tableNames: []string{"users", "orders"},
				schemas:    []string{"public", "public"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := newTX(t, db)

			if tt.queryToExcuteBeforeTest != "" {
				_, err = tx.ExecContext(ctx, tt.queryToExcuteBeforeTest)
				if err != nil {
					t.Fatal(err)
				}
			}

			dbHandler := getDBHander(t, tx)
			data, err := dbHandler.ListTables(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			tablesGot := make([]string, len(data.Body))
			schemasGot := make([]string, len(data.Body))
			for i, table := range data.Body {
				tablesGot[i] = table.TableName
				schemasGot[i] = table.TableSchema
			}

			slices.Sort(tt.expected.tableNames)
			slices.Sort(tablesGot)
			if !reflect.DeepEqual(tt.expected.tableNames, tablesGot) {
				t.Errorf("Expected tables %v, got %v", tt.expected.tableNames, tablesGot)
			}

			slices.Sort(tt.expected.schemas)
			slices.Sort(schemasGot)
			if !reflect.DeepEqual(tt.expected.schemas, schemasGot) {
				t.Errorf("Expected schemas %v, got %v", tt.expected.schemas, schemasGot)
			}
		})
	}
}
