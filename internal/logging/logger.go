package logging

import (
	"log/slog"
	"os"
)

// New создаёт структурированный логгер
func New(level slog.Level) *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Всегда сохраняем msg
			if a.Key == slog.MessageKey {
				return a
			}

			// Сохраняем ip и browser
			if a.Key == "client_ip" || a.Key == "user_agent" {
				return a
			}

			// Все остальные поля убираем
			return slog.Attr{}
		},
	})
	return slog.New(handler)
}

// Или вариант с JSON (удобнее для продакшена / парсинга)
func NewJSON(level slog.Level) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Всегда оставляем msg
			if a.Key == slog.MessageKey {
				return a
			}

			// Оставляем ip и user_agent
			if a.Key == "client_ip" || a.Key == "user_agent" {
				return a
			}

			// Все остальные атрибуты удаляем
			return slog.Attr{}
		},
	})

	return slog.New(handler)
}
