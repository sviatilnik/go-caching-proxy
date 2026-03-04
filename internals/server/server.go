package server

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sviatilnik/go-caching-proxy/internals/cache"
	"github.com/sviatilnik/go-caching-proxy/internals/config"
	"github.com/sviatilnik/go-caching-proxy/internals/middleware"
	"github.com/sviatilnik/go-caching-proxy/internals/proxy"
)

type Server struct {
	conf *config.Config
}

func NewServer(c *config.Config) *Server {
	return &Server{
		conf: c,
	}
}

func (server *Server) Start() error {
	prx, err := proxy.NewProxy(server.conf.Pattern, server.conf.Target, cache.NewCache(cache.NewInMemoryStore(), server.conf.TTL))
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	mux := http.NewServeMux()

	httpServer := &http.Server{
		Addr:    server.conf.Addr,
		Handler: middleware.Log(mux),
	}

	mux.HandleFunc("/", prx.ServeHTTP)

	go func() {
		slog.Info("Server starting...", "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	slog.Info("Shutting down server ...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("Failed to shutdown server", "err", err.Error())
		return err
	}

	slog.Info("Server stopped...")

	return nil
}
