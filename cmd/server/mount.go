package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

func printLogo(version string) {
	gitLink := color.HiGreenString("https://github.com/biisal/rowsql")
	version = color.HiYellowString("You are using version %s", version)

	logo := fmt.Sprintf(`
                          ❤️ Thanks for using
█▀█ █▀█ █░█░█ █▀ █▀█ █░░  ⭐ Star on GitHub: %s
█▀▄ █▄█ ▀▄▀▄▀ ▄█ ▀▀█ █▄▄  %s
`, gitLink, version)
	fmt.Println(logo)
}

const (
	databaseConnectTimeout = 15 * time.Second
	readHeaderTimeout      = 10 * time.Second
	idleTimeout            = 2 * time.Minute
	shutdownTimeout        = 5 * time.Second
)

func shutdownOnSignal(ctx context.Context, server *http.Server) {
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Errorln("Failed to shut down server:", err)
	}
}

func mount(cfg *configs.Config) error {
	signalCtx, stopSignals := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignals()
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
	logger.Info("Connecting to %s database...", driver)
	connectCtx, cancelConnect := context.WithTimeout(ctx, databaseConnectTimeout)
	defer cancelConnect()
	dbConn, err := sqlx.ConnectContext(connectCtx, string(driver), cfg.DBString)
	if err != nil {
		logger.Errorln("Failed to connect to database:", err)
		return err
	}
	defer func() {
		err := dbConn.Close()
		if err != nil {
			logger.Errorln("Failed to close database connection:", err)
		}
	}()

	logger.Info("Connected to %s database successfully", driver)

	queryBuilder, err := queries.NewBuilder(cfg.Driver, cfg.MaxItemsPerPage)
	if err != nil {
		return err
	}

	dbRepo := repo.New(dbConn, cfg.Driver, queryBuilder, cfg.MaxItemsPerPage)

	if err = dbRepo.Init(ctx); err != nil {
		logger.Errorln("Failed to initialize database repository:", err)
		return err
	}

	dbService := service.NewService(dbRepo, cfg.MaxItemsPerPage)
	dbHandler := router.NewHandler(dbService, cfg.MaxItemsPerPage)

	mux, err := router.MountRouter(dbHandler, cfg)
	if err != nil {
		logger.Errorln("Failed to mount router:", err)
		return err
	}

	sessionToken := ""
	if cfg.Env != configs.EnvDevelopment {
		sessionToken, err = router.NewSessionToken()
		if err != nil {
			return err
		}
	}
	guardedMux := router.RequireLocalAccess(sessionToken, cfg.Host)(mux)

	server := http.Server{
		Handler:           guardedMux,
		ReadHeaderTimeout: readHeaderTimeout,
		IdleTimeout:       idleTimeout,
	}
	ln, addr, err := listen(cfg.Host, cfg.Port)
	if err != nil {
		logger.Errorln("Failed to listen on port:", err)
		return err
	}
	logger.Success("Listening on %s", addr)
	if cfg.Env != configs.EnvDevelopment {
		appURL := fmt.Sprintf("http://%s/?token=%s", addr, sessionToken)
		fmt.Printf("Open %s if your browser does not start automatically\n", appURL)
		go openBrowser(appURL)
	}
	go shutdownOnSignal(signalCtx, &server)
	if err := server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Errorln("Failed to start server:", err)
		return err
	}
	logger.Info("Server stopped")
	return nil
}
