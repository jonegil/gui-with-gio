package main

import (
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"math"
	"os"

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
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type C = layout.Context
type D = layout.Dimensions

var speech string

func main() {

	// read from file
	dat, err := ioutil.ReadFile("speech.txt")
	if err == nil {
		speech = string(dat)
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
	// y-position for scroll text
	var scrollY float32 = 0

	// y-position for red highlight bar
	var highlightY float32 = 78

	// width of text area
	var textWidth float32 = 300

	// fontSize
	var fontSize float32 = 35

	// editor is a text field
	//var editor widget.Editor
	//editor.SetText(speech)

	// th defines the material design style
	th := material.NewTheme(gofont.Collection())

	// listen for events in the window.
	for e := range w.Events() {

		// detect what type of event
		switch e := e.(type) {

		case key.Event:
			if e.State == key.Press {
				// To set increment
				stepSize := float32(10)
				if e.Modifiers == key.ModShift {
					stepSize = 1
				}

				// To scroll text
				if e.Name == key.NameDownArrow || e.Name == key.NameSpace {
					scrollY = scrollY - stepSize*4*1.025
					highlightY = highlightY + stepSize*4*1.025
				}
				if e.Name == key.NameUpArrow {
					scrollY = scrollY + stepSize*4*1.025
					highlightY = highlightY - stepSize*4*1.025
				}

				// To adjust the highlighter
				if e.Name == "U" || e.Name == "K" {
					highlightY = highlightY - stepSize
				}
				if e.Name == "D" || e.Name == "J" {
					highlightY = highlightY + stepSize
				}

				// To adjust margin width
				if e.Name == "W" || e.Name == "L" {
					textWidth = textWidth + stepSize
				}
				if e.Name == "N" || e.Name == "H" {
					textWidth = textWidth - stepSize
				}

				// To adjust fontsize
				// + and - are unmodified
				// Shift and + is ?
				// Shitf and - is _
				if e.Name == "+" || e.Name == "?" {
					fontSize = fontSize + stepSize
				}
				if e.Name == "-" || e.Name == "_" {
					fontSize = fontSize - stepSize
				}

				w.Invalidate()
			}

		case pointer.Event:
			if e.Type == pointer.Scroll {
				stepSize := e.Scroll.Y
				scrollY = scrollY - stepSize
				highlightY = highlightY + stepSize

				w.Invalidate()
			}

		// this is sent when the application should re-render.
		case system.FrameEvent:
			// ops are the operations from the UI
			var ops op.Ops

			// Graphical context
			gtx := layout.NewContext(&ops, e)

			// Background
			paint.Fill(&ops, color.NRGBA{R: 0xff, G: 0xfe, B: 0xe0, A: 0xff})

			// Margins
			marginWidth := (float32(gtx.Constraints.Max.X) - textWidth) / 2.0

			// Text
			speech := material.Label(th, unit.Dp(float32(fontSize)), speech)
			speech.Alignment = text.Middle

			// Save our current drawing offset and constraints before
			// transforming them.
			saved := op.Save(&ops)
			originalConstraints := gtx.Constraints
			// Center the text by offsetting it from the left edge by marginWidth,
			// and offset it vertically by the current scroll position.
			op.Offset(f32.Pt(float32(marginWidth), float32(scrollY))).Add(&ops)
			gtx.Constraints.Max.X = int(textWidth)
			// Set the vertical constraints to be completely unlimited. Otherwise,
			// the material.LabelStyle will stop drawing text once it is confident
			// that the next line won't fit within the provided constraints.
			gtx.Constraints.Max.Y = math.MaxInt
			// Actually lay out the text.
			speech.Layout(gtx)
			// Reset our offset and constraints.
			saved.Load()
			gtx.Constraints = originalConstraints

			// Draw a transparent red rectangle.
			stack := op.Save(&ops)
			op.Offset(f32.Pt(0, float32(highlightY))).Add(&ops)
			clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, 50)}.Add(&ops)
			paint.ColorOp{Color: color.NRGBA{R: 0xff, A: 0x66}}.Add(&ops)
			paint.PaintOp{}.Add(&ops)
			stack.Load()

			e.Frame(&ops)

		case system.DestroyEvent:
			return e.Err

		}
	}
	return nil
}
