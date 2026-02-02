package mcp

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"
)

func StartServer() error {
	s := server.NewMCPServer(
		"gqmd",
		"0.1.0",
		server.WithToolCapabilities(true),
	)

	if err := registerTools(s); err != nil {
		return fmt.Errorf("failed to register tools: %w", err)
	}

	return server.ServeStdio(s)
}
