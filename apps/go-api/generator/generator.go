package generator

import (
	"api/config"
	"log/slog"
	"net/http"
)

type Generator interface {
	GenerateModule(input GeneratorInput) error
	ZipOutput() error
}

type GeneratorInput struct {
	Scope    string `json:"scope" jsonschema:"the scope for generating the code"` 
	SkipEval bool   `json:"skipEval"`
}

func NewGenerator(logger *slog.Logger, httpClient *http.Client, cfg *config.Config) Generator {
	return &generatorImpl{
		logger:     logger,
		httpClient: httpClient,
		config:     cfg,
	}
}
