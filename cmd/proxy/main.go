package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"reverse-proxy/internal/logging"
	"reverse-proxy/internal/proxy"
	"reverse-proxy/internal/server"
)

func main() {
	// создаём логгер через твой модуль logging
	logger := logging.New(slog.LevelInfo)

	const upstream = "https://httpbin.org"

	// создаём proxy, теперь передаём logger как второй аргумент
	p, err := proxy.New(upstream, logger)
	if err != nil {
		logger.Error("не удалось создать прокси", "error", err)
		os.Exit(1)
	}

	// middleware для логирования запросов
	handler := logging.RequestLogger(p, logger)

	logger.Info("reverse proxy запущен",
		slog.String("listen", ":8080"),
		slog.String("upstream", upstream),
	)

	// ← создаём сервер
	srv := server.New(handler, logger)

	// ← запускаем сервер в фоне
	go func() {
		logger.Info("сервер стартует", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("сервер упал", "error", err)
			os.Exit(1)
		}
	}()

	// ← ждём сигнал на остановку
	server.Graceful(context.Background(), srv, logger)

	logger.Info("приложение полностью остановлено")
}
