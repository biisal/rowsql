package main

import (
	"log"

	"github.com/biisal/db-gui/configs"
)

func main() {
	cfg := configs.MustLoad()
	if err := mount(cfg); err != nil {
		log.Fatal(err)
	}
}
