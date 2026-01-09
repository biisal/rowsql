package router

import (
	"net/http"

	"github.com/biisal/db-gui/internal/logger"
	"github.com/biisal/db-gui/internal/response"
)

func CORS() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			if r.Method == "OPTIONS" {
				if origin != "" {
					logger.Debug("CORS preflight request from origin: %s", origin)
				}
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (h *DBHandler) withTable(handlerFunc http.HandlerFunc) http.Handler {
	return h.middlewareCheckTableExists(http.HandlerFunc(handlerFunc))
}

func (h *DBHandler) middlewareCheckTableExists(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tableName := r.PathValue("tableName")
		if err := h.service.CheckTableExits(r.Context(), tableName); err != nil {
			logger.Error("%s", err)
			response.Error(w, http.StatusNotFound, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}
