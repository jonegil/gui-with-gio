// eclipse illustrates the eclipse
package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/ajstarks/giocanvas"
)

func draw(w *app.Window) error {

	var ops op.Ops

	// Origo for eclipse
	eclipseX := float32(0)
	eclipseY := float32(0)

	// Some colors
	black := color.NRGBA{0, 0, 0, 255}
	white := color.NRGBA{255, 255, 255, 255}
	blue := color.NRGBA{0, 0, 255, 255}
	cheese := color.NRGBA{0xFF, 0xa6, 0, 255}
	moon := white

	// button is a clickable widget
	var button widget.Clickable

	// th defines the material design style
	th := material.NewTheme()

	for {
		// listen for events in the window
		switch e := w.NextEvent().(type) {
		case app.DestroyEvent:
			os.Exit(0)

		case app.FrameEvent:
			// fmt.Printf("Main: %#v\n", gtx.Constraints)
			gtx := app.NewContext(&ops, e)
			if button.Clicked(gtx) {
				switch moon {
				case white:
					moon = blue
				case blue:
					moon = cheese
				case cheese:
					moon = white
				}
			}

			// Let's try out the flexbox layout:
			layout.Flex{
				// Vertical alignment, from top to bottom
				Axis: layout.Vertical,
				// Empty space is left at the start, i.e. at the top
				Spacing: layout.SpaceStart,
			}.Layout(gtx,

				// We insert two rigid elements:
				// ... first one for the canvas
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						canvas := giocanvas.Canvas{
							Width:   float32(gtx.Constraints.Max.X),
							Height:  float32(500),
							Context: gtx,
						}
						// fmt.Printf("  Canvas: %#v\n", canvas.Context.Constraints)

						canvas.CenterRect(50, 50, 95, 95, black)
						var x float32 = 10.0
						var r float32 = 5.0
						var y float32 = 50.0
						for x = 25.0; x < 100.0; x += 25 {
							canvas.Circle(x, 50, r+0.5, moon)
							canvas.Circle(eclipseX+x-1, eclipseY+y+2, r, black)
							y -= 2
						}

						return layout.Dimensions{
							Size: image.Point{X: int(canvas.Width), Y: int(canvas.Height)},
						}
					},
				),
				// ... then one for the button
				layout.Rigid(
					func(gtx layout.Context) layout.Dimensions {
						// fmt.Printf("  Button: %#v\n", gtx.Constraints)
						btn := material.Button(th, &button, "Loco luna")
						btn.CornerRadius = unit.Dp(0)
						return btn.Layout(gtx)
					},
				),
			)
			e.Frame(gtx.Ops)

		case key.Event:
			fmt.Printf("  Key: %#v\n", e)
			if e.State == key.Press {
				switch e.Name {
				case "Q", key.NameEscape:
					os.Exit(0)
				case key.NameDownArrow:
					eclipseY--
				case key.NameUpArrow:
					eclipseY++
				case key.NameRightArrow:
					eclipseX++
				case key.NameLeftArrow:
					eclipseX--
				}
				w.Invalidate()
			}

		case pointer.Event:
			fmt.Printf("  Pointer: %#v\n", e)
			if e.Kind == pointer.Scroll {
				eclipseX += e.Scroll.X / 100
				eclipseY += e.Scroll.Y / 100
				w.Invalidate()
			}
		}
	}
}

func main() {
	go func() {
		// create new window
		w := app.NewWindow(
			app.Title("Eclipse"),
			app.Size(unit.Dp(800), unit.Dp(500)),
		)
		if err := draw(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
