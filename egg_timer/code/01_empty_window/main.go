package main

import (
	"gioui.org/app"
)

func main() {
	go func() {
		// create new window
		w := new(app.Window)

		// listen for events in the window
		for {
			w.Event()
		}
	}()
	app.Main()
}
