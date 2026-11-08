package queries

import (
	"errors"
	"testing"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database/models"
)

func column(name, dataType string, mutate func(*models.ColValue)) models.ColValue {
	definition := models.ColValue{ColumnName: name, ColumnType: models.ColType{DataType: dataType}}
	if mutate != nil {
		mutate(&definition)
	}
	return definition
}

func TestCreateTable(t *testing.T) {
	autoIncrementKey := func(c *models.ColValue) {
		c.ColumnType.IsPk = true
		c.ColumnType.HasAutoIncrement = true
	}
	tests := []struct {
		name    string
		driver  configs.Driver
		columns []models.ColValue
		want    string
	}{
		{
			name: "postgres auto increment uses serial type", driver: configs.DriverPostgres,
			columns: []models.ColValue{column("id", "integer", autoIncrementKey)},
			want:    `CREATE TABLE "users" ("id" SERIAL NOT NULL PRIMARY KEY)`,
		},
		{
			name: "postgres bigint becomes bigserial", driver: configs.DriverPostgres,
			columns: []models.ColValue{column("id", "bigint", autoIncrementKey)},
			want:    `CREATE TABLE "users" ("id" BIGSERIAL NOT NULL PRIMARY KEY)`,
		},
		{
			name: "mysql auto increment", driver: configs.DriverMySQL,
			columns: []models.ColValue{column("id", "int", autoIncrementKey)},
			want:    "CREATE TABLE `users` (`id` INT NOT NULL PRIMARY KEY AUTO_INCREMENT)",
		},
		{
			name: "sqlite auto increment", driver: configs.DriverSQLite,
			columns: []models.ColValue{column("id", "integer", autoIncrementKey)},
			want:    `CREATE TABLE "users" ("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT)`,
		},
		{
			name: "size nullable unique", driver: configs.DriverPostgres,
			columns: []models.ColValue{column("email", "varchar", func(c *models.ColValue) {
				c.Size = 255
				c.ColumnType.IsNull = true
				c.ColumnType.IsUnique = true
			})},
			want: `CREATE TABLE "users" ("email" VARCHAR(255) UNIQUE)`,
		},
		{
			name:   "string default is a literal",
			driver: configs.DriverPostgres,
			columns: []models.ColValue{
				column("status", "text", func(c *models.ColValue) { c.DefaultValue = "it's new" }),
			},
			want: `CREATE TABLE "users" ("status" TEXT NOT NULL DEFAULT 'it''s new')`,
		},
		{
			name:   "number default",
			driver: configs.DriverSQLite,
			columns: []models.ColValue{
				column("score", "integer", func(c *models.ColValue) { c.DefaultValue = float64(0) }),
			},
			want: `CREATE TABLE "users" ("score" INTEGER NOT NULL DEFAULT 0)`,
		},
		{
			name:   "expression default",
			driver: configs.DriverSQLite,
			columns: []models.ColValue{
				column("created", "datetime", func(c *models.ColValue) { c.DefaultValue = "current_timestamp" }),
			},
			want: `CREATE TABLE "users" ("created" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP)`,
		},
		{
			name: "column named like a reserved word is quoted", driver: configs.DriverPostgres,
			columns: []models.ColValue{column("order", "text", nil)},
			want:    `CREATE TABLE "users" ("order" TEXT NOT NULL)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newTestBuilder(t, tt.driver).CreateTable("users", tt.columns)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("got  %s\nwant %s", got, tt.want)
			}
		})
	}
}

func TestCreateTableErrors(t *testing.T) {
	tests := []struct {
		name    string
		driver  configs.Driver
		table   string
		columns []models.ColValue
		want    error
	}{
		{
			"hostile data type",
			configs.DriverPostgres,
			"t",
			[]models.ColValue{column("a", "int); DROP TABLE x; --", nil)},
			apperr.ErrorUnsupportedDataType,
		},
		{
			"unknown data type",
			configs.DriverSQLite,
			"t",
			[]models.ColValue{column("a", "geometry", nil)},
			apperr.ErrorUnsupportedDataType,
		},
		{
			"invalid table name",
			configs.DriverPostgres,
			"bad name",
			[]models.ColValue{column("a", "text", nil)},
			apperr.ErrorInvalidCreateTableName,
		},
		{"no columns", configs.DriverPostgres, "t", nil, apperr.ErrorNoColumnsInArgs},
		{
			"auto increment without primary key",
			configs.DriverMySQL,
			"t",
			[]models.ColValue{column("a", "int", func(c *models.ColValue) { c.ColumnType.HasAutoIncrement = true })},
			apperr.ErrorInvalidAutoIncrement,
		},
		{
			"sqlite auto increment on non integer",
			configs.DriverSQLite,
			"t",
			[]models.ColValue{column("a", "text", func(c *models.ColValue) {
				c.ColumnType.HasAutoIncrement = true
				c.ColumnType.IsPk = true
			})},
			apperr.ErrorInvalidAutoIncrement,
		},
		{
			"unsupported default",
			configs.DriverPostgres,
			"t",
			[]models.ColValue{column("a", "text", func(c *models.ColValue) { c.DefaultValue = []string{"x"} })},
			apperr.ErrorUnsupportedDefault,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := newTestBuilder(t, tt.driver).CreateTable(tt.table, tt.columns); !errors.Is(err, tt.want) {
				t.Fatalf("error = %v, want %v", err, tt.want)
			}
		})
	}
}
