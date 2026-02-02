package cli

import (
	"fmt"

	"github.com/NOTAschool/gqmd/internal/store"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan [name]",
	Short: "Scan and index a collection",
	Long:  `Scan a collection directory and index all matching documents.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := store.Open()
		if err != nil {
			return err
		}
		defer db.Close()

		if len(args) == 0 {
			// Scan all collections
			cols, err := db.ListCollections()
			if err != nil {
				return err
			}
			for _, col := range cols {
				fmt.Printf("Scanning %s...\n", col.Name)
				result, err := db.ScanCollection(col.Name)
				if err != nil {
					fmt.Printf("  Error: %v\n", err)
					continue
				}
				fmt.Printf("  Added: %d\n", result.Added)
			}
			return nil
		}

		// Scan specific collection
		name := args[0]
		fmt.Printf("Scanning %s...\n", name)
		result, err := db.ScanCollection(name)
		if err != nil {
			return err
		}
		fmt.Printf("Added: %d, Errors: %d\n", result.Added, result.Errors)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
