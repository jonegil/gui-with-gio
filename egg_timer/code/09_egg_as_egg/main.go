package main

import (
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Define the progress variables, a channel and a variable
var progressIncrementer chan float32
var progress float32

func main() {
	// Setup a separate channel to provide ticks to increment progress
	progressIncrementer = make(chan float32)
	go func() {
		for {
			time.Sleep(time.Second / 25)
			progressIncrementer <- 0.004
		}
	}()

	go func() {
		// create new window
		w := app.NewWindow(
			app.Title("Egg timer"),
			app.Size(unit.Dp(400), unit.Dp(600)),
		)
		if err := draw(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

type C = layout.Context
type D = layout.Dimensions

func draw(w *app.Window) error {
	// ops are the operations from the UI
	var ops op.Ops

	// startButton is a clickable widget
	var startButton widget.Clickable

	// is the egg boiling?
	var boiling bool

	// th defines the material design style
	th := material.NewTheme(gofont.Collection())

	for {
		select {
		// listen for events in the window.
		case e := <-w.Events():

			// detect what type of event
			switch e := e.(type) {

			// this is sent when the application should re-render.
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				// Let's try out the flexbox layout concept
				if startButton.Clicked() {
					boiling = !boiling
				}

				layout.Flex{
					// Vertical alignment, from top to bottom
					Axis: layout.Vertical,
					// Empty space is left at the start, i.e. at the top
					Spacing: layout.SpaceStart,
				}.Layout(gtx,

					layout.Rigid(
						func(gtx C) D {
							// Draw a custom path, shaped like an egg
							var eggPath clip.Path
							op.Offset(f32.Pt(200, 150)).Add(gtx.Ops)
							eggPath.Begin(gtx.Ops)
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
								// Finally the point on the outline
								p := f32.Pt(float32(x), float32(y))
								// Draw the line to this point
								eggPath.LineTo(p)
							}
							// Close the path
							eggPath.Close()

							// Get hold of the actual clip
							eggArea := clip.Outline{Path: eggPath.End()}.Op()

							// Fill the shape
							// color := color.NRGBA{R: 255, G: 239, B: 174, A: 255}
							color := color.NRGBA{R: 255, G: uint8(239 * (1 - progress)), B: uint8(174 * (1 - progress)), A: 255}
							paint.FillShape(gtx.Ops, color, eggArea)

							d := image.Point{Y: 375}
							return layout.Dimensions{Size: d}
						},
					),

					layout.Rigid(
						func(gtx C) D {
							bar := material.ProgressBar(th, progress)
							return bar.Layout(gtx)
						},
					),

					layout.Rigid(
						func(gtx C) D {
							// We start by defining a set of margins
							margins := layout.Inset{
								Top:    unit.Dp(25),
								Bottom: unit.Dp(25),
								Right:  unit.Dp(35),
								Left:   unit.Dp(35),
							}
							// Then we lay out within those margins ...
							return margins.Layout(gtx,
								// ...the same function we earlier used to create a button
								func(gtx C) D {
									var text string
									if !boiling {
										text = "Start"
									} else {
										text = "Stop"
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

		// listen for events in the incrementor channel
		case p := <-progressIncrementer:
			if boiling && progress < 1 {
				progress += p
				w.Invalidate()
			}
		}
	}
}
