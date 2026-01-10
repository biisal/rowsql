// Package router contains the handler for the database
package router

import (
	"fmt"
	"net/http"

	"github.com/biisal/rowsql/frontend"
)

const apiPrefix = "/api/v1"

func route(method methodType, path string) string {
	return fmt.Sprintf("%s %s%s", method, apiPrefix, path)
}

func MountRouter(handler DBHandler) (*http.ServeMux, error) {
	mux := http.NewServeMux()

	mux.Handle("GET /", frontend.ReactHandler("/"))

	mux.HandleFunc(route(GET, "/tables"), handler.ListTables)
	mux.Handle(route(GET, "/tables/{tableName}"), handler.withTable(handler.ListRows))
	mux.Handle(route(GET, "/tables/{tableName}/form"), handler.withTable(handler.RowInsertForm))
	mux.Handle(route(GET, "/tables/{tableName}/columns"), handler.withTable(handler.ListColumns))
	mux.Handle(route(POST, "/tables/{tableName}/form"), handler.withTable(handler.InsertOrUpdateRow))
	mux.Handle(route(DELETE, "/tables/{tableName}/row/{hash}"), handler.withTable(handler.DeleteRow))

	mux.HandleFunc(route(GET, "/tables/form/new"), handler.NewTableFormFileds)
	mux.HandleFunc(route(POST, "/tables/form/new"), handler.CreeteNewTable)
	mux.HandleFunc(route(DELETE, "/tables"), handler.DeleteTable)
	mux.HandleFunc(route(GET, "/history"), handler.ListHistory)
	mux.HandleFunc(route(GET, "/history/recent"), handler.ListRecentHistory)

	// fs := http.FileServer(http.Dir("frontend/static"))
	// mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	return mux, nil
}
