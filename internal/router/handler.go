package router

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/biisal/db-gui/internal/database"
	"github.com/biisal/db-gui/internal/database/repo"
	resopnse "github.com/biisal/db-gui/internal/response"
)

type DbHandler struct {
	service DbService
}

type BaseHtmlData struct {
	Tables      []repo.ListTablesRow
	Cols        []repo.ListDataCol
	ActiveTable string
}

type ErrorMessage struct {
	Message string
}

func NewHandler(service DbService) DbHandler {
	return DbHandler{
		service,
	}
}

func (h DbHandler) getBaseData(ctx context.Context, tableName ...string) (*BaseHtmlData, error) {
	tables, err := h.service.ListTables(ctx)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	if len(tableName) == 0 {
		return &BaseHtmlData{Tables: tables}, nil
	}
	cols, err := h.service.ListCols(ctx, tableName[0])
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return &BaseHtmlData{Tables: tables, Cols: cols, ActiveTable: tableName[0]}, nil
}

func (h DbHandler) ListTables(w http.ResponseWriter, r *http.Request) {
	tables, err := h.service.ListTables(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resopnse.Success(w, http.StatusOK, tables)
}

func (h DbHandler) TableRows(w http.ResponseWriter, r *http.Request) {
	tableName := r.PathValue("tableName")
	page := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	pageInt = max(pageInt, 1)
	rows, err := h.service.ListRows(r.Context(), tableName, pageInt)
	if err != nil {
		slog.Error(err.Error())
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}

	cols, err := h.service.ListCols(r.Context(), tableName)
	if err != nil {
		slog.Error(err.Error())
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}

	data := struct {
		Page        int
		Rows        repo.ListDataRow
		Cols        []repo.ListDataCol
		ActiveTable string
	}{
		pageInt,
		rows,
		cols,
		tableName,
	}
	resopnse.Success(w, http.StatusOK, data)
}

func (h DbHandler) TableInsertForm(w http.ResponseWriter, r *http.Request) {
	tableName := r.PathValue("tableName")
	tables, err := h.service.ListTables(r.Context())
	if err != nil {
		slog.Error(err.Error())
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}
	found := false
	for _, table := range tables {
		if table.TableName == strings.TrimSpace(tableName) {
			found = true
		}
	}
	if !found {
		resopnse.Error(w, http.StatusNotFound, fmt.Errorf("table %s not found", tableName))
		return
	}
	action := "Insert"
	hash := strings.TrimSpace(r.URL.Query().Get("hash"))
	page := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	var initialRow []any
	if hash != "" {
		var err error
		action = "Update"
		initialRow, err = h.service.GetRow(r.Context(), tableName, hash, pageInt)
		if err != nil {
			slog.Error(err.Error())
			resopnse.Error(w, http.StatusInternalServerError, err)
			return
		}
	}

	basseData, err := h.getBaseData(r.Context(), tableName)
	if err != nil {
		slog.Error(err.Error())
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
		BaseHtmlData
	}{
		action,
		*basseData,
	}
	resopnse.Success(w, http.StatusOK, data)
}

func (h DbHandler) InsertOrUpdateRow(w http.ResponseWriter, r *http.Request) {
	tableName := r.PathValue("tableName")
	var form = make(map[string]repo.FormValue)
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		slog.Error(err.Error())
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}

	page := r.URL.Query().Get("page")
	hash := strings.TrimSpace(r.URL.Query().Get("hash"))

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}

	if hash != "" {
		if err := h.service.UpdateRow(r.Context(), form, tableName, hash, pageInt); err != nil {
			slog.Error(err.Error())
			resopnse.Error(w, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	if err := h.service.InsertRow(r.Context(), repo.InsertDataProps{
		TableName: tableName,
		Values:    form,
	}); err != nil {
		slog.Error(err.Error())
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h DbHandler) DeleteRow(w http.ResponseWriter, r *http.Request) {
	tableName := r.PathValue("tableName")
	hash := r.PathValue("hash")
	page := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	if err := h.service.DeleteRow(r.Context(), tableName, hash, pageInt); err != nil {
		slog.Error(err.Error())
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}
	resopnse.Success(w, http.StatusOK, nil)
}

func (h DbHandler) NewTableFormFileds(w http.ResponseWriter, r *http.Request) {
	fields := h.service.GetTableFormDataTypes()
	if fields == nil {
		resopnse.Error(w, http.StatusInternalServerError, fmt.Errorf("no data found"))
		return
	}

	resopnse.Success(w, http.StatusOK, fields)
}

func (h DbHandler) CreeteNewTable(w http.ResponseWriter, r *http.Request) {
	var req = struct {
		TableName string           `json:"tableName"`
		Inputs    []database.Input `json:"inputs"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info("Request data", "req", req)
	if err := h.service.CreateTable(r.Context(), req.TableName, req.Inputs); err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

}

type DeleteTableRequest struct {
	TableName         string `json:"tableName"`
	VerificationQuiry string `json:"verificationQuery"`
}

func (h *DbHandler) DeleteTable(w http.ResponseWriter, r *http.Request) {
	var req DeleteTableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error(err.Error())
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}
	err := h.service.DeleteTable(r.Context(), req.TableName, req.VerificationQuiry)
	if err != nil {
		slog.Error(err.Error())
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *DbHandler) ListHistory(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	pageInt = max(pageInt, 1)

	history, err := h.service.ListHistory(r.Context(), pageInt)
	if err != nil {
		slog.Error("Failed to list history", "err", err)
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}

	resopnse.Success(w, http.StatusOK, history)
}

func (h *DbHandler) ListRecentHistory(w http.ResponseWriter, r *http.Request) {
	history, err := h.service.ListHistory(r.Context(), 1)
	if err != nil {
		slog.Error("Failed to list recent history", "err", err)
		resopnse.Error(w, http.StatusInternalServerError, err)
		return
	}

	if len(history) > 10 {
		history = history[:10]
	}

	resopnse.Success(w, http.StatusOK, history)
}
