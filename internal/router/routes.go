// Package router contains the handler for the database
package router

import (
	"fmt"
	"net/http"

	"github.com/biisal/rowsql/frontend"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

const apiPrefix = "/api/v1"

func route(method methodType, path string) string {
	return fmt.Sprintf("%s %s%s", method, apiPrefix, path)
}

func MountRouter(handler DBHandler) (*http.ServeMux, error) {
	mux := http.NewServeMux()

	mux.Handle("GET /", frontend.ReactHandler("/"))
	api := humago.New(mux, huma.DefaultConfig("RowSQL API", "1.0.0"))

	registerRoutes(api, handler)

	// fs := http.FileServer(http.Dir("frontend/static"))
	// mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	return mux, nil
}

func registerRoutes(api huma.API, h DBHandler) {
	huma.Register(api, huma.Operation{
		OperationID: "listRows",
		Method:      http.MethodGet,
		Path:        "/api/v1/tables/{tableName}",
		Summary:     "List rows of a table",
		Tags:        []string{"rows"},
	}, h.ListRows)

	huma.Register(api, huma.Operation{
		OperationID: "listTables",
		Method:      http.MethodGet,
		Path:        "/api/v1/tables",
		Summary:     "List of all tables",
		Tags:        []string{"tables"},
	}, h.ListTables)

	huma.Register(api, huma.Operation{
		OperationID: "listColumns",
		Method:      http.MethodGet,
		Path:        "/api/v1/tables/{tableName}/columns",
		Summary:     "List columns of a table",
		Tags:        []string{"columns"},
	}, h.ListColumns)

	huma.Register(api, huma.Operation{
		OperationID: "rowInsertOrUpdateForm",
		Method:      http.MethodGet,
		Path:        "/api/v1/tables/{tableName}/form",
		Summary:     "Get row insert/update form metadata",
		Tags:        []string{"tables"},
	}, h.RowInsertOrUpdateForm)

	huma.Register(api, huma.Operation{
		OperationID: "insertOrUpdateRow",
		Method:      http.MethodPost,
		Path:        "/api/v1/tables/{tableName}/form",
		Summary:     "Insert or update a row",
		Tags:        []string{"rows"},
	}, h.InsertOrUpdateRow)

	huma.Register(api, huma.Operation{
		OperationID: "deleteRow",
		Method:      http.MethodDelete,
		Path:        "/api/v1/tables/{tableName}/row/{hash}",
		Summary:     "Delete a row",
		Tags:        []string{"rows"},
	}, h.DeleteRow)

	huma.Register(api, huma.Operation{
		OperationID: "newTableFormFileds",
		Method:      http.MethodGet,
		Path:        "/api/v1/tables/form/new",
		Summary:     "Get data types for new table form",
		Tags:        []string{"tables"},
	}, h.NewTableFormFileds)

	huma.Register(api, huma.Operation{
		OperationID: "createNewTable",
		Method:      http.MethodPost,
		Path:        "/api/v1/tables/form/new",
		Summary:     "Create a new table",
		Tags:        []string{"tables"},
	}, h.CreeteNewTable)

	huma.Register(api, huma.Operation{
		OperationID: "deleteTable",
		Method:      http.MethodDelete,
		Path:        "/api/v1/tables",
		Summary:     "Delete a table",
		Tags:        []string{"tables"},
	}, h.DeleteTable)

	huma.Register(api, huma.Operation{
		OperationID: "listHistory",
		Method:      http.MethodGet,
		Path:        "/api/v1/history",
		Summary:     "List query history",
		Tags:        []string{"history"},
	}, h.ListHistory)

	huma.Register(api, huma.Operation{
		OperationID: "listRecentHistory",
		Method:      http.MethodGet,
		Path:        "/api/v1/history/recent",
		Summary:     "List recent history",
		Tags:        []string{"history"},
	}, h.ListRecentHistory)
}
