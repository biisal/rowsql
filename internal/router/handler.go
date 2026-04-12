package router

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/biisal/rowsql/internal/apperr"
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
	Tables []models.ListTablesRow
	Cols   []models.ColMetaData
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

// TODO : remove ths function
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
	return &BaseHTMLData{Tables: tables, Cols: cols}, nil
}

// @Summary List columns of a table
// @Description get list of all columns for a specific table
// @Tags columns
// @Accept json
// @Produce json
// @Param tableName path string true "Table Name"
// @Success 200 {object} response.Response{data=[]models.ColMetaData}
// @Router /api/v1/tables/{tableName}/columns [get]
func (h DBHandler) ListColumns(w http.ResponseWriter, r *http.Request) {
	cols, err := h.service.ListCols(r.Context(), r.PathValue("tableName"))
	if err != nil {
		logger.Error("%s", err)
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}
	resopnse.Success(w, http.StatusOK, cols)
}

// @Summary list of all tables
// @Description get list of all tables
// @Tags tables
// @Accept json
// @Produce      json
// @Success      200  {object}  map[string]any{}
// @Router       /api/v1/tables [get]
func (h DBHandler) ListTables(w http.ResponseWriter, r *http.Request) {
	tables, err := h.service.ListTables(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	logger.Info("Fetched %d tables", len(tables))
	resopnse.Success(w, http.StatusOK, tables)
}

// @Summary List rows of a table
// @Description get paginated list of rows for a specific table
// @Tags rows
// @Accept json
// @Produce json
// @Param tableName path string true "Table Name"
// @Param page query int false "Page number" default(1)
// @Param column query string false "Column to order by"
// @Param order query string false "Order direction (ASC/DESC)"
// @Success 200 {object} response.Response{data=router.ListRowsResponse}
// @Router /api/v1/tables/{tableName} [get]
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
		cols, err := h.service.ListCols(r.Context(), tableName)
		if err != nil {
			resopnse.Error(w, http.StatusInternalServerError, err)
		}
		for _, col := range cols {
			if col.Name == colParam {
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

// @Summary Get row insert/update form metadata
// @Description get metadata for columns to build an insert or update form
// @Tags tables
// @Accept json
// @Produce json
// @Param tableName path string true "Table Name"
// @Success 200 {object} response.Response{data=object{Cols=[]models.ColMetaData}}
// @Router /api/v1/tables/{tableName}/form [get]
func (h DBHandler) RowInsertOrUpdateForm(w http.ResponseWriter, r *http.Request) {
	tableName := r.PathValue("tableName")

	// hash := strings.TrimSpace(r.URL.Query().Get("hash"))
	// page := r.URL.Query().Get("page")
	// pageInt, err := strconv.Atoi(page)
	// if err != nil {
	// 	pageInt = 1
	// }
	// var initialRow []any
	// if hash != "" {
	// 	initialRow, err = h.service.GetRow(r.Context(), tableName, hash, pageInt)
	// 	if err != nil {
	// 		logger.Error("%s", err)
	// 		resopnse.Error(w, http.StatusInternalServerError, err)
	// 		return
	// 	}
	// }

	colsMeta, err := h.service.ListCols(r.Context(), tableName)
	if err != nil {
		logger.Error("%s", err)
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}
	// if len(initialRow) == len(colsMeta) {
	// 	for i := range colsMeta {
	// 		colsMeta[i].Value = initialRow[i]
	// 	}
	// }
	data := struct {
		Cols []models.ColMetaData
	}{
		Cols: colsMeta,
	}
	resopnse.Success(w, http.StatusOK, data)
}

// @Summary Insert or update a row
// @Description insert a new row or update an existing row if hash is provided
// @Tags rows
// @Accept json
// @Produce json
// @Param tableName path string true "Table Name"
// @Param hash query string false "Row hash (for update)"
// @Param page query int false "Page number (for update context)" default(1)
// @Param body body models.InsertRowRequest true "Row data"
// @Success 200 {object} response.Response "Success"
// @Success 201 {object} response.Response "Created"
// @Router /api/v1/tables/{tableName}/form [post]
func (h DBHandler) InsertOrUpdateRow(w http.ResponseWriter, r *http.Request) {
	tableName := r.PathValue("tableName")
	ctx := r.Context()
	form := models.InsertRowRequest{}

	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		logger.Error("%s", err)
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}
	logger.Info("Request data: %+v", form)

	page := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}

	hash := strings.TrimSpace(r.URL.Query().Get("hash"))
	if hash != "" {
		if err := h.service.UpdateRow(ctx, form.Data, tableName, hash, pageInt); err != nil {
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
		Values:    form.Data,
	}); err != nil {
		logger.Errorln(err.Error())
		logger.Error("Failed to insert row in table '%s'", tableName)
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}
	logger.Success("Row inserted successfully in table '%s'", tableName)
	w.WriteHeader(http.StatusCreated)
}

// @Summary Delete a row
// @Description delete a specific row by its hash
// @Tags rows
// @Accept json
// @Produce json
// @Param tableName path string true "Table Name"
// @Param hash path string true "Row hash"
// @Param page query int false "Page number" default(1)
// @Success 200 {object} response.Response "Success"
// @Router /api/v1/tables/{tableName}/row/{hash} [delete]
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

// @Summary Get data types for new table form
// @Description get available numeric and string data types for creating a new table
// @Tags tables
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=service.FormDatatype}
// @Router /api/v1/tables/form/new [get]
func (h DBHandler) NewTableFormFileds(w http.ResponseWriter, r *http.Request) {
	fields := h.service.GetTableFormDataTypes()
	if fields == nil {
		resopnse.Error(w, http.StatusInternalServerError, fmt.Errorf("no data found"))
		return
	}

	resopnse.Success(w, http.StatusOK, fields)
}

// @Summary Create a new table
// @Description create a new table with specified name and columns
// @Tags tables
// @Accept json
// @Produce json
// @Param body body object{tableName=string,inputs=[]models.ColValues} true "Table Schema"
// @Success 201 {object} response.Response "Created"
// @Router /api/v1/tables/form/new [post]
func (h DBHandler) CreeteNewTable(w http.ResponseWriter, r *http.Request) {
	req := struct {
		TableName string             `json:"tableName"`
		Inputs    []models.ColValues `json:"inputs"`
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

// @Summary Delete a table
// @Description drop a table from the database
// @Tags tables
// @Accept json
// @Produce json
// @Param body body router.DeleteTableRequest true "Delete table request"
// @Success 204 "No Content"
// @Router /api/v1/tables [delete]
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

// @Summary List query history
// @Description get paginated list of query history
// @Tags history
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Success 200 {object} response.Response{data=[]models.History}
// @Router /api/v1/history [get]
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

// @Summary List recent history
// @Description get the last 10 query history entries
// @Tags history
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]models.History}
// @Router /api/v1/history/recent [get]
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
