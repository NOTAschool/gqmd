package cli

import (
	"fmt"

	"github.com/NOTAschool/gqmd/internal/store"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show index status",
	Long:  `Show the status of the gqmd index.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := store.Open()
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		defer db.Close()

		status, err := db.GetStatus()
		if err != nil {
			return fmt.Errorf("failed to get status: %w", err)
		}

		fmt.Println("gqmd Index Status:")
		fmt.Printf("  Database: %s\n", status.DBPath)
		fmt.Printf("  Total documents: %d\n", status.TotalDocs)
		fmt.Printf("  Collections: %d\n", status.Collections)
		fmt.Printf("  Has vector index: %v\n", status.HasVectorIndex)

		return nil
	},
}
