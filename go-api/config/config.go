package config

import (
	"log/slog"
	"os"
)

// Add config variables that will be fetched from env
type Config struct {
	LogLevel  string
	LogFormat string
	Port      string

	// External resources
	LLMUrl      string
	DatabaseURL string
}

func GetConfig() *Config {
	var c = Config{}
	c.LogLevel = os.Getenv("LOG_LEVEL")
	if c.LogLevel == "" {
		slog.Info("LOG_LEVEL not set, defaulting to info")
		c.LogLevel = "info"
	}
	c.Port = os.Getenv("PORT")
	if c.Port == "" {
		slog.Info("PORT not set, defaulting to 8080")
		c.Port = "8080"
	}

	c.LLMUrl = os.Getenv("LLM_URL")
	if c.LLMUrl == "" {
		slog.Info("LLM_URL not set, defaulting to ollama-service.default.svc.cluster.local:11434")
		c.LLMUrl = "http://ollama-service.default.svc.cluster.local:11434"
	}

	c.LogFormat = os.Getenv("LOG_FORMAT")
	if c.LogFormat == "" {
		slog.Info("LOG_FORMAT not set, defaulting to PLAIN")
		c.LogFormat = "PLAIN"
	}

	return &c
}
