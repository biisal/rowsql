package queries

import (
	"testing"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database/models"
)

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
			dialect, bErr := GetDialect(tt.driver)
			if bErr != nil {
				assertErr(t, bErr, tt.err)
				return
			}
			query, args, err := dialect.WhereCluse(tt.cols, tt.rows, tt.argsIdx)

			assertErr(t, err, tt.err)
			assertQuery(t, query, tt.want)
			assertArgs(t, args, tt.args)
		})
	}
}
