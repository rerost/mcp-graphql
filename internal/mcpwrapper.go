package internal

import (
	"context"

	"github.com/cockroachdb/errors"
	mcp "github.com/ktr0731/go-mcp"
	"golang.org/x/exp/jsonrpc2"
)

type Server struct {
	handler *mcp.Handler
}

func NewServer(name string, version string) *Server {
	handler := &mcp.Handler{
		Implementation: mcp.Implementation{
			Name:    name,
			Version: version,
		},
		Capabilities: mcp.ServerCapabilities{
			Tools:   &mcp.ToolCapability{},
			Logging: &mcp.LoggingCapability{},
		},
	}
	
	RegisterTools(handler)
	
	return &Server{
		handler: handler,
	}
}

func (s *Server) ServeStdio() error {
	ctx, listener, binder := mcp.NewStdioTransport(context.Background(), s.handler, nil)
	srv, err := jsonrpc2.Serve(ctx, listener, binder)
	if err != nil {
		return errors.WithStack(err)
	}
	srv.Wait()
	return nil
}

func (s *Server) AddTool(t *Tool) {
}

type Tool struct {
	Name        string
	Description string
	Handler     func(ctx context.Context, request map[string]interface{}) (interface{}, error)
}
