package main

import (
	"context"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	impl := &mcp.Implementation{
		Name:    "MCP Server",
		Version: "v1.0.0",
	}

	// Options for server
	opts := &mcp.ServerOptions{}

	// Create a server.
	server := mcp.NewServer(impl, opts)

	// Add tools
	mcp.AddTool(server, &mcp.Tool{Name: "Create file", Description: "Create a file with a given content at given file path"}, createFileTool)
	mcp.AddTool(server, &mcp.Tool{Name: "Read file", Description: "Reads file at given file path and returns the content"}, readFileTool)

	// Run the server over stdin/stdout, until the client disconnects.
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}

type Args struct {
	Content []byte `json:"content" jsonschema:"the content to write to the file"`
	Path    string `json:"path" jsonschema:"the path to the file"`
}

func readFileTool(ctx context.Context, req *mcp.CallToolRequest, args Args) (*mcp.CallToolResult, any, error) {

	file, err := os.ReadFile(args.Path)
	if err != nil {
		return nil, nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(file)},
		},
	}, nil, nil
}
func createFileTool(ctx context.Context, req *mcp.CallToolRequest, args Args) (*mcp.CallToolResult, any, error) {

	os.WriteFile(args.Path, args.Content, 0777)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: args.Path},
		},
	}, nil, nil
}
