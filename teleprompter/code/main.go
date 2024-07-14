package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// Command line input variables
var filename *string

// Define context and dimension types, just for shorthand comfort
type C = layout.Context
type D = layout.Dimensions

// A []string to hold the speech as a list of paragraphs
var paragraphList []string

// Colors
type colorMode struct {
	background color.NRGBA
	foreground color.NRGBA
	focusbar   color.NRGBA
}

func main() {
	// Step 1 - Read input from command line
	filename = flag.String("file", "speech.txt", "Which .txt file shall I present?")
	flag.Parse()

	// Step 2 - Read from file
	paragraphList = readText(filename)

	// Step 3 - Start the GUI
	go func() {
		// create new window
		w := new(app.Window)
		w.Option(app.Title("Teleprompter"))
		w.Option(app.Size(unit.Dp(650), unit.Dp(600)))
		// draw on screen
		if err := draw(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

func readText(filename *string) []string {
	f, err := os.ReadFile(*filename)
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

// The main draw function
func draw(w *app.Window) error {
	// y-position for text
	var scrollY unit.Dp = 0

	// y-position for red focusBar
	var focusBarY unit.Dp = 170

	// width of text area
	var textWidth unit.Dp = 550

	// fontSize
	var fontSize unit.Sp = 35

	// Are we auto scrolling?
	var autoscroll bool = false
	var autospeed unit.Dp = 1

	// th defines the material design style
	th := material.NewTheme()

	// ops are the operations from the UI
	var ops op.Ops

	// Define a tag for input routing
	var tag = "My Input Routing Tag - which could be this silly string, or an int/float/address, or anything else"

	// Colors
	colorDark := colorMode{
		background: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff},
		foreground: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		focusbar:   color.NRGBA{R: 0xff, G: 0x00, B: 0x00, A: 0x33},
	}

	colorLight := colorMode{
		background: color.NRGBA{R: 0xff, G: 0xfe, B: 0xe0, A: 0xff},
		foreground: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff},
		focusbar:   color.NRGBA{R: 0xff, A: 0x66},
	}

	// Define a color to start with. We like dark
	myColor := colorDark

	for {

		// listen for events in the window
		switch winE := w.Event().(type) {

		// Should we draw a new frame?
		case app.FrameEvent:
			gtx := app.NewContext(&ops, winE)

			// ---------- Handle input ----------
			// Time to deal with inputs since last frame.

			// Scrolled a mouse wheel?
			for {
				ev, ok := gtx.Event(
					pointer.Filter{
						Target:  tag,
						Kinds:   pointer.Scroll,
						ScrollY: pointer.ScrollRange{Min: -1, Max: +1},
					},
				)
				if !ok {
					break
				}
				fmt.Printf("SCROLL: %+v\n", ev)
				scrollY = scrollY + unit.Dp(ev.(pointer.Event).Scroll.Y*float32(fontSize))
				if scrollY < 0 {
					scrollY = 0
				}
			}

			// Pressed a mouse button?
			for {
				ev, ok := gtx.Event(
					pointer.Filter{
						Target: tag,
						Kinds:  pointer.Press,
					},
				)
				if !ok {
					break
				}
				fmt.Printf("PRESS : %+v\n", ev)
				// Start / stop
				autoscroll = !autoscroll
			}

			// Pressed a key?
			for {
				ev, ok := gtx.Event(
					key.Filter{Name: key.NameSpace},
					key.Filter{Optional: key.ModShift, Name: "U"},
					key.Filter{Optional: key.ModShift, Name: "D"},
					key.Filter{Optional: key.ModShift, Name: "J"},
					key.Filter{Optional: key.ModShift, Name: "K"},
					key.Filter{Optional: key.ModShift, Name: key.NameUpArrow},
					key.Filter{Optional: key.ModShift, Name: key.NameDownArrow},
					key.Filter{Optional: key.ModShift, Name: key.NamePageUp},
					key.Filter{Optional: key.ModShift, Name: key.NamePageDown},
					key.Filter{Optional: key.ModShift, Name: "F"},
					key.Filter{Optional: key.ModShift, Name: "S"},
					key.Filter{Optional: key.ModShift, Name: "+"},
					key.Filter{Optional: key.ModShift, Name: "-"},
					key.Filter{Optional: key.ModShift, Name: "W"},
					key.Filter{Optional: key.ModShift, Name: "N"},
					key.Filter{Optional: key.ModShift, Name: "C"},
				)
				if !ok {
					break
				}
				fmt.Printf("KEY   : %+v\n", ev)
				if ev.(key.Event).State == key.Press {
					name := ev.(key.Event).Name
					mod := ev.(key.Event).Modifiers

					// Set stepsize
					var stepSize unit.Dp = 1
					if mod == key.ModShift {
						stepSize = 5
					}

					// Start / stop
					if name == key.NameSpace {
						autoscroll = !autoscroll
						if autoscroll && autospeed <= 0 {
							autospeed = stepSize
						}
					}

					// Move the focusBar Up
					if name == "U" {
						focusBarY = focusBarY - stepSize
					}

					// Move the focusBar Down
					if name == "D" {
						focusBarY = focusBarY + stepSize
					}

					// Scroll up
					if name == "K" || name == key.NameUpArrow {
						scrollY = scrollY - stepSize*4
					}
					if name == key.NamePageUp {
						scrollY = scrollY - stepSize*100
					}
					if scrollY < 0 {
						scrollY = 0
					}

					// Scroll down
					if name == "J" || name == key.NameDownArrow || name == key.NamePageDown {
						scrollY = scrollY + stepSize*4
					}
					if name == key.NamePageDown {
						scrollY = scrollY + stepSize*100
					}

					// Faster scrollspeed
					if name == "F" {
						autoscroll = true
						autospeed += stepSize
					}

					// Slower scrollspeed
					if name == "S" {
						if autospeed > 0 {
							autospeed -= stepSize
						}
						if autospeed <= 0 {
							autospeed = 0
							autoscroll = false
						}
					}

					// To increase the fontsize
					if name == "+" {
						fontSize = fontSize + unit.Sp(stepSize)
					}

					// To decrease the fontsize
					if name == "-" {
						fontSize = fontSize - unit.Sp(stepSize)
					}

					// Widen text to be displayed
					if name == "W" {
						textWidth = textWidth + stepSize*10
					}
					// Narrow text to be displayed
					if name == "N" {
						textWidth = textWidth - stepSize*10
					}

					// Swhich Colormode
					if name == "C" {
						if myColor == colorDark {
							myColor = colorLight
						} else {
							myColor = colorDark
						}
					}
				}
			}

			// ---------- LAYOUT ----------
			// First we layout the user interface.
			// Let's start with a background color
			paint.Fill(&ops, myColor.background)

			// ---------- THE SCROLLING TEXT ----------
			// First, check if we should autoscroll
			// That's done by increasing the value of scrollY
			if autoscroll {
				if autospeed < 0 {
					autospeed = 0
				}
				scrollY = scrollY + autospeed
				// Invalidate 50 times per second
				inv := op.InvalidateCmd{At: gtx.Now.Add(time.Second / 50)}
				gtx.Execute(inv)
			}
			// Then we use scrollY to control the distance from the top of the screen to the first element.
			// We visualize the text using a list where each paragraph is a separate item.
			var vizList = layout.List{
				Axis: layout.Vertical,
				Position: layout.Position{
					Offset: int(scrollY),
				},
			}

			// ---------- MARGINS ----------
			// Margins
			var marginWidth unit.Dp
			marginWidth = (unit.Dp(gtx.Constraints.Max.X) - textWidth) / 3
			margins := layout.Inset{
				Left:   marginWidth,
				Right:  marginWidth,
				Top:    unit.Dp(0),
				Bottom: unit.Dp(0),
			}

			// ---------- LIST WITHIN MARGINS ----------
			// 1) First the margins ...
			margins.Layout(gtx,
				func(gtx C) D {
					// 2) ... then the list inside those margins ...
					return vizList.Layout(gtx, len(paragraphList),
						// 3) ... where each paragraph is a separate item
						func(gtx C, index int) D {
							// One label per paragraph
							paragraph := material.Label(th, unit.Sp(float32(fontSize)), paragraphList[index])
							// The text is centered
							paragraph.Alignment = text.Middle
							// Set color
							paragraph.Color = myColor.foreground
							// Return the laid out paragraph
							return paragraph.Layout(gtx)
						},
					)
				},
			)

			// ---------- THE FOCUS BAR ----------
			// Draw the transparent red focus bar.
			focusBar := clip.Rect{
				Min: image.Pt(0, int(focusBarY)),
				Max: image.Pt(gtx.Constraints.Max.X, int(focusBarY)+int(fontSize*1.5)),
			}.Push(&ops)
			paint.ColorOp{Color: myColor.focusbar}.Add(&ops)
			paint.PaintOp{}.Add(&ops)
			focusBar.Pop()

			// ---------- REGISTERING EVENTS ----------
			// registering events here work
			event.Op(&ops, tag)

			// ---------- FINALIZE ----------
			// Frame completes the FrameEvent by drawing the graphical operations from ops into the window.
			winE.Frame(&ops)

			// Should we shut down?
		case app.DestroyEvent:
			return winE.Err
		}
	}
}
