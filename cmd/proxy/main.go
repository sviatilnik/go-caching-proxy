package main

import (
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/sviatilnik/go-caching-proxy/internals/config"
	"github.com/sviatilnik/go-caching-proxy/internals/server"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		slog.Error("No .env file found, relying on OS environment variables")
	}

	server := server.NewServer(config.NewConfig())
	server.Start()
}
