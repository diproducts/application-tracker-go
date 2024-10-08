package main

import (
	"flag"
	"github.com/diproducts/application-tracker-go/internal/config"
	"github.com/diproducts/application-tracker-go/internal/migrator"
)

func main() {
	cfg := config.MustLoad()

	var migrationsPath string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.Parse()

	if migrationsPath == "" {
		panic("migrations-path flag is required")
	}

	migrator.Up(&cfg.DB, migrationsPath)
}
