package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/database/queries"
	"github.com/biisal/rowsql/internal/database/repo"
	"github.com/biisal/rowsql/internal/logger"
	"github.com/biisal/rowsql/internal/router"
	"github.com/biisal/rowsql/internal/service"
	"github.com/biisal/rowsql/internal/utils"
	"github.com/fatih/color"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func printLogo() {
	gitLink := color.HiGreenString("https://github.com/biisal/rowsql")
	logo := fmt.Sprintf(`
█▀█ █▀█ █░█░█ █▀ █▀█ █░░  ❤️ Thanks for using
█▀▄ █▄█ ▀▄▀▄▀ ▄█ ▀▀█ █▄▄  ⭐ Star on GitHub: %s
	`, gitLink)
	fmt.Println(logo)
}

func mount(cfg *configs.Config) error {
	ctx := context.Background()

	logFilePath, err := utils.ReplaceTildeWithHomeDir(cfg.LogFilePath)
	if err != nil {
		return err
	}
	cfg.LogFilePath = logFilePath

	if err = logger.SetupFile(cfg.LogFilePath); err != nil {
		return err
	}
	defer logger.Close()
	if cfg.Env == configs.EnvDevelopment {
		logger.SetLogLevel(logger.LevelDebug)
	}

	logger.Info("All logs will be written in %s", cfg.LogFilePath)

	driver, err := utils.DetectDriver(&cfg.DBString)
	if err != nil {
		return err
	}
	logger.Info("Database driver detected: %s", driver)
	cfg.Driver = driver
	dbConn, err := sqlx.ConnectContext(ctx, driver, cfg.DBString)
	if err != nil {
		logger.Errorln("Failed to connect to database:", err)
		return err
	}

	queiryBuilder := queries.NewBuilder(cfg.Driver)

	dbRepo := repo.New(dbConn, driver, queiryBuilder, cfg.MaxItemsPerPage)

	if err = dbRepo.Init(ctx); err != nil {
		logger.Errorln("Failed to initialize database repository:", err)
		return err
	}

	dbService := service.NewService(dbRepo, cfg.MaxItemsPerPage)
	dbHandler := router.NewHandler(dbService, cfg.MaxItemsPerPage)

	mux, err := router.MountRouter(dbHandler)
	corsMux := router.CORS()(mux)
	if err != nil {
		logger.Errorln("Failed to mount router:", err)
		return err
	}

	server := http.Server{
		Addr:    cfg.Server.Port,
		Handler: corsMux,
	}
	logger.Success("Running server on port %s", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		logger.Errorln("Failed to start server:", err)
		return err
	}
	return nil
}
