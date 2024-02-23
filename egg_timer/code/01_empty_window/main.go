package main

import (
	"gioui.org/app"
)

func main() {
	go func() {
		// create new window
		w := app.NewWindow()

		// listen for events in the window
		for {
			w.NextEvent()
		}
	}()
	app.Main()
}
