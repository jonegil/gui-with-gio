package main

import (
	"fmt"
	"image"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/op"
)

func main() {
	// Step 3 - Start the GUI
	go func() {
		// create new window
		w := app.NewWindow(
			app.Title("Teleprompter"),
		)
		// draw on screen
		if err := draw(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

// The main draw function
func draw(w *app.Window) error {

	var ops op.Ops

	// Define a tag for input routing
	var tag = "My Input Routing Tag - which could be this silly string, or an int/float/address, or anything else"

	for {

		// listen for events in the window
		switch winE := w.NextEvent().(type) {

		// Should we draw a new frame?
		case app.FrameEvent:
			gtx := app.NewContext(&ops, winE)

			// ---------- Handle input ----------
			// Time to deal with inputs since last frame.

			// Scrolled a mouse wheel?
			for {
				ev, ok := gtx.Event(
					pointer.Filter{
						Target: tag,
						Kinds:  pointer.Scroll,
						ScrollBounds: image.Rectangle{
							Min: image.Point{X: 0, Y: -1},
							Max: image.Point{X: 0, Y: +1},
						},
					},
				)
				if !ok {
					break
				}
				fmt.Printf("SCROLL: %+v\n", ev)
			}

			// Pressed a mouse button?
			for {
				ev, ok := gtx.Event(
					pointer.Filter{
						Target: tag,
						Kinds:  pointer.Press,
					},
				)
				if !ok {
					break
				}
				fmt.Printf("PRESS : %+v\n", ev)
			}

			// ---------- REGISTERING EVENTS ----------
			// registering events here work
			event.Op(&ops, tag)

			// ---------- FINALIZE ----------
			// Frame completes the FrameEvent by drawing the graphical operations from ops into the window.
			winE.Frame(&ops)

			// Should we shut down?
		case app.DestroyEvent:
			return winE.Err
		}
	}
}
