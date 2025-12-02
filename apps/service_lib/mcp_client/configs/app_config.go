package config

import (
	"log/slog"
	"os"
)

// Add config variables that will be fetched from env
type Config struct {
	// General Config
	Port string

	// Generator Config
	OutputDir string
	LLMUrl    string
}

func GetConfig(logger *slog.Logger) *Config {
	var c = Config{}

	c.Port = os.Getenv("PORT")
	if c.Port == "" {
		logger.Info("PORT not set, defaulting to 8080")
		c.Port = "8082"
	}

	c.LLMUrl = os.Getenv("LLM_URL")
	if c.LLMUrl == "" {
		logger.Info("LLM_URL not set, defaulting to ollama-service.default.svc.cluster.local:11434")
		c.LLMUrl = "http://ollama-service.default.svc.cluster.local:11434"
	}

	c.OutputDir = os.Getenv("OUTPUT_DIR")
	if c.OutputDir == "" {
		logger.Info("OUTPUT_DIR not set, defaulting to ./output")
		c.OutputDir = "./output"
	}

	return &c
}
