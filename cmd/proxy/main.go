package main

import (
	"log/slog"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/sviatilnik/go-caching-proxy/internals/config"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		slog.Error("No .env file found, relying on OS environment variables")
	}

	c := config.NewConfig()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Hello, World!"))
	})

	if err := http.ListenAndServe(c.Port, mux); err != nil {
		slog.Error(err.Error())
	}
}
