package repo

import (
	"context"

	"github.com/biisal/rowsql/internal/database/models"
)

type Querier interface {
	ListTables(ctx context.Context) ([]models.ListTablesRow, error)
	ListCols(ctx context.Context, tableName string) ([]models.ListDataCol, error)
	ListRows(ctx context.Context, props models.ListDataProps) (models.ListDataRow, error)
	InsertRow(ctx context.Context, props models.InsertDataProps) error
	GetRow(ctx context.Context, tableName, hash string, offset, limit int) ([]any, error)
	GetDriver() string
}

var _ Querier = (*Queries)(nil)
