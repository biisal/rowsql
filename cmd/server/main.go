package main

import (
	"errors"
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

	flags := cmd.ParseFlags(command)
	configService := configs.NewConfigService(&configs.PrompterImpl{})

	loadConfig := configService.LoadConfig
	if flags.DevMode {
		loadConfig = configService.LoadDevConfig
	}
	cfg, err := loadConfig(flags.EnvPath)
	if err != nil {
		logger.Error("Failed to load config: %v", err)
		os.Exit(1)
	}

	if err := runAutoUpdate(command, version, cfg.DisableAutoUpdate); err != nil {
		if errors.Is(err, errRestartRequired) {
			return
		}
		logger.ErrorWriteOnlyFile("Error while updating: %s", err)
	}

	if err := mount(&cfg); err != nil {
		log.Fatal("Failed to mount app:", err)
		return
	}
}
