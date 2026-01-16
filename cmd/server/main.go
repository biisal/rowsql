package main

import (
	"log"
	"os"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/logger"
)

// go build -ldflags="-X main.version=$(date +%d-%m-%Y)"
var (
	version = "dev"
)

func main() {
	command := os.Args[0]

	envPath := perseFlags(command)
	printLogo(version)
	cfg := configs.MustLoad(envPath)

	if err := runAutoUpdate(command, version, &cfg.Update); err != nil {
		logger.Error("Failed to check for updates: %s", err)
		logger.Info("Continuing with current version...")
	}

	if err := mount(cfg); err != nil {
		log.Fatal("Failed to mount app:", err)
		return
	}
}
