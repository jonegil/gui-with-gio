package main

import (
	"flag"
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

// Command line input variables
var filename *string

// A []string to hold the speech as a list of paragraphs
var paragraphList []string

func main() {
	// Part 1 - Read input from command line
	filename = flag.String("file", "speech.txt", "Which .txt file shall I present?")
	flag.Parse()

	// Part 2 - Read from file
	paragraphList = readText(filename)

	// Part 3 - Start the GUI
	go func() {
		// create new window
		w := app.NewWindow(
			app.Title("Teleprompter"),
			app.Size(unit.Dp(650), unit.Dp(600)),
		)
		// draw on screen
		if err := draw(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

func readText(filename *string) []string {
	f, err := ioutil.ReadFile(*filename)
	text := []string{}
	if err != nil {
		log.Fatal("Error when reading file:\n  ", err)
	}
	if err == nil {
		// Convert whole text into a slice of strings.
		text = strings.Split(string(f), "\n")
		// Add extra empty lines a the end. Simple trick to ensure
		// the last line of the speech scrolls out of the screen
		for i := 1; i <= 10; i++ {
			text = append(text, "")
		}
	}

	// Alternative to reading from file, we can generate paragraphs programatically
	// Handy for debugging
	//for i := 1; i <= 2500; i++ {
	//	text = append(text, fmt.Sprintf("Eloquent speech, interesting phrase %d", i))
	//}
	return text
}

func draw(w *app.Window) error {
	// y-position for text
	var scrollY unit.Dp = 0

	// y-position for red focusBar
	var focusBarY unit.Dp = 78

	// width of text area
	var textWidth unit.Dp = 550

	// fontSize
	var fontSize unit.Sp = 35

	// Are we auto scrolling?
	var autoscroll bool = false
	var autospeed unit.Dp = 1

	// th defines the material design style
	th := material.NewTheme(gofont.Collection())

	// ops are the operations from the UI
	var ops op.Ops

	// listen for events in the window.
	for windowEvent := range w.Events() {
		switch e := windowEvent.(type) {

		// FrameEvent?
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			// ---------- Handle input ----------
			// Since we use the window w as the event routing tag,
			// we here call gtx.Events(w) to get these events.

			// To set increment
			var stepSize unit.Dp = 1

			for _, gtxEvent := range gtx.Events(w) {
				switch e := gtxEvent.(type) {

				case key.EditEvent:
					e.Text = strings.ToUpper(e.Text)
					// To increase the fontsize
					if e.Text == "+" {
						fontSize = fontSize + unit.Sp(stepSize)
					}
					// To decrease the fontsize
					if e.Text == "-" {
						fontSize = fontSize - unit.Sp(stepSize)
					}

				case key.Event:
					// For better controll, we only care about pressing the key down, not releasing it up
					if e.State.String() == "Press" {
						if e.Modifiers.String() == "Shift" {
							stepSize = stepSize * 3
						}
						// Start/Stop
						if e.Name == "Space" {
							autoscroll = !autoscroll
							if autospeed == 0 {
								autoscroll = true
								autospeed++
							}
						}
						// Scroll up
						if e.Name == "K" {
							scrollY = scrollY - stepSize*4
							if scrollY < 0 {
								scrollY = 0
							}
						}
						// Scroll down
						if e.Name == "J" {
							scrollY = scrollY + stepSize*4
						}
						// Faster scrollspeed
						if e.Name == "F" {
							autoscroll = true
							autospeed++
						}
						// Slower scrollspeed
						if e.Name == "S" {
							if autospeed > 0 {
								autospeed--
							}
							if autospeed == 0 {
								autoscroll = false
							}
						}
						// Wider text to be displayed
						if e.Name == "W" {
							textWidth = textWidth + stepSize*10
						}
						// Narrow text to be displayed
						if e.Name == "N" {
							textWidth = textWidth - stepSize*10
						}
						// Move the focusBar Up
						if e.Name == "U" {
							focusBarY = focusBarY - stepSize
						}
						// Move the focusBar Down
						if e.Name == "D" {
							focusBarY = focusBarY + stepSize
						}
					}

				case pointer.Event:
					//fmt.Printf("  pointer: %#+v \n", e)
					if e.Type == pointer.Scroll {
						if e.Modifiers == key.ModShift {
							stepSize = 3
						}
						// By how much should the user scroll this time?
						thisScroll := unit.Dp(e.Scroll.Y)

						// Increment scrollY with that distance
						scrollY = scrollY + thisScroll*stepSize
						if scrollY < 0 {
							scrollY = 0
						}
					}

				default:
					fmt.Printf("gtxEvent: %#+v \n", e)
				}
			}

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
			var marginWidth unit.Dp
			marginWidth = (unit.Dp(gtx.Constraints.Max.X) - textWidth) / 2
			margins := layout.Inset{
				Left:   marginWidth,
				Right:  marginWidth,
				Top:    unit.Dp(0),
				Bottom: unit.Dp(0),
			}

			// Visualisation of the text, using a list where each paragraph is a separate item.
			// Offset is the distance from the top of the screen to the first element.
			// I.e. it controls how far we have scrolled.
			var visList = layout.List{
				Axis: layout.Vertical,
				Position: layout.Position{
					Offset: int(scrollY),
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
			op.Offset(image.Pt(0, int(focusBarY))).Add(&ops)
			stack := clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, 50)}.Push(&ops)
			paint.ColorOp{Color: color.NRGBA{R: 0xff, A: 0x66}}.Add(&ops)
			paint.PaintOp{}.Add(&ops)
			stack.Pop()

			// ---------- COLLECT INPUT ----------
			// Create a clip area the size of the window.
			// Note the Tag: w, as discussed above
			eventArea := clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Push(gtx.Ops)

			// pointer input
			pointer.InputOp{
				Types: pointer.Scroll,
				Tag:   w,
				ScrollBounds: image.Rectangle{
					Min: image.Point{
						X: 0,
						Y: -int(unit.Dp(200)),
					},
					Max: image.Point{
						X: 0,
						Y: int(unit.Dp(200)),
					},
				},
			}.Add(gtx.Ops)

			// keyboard focus, needed for general keybaord output, except the ones defined in key.InputOp
			key.FocusOp{
				Tag: w, // Use the window as the event routing tag. This means we can call gtx.Events(w) and get these events.
			}.Add(gtx.Ops)

			// Specify keys for key.Event
			// Other keys are caught as key.EditEvent
			key.InputOp{
				Keys: key.Set("(Shift)-F|(Shift)-S|(Shift)-U|(Shift)-D|(Shift)-J|(Shift)-K|(Shift)-W|(Shift)-N|Space"),
				Tag:  w, // Use the window as the event routing tag. This means we can call gtx.Events(w) and get these events.
			}.Add(gtx.Ops)

			eventArea.Pop()

			// ---------- FINALIZE ----------
			// Frame completes the FrameEvent by drawing the graphical operations from ops into the window.
			e.Frame(&ops)

			// Shutdown?
		case system.DestroyEvent:
			return e.Err

		}
	}
	return nil
}
