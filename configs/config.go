// Package configs contains the configuration for the application
package configs

import (
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/biisal/db-gui/internal/logger"
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
	LogFilePath     string `env:"LOG_FILE_PATH" env-default:"rowsql.log"`
}

func MustLoad() *Config {
	var cfg Config

	if err := godotenv.Load(".env"); err != nil {
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
