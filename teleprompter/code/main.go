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
	"gioui.org/f32"
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

var speechList []string

func main() {

	// Read from file
	f, err := ioutil.ReadFile("speech.txt")
	//f, err := ioutil.ReadFile("shakespeare_complete.txt")
	// Shakespeare has issues from "Painting my age" or the start of sonnet 63
	// Sometimes a few sonnets earlier too
	//f, err := ioutil.ReadFile("gatsby.txt")
	if err == nil {
		// Convert whole text into a slice of strings,
		// instead of one huge string. Useful to later
		// only show the lines which are visible at the moment
		speechList = strings.Split(string(f), "\n")
		// Add extra empty lines a the end. Easy trick to ensure
		// the last line of the speech scrolls up and out of the
		// screen
		for i := 1; i <= 10; i++ {
			speechList = append(speechList, "")
		}
	}

	// GUI
	go func() {
		// create new window
		w := app.NewWindow(
			app.Title("Teleprompter"),
			app.Size(unit.Dp(350), unit.Dp(300)),
			//app.Fullscreen,
		)

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

	// y-position for red highlight bar
	var highlightY int = 78

	// width of text area
	var textWidth int = 300

	// fontSize
	var fontSize int = 35

	// Are we auato scrolling?
	var autoscroll bool = false
	var autospeed int = 1

	// th defines the material design style
	th := material.NewTheme(gofont.Collection())

	// listen for events in the window.
	for e := range w.Events() {

		// Detect what type of event
		switch e := e.(type) {

		// A keypress?
		case key.Event:
			if e.State == key.Press {
				// To set increment
				var stepSize int = 10
				if e.Modifiers == key.ModShift {
					stepSize = 1
				}

				// To scroll text
				if e.Name == key.NameDownArrow || e.Name == "J" {
					scrollY = scrollY + stepSize*4
					if scrollY < 0 {
						scrollY = 0
					}
				}
				if e.Name == key.NameUpArrow || e.Name == "K" {
					scrollY = scrollY - stepSize*4
				}

				// To turn on/off autoscroll, and set the scrollspeed
				if e.Name == key.NameSpace {
					autoscroll = !autoscroll
				}
				if e.Name == "F" {
					if autoscroll {
						autospeed++
					}
					if !autoscroll {
						autoscroll = true
					}
				}
				if e.Name == "S" {
					autospeed--
					if autospeed < 0 {
						autospeed = 0
					}
					if autospeed == 0 {
						autoscroll = false
					}
				}

				// To adjust text width
				if e.Name == "W" {
					textWidth = textWidth + stepSize
				}
				if e.Name == "N" {
					textWidth = textWidth - stepSize
				}

				// To adjust fontsize
				if e.Name == "+" {
					fontSize = fontSize + stepSize
				}
				if e.Name == "-" {
					fontSize = fontSize - stepSize
				}

				// To adjust the highlighter
				if e.Name == "U" {
					highlightY = highlightY - stepSize
				}
				if e.Name == "D" {
					highlightY = highlightY + stepSize
				}

				w.Invalidate()
			}

		// A mouse event?
		case pointer.Event:
			if e.Type == pointer.Scroll {
				// How far did the pointer scroll?
				step := e.Scroll.Y
				// Increment the Y-scroll with that distance
				scrollY = scrollY + int(step)
				if scrollY < 0 {
					scrollY = 0
				}
				w.Invalidate()
			}

		// A re-render request
		case system.FrameEvent:
			// ops are the operations from the UI
			var ops op.Ops

			// Graphical context
			gtx := layout.NewContext(&ops, e)

			// Bacground
			paint.Fill(&ops, color.NRGBA{R: 0xff, G: 0xfe, B: 0xe0, A: 0xff})

			// Textscroll
			if autoscroll {
				scrollY = scrollY + autospeed
				op.InvalidateOp{At: gtx.Now.Add(time.Second / 25)}.Add(&ops)
			}

			// Margins
			marginWidth := (gtx.Constraints.Max.X - textWidth) / 2
			margins := layout.Inset{
				Left:   unit.Dp(float32(marginWidth)),
				Right:  unit.Dp(float32(marginWidth)),
				Top:    unit.Dp(float32(0)),
				Bottom: unit.Dp(float32(0)),
			}

			// Visualisation of the speech, using a list
			// The offset is the pixel-distance from the top edge to the
			// first element in our speechlist
			var speechViz = layout.List{
				Axis: layout.Vertical,
				Position: layout.Position{
					Offset: scrollY,
				},
			}

			// Layout the list inside the margins
			margins.Layout(gtx,
				func(gtx C) D {
					return speechViz.Layout(gtx, len(speechList),
						func(gtx C, index int) D {
							fmt.Println(index)
							// One label per paragraph
							l := material.Label(th, unit.Dp(float32(fontSize)), speechList[index])
							// The label is centered
							l.Alignment = 2
							// Return the laid out label
							return l.Layout(gtx)
						},
					)
				},
			)

			// Draw a transparent red rectangle.
			path := new(clip.Path)
			stack := op.Save(&ops)
			path.Begin(&ops)
			path.MoveTo(f32.Pt(0, 0))
			path.End()
			op.Offset(f32.Pt(0, float32(highlightY))).Add(&ops)
			clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, 50)}.Add(&ops)
			paint.ColorOp{Color: color.NRGBA{R: 0xff, A: 0x66}}.Add(&ops)
			paint.PaintOp{}.Add(&ops)
			stack.Load()

			e.Frame(&ops)

		// Shutdown?
		case system.DestroyEvent:
			return e.Err
		}
	}
	return nil
}
