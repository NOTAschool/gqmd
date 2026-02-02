package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerTools(s *server.MCPServer) error {
	// status tool
	statusTool := mcp.NewTool("status",
		mcp.WithDescription("Show index status and health information"),
	)
	s.AddTool(statusTool, statusHandler)

	// search tool
	searchTool := mcp.NewTool("search",
		mcp.WithDescription("Search documents using FTS5 full-text search"),
		mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
		mcp.WithNumber("limit", mcp.Description("Max results (default 10)")),
	)
	s.AddTool(searchTool, searchHandler)

	// get tool
	getTool := mcp.NewTool("get",
		mcp.WithDescription("Get a document by collection and path"),
		mcp.WithString("collection", mcp.Required(), mcp.Description("Collection name")),
		mcp.WithString("path", mcp.Required(), mcp.Description("Document path")),
	)
	s.AddTool(getTool, getHandler)

	// multi_get tool
	multiGetTool := mcp.NewTool("multi_get",
		mcp.WithDescription("Get multiple documents by paths"),
		mcp.WithArray("paths", mcp.Required(), mcp.Description("Array of collection/path strings")),
		mcp.WithNumber("max_bytes", mcp.Description("Max total bytes (default 10KB)")),
	)
	s.AddTool(multiGetTool, multiGetHandler)

	return nil
}
