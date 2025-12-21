package configs

import (
	"log"
	"reflect"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	DRIVER_POSTGRES = "pgx"
	DRIVER_MYSQL    = "mysql"
	DRIVER_SQLITE   = "sqlite3"
)

type ServerConfig struct {
	Host string `env:"HOST"`
	Port string `env:"PORT" env-required:"true"`
}

type Config struct {
	DBSTRING        string `env:"DBSTRING" env-required:"true"`
	Server          ServerConfig
	Driver          string
	MaxItemsPerPage int `env:"MAX_ITEMS_PER_PAGE" env-default:"10"`
}

func MustLoad() *Config {
	var cfg Config

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal(err)
	}
	if cfg.DBSTRING == "" {
		t := reflect.TypeFor[Config]()
		field, _ := t.FieldByName("DBSTRING")
		tagVal := field.Tag.Get("env")
		log.Fatalf("%s not found in .env", tagVal)
	}
	if !strings.HasPrefix(cfg.Server.Port, ":") {
		cfg.Server.Port = ":" + cfg.Server.Port
	}
	return &cfg
}
