// Package configs contains the configuration for the application
package configs

import (
	"encoding/json"
	"os"

	"github.com/biisal/rowsql/internal/apperr"
	"github.com/biisal/rowsql/internal/logger"
)

type Driver string

const (
	DriverPostgres Driver = "pgx"
	DriverMySQL    Driver = "mysql"
	DriverSQLite   Driver = "sqlite"
	EnvDevelopment string = "development"
	EnvProduction  string = "production"
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
	DisableAutoUpdate bool `json:"disable_auto_update,omitempty"`
	Driver            Driver

	MaxItemsPerPage int    `json:"max_items_per_page,omitempty"`
	LogFilePath     string `json:"log_file_path,omitempty"`
}

type ConfigFileFileds struct {
	Connections []ConnectionConfig `json:"connections"`
	AppConfig
}

type Config struct {
	ConnectionConfig
	AppConfig
}

func getCofigPath() string {
	userHome, err := os.UserHomeDir()
	if err != nil {
		logger.Error("Error getting user home directory: %s", err)
		os.Exit(1)
	}
	path := userHome + "/.rowsql"
	if err = os.MkdirAll(path, 0o755); err != nil {
		logger.Error("Error creating .rowsql directory: %s", err)
		os.Exit(1)
	}
	fileName := "config.json"
	fullPath := path + "/" + fileName
	_, err = os.OpenFile(fullPath, os.O_RDONLY, 0o644)
	if err != nil {
		if os.IsNotExist(err) {
			promptForDefaultEnv(path, fileName)
		} else {
			logger.Error("Error opening .env file: %s", err)
			os.Exit(1)
		}
	}
	return fullPath
}

func MustLoad(configPath ...string) *Config {
	var cfg Config

	var path string
	if len(configPath) > 0 && configPath[0] != "" {
		path = configPath[0]
	} else {
		path = getCofigPath()
	}
	configFileFileds, err := listConfig(path)
	if err != nil {
		logger.Errorln(err)
		os.Exit(1)
	}

	appConfig, err := askConfig(configFileFileds.Connections)
	if err != nil {
		logger.Errorln(err)
		os.Exit(1)
	}

	cfg.ConnectionConfig = *appConfig
	if cfg.Env == "" {
		cfg.Env = EnvProduction
	}
	if cfg.LogFilePath == "" {
		userHome, err := os.UserHomeDir()
		if err != nil {
			logger.Error("Error getting user home directory: %s", err)
			os.Exit(1)
		}
		cfg.LogFilePath = userHome + "/.rowsql/rowsql.log"
	}

	if cfg.MaxItemsPerPage < 1 {
		cfg.MaxItemsPerPage = 1
	}

	return &cfg
}

func listConfig(path string) (ConfigFileFileds, error) {
	var configFileFields ConfigFileFileds
	fileData, err := os.ReadFile(path)
	if err != nil {
		return configFileFields, err
	}
	if err := json.Unmarshal(fileData, &configFileFields); err != nil {
		return configFileFields, apperr.ErrorInvalidJSONFile
	}

	for _, conn := range configFileFields.Connections {
		if conn.Port == 0 || conn.DBString == "" {
			return ConfigFileFileds{}, apperr.ErrorInvalidJSONFile
		}
	}
	return configFileFields, nil
}
