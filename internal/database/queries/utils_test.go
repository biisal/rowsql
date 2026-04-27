package queries

import (
	"testing"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database/models"
)

func TestFormatColumnDefinition(t *testing.T) {
	tests := []struct {
		name   string
		input  models.ColValue
		want   string
		err    error
		driver configs.Driver
	}{
		{
			driver: configs.DriverPostgres,
			name:   "psql unique witn not null",
			input: models.ColValue{
				ColumnName: "id",
				ColumnType: models.ColType{
					DataType: "integer",
					IsUnique: true,
					IsPk:     true,
				},
			},
			want: "id integer UNIQUE NOT NULL PRIMARY KEY",
		},
		{
			driver: configs.DriverPostgres,
			name:   "psql unique with null",
			input: models.ColValue{
				ColumnName: "id",
				ColumnType: models.ColType{
					DataType: "integer",
					IsUnique: true,
					IsPk:     true,
					IsNull:   true,
				},
			},
			want: "id integer UNIQUE PRIMARY KEY",
		},
		{
			driver: configs.DriverPostgres,
			name:   "psql primary key without unique",
			input: models.ColValue{
				ColumnName: "id",
				ColumnType: models.ColType{
					DataType: "integer",
					IsPk:     true,
					IsNull:   true,
					IsUnique: false,
				},
			},
			want: "id integer PRIMARY KEY",
		},
		{
			driver: configs.DriverMySQL,
			name:   "mysql auto increment pkey",
			input: models.ColValue{
				ColumnName: "id",
				ColumnType: models.ColType{
					DataType:         "int",
					IsPk:             true,
					HasAutoIncrement: true,
				},
			},
			want: "id int NOT NULL PRIMARY KEY AUTO_INCREMENT",
		},
		{
			driver: configs.DriverSQLite,
			name:   "sqlite auto increment pkey",
			input: models.ColValue{
				ColumnName: "id",
				ColumnType: models.ColType{
					DataType:         "INTEGER",
					IsPk:             true,
					HasAutoIncrement: true,
				},
			},
			want: "id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT",
		},
		{
			driver: configs.DriverPostgres,
			name:   "psql column size",
			input: models.ColValue{
				ColumnName: "name",
				Size:       255,
				ColumnType: models.ColType{
					DataType: "varchar",
					HasSize:  true,
				},
			},
			want: "name varchar(255) NOT NULL",
		},
		{
			driver: configs.DriverPostgres,
			name:   "psql default value string",
			input: models.ColValue{
				ColumnName:   "status",
				DefaultValue: "active",
				ColumnType: models.ColType{
					DataType: "text",
				},
			},
			want: "status text NOT NULL DEFAULT active",
		},
		{
			driver: configs.DriverPostgres,
			name:   "psql default value int",
			input: models.ColValue{
				ColumnName:   "age",
				DefaultValue: 18,
				ColumnType: models.ColType{
					DataType: "integer",
				},
			},
			want: "age integer NOT NULL DEFAULT 18",
		},
		{
			driver: configs.DriverPostgres,
			name:   "psql name with space",
			input: models.ColValue{
				ColumnName: "group id",
				ColumnType: models.ColType{
					DataType: "integer",
				},
			},
			want: "\"group id\" integer NOT NULL",
		},
		{
			driver: configs.DriverPostgres,
			name:   "psql auto increment error (not pk)",
			input: models.ColValue{
				ColumnName: "id",
				ColumnType: models.ColType{
					DataType:         "integer",
					HasAutoIncrement: true,
					IsPk:             false,
				},
			},
			err: apperr.ErrorInvalidAutoIncrement,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dialect, err := GetDialect(tt.driver)
			if err != nil {
				assertErr(t, err, tt.err)
				return
			}
			b := &builder{
				driver:   tt.driver,
				maxLimit: 10,
				dialect:  dialect,
			}
			got, err := b.formatColumnDefinition(tt.input)
			if err != nil {
				assertErr(t, err, tt.err)
				return
			}
			if got != tt.want {
				t.Errorf("got = %q, want %q", got, tt.want)
			}
		})
	}

}
