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

func searchHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query := req.GetString("query", "")
	if query == "" {
		return mcp.NewToolResultError("query is required"), nil
	}

	limit := req.GetInt("limit", 10)

	db, err := store.Open()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to open database: %v", err)), nil
	}
	defer db.Close()

	results, err := db.Search(query, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("search failed: %v", err)), nil
	}

	if len(results) == 0 {
		return mcp.NewToolResultText("No results found"), nil
	}

	var text string
	for i, r := range results {
		text += fmt.Sprintf("%d. %s/%s\n   Title: %s\n   %s\n\n",
			i+1, r.Collection, r.Path, r.Title, r.Snippet)
	}

	return mcp.NewToolResultText(text), nil
}

func getHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	collection := req.GetString("collection", "")
	path := req.GetString("path", "")

	if collection == "" || path == "" {
		return mcp.NewToolResultError("collection and path are required"), nil
	}

	db, err := store.Open()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to open database: %v", err)), nil
	}
	defer db.Close()

	doc, content, err := db.Get(collection, path)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("document not found: %v", err)), nil
	}

	text := fmt.Sprintf("# %s\n\nPath: %s/%s\nModified: %s\n\n---\n\n%s",
		doc.Title, doc.Collection, doc.Path, doc.ModifiedAt, content)

	return mcp.NewToolResultText(text), nil
}

func multiGetHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	paths := req.GetStringSlice("paths", nil)
	if len(paths) == 0 {
		return mcp.NewToolResultError("paths array is required"), nil
	}

	maxBytes := req.GetInt("max_bytes", 10*1024)

	db, err := store.Open()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to open database: %v", err)), nil
	}
	defer db.Close()

	results, err := db.MultiGet(paths, maxBytes)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("multi_get failed: %v", err)), nil
	}

	if len(results) == 0 {
		return mcp.NewToolResultText("No documents found"), nil
	}

	var text string
	for _, r := range results {
		text += fmt.Sprintf("## %s\n\nPath: %s/%s\n\n%s\n\n---\n\n",
			r.Document.Title, r.Document.Collection, r.Document.Path, r.Content)
	}

	return mcp.NewToolResultText(text), nil
}
