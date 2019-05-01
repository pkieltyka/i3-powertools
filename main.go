package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"go.i3wm.org/i3"
)

const VERSION = "v0.2"

var rootCmd = &cobra.Command{
	Use:   "i3-powertools",
	Short: "i3-powertools",
	Args:  cobra.MinimumNArgs(1),
}

func init() {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("i3-powertools", VERSION)
		},
	}

	// rootCmd.AddCommand(&cobra.Command{
	// 	Use: "bus",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		recv := i3.Subscribe(i3.WindowEventType)
	// 		go func() {
	// 			for recv.Next() {
	// 				ev := recv.Event().(*i3.WindowEvent)
	// 				spew.Dump(ev)
	// 			}
	// 		}()
	// 		time.Sleep(100 * time.Second)
	// 	},
	// })

	rootCmd.AddCommand(versionCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func getFocusedNode() *i3.Node {
	tree, err := i3.GetTree()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	n := tree.Root.FindFocused(func(n *i3.Node) bool {
		return n.Focused == true
	})

	return n
}
