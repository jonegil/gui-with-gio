// eclipse illustrates the eclipse
package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
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
	th := material.NewTheme(gofont.Collection())

	for e := range w.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			os.Exit(0)

		case system.FrameEvent:

			gtx := layout.NewContext(&ops, e)

			fmt.Printf("Main: %#v\n", gtx.Constraints)

			if button.Clicked() {
				switch moon {
				case white:
					moon = blue
				case blue:
					moon = cheese
				case cheese:
					moon = white
				}
			}

			fmt.Printf("Main: %#v\n", gtx.Constraints)
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
						fmt.Printf("  Canvas: %#v\n", canvas.Context.Constraints)

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
						fmt.Printf("  Button: %#v\n", gtx.Constraints)
						btn := material.Button(th, &button, "Loco luna")
						btn.CornerRadius = unit.Dp(0)
						return btn.Layout(gtx)
					},
				),
			)
			e.Frame(gtx.Ops)

		case key.Event:
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
			if e.Type == pointer.Scroll {
				eclipseX += e.Scroll.X / 100
				eclipseY += e.Scroll.Y / 100
				w.Invalidate()
			}
		}

	}
	return nil
}

func main() {
	go func() {
		// create new window
		w := app.NewWindow(
			app.Title("Eclipse"),
			app.Size(unit.Dp(400), unit.Dp(600)),
		)
		if err := draw(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
