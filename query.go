package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
	"go.i3wm.org/i3/v4"
)

func init() {
	var queryCmd = &cobra.Command{
		Use:   "query",
		Short: "search the i3 tree for a nodeid",
		Run: func(cmd *cobra.Command, args []string) {
			tree, err := i3.GetTree()
			if err != nil {
				log.Fatal(err)
			}
			spew.Dump(tree)
		},
	}
	// queryCmd.Flags().String("instance", "", "filter by process name")

	rootCmd.AddCommand(queryCmd)
}
