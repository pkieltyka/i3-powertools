package main

import (
	"github.com/davecgh/go-spew/spew"
	"go.i3wm.org/i3/v4"
)

// TODO ..
func events() {
	recv := i3.Subscribe(i3.WindowEventType)
	for recv.Next() {
		ev := recv.Event().(*i3.WindowEvent)
		spew.Dump(ev)
	}
}
