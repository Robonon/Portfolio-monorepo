package generator

type (
	Job struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Result string `json:"result,omitempty"`
		Error  string `json:"error,omitempty"`
		Stage  string `json:"stage,omitempty"`
	}

	ollamaRequest struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
		Stream bool   `json:"stream"`
		Format any    `json:"format,omitempty"`
	}

	ollamaResponse struct {
		Model     string `json:"model"`
		CreatedAt string `json:"created_at"`
		Response  string `json:"response"`
		Done      bool   `json:"done"`
	}
)
