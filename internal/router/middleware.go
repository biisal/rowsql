package router

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"net"
	"net/http"
	"strings"
)

const (
	apiPathPrefix    = "/api/"
	sessionTokenSize = 32
	bearerPrefix     = "Bearer "
)

func NewSessionToken() (string, error) {
	tokenBytes := make([]byte, sessionTokenSize)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(tokenBytes), nil
}

func RequireLocalAccess(sessionToken string, extraAllowedHosts ...string) func(http.Handler) http.Handler {
	allowedHosts := map[string]bool{"localhost": true, "127.0.0.1": true, "::1": true}
	for _, host := range extraAllowedHosts {
		allowedHosts[host] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !allowedHosts[hostnameOf(r.Host)] {
				http.Error(w, "misdirected request", http.StatusMisdirectedRequest)
				return
			}
			if sessionToken != "" && strings.HasPrefix(r.URL.Path, apiPathPrefix) && !hasValidToken(r, sessionToken) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func hostnameOf(hostWithPort string) string {
	hostname, _, err := net.SplitHostPort(hostWithPort)
	if err != nil {
		return strings.Trim(hostWithPort, "[]")
	}
	return hostname
}

func hasValidToken(r *http.Request, sessionToken string) bool {
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, bearerPrefix) {
		return false
	}
	givenToken := strings.TrimPrefix(header, bearerPrefix)
	return subtle.ConstantTimeCompare([]byte(givenToken), []byte(sessionToken)) == 1
}
