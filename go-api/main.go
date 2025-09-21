package main

import (
	"api/calculations"
	"api/config"
	"api/generator"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	apiLogger := setupLogger()
	slog.SetDefault(apiLogger)
	cfg := config.GetConfig(apiLogger)

	httpClient := &http.Client{}

	// Initialize modules
	var gen = generator.NewGenerator(apiLogger, httpClient, cfg)

	// Define routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Go API!"))
	})
	http.HandleFunc("/calculations/max", calculations.MaxHandler(apiLogger))
	http.HandleFunc("/calculations/sum", calculations.SumHandler(apiLogger))
	http.HandleFunc("/calculations/reverse", calculations.ReverseHandler(apiLogger))
	http.HandleFunc("/calculations/countUnique", calculations.CountUniqueHandler(apiLogger))

	// Generator related endpoints
	http.HandleFunc("/generate-module", func(w http.ResponseWriter, r *http.Request) {
		apiLogger.Info("Received request to /generate-module")
		var input generator.GeneratorInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		// Manual validation
		if input.Scope == "" {
			http.Error(w, "Schema validation failed: 'scope' is required", http.StatusBadRequest)
			return
		}
		// You can add more checks for Options if needed

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "Module generation started",
		})

		go func() {
			apiLogger.Info("Starting module generation", "scope", input.Scope)
			if err := gen.GenerateModule(input); err != nil {
				apiLogger.Error("Failed to generate module", "error", err)
			}
		}()
	})

	http.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir(cfg.OutputDir))))

	apiLogger.Info(fmt.Sprintf("Starting server on :%v", cfg.Port), "port", cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), nil)
}

func setupLogger() *slog.Logger {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	var level slog.Level
	switch os.Getenv("LOG_LEVEL") {
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

	handler := &apiLogger{Handler: slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})}
	apislog := slog.New(handler)
	apislog.Info("logger configured", "level", level, "pod", "API", "hostname", hostname)
	return apislog
}

type apiLogger struct {
	slog.Handler
}

func colorForLevel(level string) string {
	switch level {
	case "DEBUG":
		return "\033[36m" // Cyan
	case "INFO":
		return "\033[32m" // Green
	case "WARN":
		return "\033[33m" // Yellow
	case "ERROR":
		return "\033[31m" // Red
	default:
		return "\033[0m" // Reset
	}
}

func (h *apiLogger) Handle(ctx context.Context, r slog.Record) error {
	// Color codes
	timestampColor := "\033[90m" // Bright black (gray)
	podColor := "\033[35m"       // Magenta
	levelColor := colorForLevel(r.Level.String())
	msgColor := "\033[97m" // Bright white
	attrColor := "\033[2m" // Dim
	pod := "API"

	reset := "\033[0m"

	timestamp := r.Time.Format("2006-01-02 15:04")
	level := r.Level.String()
	msg := r.Message

	attrs := ""
	r.Attrs(func(a slog.Attr) bool {
		attrs += fmt.Sprintf(" %s=%v", a.Key, a.Value)
		return true
	})

	// Print each field with its color and clear separation
	left := fmt.Sprintf("%s%s%s | %s[%s]%s | %s%s%s | %s%s%s",
		timestampColor, timestamp, reset,
		podColor, pod, reset,
		levelColor, level, reset,
		msgColor, msg, reset,
	)
	right := fmt.Sprintf("%s%s%s\n", attrColor, attrs, reset)

	fmt.Printf("%-150s |", left)
	fmt.Printf("%s", right)
	return nil
}
