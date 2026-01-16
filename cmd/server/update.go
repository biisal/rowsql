package main

import (
	"context"
	"fmt"
	"os"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/logger"
	"github.com/creativeprojects/go-selfupdate"
	"github.com/fatih/color"
)

func runAutoUpdate(cmd string, currentVersion string, cfg *configs.AutoUpdateConfig) error {
	if cfg.DisableAutoUpdate {
		logger.Info("Auto-update disabled using DISABLE_AUTO_UPDATE=true")
		return nil
	}

	if currentVersion == "" || currentVersion == "dev" || currentVersion == "latest" {
		logger.Info("Development build: skipping auto-update")
		return nil
	}

	logger.Info("Checking for updates...You can disable auto-update by setting DISABLE_AUTO_UPDATE=true")

	repo := selfupdate.ParseSlug("biisal/rowsql")

	latest, err := selfupdate.UpdateSelf(context.Background(), currentVersion, repo)
	if err != nil {
		return fmt.Errorf("error occurred while updating binary: %w", err)
	}

	if latest.Version() == currentVersion {
		logger.Success("Already on the latest version (%s)", color.HiGreenString(currentVersion))
		return nil
	}

	logger.Success("Successfully updated to version %s", color.HiGreenString(latest.Version()))
	color.Cyan("Please run %s to restart the server", cmd)
	os.Exit(0)
	return nil
}
