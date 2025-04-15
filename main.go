package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/cockroachdb/errors"
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
	
	internal.SetEndpoint(endpoint)
	internal.SetDefaultHeaders(defaultHeader.Get())
}

func RunQueryHandler(ctx context.Context, request map[string]interface{}) (interface{}, error) {
	query := request["query"].(string)
	
	var variables *string
	if v, ok := request["variables"].(string); ok && v != "" {
		variables = &v
	}
	
	var headers *string
	if v, ok := request["headers"].(string); ok && v != "" {
		headers = &v
	}

	// Merge default headers
	headersMap := make(map[string]string)
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

	response, err := internal.CallGraphQL(ctx, *endpoint, query, variables, headersMap)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return response, nil
}

func main() {
	s := internal.NewServer("MCP GraphQL", "0.1.0")

	s.AddTool(&internal.Tool{
		Name:        "run-query",
		Description: "Run a GraphQL query",
		Handler:     RunQueryHandler,
	})

	if err := s.ServeStdio(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func Debug(args ...interface{}) {
	fmt.Fprintln(os.Stderr, "DEBUG", args)
}
