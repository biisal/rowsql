// Package router contains the handler for the database
package router

import (
	"fmt"
	"net/http"

	"github.com/biisal/db-gui/frontend"
)

const apiPrefix = "/api/v1"

func route(method methodType, path string) string {
	return fmt.Sprintf("%s %s%s", method, apiPrefix, path)
}

func MountRouter(handler DBHandler) (*http.ServeMux, error) {
	mux := http.NewServeMux()

	mux.Handle("/", frontend.ReactHandler("/"))

	mux.HandleFunc(route(GET, "/tables"), handler.ListTables)
	mux.HandleFunc(route(GET, "/table/{tableName}"), handler.TableRows)
	mux.HandleFunc(route(GET, "/table/{tableName}/form"), handler.RowInsertForm)
	mux.HandleFunc(route(POST, "/table/{tableName}/form"), handler.InsertOrUpdateRow)
	mux.HandleFunc(route(DELETE, "/table/{tableName}/row/{hash}"), handler.DeleteRow)
	mux.HandleFunc(route(GET, "/table/form/new"), handler.NewTableFormFileds)
	mux.HandleFunc(route(POST, "/table/form/new"), handler.CreeteNewTable)
	mux.HandleFunc(route(DELETE, "/table"), handler.DeleteTable)
	mux.HandleFunc(route(GET, "/history"), handler.ListHistory)
	mux.HandleFunc(route(GET, "/history/recent"), handler.ListRecentHistory)

	// fs := http.FileServer(http.Dir("frontend/static"))
	// mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	return mux, nil
}
