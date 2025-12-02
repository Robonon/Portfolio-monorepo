package configs

import (
	"fmt"
	"log/slog"
	"os"
)

type Config struct {
	LogLevel string
	Port     string
}

func NewConfig(log *slog.Logger) *Config {
	return &Config{
		LogLevel: getenvOrDefault(log, "LOG_LEVEL", "DEBUG"),
		Port:     getenvOrDefault(log, "PORT", "8082"),
	}
}

func getenvOrDefault(log *slog.Logger, key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Info(fmt.Sprintf("%s not set, defaulting to %s", key, def))
		return def
	}
	return val
}
