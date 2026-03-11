package proxy

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Proxy — наш reverse proxy с возможностью расширения
type Proxy struct {
	ReverseProxy *httputil.ReverseProxy
	Target       *url.URL
	Logger       *slog.Logger
}

// Метрики Prometheus

func init() {
	prometheus.MustRegister(RequestCount)
	prometheus.MustRegister(RequestDuration)
}

// New создаёт новый reverse proxy, направленный на указанный target
func New(targetURL string, logger *slog.Logger) (*Proxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	rp := httputil.NewSingleHostReverseProxy(target)

	// Настройка Director
	originalDirector := rp.Director
	rp.Director = func(req *http.Request) {
		originalDirector(req)

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host

		// Генерация request_id
		reqID := generateRequestID()
		req.Header.Set("X-Request-ID", reqID)

		// X-Forwarded-For для реального IP
		clientIP := req.RemoteAddr
		if forwarded := req.Header.Get("X-Forwarded-For"); forwarded != "" {
			clientIP = forwarded + ", " + req.RemoteAddr
		}
		req.Header.Set("X-Forwarded-For", clientIP)

		// Сохраняем время начала запроса
		req.Header.Set("X-Start-Time", time.Now().Format(time.RFC3339Nano))

		// Логирование начала запроса
		logger.Info("request started",
			slog.String("req_id", reqID),
			slog.String("client_ip", clientIP),
			slog.String("user_agent", req.UserAgent()),
			slog.String("method", req.Method),
			slog.String("path", req.URL.Path),
		)
	}

	// Логирование завершения запроса и сбор метрик
	rp.ModifyResponse = func(resp *http.Response) error {
		req := resp.Request
		reqID := req.Header.Get("X-Request-ID")
		clientIP := req.Header.Get("X-Forwarded-For")

		startTime, err := time.Parse(time.RFC3339Nano, req.Header.Get("X-Start-Time"))
		duration := time.Since(startTime)
		if err != nil {
			duration = 0
		}

		statusCode := resp.StatusCode
		method := req.Method
		path := req.URL.Path

		// Логирование
		logger.Info("request completed",
			slog.String("req_id", reqID),
			slog.String("client_ip", clientIP),
			slog.String("user_agent", req.UserAgent()),
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", statusCode),
			slog.String("duration", duration.String()),
		)

		// Сбор метрик
		RequestCount.WithLabelValues(method, path, http.StatusText(statusCode)).Inc()
		RequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())

		return nil
	}

	// Логирование ошибок прокси
	rp.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
		reqID := req.Header.Get("X-Request-ID")
		clientIP := req.Header.Get("X-Forwarded-For")

		logger.Error("proxy error",
			slog.String("req_id", reqID),
			slog.String("client_ip", clientIP),
			slog.String("method", req.Method),
			slog.String("path", req.URL.Path),
			slog.String("error", err.Error()),
		)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}

	return &Proxy{
		ReverseProxy: rp,
		Target:       target,
		Logger:       logger,
	}, nil
}

// ServeHTTP реализует интерфейс http.Handler
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.ReverseProxy.ServeHTTP(w, r)
}

// Генерация безопасного request_id
func generateRequestID() string {
	timestamp := time.Now().Format("20060102-150405")
	randomPart := randomStringSecure(6)
	return timestamp + "-" + randomPart
}

// Безопасная случайная строка
func randomStringSecure(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return strings.ToLower(hex.EncodeToString(b)[:n])
}
