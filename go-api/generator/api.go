package generator

import (
	"api/config"
	"log/slog"
	"net/http"
)

type Generator interface {
	GenerateModule(input GeneratorInput) error
}

type GeneratorInput struct {
	Scope string `json:"scope"`
}

func NewGenerator(logger *slog.Logger, httpClient *http.Client, cfg *config.Config) Generator {
	return &generatorImpl{
		logger:     logger,
		httpClient: httpClient,
		config:     cfg,
	}
}
