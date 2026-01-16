package main

import (
	"log"

	"github.com/biisal/rowsql/configs"
)

// go build -ldflags="-X main.version=$(date +%d-%m-%Y)"
var (
	version = "dev"
)

func main() {
	printLogo(version)
	envPath := perseFlags()
	cfg := configs.MustLoad(envPath)
	if err := mount(cfg); err != nil {
		log.Fatal("Failed to mount app:", err)
		return
	}
}
