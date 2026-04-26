package cmd

import (
	"flag"
	"fmt"

	"github.com/fatih/color"
)

func ParseFlags(command string) (envPath string) {
	pathInstruction := color.CyanString(fmt.Sprintf("Path to the environment file\nExample: %s -env=./env", command))
	envPath = *flag.String("env", "", pathInstruction)
	flag.Parse()
	return
}
