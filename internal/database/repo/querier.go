package repo

import (
	"context"
)

type Querier interface {
	ListTables(ctx context.Context) ([]ListTablesRow, error)
	ListCols(ctx context.Context, tableName string) ([]ListDataCol, error)
	ListRows(ctx context.Context, props ListDataProps) (ListDataRow, error)
	InsertRow(ctx context.Context, props InsertDataProps) error
	GetRow(ctx context.Context, tableName, hash string, offset, limit int) ([]any, error)
	GetDriver() string
}

var _ Querier = (*Queries)(nil)
