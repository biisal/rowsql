package repo

import (
	"context"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/database/models"
)

type Querier interface {
	ListTables(ctx context.Context) ([]models.ListTablesRow, error)
	ListColsMetaData(ctx context.Context, tableName string) ([]models.ColValue, error)
	ListRows(ctx context.Context, props models.ListDataProps) (models.RowSet, error)
	InsertRow(ctx context.Context, tableName string, colValues []models.ColValue) error
	UpdateRow(ctx context.Context, props UpdateOrDeleteRowProps) error
	DeleteRow(ctx context.Context, props UpdateOrDeleteRowProps) error
	GetRow(ctx context.Context, tableName, hash string, offset, limit int) (models.DataRow, error)
	GetDriver() configs.Driver
}

var _ Querier = (*Queries)(nil)
