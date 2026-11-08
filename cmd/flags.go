package cmd

import (
	"flag"
	"fmt"

	"github.com/fatih/color"
)

type Flags struct {
	EnvPath string
	DevMode bool
}

func ParseFlags(command string) Flags {
	pathInstruction := color.CyanString(fmt.Sprintf("Path to the environment file\nExample: %s -env=./env", command))
	envPath := flag.String("env", "", pathInstruction)
	devMode := flag.Bool(
		"dev",
		false,
		"Run in development mode: uses the first connection in the config and skips the session token",
	)
	flag.Parse()
	return Flags{EnvPath: *envPath, DevMode: *devMode}
}
