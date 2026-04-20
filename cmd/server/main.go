package main

import (
	"log"
	"os"

	"github.com/biisal/rowsql/configs"
	"github.com/biisal/rowsql/internal/logger"
)

var (
	version = "dev"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8000
// @BasePath
func main() {
	command := os.Args[0]

	envPath := perseFlags(command)
	printLogo(version)
	cfg := configs.MustLoad(envPath)

	if err := runAutoUpdate(command, version, &cfg.Update); err != nil {
		logger.ErrorWriteOnlyFile("Error while updating: %s", err)
	}

	if err := mount(cfg); err != nil {
		log.Fatal("Failed to mount app:", err)
		return
	}
}
