package server

import (
	"log/slog"
	"net/http"
	"time"
	// если нужно
)

func New(handler http.Handler, logger *slog.Logger) *http.Server {
	return &http.Server{
		Addr:              ":8080", // потом из конфига
		Handler:           handler,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
}
