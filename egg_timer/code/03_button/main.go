package main

import (
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		// create new window
		w := new(app.Window)
		w.Option(app.Title("Egg timer"))
		w.Option(app.Size(unit.Dp(400), unit.Dp(600)))

		// ops are the operations from the UI
		var ops op.Ops

		// startButton is a clickable widget
		var startButton widget.Clickable

		// th defines the material design style
		th := material.NewTheme()

		// listen for events in the window.
		for {
			// first grab the event
			evt := w.Event()

			// then detect the type
			switch typ := evt.(type) {

			// this is sent when the application should re-render
			case app.FrameEvent:
				gtx := app.NewContext(&ops, typ)
				btn := material.Button(th, &startButton, "Start")
				btn.Layout(gtx)
				typ.Frame(gtx.Ops)

			// and this is sent when the application should exit
			case app.DestroyEvent:
				os.Exit(0)
			}
		}
	}()
	app.Main()
}
