package config

import (
	"time"
)

type Config struct {
	ListenAddr      string        `yaml:"listen_addr" env:"LISTEN_ADDR"`
	Upstream        string        `yaml:"upstream" env:"UPSTREAM_URL"`
	LogLevel        string        `yaml:"log_level" env:"LOG_LEVEL"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	// Добавь потом: TLS, Fingerprint rules, MetricsPort и т.д.
}
