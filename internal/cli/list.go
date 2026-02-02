package cli

import (
	"fmt"

	"github.com/NOTAschool/gqmd/internal/store"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List collections",
	Long:  `List all registered collections.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := store.Open()
		if err != nil {
			return err
		}
		defer db.Close()

		cols, err := db.ListCollections()
		if err != nil {
			return err
		}

		if len(cols) == 0 {
			fmt.Println("No collections")
			return nil
		}

		for _, c := range cols {
			fmt.Printf("%s -> %s (%s)\n", c.Name, c.Path, c.Pattern)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
