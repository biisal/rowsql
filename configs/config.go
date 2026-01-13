// Package configs contains the configuration for the application
package configs

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/biisal/rowsql/internal/logger"
	"github.com/fatih/color"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	DriverPostgres        = "pgx"
	DriverMySQL           = "mysql"
	DriverSQLite          = "sqlite"
	EnvDevelopment string = "development"
	EnvProduction  string = "production"
)

type ServerConfig struct {
	Host string `env:"HOST"`
	Port string `env:"PORT" env-required:"true"`
}

type Config struct {
	DBString        string `env:"DBSTRING" env-required:"true"`
	Server          ServerConfig
	Driver          string
	MaxItemsPerPage int    `env:"MAX_ITEMS_PER_PAGE" env-default:"10"`
	Env             string `env:"ENV" env-default:"production"`
	LogFilePath     string `env:"LOG_FILE_PATH" env-default:"~/.rowsql/rowsql.log"`
}

func promptForDefaultEnv(path string) {
	color.Cyan("No .env found in %s\nDo you want to create one with default values? (y/n): ", path)
	var choice string
	if _, err := fmt.Scan(&choice); err != nil {
		logger.Errorln(err)
		os.Exit(0)
	}
	if strings.ToLower(choice) == "y" {
		file, err := os.Create(path)
		if err != nil {
			logger.Error("Error creating .env file: %s", err)
			os.Exit(1)
		}
		if _, err = file.WriteString("DBSTRING=test.db\nPORT=8000"); err != nil {
			logger.Errorln(err)
			os.Exit(0)
		}
		defer func() {
			if err := file.Close(); err != nil {
				logger.Errorln(err)
			}
		}()
	} else {
		logger.Error("No .env file found")
		os.Exit(1)
	}
}

func getEnvPath() string {
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
	fullPath := path + "/.env"
	_, err = os.OpenFile(fullPath, os.O_RDONLY, 0o644)
	if err != nil {
		if os.IsNotExist(err) {
			promptForDefaultEnv(fullPath)
		} else {
			logger.Error("Error opening .env file: %s", err)
			os.Exit(1)
		}
	}
	return fullPath
}

func MustLoad(envPath ...string) *Config {
	var cfg Config

	var path string
	if len(envPath) > 0 && envPath[0] != "" {
		path = envPath[0]
	} else {
		path = getEnvPath()
	}

	if err := godotenv.Load(path); err != nil {
		log.Fatal(err)
	}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal(err)
	}
	if cfg.DBString == "" {
		t := reflect.TypeFor[Config]()
		field, _ := t.FieldByName("DBSTRING")
		tagVal := field.Tag.Get("env")
		logger.Error("%s not found in .env", tagVal)
		os.Exit(1)
	}
	if !strings.HasPrefix(cfg.Server.Port, ":") {
		cfg.Server.Port = ":" + cfg.Server.Port
	}
	if cfg.Env != string(EnvDevelopment) && cfg.Env != string(EnvProduction) {
		logger.Error("%s env can't be set! Make sure it's '%s' or '%s', Default '%s'",
			cfg.Env, EnvDevelopment, EnvProduction, EnvProduction)
		os.Exit(1)
	}
	return &cfg
}
