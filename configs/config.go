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
	DisableAutoUpdate bool `json:"disable_auto_update,omitempty"`
	Driver            Driver

	MaxItemsPerPage int    `json:"max_items_per_page,omitempty"`
	MinItemsPerPage int    `json:"min_items_per_page,omitempty"`
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

type ConfigService interface {
	ParseConfig(path string) (ConfigFileFileds, error)
	LoadConfig(path ...string) *Config
	GetConfigPath() string
}

type ConfigServiceImpl struct {
	prompter Prompter
}

func NewConfigService() ConfigService {
	return &ConfigServiceImpl{
		prompter: &PrompterImpl{},
	}
}

func (c *ConfigServiceImpl) LoadConfig(configPath ...string) *Config {
	var cfg Config

	var path string
	if len(configPath) > 0 && configPath[0] != "" {
		path = configPath[0]
	} else {
		path = c.GetConfigPath()
	}
	configFileFileds, err := c.ParseConfig(path)
	if err != nil {
		logger.Errorln(err)
		os.Exit(1)
	}

	cfg.AppConfig = configFileFileds.AppConfig

	var appConfig *ConnectionConfig
	if len(configFileFileds.Connections) == 1 {
		appConfig = &configFileFileds.Connections[0]
	} else {
		appConfig, err = c.prompter.AskConnection(configFileFileds.Connections)
		if err != nil {
			logger.Errorln(err)
			os.Exit(1)
		}
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

	if cfg.MaxItemsPerPage == 0 {
		logger.Warning("max_items_per_page not set or set to 0, defaulting to %d", MaxItemsPerPage)
		cfg.MaxItemsPerPage = MaxItemsPerPage
	}

	if cfg.MinItemsPerPage == 0 {
		logger.Warning("min_items_per_page not set or set to 0, defaulting to %d", MinItemsPerPage)
		cfg.MinItemsPerPage = MinItemsPerPage
	}

	return &cfg
}

func (c *ConfigServiceImpl) GetConfigPath() string {
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

func (c *ConfigServiceImpl) ParseConfig(path string) (ConfigFileFileds, error) {
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
