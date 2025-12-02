package main

import (
	c "api/configs"
	l "api/logger"
	"context"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	logger := l.NewLogger("MCP_CLIENT")
	cfg := c.NewConfig(logger)
	mcpClient := NewClient()
	httpClient := http.DefaultClient
	ctx := context.Background()

	http.HandleFunc("/ollama", func(w http.ResponseWriter, r *http.Request) {
		response, err := OllamaHandler(logger, httpClient)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		logger.Info("Ollama response", "response", response)

		t := &mcp.CommandTransport{Command: exec.Command("server")}
		session, err := mcpClient.Connect(ctx, t, nil)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer session.Close()

		params := &mcp.CallToolParams{
			Name: "Create file",
			Arguments: map[string]any{
				"path":    "/path/to/file.txt",
				"content": "file content",
			},
		}

		res, err := session.CallTool(ctx, params)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if res.IsError {
			logger.Error("Tool call error", "tool", params.Name)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		logger.Info(fmt.Sprintf("tool response: %v", res.Content))

	})

	http.ListenAndServe(":"+cfg.Port, nil)
}
