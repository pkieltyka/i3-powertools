package main

import (
	"flag"
	"fmt"
	"log"

	"go.i3wm.org/i3"
)

var workspaceArg = flag.String("workspace", "", "workspace command, args: [next, prev]")

func main() {
	flag.Parse()

	var workspaceOp WorkspaceOp

	switch *workspaceArg {
	case "next":
		workspaceOp = WS_NEXT
	case "prev":
		workspaceOp = WS_PREV
	default:
		log.Fatal("invalid workspace argument")
	}

	switchWorkspace(workspaceOp)
}

type WorkspaceOp uint

const (
	WS_NEXT WorkspaceOp = iota // right
	WS_PREV                    // left
)

// switchWorkspace will move to the next or previous workspace, while staying
// on the same monitor, aka output device. UX+1
func switchWorkspace(op WorkspaceOp) {
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
	case WS_NEXT:
		if i == len(spaces)-1 {
			gotoSpace = spaces[0].Name
		} else {
			gotoSpace = spaces[i+1].Name
		}
	case WS_PREV:
		if i == 0 {
			gotoSpace = spaces[len(spaces)-1].Name
		} else {
			gotoSpace = spaces[i-1].Name
		}
	}

	i3.RunCommand(fmt.Sprintf("workspace %s", gotoSpace))
}
