package router

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/database"
	"github.com/biisal/rowsql/internal/database/models"
	"github.com/biisal/rowsql/internal/logger"
	resopnse "github.com/biisal/rowsql/internal/response"
	"github.com/biisal/rowsql/internal/service"
)

type DBHandler struct {
	service    service.DBService
	itemsLimit int
}

type BaseHTMLData struct {
	Tables      []models.ListTablesRow
	Cols        []models.ListDataCol
	ActiveTable string
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

func (h DBHandler) getBaseData(ctx context.Context, tableName ...string) (*BaseHTMLData, error) {
	tables, err := h.service.ListTables(ctx)
	if err != nil {
		logger.Error("%s", err)
		return nil, err
	}
	if len(tableName) == 0 {
		return &BaseHTMLData{Tables: tables}, nil
	}
	cols, err := h.service.ListCols(ctx, tableName[0])
	if err != nil {
		logger.Error("%s", err)
		return nil, err
	}
	return &BaseHTMLData{Tables: tables, Cols: cols, ActiveTable: tableName[0]}, nil
}

func (h DBHandler) ListColumns(w http.ResponseWriter, r *http.Request) {
	cols, err := h.service.ListCols(r.Context(), r.PathValue("tableName"))
	if err != nil {
		logger.Error("%s", err)
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}
	resopnse.Success(w, http.StatusOK, cols)
}

func (h DBHandler) ListTables(w http.ResponseWriter, r *http.Request) {
	tables, err := h.service.ListTables(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	logger.Info("Fetched %d tables", len(tables))
	resopnse.Success(w, http.StatusOK, tables)
}

func (h DBHandler) ListRows(w http.ResponseWriter, r *http.Request) {
	tableName := r.PathValue("tableName")
	page := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	pageInt = max(pageInt, 1)

	colParam := strings.TrimSpace(r.URL.Query().Get("column"))
	order := r.URL.Query().Get("order")

	colFound := false
	if colParam != "" {
		var cols []models.ListDataCol
		cols, err = h.service.ListCols(r.Context(), tableName)
		if err != nil {
			resopnse.Error(w, http.StatusInternalServerError, err)
		}
		for _, col := range cols {
			if col.ColumnName == colParam {
				colFound = true
				break
			}
		}
		if !colFound {
			logger.Error("error: %s, table: %s, requested column: %s", apperr.ErrorInvalidColumn, tableName, colParam)
			resopnse.Error(w, http.StatusBadRequest, apperr.ErrorInvalidColumn)
			return
		}

	}
	rows, err := h.service.ListRows(r.Context(), tableName, pageInt, colParam, order)
	if err != nil {
		logger.Error("Failed to fetch rows from table '%s'", tableName)
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}

	cols, err := h.service.ListCols(r.Context(), tableName)
	if err != nil {
		logger.Error("%s", err)
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}

	count, err := h.service.GetRowCount(r.Context(), tableName)
	if err != nil {
		logger.Error("%s", err)
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}

	logger.Debug("Loaded page %d for table '%s'", pageInt, tableName)
	resopnse.Success(w, http.StatusOK,
		ListRowsResponse{
			Page:        pageInt,
			Rows:        rows,
			Cols:        cols,
			RowCount:    count,
			ActiveTable: tableName,
			HasNextPage: h.service.HasNextPage(r.Context(), count, pageInt),
			TotalPages:  count / h.itemsLimit,
		},
	)
}

func (h DBHandler) RowInsertForm(w http.ResponseWriter, r *http.Request) {
	tableName := r.PathValue("tableName")

	action := "Insert"
	hash := strings.TrimSpace(r.URL.Query().Get("hash"))
	page := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	var initialRow []any
	if hash != "" {
		action = "Update"
		initialRow, err = h.service.GetRow(r.Context(), tableName, hash, pageInt)
		if err != nil {
			logger.Error("%s", err)
			resopnse.Error(w, http.StatusInternalServerError, err)
			return
		}
	}

	basseData, err := h.getBaseData(r.Context(), tableName)
	if err != nil {
		logger.Error("%s", err)
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}
	if len(initialRow) == len(basseData.Cols) {
		for i := range basseData.Cols {
			basseData.Cols[i].Value = initialRow[i]
		}
	}
	data := struct {
		Action string
		BaseHTMLData
	}{
		action,
		*basseData,
	}
	resopnse.Success(w, http.StatusOK, data)
}

func (h DBHandler) InsertOrUpdateRow(w http.ResponseWriter, r *http.Request) {
	tableName := r.PathValue("tableName")
	ctx := r.Context()
	form := make(map[string]models.FormValue)
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		logger.Error("%s", err)
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}

	page := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}

	hash := strings.TrimSpace(r.URL.Query().Get("hash"))
	if hash != "" {
		if err := h.service.UpdateRow(ctx, form, tableName, hash, pageInt); err != nil {
			logger.Error("%s", err)
			logger.Error("Failed to update row in table '%s'", tableName)
			resopnse.Error(w, http.StatusInternalServerError, err)
			return
		}
		logger.Success("Row updated successfully in table '%s'", tableName)
		w.WriteHeader(http.StatusOK)
		return
	}

	if err := h.service.InsertRow(ctx, models.InsertDataProps{
		TableName: tableName,
		Values:    form,
	}); err != nil {
		logger.Errorln(err.Error())
		logger.Error("Failed to insert row in table '%s'", tableName)
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}
	logger.Success("Row inserted successfully in table '%s'", tableName)
	w.WriteHeader(http.StatusCreated)
}

func (h DBHandler) DeleteRow(w http.ResponseWriter, r *http.Request) {
	tableName := r.PathValue("tableName")
	hash := r.PathValue("hash")
	page := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	if err := h.service.DeleteRow(r.Context(), tableName, hash, pageInt); err != nil {
		logger.Error("%s", err)
		logger.Error("Failed to delete row from table '%s'", tableName)
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}
	logger.Success("Row deleted successfully from table '%s'", tableName)
	resopnse.Success(w, http.StatusOK, nil)
}

func (h DBHandler) NewTableFormFileds(w http.ResponseWriter, r *http.Request) {
	fields := h.service.GetTableFormDataTypes()
	if fields == nil {
		resopnse.Error(w, http.StatusInternalServerError, fmt.Errorf("no data found"))
		return
	}

	resopnse.Success(w, http.StatusOK, fields)
}

func (h DBHandler) CreeteNewTable(w http.ResponseWriter, r *http.Request) {
	req := struct {
		TableName string           `json:"tableName"`
		Inputs    []database.Input `json:"inputs"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Info("Request data: %+v", req)
	if err := h.service.CreateTable(r.Context(), req.TableName, req.Inputs); err != nil {
		logger.Error("%s", err)
		logger.Error("Failed to create table '%s'", req.TableName)
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}
	logger.Success("Table '%s' created successfully with %d columns", req.TableName, len(req.Inputs))
	resopnse.Success(w, http.StatusCreated, nil)
}

type DeleteTableRequest struct {
	TableName         string `json:"tableName"`
	VerificationQuiry string `json:"verificationQuery"`
}

func (h *DBHandler) DeleteTable(w http.ResponseWriter, r *http.Request) {
	var req DeleteTableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("%s", err)
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}
	err := h.service.DeleteTable(r.Context(), req.TableName, req.VerificationQuiry)
	if err != nil {
		logger.Error("%s", err)
		logger.Error("Failed to delete table '%s'", req.TableName)
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}
	logger.Success("Table '%s' deleted successfully", req.TableName)
	w.WriteHeader(http.StatusNoContent)
}

func (h *DBHandler) ListHistory(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	pageInt = max(pageInt, 1)

	history, err := h.service.ListHistory(r.Context(), pageInt)
	if err != nil {
		logger.Error("Failed to list history: %v", err)
		logger.Error("Failed to fetch query history")
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}

	logger.Debug("Retrieved %d history entries (page %d)", len(history), pageInt)
	resopnse.Success(w, http.StatusOK, history)
}

func (h *DBHandler) ListRecentHistory(w http.ResponseWriter, r *http.Request) {
	history, err := h.service.ListHistory(r.Context(), 1)
	if err != nil {
		logger.Error("Failed to list recent history: %v", err)
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}

	if len(history) > 10 {
		history = history[:10]
	}

	resopnse.Success(w, http.StatusOK, history)
}
