package main

import (
	"log"
	"os"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/logger"
)

var (
	version = "dev"
)

func main() {
	command := os.Args[0]

	envPath := perseFlags(command)
	printLogo(version)
	cfg := configs.MustLoad(envPath)

	if err := runAutoUpdate(command, version, &cfg.Update); err != nil {
		logger.ErrorWriteOnlyFile("Error while updating: %s", err)
	}

	if err := mount(cfg); err != nil {
		log.Fatal("Failed to mount app:", err)
		return
	}
}
