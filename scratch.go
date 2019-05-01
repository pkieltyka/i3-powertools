package main

import "github.com/spf13/cobra"

func init() {
	scratchPad := &AppCmd{
		Command: &cobra.Command{
			Use:   "scratch [command]",
			Short: "scratch pad thing",
			Args:  cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {

			},
		},
	}

	rootCmd.AddCommand(scratchPad.Command)
}
