package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/biisal/rowsql/configs"
	"github.com/fatih/color"
)

func main() {
	command := os.Args[0]
	pathInstruction := color.CyanString(fmt.Sprintf("Path to the environment file\nExample: %s -env=./env", command))
	envPath := flag.String("env", "", pathInstruction)
	flag.Parse()

	cfg := configs.MustLoad(*envPath)
	if err := mount(cfg); err != nil {
		log.Fatal("Failed to mount app:", err)
		return
	}
}
