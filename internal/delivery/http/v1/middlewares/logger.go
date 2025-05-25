package middleware

import (
	"log/slog"
	"net/http"
	"service/internal/pkg/logs"
	"time"
)

var opLogger = "middleware.Logger"

func Logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logs.Info(
			r.Context(),
			"request to server",
			opLogger,
			slog.Group(
				"request",
				"method", r.Method,
				"path", r.URL.Path,
			),
		)
		
		start := time.Now()

		handler.ServeHTTP(w, r)

		logs.Info(
			r.Context(),
			"request completed",
			opLogger,
			slog.String("request processing time", time.Since(start).String()),
		)
	})
}