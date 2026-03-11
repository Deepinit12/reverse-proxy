package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"reverse-proxy/internal/config" // ← добавим позже
	"reverse-proxy/internal/logging"
	"reverse-proxy/internal/proxy"
	"reverse-proxy/internal/server"
)

type App struct {
	Logger *slog.Logger
	Config *config.Config // ← пока можно без, потом добавим
	Proxy  *proxy.Proxy
	Server *http.Server
	// Metrics  *metrics.Registry  // если добавишь prometheus
}

func New(logger *slog.Logger) (*App, error) {
	// upstream пока хардкод, потом из config
	const upstream = "https://httpbin.org"

	p, err := proxy.New(upstream)
	if err != nil {
		return nil, err
	}

	return &App{
		Logger: logger,
		Proxy:  p,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	handler := logging.RequestLogger(a.Proxy, a.Logger)

	// Создаём сервер
	a.Server = server.New(handler, a.Logger)

	// Запуск в фоне
	go func() {
		a.Logger.Info("сервер запускается", "addr", a.Server.Addr)
		if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Logger.Error("сервер упал", "error", err)
			os.Exit(1)
		}
	}()

	// Ожидание сигнала
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	sig := <-quit

	a.Logger.Warn("сигнал завершения", "signal", sig)

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return a.Server.Shutdown(shutdownCtx)
}
