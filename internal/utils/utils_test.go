package utils

import (
	"os"
	"strings"
	"testing"

	"github.com/biisal/rowsql/configs"
)

func TestDetectDriver(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		wantDriver     configs.Driver
		wantConnection string
		wantContains   []string
	}{
		{
			name:           "postgres url",
			input:          "postgres://u:p@localhost:5432/db",
			wantDriver:     configs.DriverPostgres,
			wantConnection: "postgres://u:p@localhost:5432/db",
		},
		{name: "sqlite file", input: "data.db", wantDriver: configs.DriverSQLite, wantConnection: "data.db"},
		{name: "sqlite url", input: "sqlite://data.db", wantDriver: configs.DriverSQLite, wantConnection: "data.db"},
		{
			name:         "mysql dsn",
			input:        "u:p@tcp(localhost:3306)/db",
			wantDriver:   configs.DriverMySQL,
			wantContains: []string{"u:p@tcp(localhost:3306)/db", "parseTime=true"},
		},
		{
			name:         "mysql url",
			input:        "mysql://u:p@localhost:3307/db?charset=utf8mb4",
			wantDriver:   configs.DriverMySQL,
			wantContains: []string{"u:p@tcp(localhost:3307)/db", "charset=utf8mb4", "parseTime=true"},
		},
		{
			name:         "mysql url default port",
			input:        "mysql://u:p@localhost/db",
			wantDriver:   configs.DriverMySQL,
			wantContains: []string{"tcp(localhost:3306)/db"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connection := tt.input
			driver, err := DetectDriver(&connection)
			if err != nil {
				t.Fatal(err)
			}
			if driver != tt.wantDriver {
				t.Fatalf("driver = %s, want %s", driver, tt.wantDriver)
			}
			if tt.wantConnection != "" && connection != tt.wantConnection {
				t.Fatalf("connection = %q, want %q", connection, tt.wantConnection)
			}
			for _, fragment := range tt.wantContains {
				if !strings.Contains(connection, fragment) {
					t.Fatalf("connection %q does not contain %q", connection, fragment)
				}
			}
		})
	}
}

func TestDetectDriverRejectsEmptyAndUnknown(t *testing.T) {
	for _, input := range []string{"", "something-unknown"} {
		connection := input
		if _, err := DetectDriver(&connection); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
}

func TestReplaceTildeWithHomeDir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	tests := map[string]string{
		"~":              home,
		"~/data.db":      home + "/data.db",
		"~bob/data.db":   "~bob/data.db",
		"/abs/~/data.db": "/abs/~/data.db",
		"relative.db":    "relative.db",
	}
	for input, want := range tests {
		got, err := ReplaceTildeWithHomeDir(input)
		if err != nil || got != want {
			t.Errorf("ReplaceTildeWithHomeDir(%q) = %q, %v; want %q", input, got, err, want)
		}
	}
}
