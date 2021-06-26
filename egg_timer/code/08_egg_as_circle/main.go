package main

import (
	"image"
	"image/color"
	"log"
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

	// is the egg boiling? When did it start? Used for progress
	var boiling bool
	var boilStart time.Time

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
				boilStart = time.Now()
			}

			layout.Flex{
				// Vertical alignment, from top to bottom
				Axis: layout.Vertical,
				//Emtpy space is left at the start, i.e. at the top
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				layout.Rigid(
					func(gtx C) D {
						clip.Circle{
							// Hard coding the x coordinate. Try resizing the window
							Center: f32.Point{X: 200, Y: 200},
							// Soft coding the x coordinate. Try resizing the window
							//Center: f32.Point{X: float32(gtx.Constraints.Max.X) / 2, Y: 200},
							Radius: 120,
						}.Add(gtx.Ops)
						color := color.NRGBA{R: 200, A: 255}
						paint.Fill(gtx.Ops, color)
						d := image.Point{Y: 500}
						return layout.Dimensions{Size: d}
					},
				),

				layout.Rigid(
					func(gtx C) D {
						progress := float32(0)
						if boiling {
							boilTime := time.Since(boilStart)
							progress = float32(boilTime.Seconds() / 10)
							if progress < 1 {
								// The progress bar hasn’t yet finished animating.
								op.InvalidateOp{}.Add(&ops)
							} else {
								progress = 1
							}
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
	}
	return nil
}
