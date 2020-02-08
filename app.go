package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"go.i3wm.org/i3/v4"
)

func init() {
	appCmd := &AppCmd{
		Command: &cobra.Command{
			Use:   "app [command]",
			Short: "app launcher/switcher",
			Args:  cobra.MinimumNArgs(1),
		},
	}

	openCmd := &cobra.Command{
		Use: "open",
		Run: appCmd.Open,
	}

	openCmd.Flags().String("focus.class", "", "focus window by class name")
	openCmd.Flags().String("focus.instance", "", "focus window by instance name")
	openCmd.Flags().String("focus.title", "", "focus window by title")
	openCmd.Flags().String("focus.mark", "", "focus window by mark")

	openCmd.Flags().String("new.exec", "", "exec command to open new instance")
	openCmd.Flags().String("new.workspace", "", "workspace to open new instance in")
	openCmd.Flags().String("new.mark", "", "set mark on new window")
	openCmd.Flags().Bool("new.floating", false, "toggle new window to start as floating")

	appCmd.AddCommand(openCmd)

	rootCmd.AddCommand(appCmd.Command)
}

type AppCmd struct {
	*cobra.Command
}

// TODO: next, lets track history, so we can go back to last window
// of a match..
// and lets isolate marked windows as their own thing. ie. chrome

// TODO: search for windows within workspace, then adjacent..

func (c *AppCmd) Open(cmd *cobra.Command, args []string) {
	var err error

	focusClassFlag, _ := cmd.Flags().GetString("focus.class")
	focusInstanceFlag, _ := cmd.Flags().GetString("focus.instance")
	focusTitleFlag, _ := cmd.Flags().GetString("focus.title")
	focusMarkFlag, _ := cmd.Flags().GetString("focus.mark")
	newMarkFlag, _ := cmd.Flags().GetString("new.mark")
	newExecFlag, _ := cmd.Flags().GetString("new.exec")
	newWorkspaceFlag, _ := cmd.Flags().GetString("new.workspace")
	newFloatingFlag, _ := cmd.Flags().GetBool("new.floating")

	if focusClassFlag == "" && focusInstanceFlag == "" && focusTitleFlag == "" && focusMarkFlag == "" && newExecFlag == "" {
		cmd.Help()
		return
	}

	// Focus app if open
	focused := false
	if focusMarkFlag != "" {
		focused = c.openByMark(focusMarkFlag)
	} else if focusClassFlag != "" || focusInstanceFlag != "" || focusTitleFlag != "" {
		focused = c.openByFilters(focusClassFlag, focusInstanceFlag, focusTitleFlag)
	}

	// Open new window if couldn't find existing app
	if !focused && newExecFlag != "" {
		// switch to workspace
		if newWorkspaceFlag != "" {
			_, err = i3.RunCommand(fmt.Sprintf(`workspace "%s"`, newWorkspaceFlag))
			if err != nil {
				log.Fatal(err)
				return
			}
		}

		beforeFocusedNode := getFocusedNode()

		// Exec
		_, err = i3.RunCommand(fmt.Sprintf(`exec %s`, newExecFlag))
		if err != nil {
			log.Fatal(err)
			return
		}

		// can we give name on exec? would be nice
		_ = newFloatingFlag

		// Mark window
		if newMarkFlag != "" {
			var setMark = false
			var num = 15
			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				defer wg.Done()

				for i := 0; i < num; i++ {
					n := getFocusedNode()
					if n.ID != beforeFocusedNode.ID {
						if focusInstanceFlag == "" || n.WindowProperties.Instance == focusInstanceFlag {
							setMark = true
							return
						}
					}
					time.Sleep(100 * time.Millisecond)
				}

			}()
			wg.Wait()

			if setMark {
				_, err = i3.RunCommand(fmt.Sprintf("mark %s", newMarkFlag))
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func (c *AppCmd) openByMark(mark string) bool {
	_, err := i3.RunCommand(fmt.Sprintf(`[con_mark="%s"] focus`, mark))
	return err == nil
}

func (c *AppCmd) openByFilters(class, instance, title string) bool {
	tree, err := i3.GetTree()
	if err != nil {
		log.Fatal(err)
		return false
	}

	app := tree.Root.FindChild(func(n *i3.Node) bool {
		var matchClass, matchInstance, matchTitle bool

		if class != "" {
			matchClass = n.WindowProperties.Class == class
		}
		if instance != "" {
			matchInstance = n.WindowProperties.Instance == instance
		}
		if title != "" {
			matchTitle = strings.Index(n.WindowProperties.Title, title) >= 0
		}

		if class != "" && instance != "" && title != "" {
			return matchClass && matchInstance && matchTitle
		} else if class != "" && instance != "" && title == "" {
			return matchClass && matchInstance
		} else if class != "" && instance == "" && title == "" {
			return matchClass
		} else if class == "" && instance != "" && title == "" {
			return matchInstance
		} else if class == "" && instance != "" && title != "" {
			return matchInstance && matchTitle
		} else if class != "" && instance == "" && title != "" {
			return matchClass && matchTitle
		} else if class == "" && instance == "" && title != "" {
			return matchTitle
		} else {
			return false
		}
	})

	if app == nil {
		return false
	}

	_, err = i3.RunCommand(fmt.Sprintf(`[con_id="%d"] focus`, app.ID))
	return err == nil
}
