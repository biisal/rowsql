package repo

import (
	"context"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/database/models"
)

type Querier interface {
	ListTables(ctx context.Context) ([]models.ListTablesRow, error)
	ListColsMetaData(ctx context.Context, tableName string) ([]models.ColValue, error)
	ListRows(ctx context.Context, props ListDataProps) ([]models.RowSet, error)
	InsertRow(ctx context.Context, tableName string, colValues []models.ColValue) error
	UpdateRow(ctx context.Context, tableName, key string, values []models.ColValue) error
	DeleteRow(ctx context.Context, tableName, key string) error
	GetRow(ctx context.Context, tableName, key string) ([]models.ColValue, error)
	GetDriver() configs.Driver
	RunSQLQuery(ctx context.Context, query string) (*models.RunSQLQueryOutput, error)
}

var _ Querier = (*Queries)(nil)
