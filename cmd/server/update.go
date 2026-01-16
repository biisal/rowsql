package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/biisal/rowsql/internal/logger"
	"github.com/creativeprojects/go-selfupdate"
	"github.com/fatih/color"
)

func runAutoUpdate(cmd string, currentVersion string) error {
	if currentVersion == "dev" {
		logger.Info("Development build: skipping auto-update")
		return nil
	}

	colordVersion := color.HiGreenString(currentVersion)
	ctx := context.Background()
	latest, found, err := selfupdate.DetectLatest(ctx, selfupdate.ParseSlug("biisal/rowsql"))
	if err != nil {
		return fmt.Errorf("error occurred while detecting version: %w", err)
	}
	if !found {
		return fmt.Errorf("latest version for %s/%s could not be found from github repository", runtime.GOOS, runtime.GOARCH)
	}
	if latest.LessOrEqual(colordVersion) {
		logger.Info("Current version (%s) is the latest", colordVersion)
		return nil
	}
	exe, err := selfupdate.ExecutablePath()
	if err != nil {
		return errors.New("could not locate executable path")
	}
	if err := selfupdate.UpdateTo(context.Background(), latest.AssetURL, "rowsql", exe); err != nil {
		return fmt.Errorf("error occurred while updating binary: %w", err)
	}

	logger.Success("Successfully updated to version %s", color.HiGreenString(latest.Version()))

	color.Cyan("Please run %s to restart the server", cmd)
	os.Exit(1)
	return nil
}
