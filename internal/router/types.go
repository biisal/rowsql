package router

import "github.com/biisal/db-gui/internal/database/repo"

type methodType string

const (
	GET    methodType = "GET"
	POST   methodType = "POST"
	PUT    methodType = "PUT"
	DELETE methodType = "DELETE"
)

type ListRowsResponse struct {
	Page        int                `json:"page"`
	Rows        repo.ListDataRow   `json:"rows"`
	Cols        []repo.ListDataCol `json:"cols"`
	RowCount    int                `json:"rowCount"`
	ActiveTable string             `json:"activeTable"`
}
