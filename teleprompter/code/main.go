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
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type C = layout.Context
type D = layout.Dimensions

var speechList []string

func main() {

	// read from file
	f, err := ioutil.ReadFile("speech.txt")
	if err == nil {
		speechList = strings.Split(string(f), "\n")
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
	var scrollY float32 = 0

	// y-position for red highlight bar
	var highlightY float32 = 78

	// width of text area
	var textWidth float32 = 300

	// fontSize
	var fontSize float32 = 35

	// Are we auato scrolling?
	var autoscroll bool = false
	var autospeed float32 = 1

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
				if e.Modifiers == key.ModCommand {
					stepSize = 1
				}

				// To scroll text
				if e.Name == key.NameDownArrow || e.Name == "J" {
					scrollY = scrollY + stepSize*4*1.025
					if scrollY < 0 {
						scrollY = 0
					}
				}
				if e.Name == key.NameUpArrow || e.Name == "K" {
					scrollY = scrollY - stepSize*4*1.025
				}

				// To turn on/off autoscroll, and set the scrollspeed
				if e.Name == key.NameSpace {
					autoscroll = !autoscroll
				}
				if e.Name == "F" {
					autospeed++
				}
				if e.Name == "S" {
					autospeed--
					if autospeed < 0 {
						autospeed = 0
					}
				}

				// To adjust the highlighter
				if e.Name == "U" {
					highlightY = highlightY - stepSize
				}
				if e.Name == "D" {
					highlightY = highlightY + stepSize
				}

				// To adjust margin width
				if e.Name == "W" {
					textWidth = textWidth + stepSize
				}
				if e.Name == "N" {
					textWidth = textWidth - stepSize
				}

				// To adjust fontsize
				// + and - are unmodified
				if e.Name == "+" {
					fontSize = fontSize + stepSize
				}
				if e.Name == "-" {
					fontSize = fontSize - stepSize
				}

				w.Invalidate()
			}

		case pointer.Event:
			if e.Type == pointer.Scroll {
				stepSize := e.Scroll.Y
				scrollY = scrollY + stepSize
				if scrollY < 0 {
					scrollY = 0
				}
				w.Invalidate()
			}

		// this is sent when the application should re-render.
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

			// Text
			wl := &widget.List{
				//Scrollbar: widget.Scrollbar{},
				List: layout.List{
					Axis: layout.Vertical,
					//ScrollToEnd: false,
					//Alignment:   0,
					Position: layout.Position{
						//BeforeEnd:  true,
						//First:  0,
						Offset: int(scrollY),
						//OffsetLast: 	0,
						//Count:      0,
						//Length: 0,
					},
				},
			}

			// Margins
			marginWidth := (float32(gtx.Constraints.Max.X) - textWidth) / 2.0
			margins := layout.Inset{
				Left:   unit.Dp(float32(marginWidth)),
				Right:  unit.Dp(float32(marginWidth)),
				Top:    unit.Dp(float32(0)),
				Bottom: unit.Dp(float32(0)),
			}

			//Layout within margins
			margins.Layout(gtx,
				func(gtx C) D {
					return material.List(th, wl).Layout(gtx, len(speechList),
						func(gtx layout.Context, index int) layout.Dimensions {
							line := speechList[index]
							speechLine := material.Label(th, unit.Dp(float32(fontSize)), line)
							speechLine.Alignment = 2
							return speechLine.Layout(gtx)
						},
					)
				},
			)
			op.Save(&ops).Load()

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

		case system.DestroyEvent:
			return e.Err
		}
	}
	return nil
}
