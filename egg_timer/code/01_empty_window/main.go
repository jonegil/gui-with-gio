package main

import (
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
)

func main() {
	go func() {
		// create new window
		w := app.NewWindow()
		var ops op.Ops

		// listen for events in the window.
		for e := range w.Events() {

			// detect which type of event
			switch e := e.(type) {

			// Is it a FrameEvent? Those are sent when the application should re-render.
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				e.Frame(gtx.Ops)
			}
		}
		os.Exit(0)
	}()
	app.Main()
}
