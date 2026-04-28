package router

import (
	"context"
	"log/slog"
	"strings"

	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database/models"
	"github.com/biisal/rowsql/internal/logger"
	"github.com/biisal/rowsql/internal/service"
	"github.com/danielgtaylor/huma/v2"
)

type DBHandler struct {
	service    service.DBService
	itemsLimit int
}

type BaseHTMLData struct {
	Tables []models.ListTablesRow
	Cols   []models.ColValue
}

type ErrorMessage struct {
	Message string
}

func NewHandler(service service.DBService, itemsLimit int) DBHandler {
	return DBHandler{
		service,
		itemsLimit,
	}
}

func (h *DBHandler) checkTableExists(ctx context.Context, tableName string) error {
	if err := h.service.CheckTableExits(ctx, tableName); err != nil {
		logger.Error("%s", err)
		return huma.Error404NotFound("table not found", err)
	}
	return nil
}

func (h DBHandler) ListTables(ctx context.Context, input *struct{}) (*struct{ Body []models.ListTablesRow }, error) {
	tables, err := h.service.ListTables(ctx)
	if err != nil {
		logger.Error("%s", err)
		return nil, err
	}

	logger.Info("Fetched %d tables", len(tables))
	return &struct{ Body []models.ListTablesRow }{Body: tables}, nil
}

type ListColumnsInput struct {
	TableName string `path:"tableName"`
}

func (h DBHandler) ListColumns(ctx context.Context, input *ListColumnsInput) (*struct{ Body []models.ColValue }, error) {
	if err := h.checkTableExists(ctx, input.TableName); err != nil {
		return nil, err
	}
	cols, err := h.service.ListCols(ctx, input.TableName)
	if err != nil {
		logger.Error("%s", err)
		return nil, err
	}
	return &struct{ Body []models.ColValue }{Body: cols}, nil
}

type ListRowsInput struct {
	TableName string `path:"tableName"`
	Page      int    `query:"page" default:"1" minimum:"1"`
	Column    string `query:"column"`
	Order     string `query:"order" enum:"ASC,DESC"`
}
type ListRowsOutput struct {
	Body ListRowsResponse
}

func (h DBHandler) ListRows(ctx context.Context, input *ListRowsInput) (*ListRowsOutput, error) {
	tableName := input.TableName
	if err := h.checkTableExists(ctx, tableName); err != nil {
		return nil, err
	}

	colParam := strings.TrimSpace(input.Column)
	order := input.Order

	var err error

	colFound := false
	if colParam != "" {
		var cols []models.ColValue
		cols, err = h.service.ListCols(ctx, tableName)
		if err != nil {
			return nil, err
		}
		for _, col := range cols {
			if col.ColumnName == colParam {
				colFound = true
				break
			}
		}
		if !colFound {
			return nil, apperr.ErrorInvalidColumn
		}
	}

	rows, err := h.service.ListRows(ctx, tableName, input.Page, colParam, order)
	if err != nil {
		return nil, err
	}

	cols, err := h.service.ListCols(ctx, tableName)
	if err != nil {
		return nil, err
	}

	count, err := h.service.GetRowCount(ctx, tableName)
	if err != nil {
		return nil, err
	}

	return &ListRowsOutput{
		Body: ListRowsResponse{
			Page:        input.Page,
			Rows:        rows,
			Cols:        cols,
			RowCount:    count,
			ActiveTable: tableName,
			HasNextPage: h.service.HasNextPage(ctx, count, input.Page),
			TotalPages:  count / h.itemsLimit,
		},
	}, nil
}

type RowInsertOrUpdateFormInput struct {
	TableName string `path:"tableName"`
	Page      int    `query:"page" default:"1" minimum:"1"`
	Hash      string `query:"hash"`
}

type RowInsertOrUpdateFormOutput struct {
	Body struct {
		Cols []models.ColValue `json:"cols"`
	}
}

func (h DBHandler) RowInsertOrUpdateForm(ctx context.Context, input *RowInsertOrUpdateFormInput) (*RowInsertOrUpdateFormOutput, error) {
	if err := h.checkTableExists(ctx, input.TableName); err != nil {
		return nil, err
	}
	colsMeta, err := h.service.ListCols(ctx, input.TableName)
	if err != nil {
		logger.Error("%s", err)
		return nil, err
	}
	if input.Hash != "" {
		initialRow, err := h.service.GetRow(ctx, input.TableName, input.Hash, input.Page)
		if err != nil {
			return nil, err
		}
		slog.Info("Initial row", "data", initialRow)
		return &RowInsertOrUpdateFormOutput{
			Body: struct {
				Cols []models.ColValue `json:"cols"`
			}{
				Cols: initialRow,
			},
		}, nil
	}

	return &RowInsertOrUpdateFormOutput{
		Body: struct {
			Cols []models.ColValue `json:"cols"`
		}{
			Cols: colsMeta,
		},
	}, nil
}

type InsertOrUpdateRowInput struct {
	TableName string            `path:"tableName"`
	Hash      string            `query:"hash"`
	Page      int               `query:"page" default:"1"`
	Body      []models.ColValue `json:"body"`
}

func (h DBHandler) InsertOrUpdateRow(ctx context.Context, input *InsertOrUpdateRowInput) (*struct{}, error) {
	tableName := input.TableName
	if err := h.checkTableExists(ctx, tableName); err != nil {
		return nil, err
	}
	form := input.Body

	logger.Info("Request data: %+v", form)

	pageInt := input.Page
	hash := strings.TrimSpace(input.Hash)

	if hash != "" {
		if err := h.service.UpdateRow(ctx, form, tableName, hash, pageInt); err != nil {
			logger.Error("%s", err)
			logger.Error("Failed to update row in table '%s'", tableName)
			return nil, err
		}
		logger.Success("Row updated successfully in table '%s'", tableName)
		return nil, nil
	}

	if err := h.service.InsertRow(ctx, tableName, form); err != nil {
		logger.Errorln(err.Error())
		logger.Error("Failed to insert row in table '%s'", tableName)
		return nil, err
	}
	logger.Success("Row inserted successfully in table '%s'", tableName)
	return nil, nil
}

type DeleteRowInput struct {
	TableName string `path:"tableName"`
	Hash      string `path:"hash"`
	Page      int    `query:"page" default:"1"`
}

func (h DBHandler) DeleteRow(ctx context.Context, input *DeleteRowInput) (*struct{}, error) {
	if err := h.checkTableExists(ctx, input.TableName); err != nil {
		return nil, err
	}
	if err := h.service.DeleteRow(ctx, input.TableName, input.Hash, input.Page); err != nil {
		logger.Error("%s", err)
		logger.Error("Failed to delete row from table '%s'", input.TableName)
		return nil, err
	}
	logger.Success("Row deleted successfully from table '%s'", input.TableName)
	return nil, nil
}

func (h DBHandler) NewTableFormFileds(ctx context.Context, input *struct{}) (*struct{ Body *service.FormDatatype }, error) {
	fields := h.service.GetTableFormDataTypes()
	if fields == nil {
		return nil, huma.Error500InternalServerError("no data found")
	}

	return &struct{ Body *service.FormDatatype }{Body: fields}, nil
}

type CreeteNewTableInput struct {
	Body struct {
		TableName string            `json:"tableName"`
		Inputs    []models.ColValue `json:"inputs"`
	}
}

func (h DBHandler) CreeteNewTable(ctx context.Context, input *CreeteNewTableInput) (*struct{}, error) {
	req := input.Body
	logger.Info("Request data: %+v", req)
	if err := h.service.CreateTable(ctx, req.TableName, req.Inputs); err != nil {
		logger.Error("%s", err)
		logger.Error("Failed to create table '%s'", req.TableName)
		return nil, err
	}
	logger.Success("Table '%s' created successfully with %d columns", req.TableName, len(req.Inputs))
	return nil, nil
}

type DeleteTableRequest struct {
	TableName         string `json:"tableName"`
	VerificationQuiry string `json:"verificationQuery"`
}

func (h *DBHandler) DeleteTable(ctx context.Context, input *struct{ Body DeleteTableRequest }) (*struct{}, error) {
	req := input.Body
	err := h.service.DeleteTable(ctx, req.TableName, req.VerificationQuiry)
	if err != nil {
		logger.Error("%s", err)
		logger.Error("Failed to delete table '%s'", req.TableName)
		return nil, err
	}
	logger.Success("Table '%s' deleted successfully", req.TableName)
	return nil, nil
}

type ListHistoryInput struct {
	Page int `query:"page" default:"1"`
}

func (h *DBHandler) ListHistory(ctx context.Context, input *ListHistoryInput) (*struct{ Body []models.History }, error) {
	pageInt := max(input.Page, 1)

	history, err := h.service.ListHistory(ctx, pageInt)
	if err != nil {
		logger.Error("Failed to list history: %v", err)
		return nil, err
	}

	logger.Debug("Retrieved %d history entries (page %d)", len(history), pageInt)
	return &struct{ Body []models.History }{Body: history}, nil
}

func (h *DBHandler) ListRecentHistory(ctx context.Context, input *struct{}) (*struct{ Body []models.History }, error) {
	history, err := h.service.ListHistory(ctx, 1)
	if err != nil {
		logger.Error("Failed to list recent history: %v", err)
		return nil, err
	}

	if len(history) > 10 {
		history = history[:10]
	}

	return &struct{ Body []models.History }{Body: history}, nil
}

type RunSQLQueryInput struct {
	Body struct {
		Query string `json:"query"`
	}
}

type RunSQLQueryOutput struct {
	Body struct {
		*models.RunSQLQueryOutput
	}
}

func (h *DBHandler) HandleRunSQLQuery(ctx context.Context, input *RunSQLQueryInput) (*RunSQLQueryOutput, error) {
	result, err := h.service.RunSQLQuery(ctx, input.Body.Query)
	if err != nil {
		logger.Error("Failed to run SQL query: %v", err)
		return nil, err
	}
	return &RunSQLQueryOutput{
		Body: struct {
			*models.RunSQLQueryOutput
		}{
			RunSQLQueryOutput: result,
		},
	}, nil
}
