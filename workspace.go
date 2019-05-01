package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"go.i3wm.org/i3"
)

func init() {
	// workspace
	var workspaceCmd = &cobra.Command{
		Use:   "workspace [command]",
		Short: "workspace switcher",
		Args:  cobra.MinimumNArgs(1),
	}

	var wsSwitchCmd = &cobra.Command{
		Use:  "switch [command]",
		Args: cobra.MinimumNArgs(1),
	}
	workspaceCmd.AddCommand(wsSwitchCmd)

	wsSwitchCmd.AddCommand(&cobra.Command{
		Use: "next",
		Run: func(cmd *cobra.Command, args []string) {
			workspaceSwitch("next")
		},
	})
	wsSwitchCmd.AddCommand(&cobra.Command{
		Use: "prev",
		Run: func(cmd *cobra.Command, args []string) {
			workspaceSwitch("prev")
		},
	})

	// TODO: add "last"
	// workspace open, so we go back to the last workspace we visited
	// we might need PowerTools instance and track the events/stuff for this..

	rootCmd.AddCommand(workspaceCmd)
}

// workspaceSwitch will move to the next or previous workspace, while staying
// on the same monitor, aka output device. UX+1
func workspaceSwitch(op string) {
	workspaces, err := i3.GetWorkspaces()
	if err != nil {
		log.Fatal(err)
	}

	spacesByOutput := map[string][]i3.Workspace{}

	var focusedOutput string
	var focusedOutputIndex int

	for _, ws := range workspaces {
		spacesByOutput[ws.Output] = append(spacesByOutput[ws.Output], ws)
		if ws.Focused {
			focusedOutput = ws.Output
			focusedOutputIndex = len(spacesByOutput[ws.Output]) - 1
		}
	}

	var gotoSpace string

	spaces := spacesByOutput[focusedOutput]
	i := focusedOutputIndex

	switch op {
	case "next":
		if i == len(spaces)-1 {
			gotoSpace = spaces[0].Name
		} else {
			gotoSpace = spaces[i+1].Name
		}
	case "prev":
		if i == 0 {
			gotoSpace = spaces[len(spaces)-1].Name
		} else {
			gotoSpace = spaces[i-1].Name
		}
	}

	i3.RunCommand(fmt.Sprintf("workspace %s", gotoSpace))
}
