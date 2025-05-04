package main

import (
	"auth-service/config"
	"auth-service/internal/app"
)

func main() {
	cfg := config.MustLoad()

	app.Run(cfg)
}
