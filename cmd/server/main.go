package main

import (
	"log"
	"os"

	"github.com/biisal/rowsql/configs"
)

// go build -ldflags="-X main.version=$(date +%d-%m-%Y)"
var (
	version = "dev"
)

func main() {
	command := os.Args[0]

	if err := runAutoUpdate(command, version); err != nil {
		log.Fatal("Failed to run auto update:", err)
		return
	}
	printLogo(version)
	envPath := perseFlags(command)
	cfg := configs.MustLoad(envPath)
	if err := mount(cfg); err != nil {
		log.Fatal("Failed to mount app:", err)
		return
	}
}
