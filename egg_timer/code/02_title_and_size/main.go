package main

import (
	"gioui.org/app"
	"gioui.org/unit"
)

func main() {
	go func() {
		// create new window
		w := new(app.Window)
		w.Option(app.Title("Egg timer"))
		w.Option(app.Size(unit.Dp(400), unit.Dp(600)))

		// listen for events in the window
		for {
			w.Event()
		}
	}()
	app.Main()
}
