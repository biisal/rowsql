package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/biisal/rowsql/internal/apperr"
)

func createTestConfigFile(t testing.TB, configPath string, testJSON string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("Failed to create temp config directory: %v", err)
	}

	if err := os.WriteFile(configPath, []byte(testJSON), 0o644); err != nil {
		t.Fatalf("Failed to write temp config file: %v", err)
	}
}

func assertError(t testing.TB, err error, want error) {
	t.Helper()
	if err != want {
		t.Errorf("Expected error %v, got %v", want, err)
	}
}

func assertConfigs(t testing.TB, configs []ConnectionConfig, want []ConnectionConfig) {
	t.Helper()
	if !slices.Equal(configs, want) {
		t.Errorf("Expected %+v, got %+v", want, configs)
	}
}

func makeJSONString(configs []map[string]any) string {
	var sb strings.Builder

	var connections []string

	for _, config := range configs {
		connections = append(connections,
			fmt.Sprintf(`{"port": %v, "db_string": %q, "env": %q}`,
				config["port"],
				config["db_string"],
				config["env"],
			))
	}

	connString := strings.Join(connections, ",")

	fmt.Fprintf(&sb, `{
		"log_file_path": "~/.rowsql/rowsql.log",
		"connections": [%s]
	}`, connString)

	return sb.String()
}

type MockPromter struct{}

func (m *MockPromter) AskConnection(configs []ConnectionConfig) (*ConnectionConfig, error) {
	return &configs[0], nil
}

func TestListConfig(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name       string
		fileConfig []map[string]any
		wantConfig []ConnectionConfig
		wantError  error
	}{
		{
			name: "Valid config file",
			fileConfig: []map[string]any{
				{
					"port":      8000,
					"env":       "production",
					"db_string": "test.db",
				},
			},

			wantConfig: []ConnectionConfig{{Port: 8000, DBString: "test.db", Env: "production"}},
			wantError:  nil,
		},
		{
			name: "InValid config file",
			fileConfig: []map[string]any{
				{
					"port":      8000,
					"env":       "production",
					"db_string": "test.db",
				},
			},

			wantConfig: []ConnectionConfig{},
			wantError:  apperr.ErrorInvalidJSONFile,
		},
	}

	configService := &ConfigServiceImpl{
		prompter: &MockPromter{},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			configPath := filepath.Join(tempDir, "config.json")

			createTestConfigFile(t, configPath, makeJSONString(test.fileConfig))

			got, err := configService.ParseConfig(configPath)

			assertConfigs(t, got.Connections, test.wantConfig)
			assertError(t, err, test.wantError)
		})
	}
}

func assertCofig(t testing.TB, got Config, want Config) {
	t.Helper()
	if got != want {
		t.Errorf("Expected %+v, got %+v", want, got)
	}
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name       string
		wantErr    error
		configPath []string
		wantConfig *Config
		fileConfig []map[string]any
	}{
		{
			name: "valid config path",
			fileConfig: []map[string]any{
				{
					"port":      8000,
					"env":       "production",
					"db_string": "test.db",
				},
			},
			wantErr: nil,
			wantConfig: &Config{
				AppConfig: AppConfig{
					LogFilePath: "~/.rowsql/rowsql.log",
				},
				ConnectionConfig: ConnectionConfig{
					Port:     8000,
					DBString: "test.db",
					Env:      "production",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			configService := &ConfigServiceImpl{
				prompter: &MockPromter{},
			}
			configPath := ""
			if len(test.configPath) > 0 {
				configPath = test.configPath[0]
			}
			got := configService.LoadConfig(configPath)
			assertCofig(t, *got, *test.wantConfig)
		})
	}
}
