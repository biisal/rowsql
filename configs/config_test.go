package configs

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
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
	if want == nil {
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		return
	}
	if !errors.Is(err, want) {
		t.Errorf("Expected error %v, got %v", want, err)
	}
}

func assertConfigs(t testing.TB, configs []ConnectionConfig, want []ConnectionConfig) {
	t.Helper()
	if !slices.Equal(configs, want) {
		t.Errorf("Expected %+v, got %+v", want, configs)
	}
}

type MockPrompter struct{}

func (m *MockPrompter) AskConnection(configs []ConnectionConfig) (*ConnectionConfig, error) {
	return &configs[0], nil
}

func toJSON(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	data, _ := json.Marshal(v)
	return string(data)
}

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name       string
		input      any
		wantConfig []ConnectionConfig
		wantError  error
	}{
		{
			name: "Valid config file",
			input: map[string]any{
				"connections": []any{
					map[string]any{
						"port":      8000,
						"db_string": "test.db",
						"env":       "production",
					},
				},
			},
			wantConfig: []ConnectionConfig{{
				Port:     8000,
				DBString: "test.db",
				Env:      "production",
			}},
			wantError: nil,
		},
		{
			name:       "Invalid config file",
			input:      `{invalid json}`,
			wantConfig: []ConnectionConfig{},
			wantError:  apperr.ErrorInvalidJSONFile,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tempDir := t.TempDir()
			configPath := filepath.Join(tempDir, "config.json")
			createTestConfigFile(t, configPath, toJSON(test.input))

			configService := &ConfigServiceImpl{}
			got, err := configService.parseConfig(configPath)

			assertError(t, err, test.wantError)
			if test.wantError == nil {
				assertConfigs(t, got.Connections, test.wantConfig)
			}
		})
	}
}

func assertConfig(t testing.TB, got Config, want Config) {
	t.Helper()
	if got != want {
		t.Errorf("Expected %+v, got %+v", want, got)
	}
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name       string
		wantErr    error
		configPath string
		wantConfig Config
		input      any
	}{
		{
			name: "valid config path",
			input: map[string]any{
				"connections": []any{
					map[string]any{
						"port":      8000,
						"db_string": "test.db",
						"env":       "production",
					},
				},
				"log_file_path": "~/.rowsql/rowsql.log",
			},
			wantErr: nil,
			wantConfig: Config{
				AppConfig: AppConfig{
					LogFilePath:     "~/.rowsql/rowsql.log",
					MaxItemsPerPage: MaxItemsPerPage,
					MinItemsPerPage: MinItemsPerPage,
				},
				ConnectionConfig: ConnectionConfig{
					Host:     "127.0.0.1",
					Port:     8000,
					DBString: "test.db",
					Env:      "production",
				},
			},
		},
		{
			name: "empty connections",
			input: map[string]any{
				"connections": []any{},
			},
			wantErr: apperr.ErrorNoConfigsFound,
		},
		{
			name: "missing required connection fields",
			input: map[string]any{
				"connections": []any{
					map[string]any{"port": 8080},
				},
			},
			wantErr: apperr.ErrorInvalidJSONFile,
		},
		{
			name: "uses defaults for missing app config",
			input: map[string]any{
				"connections": []any{
					map[string]any{"port": 8080, "db_string": "default.db"},
				},
			},
			wantErr: nil,
			wantConfig: Config{
				AppConfig: DefaultConfig().AppConfig,
				ConnectionConfig: ConnectionConfig{
					Host:     "127.0.0.1",
					Port:     8080,
					DBString: "default.db",
					Env:      EnvProduction,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tempDir := t.TempDir()
			configPath := filepath.Join(tempDir, "config.json")
			createTestConfigFile(t, configPath, toJSON(test.input))

			configService := NewConfigService(&MockPrompter{})

			got, err := configService.LoadConfig(configPath)
			assertError(t, err, test.wantErr)
			assertConfig(t, got, test.wantConfig)
		})
	}
}

type failingPrompter struct{}

func (failingPrompter) AskConnection([]ConnectionConfig) (*ConnectionConfig, error) {
	return nil, errors.New("prompter must not be called in dev mode")
}

func TestLoadDevConfigUsesFirstConnectionWithoutPrompting(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	createTestConfigFile(t, configPath, `{"connections":[
		{"port":8081,"db_string":"first.db","env":"production"},
		{"port":8082,"db_string":"second.db"}
	]}`)

	cfg, err := NewConfigService(failingPrompter{}).LoadDevConfig(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DBString != "first.db" || cfg.Port != 8081 {
		t.Errorf("connection = %s on port %d, want first.db on 8081", cfg.DBString, cfg.Port)
	}
	if cfg.Env != EnvDevelopment {
		t.Errorf("env = %s, want %s", cfg.Env, EnvDevelopment)
	}
}

func TestLoadConfigStillPromptsWithoutDevMode(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	createTestConfigFile(t, configPath, `{"connections":[{"port":1,"db_string":"a.db"},{"port":2,"db_string":"b.db"}]}`)

	if _, err := NewConfigService(failingPrompter{}).LoadConfig(configPath); err == nil {
		t.Fatal("expected the prompter to be consulted when several connections exist")
	}
}
