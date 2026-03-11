package context

import (
	"context"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// ────────────────────────────────────────────────
// Request ID
// ────────────────────────────────────────────────

func WithRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, reqID)
}

func GetRequestID(ctx context.Context) string {
	if v := ctx.Value(RequestIDKey); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}

// ────────────────────────────────────────────────
// Real client IP (учитывая прокси-заголовки)
// ────────────────────────────────────────────────

func WithRealIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, RealIPKey, ip)
}

func GetRealIP(ctx context.Context) string {
	if v := ctx.Value(RealIPKey); v != nil {
		if ipStr, ok := v.(string); ok {
			return ipStr
		}
	}
	return ""
}

// ────────────────────────────────────────────────
// Удобная функция для middleware — извлечь/сгенерировать всё сразу
// ────────────────────────────────────────────────

func EnrichRequestContext(r *http.Request) *http.Request {
	ctx := r.Context()

	// Request ID
	reqID := GetRequestID(ctx)
	if reqID == "" {
		reqID = generateRequestID() // ← твоя функция, можно вынести сюда или в utils
		ctx = WithRequestID(ctx, reqID)
	}

	// Real IP — берём из заголовка или RemoteAddr
	realIP := GetRealIP(ctx)
	if realIP == "" {
		realIP = extractRealIP(r)
		ctx = WithRealIP(ctx, realIP)
	}

	return r.WithContext(ctx)
}

// extractRealIP — улучшенная версия getClientIP
func extractRealIP(r *http.Request) string {
	// X-Forwarded-For может содержать цепочку: client, proxy1, proxy2
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		if len(parts) > 0 {
			// Первый IP — клиентский (самый левый)
			return strings.TrimSpace(parts[0])
		}
	}

	// X-Real-IP (часто ставит nginx)
	if real := r.Header.Get("X-Real-IP"); real != "" {
		return strings.TrimSpace(real)
	}

	// Последний fallback
	return r.RemoteAddr
}

// generateRequestID — можно улучшить позже (uuid, snowflake и т.д.)
func generateRequestID() string {
	// Простой вариант, но читаемый
	return time.Now().Format("20060102-150405.000000") + "-" + randomString(6)
	// Лучше: github.com/google/uuid → uuid.New().String()
}

// randomString — заглушка, реализуй нормально
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
