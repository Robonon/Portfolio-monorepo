package main

import (
	"api/calculations"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type API struct {
	logger *slog.Logger
}

type Config struct {
	LogLevel string
	Port     string
}

func getConfig() (*Config, error) {
	var c = Config{}
	c.LogLevel = os.Getenv("LOG_LEVEL")
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	c.Port = os.Getenv("PORT")
	if c.Port == "" {
		c.Port = "8080"
	}
	return &c, nil
}

func main() {
	cfg, _ := getConfig()
	api := &API{
		logger: setupLogger(),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Go API!"))
	})
	http.HandleFunc("/calculations/max", calculations.MaxHandler(api.logger))
	http.HandleFunc("/calculations/sum", calculations.SumHandler(api.logger))
	http.HandleFunc("/calculations/reverse", calculations.ReverseHandler(api.logger))
	http.HandleFunc("/calculations/countUnique", calculations.CountUniqueHandler(api.logger))

	api.logger.Info(fmt.Sprintf("Starting server on :%v", cfg.Port), "port", cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), nil)
}

func setupLogger() *slog.Logger {
	jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
	apislog := slog.New(jsonHandler)
	apislog.Info("logger configured")
	return apislog
}
