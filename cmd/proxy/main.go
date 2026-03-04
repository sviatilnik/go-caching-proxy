package main

import (
	"log"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/sviatilnik/go-caching-proxy/internals/config"
	"github.com/sviatilnik/go-caching-proxy/internals/server"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		slog.Error("No .env file found, relying on OS environment variables")
	}

	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	server := server.NewServer(conf)

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
