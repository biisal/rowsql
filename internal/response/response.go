package resopnse

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Response struct {
	Error   string `json:"error,omitempty"`
	Success bool   `json:"success"`
	Data    any    `json:"data"`
}

func Success(w http.ResponseWriter, status int, data any) {
	jsonData, err := json.Marshal(Response{Success: true, Data: data})
	if err != nil {
		slog.Error("failed to marshal response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err = w.Write(jsonData); err != nil {
		slog.Error("failed to write error response", "error", err)
	}
}

func Error(w http.ResponseWriter, status int, errMsg error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	jsonData, err := json.Marshal(Response{Error: errMsg.Error()})
	if err != nil {
		slog.Error("failed to marshal response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(jsonData); err != nil {
		slog.Error("failed to write error response", "error", err)
	}
}
