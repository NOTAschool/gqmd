package mcp

import (
	"context"
	"fmt"

	"github.com/NOTAschool/gqmd/internal/store"
	"github.com/mark3labs/mcp-go/mcp"
)

func statusHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	db, err := store.Open()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to open database: %v", err)), nil
	}
	defer db.Close()

	status, err := db.GetStatus()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get status: %v", err)), nil
	}

	text := fmt.Sprintf(`gqmd Index Status:
  Database: %s
  Total documents: %d
  Collections: %d
  Has vector index: %v`,
		status.DBPath,
		status.TotalDocs,
		status.Collections,
		status.HasVectorIndex,
	)

	return mcp.NewToolResultText(text), nil
}
