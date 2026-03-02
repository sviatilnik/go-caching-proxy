package middleware

import (
	"log/slog"
	"net/http"
)

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request received: ", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
