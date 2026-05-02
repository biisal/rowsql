// Package configs contains the configuration for the application
package configs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/biisal/rowsql/internal/apperr"
)

type Driver string

const (
	DriverPostgres  Driver = "pgx"
	DriverMySQL     Driver = "mysql"
	DriverSQLite    Driver = "sqlite"
	EnvDevelopment  string = "development"
	EnvProduction   string = "production"
	MinItemsPerPage int    = 10
	MaxItemsPerPage int    = 100
)

var Drivers = map[Driver]Driver{
	"pgx":    DriverPostgres,
	"mysql":  DriverMySQL,
	"sqlite": DriverSQLite,
}

type ConnectionConfig struct {
	Port     int    `json:"port"`
	DBString string `json:"db_string"`
	Env      string `json:"env,omitempty"`
}

type AppConfig struct {
	DisableAutoUpdate bool   `json:"disable_auto_update,omitempty"`
	Driver            Driver `json:"driver,omitempty"`
	MaxItemsPerPage   int    `json:"max_items_per_page,omitempty"`
	MinItemsPerPage   int    `json:"min_items_per_page,omitempty"`
	LogFilePath       string `json:"log_file_path,omitempty"`
}

type ConfigFileFields struct {
	Connections []ConnectionConfig `json:"connections"`
	AppConfig
}

type Config struct {
	ConnectionConfig
	AppConfig
}

func DefaultConfig() Config {
	userHome, _ := os.UserHomeDir()
	return Config{
		ConnectionConfig: ConnectionConfig{
			Port:     8080,
			DBString: "test.db",
			Env:      EnvProduction,
		},
		AppConfig: AppConfig{
			MaxItemsPerPage: MaxItemsPerPage,
			MinItemsPerPage: MinItemsPerPage,
			LogFilePath:     filepath.Join(userHome, ".rowsql", "rowsql.log"),
		},
	}
}

type ConfigService interface {
	LoadConfig(path ...string) (Config, error)
	GetConfigPath() (string, error)
}

type ConfigServiceImpl struct {
	prompter Prompter
}

func NewConfigService(prompter Prompter) ConfigService {
	return &ConfigServiceImpl{
		prompter: prompter,
	}
}

func (c *ConfigServiceImpl) LoadConfig(configPath ...string) (Config, error) {
	var path string
	if len(configPath) > 0 && configPath[0] != "" {
		path = configPath[0]
	} else {
		var err error
		path, err = c.GetConfigPath()
		if err != nil {
			return Config{}, err
		}
	}

	fields, err := c.parseConfig(path)
	if err != nil {
		return Config{}, fmt.Errorf("parsing config: %w", err)
	}

	var conn *ConnectionConfig
	if len(fields.Connections) == 0 {
		return Config{}, apperr.ErrorNoConfigsFound
	} else if len(fields.Connections) == 1 {
		conn = &fields.Connections[0]
	} else {
		conn, err = c.prompter.AskConnection(fields.Connections)
		if err != nil {
			return Config{}, fmt.Errorf("selecting connection: %w", err)
		}
	}

	cfg := Config{
		ConnectionConfig: *conn,
		AppConfig:        fields.AppConfig,
	}

	c.applyDefaults(&cfg)
	return cfg, nil
}

func (c *ConfigServiceImpl) applyDefaults(cfg *Config) {
	if cfg.Env == "" {
		cfg.Env = EnvProduction
	}
	if cfg.LogFilePath == "" {
		userHome, _ := os.UserHomeDir()
		cfg.LogFilePath = filepath.Join(userHome, ".rowsql", "rowsql.log")
	}
	if cfg.MaxItemsPerPage == 0 {
		cfg.MaxItemsPerPage = MaxItemsPerPage
	}
	if cfg.MinItemsPerPage == 0 {
		cfg.MinItemsPerPage = MinItemsPerPage
	}
}

func (c *ConfigServiceImpl) GetConfigPath() (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home dir: %w", err)
	}

	configDir := filepath.Join(userHome, ".rowsql")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return "", fmt.Errorf("creating config directory: %w", err)
	}

	fileName := "config.json"
	fullPath := filepath.Join(configDir, fileName)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		if err := promptForDefaultEnv(configDir, fileName); err != nil {
			return "", err
		}
	} else if err != nil {
		return "", fmt.Errorf("checking config file: %w", err)
	}

	return fullPath, nil
}

func (c *ConfigServiceImpl) parseConfig(path string) (ConfigFileFields, error) {
	var fields ConfigFileFields
	data, err := os.ReadFile(path)
	if err != nil {
		return fields, fmt.Errorf("reading file: %w", err)
	}

	if err := json.Unmarshal(data, &fields); err != nil {
		return fields, apperr.ErrorInvalidJSONFile
	}

	for _, conn := range fields.Connections {
		if conn.Port == 0 || conn.DBString == "" {
			return ConfigFileFields{}, apperr.ErrorInvalidJSONFile
		}
	}
	return fields, nil
}
