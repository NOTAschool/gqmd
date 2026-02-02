package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/NOTAschool/gqmd/internal/store"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <name> <path>",
	Short: "Add a collection",
	Long:  `Add a new collection to index documents from a directory.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		path := args[1]

		// Resolve to absolute path
		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("invalid path: %w", err)
		}

		// Check if path exists
		info, err := os.Stat(absPath)
		if err != nil {
			return fmt.Errorf("path not found: %w", err)
		}
		if !info.IsDir() {
			return fmt.Errorf("path must be a directory")
		}

		pattern, _ := cmd.Flags().GetString("pattern")

		db, err := store.Open()
		if err != nil {
			return err
		}
		defer db.Close()

		if err := db.AddCollection(name, absPath, pattern); err != nil {
			return fmt.Errorf("failed to add collection: %w", err)
		}

		fmt.Printf("Added collection %q -> %s\n", name, absPath)
		return nil
	},
}

func init() {
	addCmd.Flags().StringP("pattern", "p", "**/*.md", "Glob pattern for files")
	rootCmd.AddCommand(addCmd)
}
