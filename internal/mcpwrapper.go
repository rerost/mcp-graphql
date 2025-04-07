package internal

import (
	"github.com/cockroachdb/errors"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Server struct {
	Server *server.MCPServer
}

func (h *Server) ServeStdio() error {
	return errors.WithStack(server.ServeStdio(h.Server))
}

func (h *Server) AddTool(t *Tool) {
	h.Server.AddTool(t.Tool, t.Handler)
}

type Tool struct {
	Tool    mcp.Tool
	Handler server.ToolHandlerFunc
}
