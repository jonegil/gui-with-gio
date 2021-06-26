package main

import (
	"gioui.org/app"
	"gioui.org/unit"
)

func main() {
	go func() {
		// create new window
		w := app.NewWindow(
			app.Title("Egg timer"),
			app.Size(unit.Dp(400), unit.Dp(600)),
		)

		// listen for events in the window.
		for range w.Events() {
		}
	}()
	app.Main()
}
