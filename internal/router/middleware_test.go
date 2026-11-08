package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireLocalAccess(t *testing.T) {
	const token = "secret-token"
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	guarded := RequireLocalAccess(token)(okHandler)

	tests := []struct {
		name       string
		host       string
		path       string
		authHeader string
		wantStatus int
	}{
		{
			name:       "api without token",
			host:       "127.0.0.1:8080",
			path:       "/api/v1/tables",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "api with wrong token",
			host:       "127.0.0.1:8080",
			path:       "/api/v1/tables",
			authHeader: "Bearer wrong",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "api with valid token",
			host:       "127.0.0.1:8080",
			path:       "/api/v1/tables",
			authHeader: "Bearer " + token,
			wantStatus: http.StatusOK,
		},
		{
			name:       "api with localhost host",
			host:       "localhost:9000",
			path:       "/api/v1/tables",
			authHeader: "Bearer " + token,
			wantStatus: http.StatusOK,
		},
		{
			name:       "api with foreign host",
			host:       "evil.example:8080",
			path:       "/api/v1/tables",
			authHeader: "Bearer " + token,
			wantStatus: http.StatusMisdirectedRequest,
		},
		{name: "static page without token", host: "127.0.0.1:8080", path: "/", wantStatus: http.StatusOK},
		{
			name:       "static page with foreign host",
			host:       "evil.example",
			path:       "/",
			wantStatus: http.StatusMisdirectedRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			request.Host = tt.host
			if tt.authHeader != "" {
				request.Header.Set("Authorization", tt.authHeader)
			}
			recorder := httptest.NewRecorder()
			guarded.ServeHTTP(recorder, request)
			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.wantStatus)
			}
		})
	}
}

func TestRequireLocalAccessWithoutTokenAllowsApi(t *testing.T) {
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	guarded := RequireLocalAccess("")(okHandler)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/tables", nil)
	request.Host = "localhost:3000"
	recorder := httptest.NewRecorder()
	guarded.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
}
