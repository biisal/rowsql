package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
)

func perseFlags() (envPath string) {
	command := os.Args[0]
	pathInstruction := color.CyanString(fmt.Sprintf("Path to the environment file\nExample: %s -env=./env", command))
	envPath = *flag.String("env", "", pathInstruction)
	flag.Parse()
	return
}
