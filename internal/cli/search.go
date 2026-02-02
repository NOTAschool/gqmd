package cli

import (
	"fmt"

	"github.com/NOTAschool/gqmd/internal/store"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search documents",
	Long:  `Search indexed documents using FTS5 full-text search.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		limit, _ := cmd.Flags().GetInt("limit")

		db, err := store.Open()
		if err != nil {
			return err
		}
		defer db.Close()

		results, err := db.Search(query, limit)
		if err != nil {
			return err
		}

		if len(results) == 0 {
			fmt.Println("No results found")
			return nil
		}

		for i, r := range results {
			fmt.Printf("%d. %s/%s\n", i+1, r.Collection, r.Path)
			fmt.Printf("   %s\n\n", r.Title)
		}
		return nil
	},
}

func init() {
	searchCmd.Flags().IntP("limit", "n", 10, "Max results")
	rootCmd.AddCommand(searchCmd)
}
