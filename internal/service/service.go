// Package service provides the service layer for the application.
package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/biisal/db-gui/configs"
	"github.com/biisal/db-gui/internal/database"
	"github.com/biisal/db-gui/internal/database/repo"
)

type DBService interface {
	CheckTableExits(ctx context.Context, tableName string) error
	ListTables(ctx context.Context) ([]repo.ListTablesRow, error)
	ListCols(ctx context.Context, tableName string) ([]repo.ListDataCol, error)
	ListRows(ctx context.Context, tableName string, page int, column string, order string) (repo.ListDataRow, error)
	InsertRow(ctx context.Context, props repo.InsertDataProps) error
	GetRow(ctx context.Context, tableName string, hash string, page int) ([]any, error)
	UpdateRow(ctx context.Context, values map[string]repo.FormValue, tableName, hash string, page int) error
	CreateTable(ctx context.Context, tableName string, inputs []database.Input) error
	GetRowCount(ctx context.Context, tableName string) (int, error)
	DeleteRow(ctx context.Context, tableName string, hash string, page int) error
	GetTableFormDataTypes() *FormDatatype
	DeleteTable(ctx context.Context, tableName, verificationQuery string) error
	ListHistory(ctx context.Context, page int) ([]repo.History, error)
	HasNextPage(ctx context.Context, total, page int) bool
}

type svc struct {
	repo  *repo.Queries
	limit int
}

func NewService(repo *repo.Queries, maxItemsPerPage int) DBService {
	return &svc{
		repo, maxItemsPerPage,
	}
}

func (s *svc) CheckTableExits(ctx context.Context, tableName string) error {
	return s.repo.CheckTableExitsInDB(ctx, tableName)
}

func (s *svc) ListTables(ctx context.Context) ([]repo.ListTablesRow, error) {
	return s.repo.ListTables(ctx)
}

func (s *svc) CreateTable(ctx context.Context, tableName string, inputs []database.Input) error {
	return s.repo.CreateTable(ctx, repo.CreateTableProps{
		TableName: tableName,
		Inputs:    inputs,
	})
}

func (s *svc) DeleteTable(ctx context.Context, tableName, verificationQuery string) error {
	q := strings.Join(strings.Fields(verificationQuery), " ")
	targetQuiry := "DROP TABLE IF EXISTS " + tableName
	if q != targetQuiry {
		return fmt.Errorf("failed to verify! input should be correct: `%s`", targetQuiry)
	}
	return s.repo.DeleteTable(ctx, tableName)
}

func (s *svc) ListCols(ctx context.Context, tableName string) ([]repo.ListDataCol, error) {
	return s.repo.ListCols(ctx, tableName)
}

func (s *svc) GetRow(ctx context.Context, tableName, hash string, page int) ([]any, error) {
	return s.repo.GetRow(ctx, tableName, hash, s.getOffset(page), s.limit)
}

func (s *svc) InsertRow(ctx context.Context, props repo.InsertDataProps) error {
	return s.repo.InsertRow(ctx, props)
}

func (s *svc) ListRows(ctx context.Context, tableName string, page int, column string, order string) (repo.ListDataRow, error) {
	return s.repo.ListRows(ctx, repo.ListDataProps{
		TableName: tableName,
		Limit:     s.limit,
		Offset:    s.getOffset(page),
		Column:    column,
		Order:     order,
	})
}

func (s *svc) UpdateRow(ctx context.Context, values map[string]repo.FormValue, tableName, hash string, page int) error {
	return s.repo.UpdateRow(ctx, repo.UpdateOrDeleteRowProps{
		TableName: tableName,
		Hash:      hash,
		Limit:     s.limit,
		Offset:    s.getOffset(page),
		Values:    values,
	})
}

func (s *svc) DeleteRow(ctx context.Context, tableName, hash string, page int) error {
	return s.repo.DeleteRow(ctx, repo.UpdateOrDeleteRowProps{
		TableName: tableName,
		Hash:      hash,
		Limit:     s.limit,
		Offset:    s.getOffset(page),
	})
}

func (s *svc) GetRowCount(ctx context.Context, tableName string) (int, error) {
	return s.repo.GetRowCount(ctx, tableName)
}

type FormDatatype struct {
	NumericDataType []database.NumericDataType `json:"numericType"`
	StringDataType  []database.StringDataType  `json:"stringType"`
}

func (s *svc) GetTableFormDataTypes() *FormDatatype {
	driver := s.repo.GetDriver()
	switch driver {
	case configs.DriverMySQL:
		return &FormDatatype{
			NumericDataType: database.MySQLNumericDataTypes,
			StringDataType:  database.MySQLStringDataTypes,
		}
	case configs.DriverPostgres:
		return &FormDatatype{
			NumericDataType: database.PostgresNumericDataTypes,
			StringDataType:  database.PostgresStringDataTypes,
		}
	case configs.DriverSQLite:
		return &FormDatatype{
			NumericDataType: database.SqliteNumericDataTypes,
			StringDataType:  database.SqliteStringDataTypes,
		}
	}
	return nil
}

func (s *svc) ListHistory(ctx context.Context, page int) ([]repo.History, error) {
	return s.repo.ListHistory(ctx, s.limit, s.getOffset(page))
}
