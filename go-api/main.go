package main

import (
	"api/calculations"
	"api/config"
	"api/llm"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type API struct {
	logger *slog.Logger
}

func main() {
	cfg := config.GetConfig()
	api := &API{
		logger: setupLogger(cfg),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Go API!"))
	})
	http.HandleFunc("/calculations/max", calculations.MaxHandler(api.logger))
	http.HandleFunc("/calculations/sum", calculations.SumHandler(api.logger))
	http.HandleFunc("/calculations/reverse", calculations.ReverseHandler(api.logger))
	http.HandleFunc("/calculations/countUnique", calculations.CountUniqueHandler(api.logger))

	// LLM related endpoints
	http.HandleFunc("/generate-module", llm.GenerateModuleHandler(api.logger))
	http.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir("/output"))))

	api.logger.Info(fmt.Sprintf("Starting server on :%v", cfg.Port), "port", cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), nil)
}

func setupLogger(cfg *config.Config) *slog.Logger {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	var level slog.Level
	switch cfg.LogLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var handler slog.Handler
	if cfg.LogFormat == "JSON" {
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	} else {
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	}

	apislog := slog.New(handler).With("pod", "API").With("hostname", hostname)
	apislog.Info("logger configured", "level", level, "pod", "API", "hostname", hostname, "format", cfg.LogFormat)
	return apislog
}
