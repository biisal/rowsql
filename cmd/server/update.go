package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/biisal/rowsql/internal/logger"
	"github.com/creativeprojects/go-selfupdate"
	"github.com/fatih/color"
)

const (
	updateCheckTimeout = 30 * time.Second
	releaseSlug        = "biisal/rowsql"
	checksumsFileName  = "checksums.txt"
)

var errRestartRequired = errors.New("restart required after update")

func runAutoUpdate(command string, currentVersion string, disableAutoUpdate bool) error {
	if disableAutoUpdate {
		logger.Info("Auto-update disabled using disable_auto_update: true in config file")
		return nil
	}

	if currentVersion == "" || currentVersion == "dev" || currentVersion == "latest" {
		logger.Info("Development build: skipping auto-update")
		return nil
	}

	logger.Info(
		"Checking for updates...You can disable auto-update by setting disable_auto_update: true in config file",
	)

	updater, err := selfupdate.NewUpdater(selfupdate.Config{
		Validator: &selfupdate.ChecksumValidator{UniqueFilename: checksumsFileName},
	})
	if err != nil {
		return fmt.Errorf("error occurred while preparing updater: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), updateCheckTimeout)
	defer cancel()

	latest, err := updater.UpdateSelf(ctx, currentVersion, selfupdate.ParseSlug(releaseSlug))
	if err != nil {
		return fmt.Errorf("error occurred while updating binary: %w", err)
	}

	if latest.Version() == currentVersion {
		logger.Success("Already on the latest version (%s)", color.HiGreenString(currentVersion))
		return nil
	}

	logger.Success("Successfully updated to version %s", color.HiGreenString(latest.Version()))
	color.Cyan("Please run %s to restart the server", command)
	return errRestartRequired
}
