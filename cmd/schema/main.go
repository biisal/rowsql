package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/biisal/rowsql/cmd"
	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/logger"
	"github.com/biisal/rowsql/internal/router"
)

func main() {
	command := os.Args[0]

	envPath := cmd.ParseFlags(command)
	cfg := configs.MustLoad(envPath)
	mux := http.NewServeMux()

	api := router.NewAPI(mux, cfg)
	router.RegisterRoutes(api, router.DBHandler{})

	b, err := json.MarshalIndent(api.OpenAPI(), "", "\t")
	if err != nil {
		panic(err)
	}
	logger.Info("Creating openapi.json in ./frontend/ ")
	if err = os.WriteFile("frontend/openapi.json", b, 0o644); err != nil {
		panic(err)
	}
	logger.Info("Done")
}
