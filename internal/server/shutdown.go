package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Graceful(ctx context.Context, srv *http.Server, logger *slog.Logger) {
	// Канал для сигналов ОС
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// Ждём сигнал
	sig := <-quit
	logger.Warn("Получен сигнал завершения", "signal", sig)

	// Контекст с таймаутом на graceful shutdown (например 30 сек)
	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Ошибка graceful shutdown", "error", err)
	} else {
		logger.Info("Сервер gracefully остановлен")
	}
}
