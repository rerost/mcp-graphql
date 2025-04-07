package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rerost/mcp-graphql/internal"
)

// Flags
var endpoint = flag.String("endpoint", "http://localhost:8080", "Required. GraphQL server URL")
var schemaFile = flag.String("schema", "", "Optional. Path to GraphQL schema file")
var defaultHeader = flag.String("headers", "{}", `Optional. Default headers for GraphQL requests. e.g. {"Authorization": "Bearer ..."}`)

func init() {
	flag.Parse()
}

// Tools
var RunQueryTool = internal.Tool{
	Tool: mcp.NewTool("run-query",
		mcp.WithDescription("Run a GraphQL query"),
		mcp.WithString("query", mcp.Required(), mcp.Description("GraphQL query to run")),
		mcp.WithString("variables", mcp.Description(`variables. JSON e.g. {"id": "123"}`)),
		mcp.WithString("headers", mcp.Description(`headers. JSON e.g. {"Content-Type": "application/json"}"`)),
	),
	Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := request.Params.Arguments["query"].(string)
		variables := request.Params.Arguments["variables"].(string)
		headers := request.Params.Arguments["headers"].(string)

		// Merge default
		// Overwrite headers
		var mergedHeader string
		{
			defaultHeadersMap := make(map[string]string)
			if err := json.Unmarshal([]byte(*defaultHeader), &defaultHeadersMap); err != nil {
				return nil, errors.WithStack(err)
			}
			headersMap := make(map[string]string)
			if err := json.Unmarshal([]byte(headers), &headersMap); err != nil {
				return nil, errors.WithStack(err)
			}

			for k, v := range headersMap {
				defaultHeadersMap[k] = v
			}

			mergedHeaderB, err := json.Marshal(defaultHeadersMap)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			mergedHeader = string(mergedHeaderB)
		}

		response, err := internal.CallGraphQL(ctx, *endpoint, query, variables, mergedHeader)
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(response), nil
	},
}

func newServer() *internal.Server {
	s := &internal.Server{
		Server: server.NewMCPServer(
			"MCP GraphQL",
			"0.1.0",
		),
	}

	s.AddTool(&RunQueryTool)

	return s
}

func main() {
	s := newServer()

	if err := s.ServeStdio(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
