package main

import (
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		// create new window
		w := app.NewWindow(
			app.Title("Egg timer"),
			app.Size(unit.Dp(400), unit.Dp(600)),
		)

		// ops are the operations from the UI
		var ops op.Ops

		// startButton is a clickable widget
		var startButton widget.Clickable

		// th defines the material design style
		th := material.NewTheme()

		// listen for events in the window.
		for {
			evt := w.NextEvent()

			// detect what type of event
			switch typ := evt.(type) {

			// this is sent when the application should re-render.
			case app.FrameEvent:
				gtx := app.NewContext(&ops, typ)
				// Let's try out the flexbox layout:
				layout.Flex{
					// Vertical alignment, from top to bottom
					Axis: layout.Vertical,
					// Empty space is left at the start, i.e. at the top
					Spacing: layout.SpaceStart,
				}.Layout(gtx,
					// We insert two rigid elements:
					// First one to hold a button ...
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(th, &startButton, "Start")
							return btn.Layout(gtx)
						},
					),
					// ... then one to hold an empty spacer
					layout.Rigid(
						// The height of the spacer is 25 Device independent pixels
						layout.Spacer{Height: unit.Dp(25)}.Layout,
					),
				)
				typ.Frame(gtx.Ops)

			case app.DestroyEvent:
				os.Exit(0)
			}
		}
	}()
	app.Main()
}
