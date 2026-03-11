package logging

import (
	"log/slog"
	"net/http"
	"time"

	"reverse-proxy/internal/context" // ← подставь своё имя модуля
)

// RequestLogger — middleware, который логирует входящий запрос и ответ
func RequestLogger(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ← Вот здесь обогащаем контекст (request id + real ip)
		r = context.EnrichRequestContext(r)

		start := time.Now()
		ctx := r.Context()

		reqID := context.GetRequestID(ctx)
		clientIP := context.GetRealIP(ctx)

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		logger.Info("request started",
			slog.String("req_id", reqID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("query", r.URL.RawQuery),
			slog.String("client_ip", clientIP),
			slog.String("user_agent", r.UserAgent()),
			slog.String("host", r.Host),
		)

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		logger.Info("request completed",
			slog.String("req_id", reqID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rw.statusCode),
			slog.Duration("duration", duration),
			slog.String("client_ip", clientIP),
		)
	})
}

// responseWriter — для перехвата статуса ответа
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	// Если статус ещё не установлен — по умолчанию 200
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	return rw.ResponseWriter.Write(b)
}

// Простая функция получения IP (можно улучшить)
func getClientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		// берём первый IP (самый левый)
		return forwarded
	}
	return r.RemoteAddr
}

// Заглушка — замени на реальную генерацию (uuid.NewString() например)
func generateRequestID() string {
	return time.Now().Format("20060102-150405.000000")
}
