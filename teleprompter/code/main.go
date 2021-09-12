package main

import (
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

var paragraphList []string

func main() {

	// Read from file
	f, err := ioutil.ReadFile("speech.txt")
	if err == nil {
		// Convert whole text into a slice of strings.
		paragraphList = strings.Split(string(f), "\n")
		// Add extra empty lines a the end. Cheap but effective trick to ensure
		// the last line of the speech scrolls out of the screen
		for i := 1; i <= 10; i++ {
			paragraphList = append(paragraphList, "")
		}
	}
	/*
		Alternative to reading from file, we can generate paragraphs
		for i := 1; i <= 2500; i++ {
			paragraphList = append(paragraphList, fmt.Sprintf("Paragraph %d", i))
		}
	*/

	// GUI
	go func() {
		// create new window
		w := app.NewWindow(
			app.Title("Teleprompter"),
			app.Size(unit.Dp(350), unit.Dp(300)),
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

	// Are we auto scrolling?
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
				}
				if e.Name == key.NameUpArrow || e.Name == "K" {
					scrollY = scrollY - stepSize*4
					if scrollY < 0 {
						scrollY = 0
					}
				}

				// To turn on/off autoscroll, and set the scrollspeed
				if e.Name == key.NameSpace {
					autoscroll = !autoscroll
					if autoscroll && autospeed == 0 {
						autospeed++
					}
				}
				if e.Name == "F" {
					autospeed++
					if !autoscroll {
						autoscroll = true
						autospeed = 1
					}
				}
				if e.Name == "S" {
					if autospeed > 0 {
						autospeed--
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

		// A re-render request?
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
				Top:    unit.Dp(0),
				Bottom: unit.Dp(0),
			}

			// Visualisation of the speech, using a list where each paragraph is a separate item.
			// Offset is the distance from the top of the screen to the first element. I.e. it controls how far we have scrolled.
			var speechViz = layout.List{
				Axis: layout.Vertical,
				Position: layout.Position{
					Offset: scrollY,
				},
			}

			// Layout the list inside the margins
			// 1) First the margins ...
			margins.Layout(gtx,
				func(gtx C) D {
					// 2) ... then the list inside those margins ...
					return speechViz.Layout(gtx, len(paragraphList),
						// 3) ... where paragraph is it's separate item
						func(gtx C, index int) D {
							// One label per paragraph
							paragraph := material.Label(th, unit.Dp(float32(fontSize)), paragraphList[index])
							// The text is centered
							paragraph.Alignment = 2
							// Return the laid out paragraph
							return paragraph.Layout(gtx)
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
