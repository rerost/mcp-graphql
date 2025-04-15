package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ktr0731/go-mcp/codegen"
)

func main() {
	outDir := "."
	if err := os.MkdirAll(outDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	f, err := os.Create(filepath.Join(outDir, "mcp.gen.go"))
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()

	def := &codegen.ServerDefinition{
		Capabilities: codegen.ServerCapabilities{
			Tools:   &codegen.ToolCapability{},
			Logging: &codegen.LoggingCapability{},
		},
		Implementation: codegen.Implementation{
			Name:    "MCP GraphQL",
			Version: "0.1.0",
		},
		Tools: []codegen.Tool{
			{
				Name:        "run-query",
				Description: "Run a GraphQL query",
				InputSchema: struct {
					Query     string `json:"query"`
					Variables string `json:"variables"`
					Headers   string `json:"headers"`
				}{},
			},
		},
	}

	if err := codegen.Generate(f, def, "internal"); err != nil {
		log.Fatalf("Failed to generate code: %v", err)
	}
}
