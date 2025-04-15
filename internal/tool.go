package internal

import (
	"context"
	"encoding/json"

	"github.com/cockroachdb/errors"
	mcp "github.com/ktr0731/go-mcp"
	"github.com/ktr0731/go-mcp/protocol"
)

var endpoint *string

func SetEndpoint(url *string) {
	endpoint = url
}

var defaultHeaders map[string]string

func SetDefaultHeaders(headers map[string]string) {
	defaultHeaders = headers
}

type ToolHandler struct{}

func (h *ToolHandler) Handle(ctx context.Context, method string, req protocol.CallToolRequestParams) (any, error) {
	switch req.Name {
	case "run-query":
		return h.handleRunQuery(ctx, req)
	default:
		return nil, errors.Newf("unknown tool: %s", req.Name)
	}
}

func (h *ToolHandler) handleRunQuery(ctx context.Context, req protocol.CallToolRequestParams) (*mcp.CallToolResult, error) {
	argsJSON, err := json.Marshal(req.Arguments)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	
	var args struct {
		Query     string `json:"query"`
		Variables string `json:"variables,omitempty"`
		Headers   string `json:"headers,omitempty"`
	}
	
	if err := json.Unmarshal(argsJSON, &args); err != nil {
		return nil, errors.WithStack(err)
	}
	
	if args.Query == "" {
		return nil, errors.New("query must be a non-empty string")
	}
	
	var variables *string
	if args.Variables != "" {
		variables = &args.Variables
	}
	
	var headers *string
	if args.Headers != "" {
		headers = &args.Headers
	}

	headersMap := make(map[string]string)
	for k, v := range defaultHeaders {
		headersMap[k] = v
	}
	
	if headers != nil {
		reqHeaders := make(map[string]string)
		if err := json.Unmarshal([]byte(*headers), &reqHeaders); err != nil {
			return nil, errors.WithStack(err)
		}

		for k, v := range reqHeaders {
			headersMap[k] = v
		}
	}

	response, err := CallGraphQL(ctx, *endpoint, args.Query, variables, headersMap)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.CallToolContent{
			mcp.TextContent{Text: response},
		},
	}, nil
}

func RegisterTools(handler *mcp.Handler) {
	tool := protocol.Tool{
		Name:        "run-query",
		Description: "Run a GraphQL query",
		InputSchema: protocol.ToolInputSchema{
			Type: "object",
			Properties: map[string]protocol.ToolInputSchemaProperty{
				"query": {
					Type:        "string",
					Description: "GraphQL query to run",
				},
				"variables": {
					Type:        "string",
					Description: "Variables in JSON format",
				},
				"headers": {
					Type:        "string",
					Description: "Headers in JSON format",
				},
			},
			Required: []string{"query"},
		},
	}

	handler.Tools = append(handler.Tools, tool)
	handler.ToolHandler = &ToolHandler{}
}
