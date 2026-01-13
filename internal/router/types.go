package router

import (
	"github.com/biisal/rowsql/internal/database/models"
)

type methodType string

const (
	GET    methodType = "GET"
	POST   methodType = "POST"
	PUT    methodType = "PUT"
	DELETE methodType = "DELETE"
)

type ListRowsResponse struct {
	Page        int                  `json:"page"`
	Rows        models.ListDataRow   `json:"rows"`
	Cols        []models.ListDataCol `json:"cols"`
	RowCount    int                  `json:"rowCount"`
	ActiveTable string               `json:"activeTable"`
	HasNextPage bool                 `json:"hasNextPage"`
	TotalPages  int                  `json:"totalPages"`
}
