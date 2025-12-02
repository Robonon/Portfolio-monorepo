package configs

import (
	"fmt"
	"log/slog"
	"os"
)

type Config struct {
	LogLevel   string
	Port       string
	APIBaseURL string
}

func NewConfig(log *slog.Logger) *Config {
	return &Config{
		LogLevel:   getenvOrDefault(log, "LOG_LEVEL", "INFO"),
		Port:       getenvOrDefault(log, "PORT", "8080"),
		APIBaseURL: getenvOrDefault(log, "API_BASE_URL", "http://localhost:8081"),
	}
}

func getenvOrDefault(log *slog.Logger, key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Info(fmt.Sprintf("%v not set, defaulting to %v", key, def))
		return def
	}
	return val
}
