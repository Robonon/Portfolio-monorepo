package main

import (
	"api/calculations"
	"api/config"
	"api/generator"
	"api/logger"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	aLog := logger.SetupLogger() // Initialize the global logger

	cfg := config.GetConfig(aLog)

	httpClient := &http.Client{}

	// Initialize modules
	var gen = generator.NewGenerator(aLog, httpClient, cfg)

	// Define routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Go API!"))
	})
	http.HandleFunc("/calculations/max", calculations.MaxHandler(aLog))
	http.HandleFunc("/calculations/sum", calculations.SumHandler(aLog))
	http.HandleFunc("/calculations/reverse", calculations.ReverseHandler(aLog))
	http.HandleFunc("/calculations/countUnique", calculations.CountUniqueHandler(aLog))

	// Generator related endpoints
	http.HandleFunc("/generate-module", func(w http.ResponseWriter, r *http.Request) {
		aLog.Info("Received request to /generate-module")
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
			aLog.Info("Starting module generation", "scope", input.Scope)
			if err := gen.GenerateModule(input); err != nil {
				aLog.Error("Failed to generate module", "error", err)
			}
		}()
	})

	// broken
	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		aLog.Info("Received request to /download")

		err := gen.ZipOutput()
		if err != nil {
			http.Error(w, "Failed to download output:", http.StatusInternalServerError)
			return
		}
		aLog.Info("Serving zip file", "path", cfg.OutputDir+".zip")
		http.ServeFile(w, r, cfg.OutputDir+".zip")
	})

	aLog.Info(fmt.Sprintf("Starting server on :%v", cfg.Port), "port", cfg.Port)
	http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), nil)
}
