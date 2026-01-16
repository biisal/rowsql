package main

import (
	"context"
	"fmt"
	"os"

	"github.com/biisal/rowsql/internal/logger"
	"github.com/creativeprojects/go-selfupdate"
	"github.com/fatih/color"
)

func runAutoUpdate(cmd string, currentVersion string) error {
	if currentVersion == "dev" || currentVersion == "" {
		logger.Info("Development build: skipping auto-update")
		return nil
	}

	ctx := context.Background()

	repo := selfupdate.ParseSlug("biisal/rowsql")

	latest, err := selfupdate.UpdateSelf(ctx, currentVersion, repo)
	if err != nil {
		return fmt.Errorf("error occurred while updating binary: %w", err)
	}

	if latest.Version() == currentVersion {
		return nil
	}

	logger.Success("Successfully updated to version %s", color.HiGreenString(latest.Version()))
	color.Cyan("Please run %s to restart the server", cmd)

	os.Exit(0)
	return nil
}
