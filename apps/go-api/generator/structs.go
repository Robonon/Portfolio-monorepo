package generator

import "encoding/json"

type (
	ollamaRequest struct {
		Model  string          `json:"model"`
		Prompt string          `json:"prompt"`
		Stream bool            `json:"stream"`
		Format json.RawMessage `json:"format,omitempty"`
	}

	ollamaResponse struct {
		Model     string `json:"model"`
		CreatedAt string `json:"created_at"`
		Response  string `json:"response"`
		Done      bool   `json:"done"`
	}
)
