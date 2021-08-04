package main

import (
	"gioui.org/app"
"fmt"
)


func main() {
	go func() {
		// create new window
		w := app.NewWindow()

		// listen for events in the window.
		for e:= range w.Events() {
fmt.Println(e)
		}
	}()
	app.Main()
}
