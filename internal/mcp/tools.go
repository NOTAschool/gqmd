package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerTools(s *server.MCPServer) error {
	statusTool := mcp.NewTool("status",
		mcp.WithDescription("Show index status and health information"),
	)

	s.AddTool(statusTool, statusHandler)

	return nil
}
