package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"github.com/biisal/db-gui/configs"
	"github.com/biisal/db-gui/internal/database/repo"
	"github.com/biisal/db-gui/internal/logger"
	"github.com/biisal/db-gui/internal/router"
	"github.com/biisal/db-gui/internal/router/middleware"
	"github.com/biisal/db-gui/internal/utils"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func mount(cfg *configs.Config) error {
	ctx := context.Background()
	logger.SetupSlog(slog.LevelDebug)
	driver, err := utils.DetectDriver(cfg.DBSTRING)
	if err != nil {
		return err
	}
	slog.Info("Database driver detected", "driver", driver)
	cfg.Driver = driver
	dbConn, err := sqlx.ConnectContext(ctx, driver, cfg.DBSTRING)
	if err != nil {
		return err
	}

	dbRepo := repo.New(dbConn, driver, cfg.MaxItemsPerPage)

	if err := dbRepo.Init(ctx); err != nil {
		return err
	}

	dbService := router.NewService(dbRepo, cfg.MaxItemsPerPage)
	dbHandler := router.NewHandler(dbService)

	mux, err := router.MountRouter(dbHandler)
	corsMux := middleware.CORS()(mux)
	if err != nil {
		return err
	}

	server := http.Server{
		Addr:    cfg.Server.Port,
		Handler: corsMux,
	}
	slog.Info("Running server on port", "port", cfg.Server.Port)
	log.Println("Running server on port", "port", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
