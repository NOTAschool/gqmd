package cli

import (
	"github.com/NOTAschool/gqmd/internal/mcp"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start MCP server (stdio transport)",
	Long:  `Start the Model Context Protocol server for AI agent integration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return mcp.StartServer()
	},
}
