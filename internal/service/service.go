// Package service provides the service layer for the application.
package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database"
	"github.com/biisal/rowsql/internal/database/models"
	"github.com/biisal/rowsql/internal/database/repo"
)

type DBService interface {
	CheckTableExists(ctx context.Context, tableName string) error
	ListTables(ctx context.Context) ([]models.ListTablesRow, error)
	ListCols(ctx context.Context, tableName string) ([]models.ColValue, error)
	ListRows(ctx context.Context, tableName string, page int, column string, order string) ([]models.RowSet, error)
	InsertRow(ctx context.Context, tableName string, form []models.ColValue) error
	GetRow(ctx context.Context, tableName, key string) ([]models.ColValue, error)
	UpdateRow(ctx context.Context, tableName, key string, values []models.ColValue) error
	CreateTable(ctx context.Context, tableName string, inputs []models.ColValue) error
	GetRowCount(ctx context.Context, tableName string) (int, error)
	DeleteRow(ctx context.Context, tableName, key string) error
	GetTableFormDataTypes() *FormDatatype
	DeleteTable(ctx context.Context, tableName, verificationQuery string) error
	ListHistory(ctx context.Context, page int) ([]models.History, error)
	HasNextPage(ctx context.Context, total, page int) bool
	TotalPages(totalRows int) int
	RunSQLQuery(ctx context.Context, query string) (*models.RunSQLQueryOutput, error)
}

type svc struct {
	repo  *repo.Queries
	limit int
}

func NewService(repo *repo.Queries, maxItemsPerPage int) DBService {
	return &svc{repo: repo, limit: maxItemsPerPage}
}

func (s *svc) CheckTableExists(ctx context.Context, tableName string) error {
	return s.repo.CheckTableExistsInDB(ctx, tableName)
}

func (s *svc) ListTables(ctx context.Context) ([]models.ListTablesRow, error) {
	return s.repo.ListTables(ctx)
}

func (s *svc) CreateTable(ctx context.Context, tableName string, inputs []models.ColValue) error {
	return s.repo.CreateTable(ctx, repo.CreateTableProps{
		TableName: tableName,
		Inputs:    inputs,
	})
}

func (s *svc) DeleteTable(ctx context.Context, tableName, verificationQuery string) error {
	typedQuery := strings.Join(strings.Fields(verificationQuery), " ")
	expectedQuery := "DROP TABLE " + tableName
	if typedQuery != expectedQuery {
		return fmt.Errorf("%w! input should be correct: `%s`", apperr.ErrorVerificationFailed, expectedQuery)
	}
	return s.repo.DeleteTable(ctx, tableName)
}

func (s *svc) ListCols(ctx context.Context, tableName string) ([]models.ColValue, error) {
	return s.repo.ListColsMetaData(ctx, tableName)
}

func (s *svc) GetRow(ctx context.Context, tableName, key string) ([]models.ColValue, error) {
	return s.repo.GetRow(ctx, tableName, key)
}

func (s *svc) InsertRow(ctx context.Context, tableName string, form []models.ColValue) error {
	return s.repo.InsertRow(ctx, tableName, form)
}

func (s *svc) ListRows(
	ctx context.Context,
	tableName string,
	page int,
	column string,
	order string,
) ([]models.RowSet, error) {
	return s.repo.ListRows(ctx, repo.ListDataProps{
		TableName: tableName,
		Limit:     s.limit,
		Offset:    s.getOffset(page),
		Column:    column,
		Order:     order,
	})
}

func (s *svc) UpdateRow(ctx context.Context, tableName, key string, values []models.ColValue) error {
	return s.repo.UpdateRow(ctx, tableName, key, values)
}

func (s *svc) DeleteRow(ctx context.Context, tableName, key string) error {
	return s.repo.DeleteRow(ctx, tableName, key)
}

func (s *svc) GetRowCount(ctx context.Context, tableName string) (int, error) {
	return s.repo.GetRowCount(ctx, tableName)
}

type FormDatatype struct {
	NumericDataType []database.VarianDataType `json:"numericType"`
	StringDataType  []database.VarianDataType `json:"stringType"`
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

func (s *svc) ListHistory(ctx context.Context, page int) ([]models.History, error) {
	return s.repo.ListHistory(ctx, s.limit, s.getOffset(page))
}

func (s *svc) RunSQLQuery(ctx context.Context, query string) (*models.RunSQLQueryOutput, error) {
	return s.repo.RunSQLQuery(ctx, query)
}
