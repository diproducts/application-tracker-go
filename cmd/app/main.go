package main

import (
	"github.com/diproducts/application-tracker-go/internal/app"
	"github.com/diproducts/application-tracker-go/internal/config"
)

func main() {
	cfg := config.MustLoad()

	app.Run(cfg)
}
