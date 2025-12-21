package router

import (
	"context"
	"fmt"
	"strings"

	"github.com/biisal/db-gui/configs"
	"github.com/biisal/db-gui/internal/database"
	"github.com/biisal/db-gui/internal/database/repo"
)

type DbService interface {
	ListTables(ctx context.Context) ([]repo.ListTablesRow, error)
	ListCols(ctx context.Context, tableName string) ([]repo.ListDataCol, error)
	ListRows(ctx context.Context, tableName string, page int) (repo.ListDataRow, error)
	InsertRow(ctx context.Context, props repo.InsertDataProps) error
	GetRow(ctx context.Context, tableName string, hash string, page int) ([]any, error)
	UpdateRow(ctx context.Context, values map[string]repo.FormValue, tableName, hash string, page int) error
	CreateTable(ctx context.Context, tableName string, inputs []database.Input) error
	DeleteRow(ctx context.Context, tableName string, hash string, page int) error
	GetTableFormDataTypes() *FormDatatype
	DeleteTable(ctx context.Context, tableName, verificationQuery string) error
	ListHistory(ctx context.Context, page int) ([]repo.History, error)
}

type svc struct {
	repo            *repo.Queries
	maxItemsPerPage int
}

func NewService(repo *repo.Queries, maxItemsPerPage int) DbService {
	return &svc{
		repo, maxItemsPerPage,
	}
}

func getLimitOffset(page int, maxItemsPerPage int) (offset, limit int) {
	offset = (page - 1) * maxItemsPerPage
	limit = offset + maxItemsPerPage
	return
}

func (s *svc) ListTables(ctx context.Context) ([]repo.ListTablesRow, error) {
	return s.repo.ListTables(ctx)
}
func (s *svc) ListCols(ctx context.Context, tableName string) ([]repo.ListDataCol, error) {
	return s.repo.ListCols(ctx, tableName)
}

func (s *svc) ListRows(ctx context.Context, tableName string, page int) (repo.ListDataRow, error) {
	offset, limit := getLimitOffset(page, s.maxItemsPerPage)

	return s.repo.ListRows(ctx, repo.ListDataProps{
		TableName: tableName,
		Limit:     limit,
		Offset:    offset,
	})
}

func (s *svc) InsertRow(ctx context.Context, props repo.InsertDataProps) error {
	return s.repo.InsertRow(ctx, props)
}

func (s *svc) GetRow(ctx context.Context, tableName, hash string, page int) ([]any, error) {
	offset, limit := getLimitOffset(page, s.maxItemsPerPage)
	return s.repo.GetRow(ctx, tableName, hash, offset, limit)
}

func (s *svc) UpdateRow(ctx context.Context, values map[string]repo.FormValue, tableName, hash string, page int) error {
	offset, limit := getLimitOffset(page, s.maxItemsPerPage)
	return s.repo.UpdateRow(ctx, repo.UpdateOrDeleteRowProps{
		TableName: tableName,
		Hash:      hash,
		Limit:     limit,
		Offset:    offset,
		Values:    values,
	})
}

func (s *svc) CreateTable(ctx context.Context, tableName string, inputs []database.Input) error {
	return s.repo.CreateTable(ctx, repo.CreateTableProps{
		TableName: tableName,
		Inputs:    inputs,
	})
}

func (s *svc) DeleteRow(ctx context.Context, tableName, hash string, page int) error {
	limit, offset := getLimitOffset(page, s.maxItemsPerPage)
	return s.repo.DeleteRow(ctx, repo.UpdateOrDeleteRowProps{
		TableName: tableName,
		Hash:      hash,
		Limit:     limit,
		Offset:    offset,
	})
}

type FormDatatype struct {
	NumericDataType []database.NumericDataType `json:"numericType"`
	StringDataType  []database.StringDataType  `json:"stringType"`
}

func (s *svc) GetTableFormDataTypes() *FormDatatype {
	var driver = s.repo.GetDriver()
	switch driver {
	case configs.DRIVER_MYSQL:
		return &FormDatatype{
			NumericDataType: database.MySqlNumericDataTypes,
			StringDataType:  database.MySqlStringDataTypes,
		}
	case configs.DRIVER_POSTGRES:
		return &FormDatatype{
			NumericDataType: database.PostgresNumericDataTypes,
			StringDataType:  database.PostgresStringDataTypes,
		}
	case configs.DRIVER_SQLITE:
		return &FormDatatype{
			NumericDataType: database.SqliteNumericDataTypes,
			StringDataType:  database.SqliteStringDataTypes,
		}
	}
	return nil
}

func (s *svc) DeleteTable(ctx context.Context, tableName, verificationQuery string) error {
	var q = strings.Join(strings.Fields(verificationQuery), " ")
	targetQuiry := "DROP TABLE IF EXISTS " + tableName
	if q != targetQuiry {
		return fmt.Errorf("failed to verify! input should be correct: `%s`", targetQuiry)
	}
	return s.repo.DeleteTable(ctx, tableName)
}

func (s *svc) ListHistory(ctx context.Context, page int) ([]repo.History, error) {
	offset, limit := getLimitOffset(page, s.maxItemsPerPage)
	return s.repo.ListHistory(ctx, limit, offset)
}
