package main

import (
	"log"

	"github.com/biisal/rowsql/configs"
)

// go build -ldflags="-X main.version=$(date +%d-%m-%Y)"
var (
	version = "14-01-2026"
)

func main() {
	log.Printf("rowsql version %s", version)
	printLogo()
	envPath := perseFlags()
	cfg := configs.MustLoad(envPath)
	if err := mount(cfg); err != nil {
		log.Fatal("Failed to mount app:", err)
		return
	}
}
