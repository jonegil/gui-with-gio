package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		// create new window
		w := app.NewWindow(
			app.Title("Egg timer"),
			app.Size(unit.Dp(400), unit.Dp(600)),
		)
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

type C = layout.Context
type D = layout.Dimensions

func loop(w *app.Window) error {

	// ops are the operations from the UI
	var ops op.Ops

	// startButton is a clickable widget
	var startButton widget.Clickable

	//durInput is a textfield to input boil duration
	var boilLengthInput widget.Editor

	// is the egg boiling? When did it start? for how long should it boil? used for progress
	var boiling bool
	var boilStart time.Time
	var boilLength float32
	progress := float32(0)

	// th defnes the material design style
	th := material.NewTheme(gofont.Collection())

	// listen for events in the window.
	for e := range w.Events() {

		// detect what type of event
		switch e := e.(type) {

		// this is sent when the application should re-render.
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			// Let's try out the flexbox layout concept
			// Here's a good reference for the main concepts
			// https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Flexible_Box_Layout/Basic_Concepts_of_Flexbox
			if startButton.Clicked() {
				boiling = !boiling
				boilStart = e.Now
				//Read from the input box
				inputString := boilLengthInput.Text()
				inputString = strings.TrimSpace(inputString)
				inputFloat, _ := strconv.ParseFloat(inputString, 32)
				boilLength = float32(inputFloat)
				//Resetting the boil
				if progress >= 1 {
					progress = 0
				}
			}

			// Moved the progress calculation out of the progressbar, so we can use it more places
			// Progressing the boil
			if boiling && progress < 1 {
				if boilLength != 0 {
					progress = float32(e.Now.Sub(boilStart).Seconds()) / boilLength
				}
			}
			// Limit the progress to between [0,1]
			// Try removing and send in negative time. Phsychedelic
			if progress > 1 {
				progress = 1
			}
			if progress < 0 {
				progress = 0
			}

			layout.Flex{
				// Vertical alignment, from top to bottom
				Axis: layout.Vertical,
				//Emtpy space is left at the start, i.e. at the top
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				layout.Rigid(
					func(gtx C) D {
						// Draw a custom path, shaped like an egg
						var egg clip.Path
						op.Offset(f32.Pt(200, 275)).Add(gtx.Ops)
						egg.Begin(gtx.Ops)
						// Rotate from 0 to 360 degrees
						for deg := 0.0; deg <= 360; deg++ {
							// Egg math (really) at this brilliant site. Thanks!
							// https://observablehq.com/@toja/egg-curve
							// Convert degrees to radians
							rad := deg / 360 * 2 * math.Pi
							// Trig gives the distance in X and Y direction
							cosT := math.Cos(rad)
							sinT := math.Sin(rad)
							// Constants to define the eggshape
							a := 110.0
							b := 150.0
							d := 20.0
							// The x/y coordinates
							x := a * cosT
							y := -(math.Sqrt(b*b-d*d*cosT*cosT) + d*sinT) * sinT
							p := f32.Pt(float32(x), float32(y))
							// Dra the line to this point
							egg.LineTo(p)
						}
						// Close the path
						egg.Close()
						clip.Outline{Path: egg.End()}.Op().Add(gtx.Ops)

						// Fill the shape
						// Using progress as a paramter for the colors. Nifty heh?
						color := color.NRGBA{R: 255, G: uint8(239 * (1 - progress)), B: uint8(174 * (1 - progress)), A: 255}
						paint.Fill(gtx.Ops, color)

						d := image.Point{Y: 450}
						return layout.Dimensions{Size: d}
					},
				),
				layout.Rigid(
					func(gtx C) D {
						ed := material.Editor(th, &boilLengthInput, "sec")
						ed.Editor.SingleLine = true
						ed.Editor.Alignment = text.Middle

						if boiling && progress < 1 {
							boilRemain := (1 - progress) * boilLength
							//Format to 1 decimal.
							//Using the good old multiply-by-10-divide-by-10 trick to get rounded values with 1 decimal
							inputStr := fmt.Sprintf("%.1f", math.Round(float64(boilRemain)*10)/10)
							boilLengthInput.SetText(inputStr)
						}

						margins := layout.Inset{
							Top:    unit.Dp(0),
							Right:  unit.Dp(170),
							Bottom: unit.Dp(40),
							Left:   unit.Dp(170),
						}
						border := widget.Border{
							Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
							CornerRadius: unit.Dp(3),
							Width:        unit.Dp(2),
						}
						return margins.Layout(gtx,
							func(gtx C) D {
								return border.Layout(gtx, ed.Layout)
							},
						)
					},
				),
				layout.Rigid(
					func(gtx C) D {
						if boiling && progress < 1 {
							// The progress bar hasn’t yet finished animating.
							op.InvalidateOp{}.Add(&ops)
						}
						//defer op.Save(&ops).Load()
						bar := material.ProgressBar(th, progress)
						return bar.Layout(gtx)
					},
				),
				layout.Rigid(
					func(gtx C) D {
						//We start by defining a set of margins
						margins := layout.Inset{
							Top:    unit.Dp(25),
							Bottom: unit.Dp(25),
							Right:  unit.Dp(35),
							Left:   unit.Dp(35),
						}
						//Then we lay out a layout within those margins ...
						return margins.Layout(gtx,
							// ...the same function we earlier used to create a button
							func(gtx C) D {
								var text string
								if !boiling {
									text = "Start"
								}
								if boiling && progress < 1 {
									text = "Stop"
								}
								if boiling && progress >= 1 {
									text = "Finished"
								}
								btn := material.Button(th, &startButton, text)
								return btn.Layout(gtx)
							},
						)
					},
				),
			)
			e.Frame(gtx.Ops)

		// this is sent when the application is closed.
		case system.DestroyEvent:
			return e.Err
		}
	}
	return nil
}
