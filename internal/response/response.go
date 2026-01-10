// Package response provides functions for creating JSON responses.
package response

import (
	"encoding/json"
	"net/http"

	"github.com/biisal/rowsql/internal/logger"
)

type Response struct {
	Error   string `json:"error,omitempty"`
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
}

func Success(w http.ResponseWriter, status int, data any) {
	jsonData, err := json.Marshal(Response{Success: true, Data: data})
	if err != nil {
		logger.Error("failed to marshal response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(jsonData); err != nil {
		logger.Error("failed to write error response: %v", err)
	}
}

func Error(w http.ResponseWriter, status int, errMsg error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	jsonData, err := json.Marshal(Response{Error: errMsg.Error()})
	if err != nil {
		logger.Error("failed to marshal response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(jsonData); err != nil {
		logger.Error("failed to write error response: %v", err)
	}
}
