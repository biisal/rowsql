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
		input  models.ColValues
		want   string
		err    error
		driver configs.Driver
	}{
		{
			driver: configs.DriverPostgres,
			name:   "psql unique witn not null",
			input: models.ColValues{
				Name:     "id",
				Type:     "integer",
				IsUnique: true,
				IsPK:     true,
			},
			want: "id integer UNIQUE NOT NULL PRIMARY KEY",
		},
		{
			driver: configs.DriverPostgres,
			name:   "psql unique with null",
			input: models.ColValues{
				Name:     "id",
				Type:     "integer",
				IsUnique: true,
				IsPK:     true,
				IsNull:   true,
			},
			want: "id integer UNIQUE PRIMARY KEY",
		},
		{
			driver: configs.DriverPostgres,
			name:   "psql primary key without unique",
			input: models.ColValues{
				Name:     "id",
				Type:     "integer",
				IsPK:     true,
				IsNull:   true,
				IsUnique: false,
			},
			want: "id integer PRIMARY KEY",
		},
		{
			driver: configs.DriverMySQL,
			name:   "mysql auto increment pkey",
			input: models.ColValues{
				Name:          "id",
				Type:          "int",
				IsPK:          true,
				AutoIncrement: true,
			},
			want: "id int NOT NULL PRIMARY KEY AUTO_INCREMENT",
		},
		{
			driver: configs.DriverSQLite,
			name:   "sqlite auto increment pkey",
			input: models.ColValues{
				Name:          "id",
				Type:          "INTEGER",
				IsPK:          true,
				AutoIncrement: true,
			},
			want: "id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT",
		},
		{
			driver: configs.DriverPostgres,
			name:   "psql column size",
			input: models.ColValues{
				Name: "name",
				Type: "varchar",
				Size: 255,
			},
			want: "name varchar(255) NOT NULL",
		},
		{
			driver: configs.DriverPostgres,
			name:   "psql default value string",
			input: models.ColValues{
				Name:    "status",
				Type:    "text",
				Default: "active",
			},
			want: "status text NOT NULL DEFAULT active",
		},
		{
			driver: configs.DriverPostgres,
			name:   "psql default value int",
			input: models.ColValues{
				Name:    "age",
				Type:    "integer",
				Default: 18,
			},
			want: "age integer NOT NULL DEFAULT 18",
		},
		{
			driver: configs.DriverPostgres,
			name:   "psql name with space",
			input: models.ColValues{
				Name: "group id",
				Type: "integer",
			},
			want: "\"group id\" integer NOT NULL",
		},
		{
			driver: configs.DriverPostgres,
			name:   "psql auto increment error (not pk)",
			input: models.ColValues{
				Name:          "id",
				Type:          "integer",
				AutoIncrement: true,
				IsPK:          false,
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
