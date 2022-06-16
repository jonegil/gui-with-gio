package main

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type C = layout.Context
type D = layout.Dimensions

// A []string to hold the speech as a list of paragraphs
var paragraphList []string

func main() {
	// Part 1 - Read from file
	f, err := ioutil.ReadFile("speech.txt")
	if err == nil {
		// Convert whole text into a slice of strings.
		paragraphList = strings.Split(string(f), "\n")
		// Add extra empty lines a the end. Simple trick to ensure
		// the last line of the speech scrolls out of the screen
		for i := 1; i <= 10; i++ {
			paragraphList = append(paragraphList, "")
		}
	}

	// Alternative to reading from file, we can generate paragraphs programatically
	// Handy for debugging
	/*
	   for i := 1; i <= 2500; i++ {
	     paragraphList = append(paragraphList, fmt.Sprintf("Eloquent speech, interesting phrase %d", i))
	   }
	*/

	// Part 2 - Start the GUI
	go func() {
		// create new window
		w := app.NewWindow(
			app.Title("Teleprompter"),
			app.Size(unit.Dp(350), unit.Dp(300)),
		)
		// draw on screen
		if err := draw(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

func draw(w *app.Window) error {
	// y-position for text
	var scrollY int = 0

	// y-position for red focusBar
	var focusBarY int = 78

	// width of text area
	var textWidth int = 300

	// fontSize
	var fontSize int = 35

	// Are we auto scrolling?
	var autoscroll bool = false
	var autospeed int = 1

	// To set increment
	var stepSize int = 1

	// th defines the material design style
	th := material.NewTheme(gofont.Collection())

	// ops are the operations from the UI
	var ops op.Ops

	// listen for events in the window.
	for windowEvent := range w.Events() {
		//fmt.Printf("+%v --\n", windowEvent)
		// Shutdown?
		switch e := windowEvent.(type) {

		case system.DestroyEvent:
			return e.Err

		case system.FrameEvent:

			// Source:
			// https://lists.sr.ht/~eliasnaur/gio/%3CCAFcc3FQNTp_UXr7oA97SsVPD7D91jSw30ZtALcT9vmopFDTeZQ%40mail.gmail.com%3E#%3CCAE_4BPB=DS9eXrmSGxkBku-VTfLZXZjp0U_VMgYrU7M3GQ7NaQ@mail.gmail.com%3E
			// https://go.dev/play/p/SDHy1LZRljf
			// https://go.dev/play/p/VDQg6sxRyA4

			// Graphical context
			gtx := layout.NewContext(&ops, e)

			// Gather and deal with all events captured by our input area since the previous frame.
			// Do eventhandling here rather than in the outer w.Events() loop
			for _, gtxEvent := range gtx.Events(w) {
				fmt.Printf("  gtx: %#+v --\n", gtxEvent)
				// Perform event handling here instead of in the outer type switch.
				switch e := gtxEvent.(type) {

				case key.EditEvent:
					fmt.Printf("    key: %#+v --\n", e)
					// Conert all input to uppercase
					e.Text = strings.ToUpper(e.Text)
					// Spacebar
					if e.Text == " " {
						autoscroll = !autoscroll
						if autospeed == 0 {
							autoscroll = true
							autospeed++
						}
					}
					// Scroll up
					if e.Text == "K" { //e.Name == key.NameUpArrow || e.Name == "K" {
						scrollY = scrollY - stepSize*4
						if scrollY < 0 {
							scrollY = 0
						}
					}
					// Scroll down
					if strings.ToUpper(e.Text) == "J" { //e.Name == key.NameDownArrow || e.Name == "J" {
						scrollY = scrollY + stepSize*4
					}
					// Faster scrollspeed
					if e.Text == "F" {
						autoscroll = true
						autospeed++
					}
					// Slower scrollspeed
					if e.Text == "S" {
						if autospeed > 0 {
							autospeed--
						}
						if autospeed == 0 {
							autoscroll = false
						}
					}
					// Wider text to be displayed
					if e.Text == "W" {
						textWidth = textWidth + stepSize*10
					}
					// Narrow text to be displayed
					if e.Text == "N" {
						textWidth = textWidth - stepSize*10
					}
					// To increase the fontsize
					if e.Text == "+" {
						fontSize = fontSize + stepSize
					}
					// To decrease the fontsize
					if e.Text == "-" {
						fontSize = fontSize - stepSize
					}
					// Move the focusBar Up
					if e.Text == "U" {
						focusBarY = focusBarY - stepSize
					}
					// Move the focusBar Down
					if e.Text == "D" {
						focusBarY = focusBarY + stepSize
					}

				case key.Event:
					fmt.Printf("    key: %#+v --\n", e)
					if e.Name == key.NameUpArrow {
						fmt.Println(e.Name, "UP")
					}
					if e.State == key.Press {
						// To set increment
						var stepSize int = 1
						if e.Modifiers == key.ModShift {
							stepSize = 10
						}
						fmt.Println(stepSize)
					}

				case pointer.Event:
					if e.Type == pointer.Scroll {
						var stepSize int = 1
						if e.Modifiers == key.ModShift {
							stepSize = 3
						}
						// By how much should the user scroll this time?
						thisScroll := int(e.Scroll.Y)
						// Increment scrollY with that distance
						scrollY = scrollY + thisScroll*stepSize
						if scrollY < 0 {
							scrollY = 0
						}
					}

				}
			}

			// ---------- AREA FOR EVENTS ----------
			// Create a clip area the size of the window.
			eventArea := clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Push(gtx.Ops)
			// Register for pointer inputs on the current clip area.
			pointer.InputOp{
				Types: pointer.Enter | pointer.Leave | pointer.Drag | pointer.Press | pointer.Release | pointer.Scroll | pointer.Move,
				Tag:   w, // Use the window as the event routing tag. This means we can call gtx.Events(w) and get these events.
			}.Add(gtx.Ops)
			// Register for keyboard focus on the current clip area.
			key.FocusOp{
				Tag: w, /// Use the window as the event routing tag. This means we can call gtx.Events(w) and get these events.
			}.Add(gtx.Ops)
			// Register for keyboard input on the current clip area.
			key.InputOp{
				Keys: key.Set("Shift-W|NameUpArrow|DownArrow"),
				Tag:  w, // Use the window as the event routing tag. This means we can call gtx.Events(w) and get these events.
			}.Add(gtx.Ops)
			// Pop the clip area to finalize it.
			eventArea.Pop()

			// ---------- LAYOUT ----------
			// Layout the interface _BEFORE_ you pop the clip area.
			// This ensures that the clip is logically the ancestor of the layout,
			// so key events unhandled by the interface will propagate upwards to it.

			// Bacground
			paint.Fill(&ops, color.NRGBA{R: 0xff, G: 0xfe, B: 0xe0, A: 0xff})

			// Textscroll
			if autoscroll {
				scrollY = scrollY + autospeed
				op.InvalidateOp{At: gtx.Now.Add(time.Second / 50)}.Add(&ops)
			}

			// Margins
			marginWidth := (gtx.Constraints.Max.X - textWidth) / 2
			margins := layout.Inset{
				Left:   unit.Dp(float32(marginWidth)),
				Right:  unit.Dp(float32(marginWidth)),
				Top:    unit.Dp(0),
				Bottom: unit.Dp(0),
			}

			// Visualisation of the text, using a list where each paragraph is a separate item.
			// Offset is the distance from the top of the screen to the first element.
			// I.e. it controls how far we have scrolled.
			var visList = layout.List{
				Axis: layout.Vertical,
				Position: layout.Position{
					Offset: scrollY,
				},
			}

			// Layout the list inside margins
			// 1) First the margins ...
			margins.Layout(gtx,
				func(gtx C) D {
					// 2) ... then the list inside those margins ...
					return visList.Layout(gtx, len(paragraphList),
						// 3) ... where each paragraph is a separate item
						func(gtx C, index int) D {
							// One label per paragraph
							paragraph := material.Label(th, unit.Sp(float32(fontSize)), paragraphList[index])
							// The text is centered
							paragraph.Alignment = 2
							// Return the laid out paragraph
							return paragraph.Layout(gtx)
						},
					)
				},
			)

			// ---------- THE FOCUS BAR ----------
			// Draw the transparent red bar.
			op.Offset(image.Pt(0, focusBarY)).Add(&ops)
			stack := clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, 50)}.Push(&ops)
			paint.ColorOp{Color: color.NRGBA{R: 0xff, A: 0x66}}.Add(&ops)
			paint.PaintOp{}.Add(&ops)
			stack.Pop()

			// ---------- FINALIZE ----------
			// Frame completes the FrameEvent by drawing the graphical operations from ops into the window.
			e.Frame(&ops)
		}
	}
	return nil
}
