package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

// Name is the name of the pod you want to be displayed in the logs
func NewLogger(podName string) *slog.Logger {
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

	handler := &Logger{
		Handler: slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}),
		PodName: podName,
	}
	logger := slog.New(handler)
	logger.Info("logger configured", "level", level, "pod", podName, "hostname", hostname)
	return logger
}

type Logger struct {
	slog.Handler
	PodName string
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

func (h *Logger) Handle(ctx context.Context, r slog.Record) error {
	// Color codes
	timestampColor := "\033[90m" // Bright black (gray)
	podColor := "\033[35m"       // Magenta
	levelColor := colorForLevel(r.Level.String())
	msgColor := "\033[97m" // Bright white
	attrColor := "\033[2m" // Dim
	pod := h.PodName

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
	left := fmt.Sprintf("%s%s%s | %s[%s]%s | %s%-5s%s | %s%s%s",
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
