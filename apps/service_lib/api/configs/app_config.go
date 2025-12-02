package configs

import (
	"fmt"
	"log/slog"
	"os"
)

type Config struct {
	LogLevel          string
	Port              string
	TplServiceBaseUrl string
}

func NewConfig(log *slog.Logger) *Config {
	return &Config{
		LogLevel:          getenvOrDefault(log, "LOG_LEVEL", "INFO"),
		Port:              getenvOrDefault(log, "PORT", "8081"),
		TplServiceBaseUrl: getenvOrDefault(log, "TPL_SERVICE_BASE_URL", "http://localhost:8082"),
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
