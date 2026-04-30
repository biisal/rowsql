package main

import (
	"log"
	"os"

	"github.com/biisal/rowsql/cmd"
	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/logger"
)

var version = "dev"

func main() {
	printLogo(version)
	command := os.Args[0]

	envPath := cmd.ParseFlags(command)
	configService := configs.NewConfigService()

	cfg := configService.LoadConfig(envPath)

	if cfg == nil {
		logger.Error("Failed to load config")
		return
	}

	if err := runAutoUpdate(command, version, cfg.DisableAutoUpdate); err != nil {
		logger.ErrorWriteOnlyFile("Error while updating: %s", err)
	}

	if err := mount(cfg); err != nil {
		log.Fatal("Failed to mount app:", err)
		return
	}
}
