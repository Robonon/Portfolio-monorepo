package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Format string `json:"format,omitempty"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func StreamPromptToFile(llmURL, model, prompt, filename, schema string) (string, error) {
	reqBody, err := json.Marshal(OllamaRequest{Model: model, Prompt: prompt, Format: schema, Stream: true})
	if err != nil {
		return "failed", fmt.Errorf("failed to marshal LLM request: %w", err)
	}

	resp, err := http.Post(llmURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "failed", fmt.Errorf("failed to send request to LLM: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "failed", fmt.Errorf("LLM returned status %d: %s", resp.StatusCode, string(body))
	}

	file, err := os.Create("/output/" + filename)
	if err != nil {
		return "failed", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var ollamaResp OllamaResponse
		if err := json.Unmarshal(scanner.Bytes(), &ollamaResp); err != nil {
			return "failed", fmt.Errorf("failed to decode LLM response: %w", err)
		}
		if _, err := file.WriteString(ollamaResp.Response); err != nil {
			return "failed", fmt.Errorf("failed to write to file: %w", err)
		}
		if ollamaResp.Done {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return "failed", fmt.Errorf("error reading LLM response: %w", err)
	}

	return "success", nil
}
