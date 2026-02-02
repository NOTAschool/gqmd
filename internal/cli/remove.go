package cli

import (
	"fmt"

	"github.com/NOTAschool/gqmd/internal/store"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a collection",
	Long:  `Remove a collection and all its indexed documents.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		db, err := store.Open()
		if err != nil {
			return err
		}
		defer db.Close()

		if err := db.RemoveCollection(name); err != nil {
			return err
		}

		fmt.Printf("Removed collection %q\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
