package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	c "mcp_client/configs"
	"net/http"
)

func OllamaHandler(logger *slog.Logger, client *http.Client) (string, error) {
	req := &ollamaGenerateRequest{
		Model:  "gemma3",
		Prompt: "",
		Stream: true,
		Format: "json",
	}

	resp, err := makeOllamaRequest(req, client, c.GetConfig(logger))
	if err != nil {
		logger.Error("Failed to make Ollama request", "error", err)
		return "", err
	}

	// pResp, err := parseOllamaResponse(resp)
	// if err != nil {
	// 	logger.Error("Failed to parse Ollama response", "error", err)
	// 	return "", err
	// }

	return resp.Response, nil
}

func makeOllamaRequest(req *ollamaGenerateRequest, httpClient *http.Client, cfg *c.Config) (*ollamaGenerateResponse, error) {
	jsonReq, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/generate", cfg.LLMUrl)
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LLM returned non-200 status: %d", resp.StatusCode)
	}

	var ollamaResp ollamaGenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, err
	}

	return &ollamaResp, nil
}

// func parseOllamaResponse(resp *ollamaGenerateResponse) (string, error) {
// 	if resp == nil {
// 		return "", fmt.Errorf("empty response")
// 	}

// 	json.Unmarshal(resp.Response)

// 	return resp.Response, nil
// }

type (
	ollamaGenerateRequest struct {
		Model     string         `json:"model"`                // required
		Prompt    string         `json:"prompt,omitempty"`     // prompt to generate a response for
		Suffix    string         `json:"suffix,omitempty"`     // text after the model response
		Images    []string       `json:"images,omitempty"`     // base64-encoded images (for multimodal models)
		Format    any            `json:"format,omitempty"`     // "json" or a JSON schema
		Options   map[string]any `json:"options,omitempty"`    // additional model parameters
		System    string         `json:"system,omitempty"`     // system message
		Template  string         `json:"template,omitempty"`   // prompt template
		Stream    bool           `json:"stream,omitempty"`     // streaming response
		Raw       bool           `json:"raw,omitempty"`        // no formatting applied to prompt
		KeepAlive string         `json:"keep_alive,omitempty"` // model memory duration (e.g., "5m")
		Context   any            `json:"context,omitempty"`    // deprecated conversational context
	}

	ollamaGenerateResponse struct {
		Model              string `json:"model"`
		CreatedAt          string `json:"created_at"`
		Response           string `json:"response"`
		Done               bool   `json:"done"`
		Context            []int  `json:"context"`
		TotalDuration      int64  `json:"total_duration"`
		LoadDuration       int64  `json:"load_duration"`
		PromptEvalCount    int    `json:"prompt_eval_count"`
		PromptEvalDuration int64  `json:"prompt_eval_duration"`
		EvalCount          int    `json:"eval_count"`
		EvalDuration       int64  `json:"eval_duration"`
	}
)
