package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rerost/mcp-graphql/internal"
)

// Flags
var endpoint = flag.String("endpoint", "http://localhost:8080", "Required. GraphQL server URL")
var schemaFile = flag.String("schema", "", "Optional. Path to GraphQL schema file")

type MapFlag map[string]string

// String はフラグの値を文字列として返します
func (m MapFlag) String() string {
	return fmt.Sprint(map[string]string(m))
}

func (m MapFlag) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	if len(parts) != 2 {
		return errors.New("引数は key=value の形式で指定してください")
	}
	key, val := parts[0], parts[1]
	m[key] = val
	return nil
}

func (m MapFlag) Get() map[string]string {
	return map[string]string(m)
}

var defaultHeader MapFlag

func init() {
	defaultHeader = make(MapFlag)
	flag.Var(&defaultHeader, "headers", "Optional. Default headers for GraphQL requests. e.g. Authorization=Bearer ...")
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
		var variables *string
		if request.Params.Arguments["variables"] != nil {
			if v, ok := request.Params.Arguments["variables"].(string); ok {
				variables = &v
			}
		}
		var headers *string
		if request.Params.Arguments["headers"] != nil {
			if v, ok := request.Params.Arguments["headers"].(string); ok {
				headers = &v
			}
		}

		// Merge default
		// Overwrite headers
		headersMap := make(map[string]string)
		{
			for k, v := range defaultHeader.Get() {
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
		}

		response, err := internal.CallGraphQL(ctx, *endpoint, query, variables, headersMap)
		if err != nil {
			return nil, errors.WithStack(err)
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

func Debug(args ...interface{}) {
	fmt.Fprintln(os.Stderr, "DEBUG", args)
}
