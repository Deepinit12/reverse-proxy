package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

func Load() (*Config, error) {
	cfg := &Config{
		ListenAddr:      ":8080",
		Upstream:        "https://httpbin.org",
		LogLevel:        "info",
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    30 * time.Second,
		ShutdownTimeout: 30 * time.Second,
	}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	// Можно добавить чтение yaml если нужно
	return cfg, nil
}
