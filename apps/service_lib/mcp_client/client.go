package main

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func NewClient() *mcp.Client {
	opts := &mcp.ClientOptions{
		ToolListChangedHandler:      toolListChangedHandler,
		CreateMessageHandler:        createMessageHandler,
		ElicitationHandler:          elicitationHandler,
		PromptListChangedHandler:    nil,
		ResourceListChangedHandler:  nil,
		ResourceUpdatedHandler:      nil,
		LoggingMessageHandler:       nil,
		ProgressNotificationHandler: nil,
		KeepAlive:                   time.Minute * 5,
	}
	return mcp.NewClient(&mcp.Implementation{Name: "Generator MCP Client", Title: "Gen MCP Client"}, opts)
}

func toolListChangedHandler(ctx context.Context, tlcr *mcp.ToolListChangedRequest) {

}

func createMessageHandler(ctx context.Context, cmr *mcp.CreateMessageRequest) (*mcp.CreateMessageResult, error) {
	return &mcp.CreateMessageResult{Content: &mcp.TextContent{}}, nil
}

func elicitationHandler(ctx context.Context, er *mcp.ElicitRequest) (*mcp.ElicitResult, error) {
	return &mcp.ElicitResult{Content: map[string]any{}}, nil
}
