package main

import (
	"flag"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"gioui.org/app"
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

// Command line input variables
var filename *string

// A []string to hold the speech as a list of paragraphs
var paragraphList []string

func main() {
	// Step 1 - Read input from command line
	filename = flag.String("file", "speech.txt", "Which .txt file shall I present?")
	flag.Parse()

	// Step 2 - Read from file
	paragraphList = readText(filename)

	// Step 3 - Start the GUI
	go func() {
		// create new window
		w := app.NewWindow(
			app.Title("Teleprompter"),
			app.Size(unit.Dp(650), unit.Dp(600)),
		)
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

	// Listen for events from the window.
	for windowEvent := range w.Events() {
		switch winE := windowEvent.(type) {

		// Should we draw a new frame?
		case system.FrameEvent:

			// ---------- Handle input ----------
			// Time to deal with inputs since last frame.
			// Since we use one global eventArea, with Tag: 0
			// we here call gtx.Events(0) to get these events.

			gtx := layout.NewContext(&ops, winE)
			for _, gtxEvent := range gtx.Events(0) {

				// To set how large each change is
				var stepSize unit.Dp = 1

				switch gtxE := gtxEvent.(type) {

				// Any key
				case key.EditEvent:
					// To increase the fontsize
					if gtxE.Text == "+" {
						fontSize = fontSize + unit.Sp(stepSize)
					}
					// To decrease the fontsize
					if gtxE.Text == "-" {
						fontSize = fontSize - unit.Sp(stepSize)
					}

				// Only specified keys, defined in key.InputOp below
				case key.Event:
					// For better control, we only care about pressing the key down, not releasing it up
					if gtxE.State == key.Press {
						// Inrease the stepSize when pressing Shift
						if gtxE.Modifiers == key.ModShift {
							stepSize = 5
						}
						// Start/Stop
						if gtxE.Name == "Space" {
							autoscroll = !autoscroll
							if autospeed == 0 {
								autoscroll = true
								autospeed = 1
							}
						}
						// Scroll up
						if gtxE.Name == "K" {
							scrollY = scrollY - stepSize*4
							if scrollY < 0 {
								scrollY = 0
							}
						}
						// Scroll down
						if gtxE.Name == "J" {
							scrollY = scrollY + stepSize*4
						}
						// Faster scrollspeed
						if gtxE.Name == "F" {
							autoscroll = true
							autospeed += stepSize
						}
						// Slower scrollspeed
						if gtxE.Name == "S" {
							if autospeed > 0 {
								autospeed -= stepSize
							}
							if autospeed == 0 {
								autoscroll = false
							}
						}
						// Wider text to be displayed
						if gtxE.Name == "W" {
							textWidth = textWidth + stepSize*10
						}
						// Narrow text to be displayed
						if gtxE.Name == "N" {
							textWidth = textWidth - stepSize*10
						}
						// Move the focusBar Up
						if gtxE.Name == "U" {
							focusBarY = focusBarY - stepSize
						}
						// Move the focusBar Down
						if gtxE.Name == "D" {
							focusBarY = focusBarY + stepSize
						}
					}

				// A mouse event?
				case pointer.Event:
					// Are we scrolling?
					if gtxE.Type == pointer.Scroll {
						if gtxE.Modifiers == key.ModShift {
							stepSize = 3
						}
						// Increment scrollY with gtxE.Scroll.Y
						scrollY = scrollY + unit.Dp(gtxE.Scroll.Y)*stepSize
						if scrollY < 0 {
							scrollY = 0
						}
					}
				}
			}

			// ---------- LAYOUT ----------
			// First we layout the user interface.
			// Afterwards we add an eventArea.
			// Let's start with a background color
			paint.Fill(&ops, color.NRGBA{R: 0xff, G: 0xfe, B: 0xe0, A: 0xff})

			// ---------- THE SCROLLING TEXT ----------
			// First, check if we should autoscroll
			// That's done by increasing the value of scrollY
			if autoscroll {
				if autospeed < 0 {
					autospeed = 0
				}
				scrollY = scrollY + autospeed
				op.InvalidateOp{At: gtx.Now.Add(time.Second * 2 / 100)}.Add(&ops)
			}
			// Then we use scrollY to control the distance from the top of the screen to the first element.
			// We visualize the text using a list where each paragraph is a separate item.
			var visList = layout.List{
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
					return visList.Layout(gtx, len(paragraphList),
						// 3) ... where each paragraph is a separate item
						func(gtx C, index int) D {
							// One label per paragraph
							paragraph := material.Label(th, unit.Sp(float32(fontSize)), paragraphList[index])
							// The text is centered
							paragraph.Alignment = text.Middle
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
				Max: image.Pt(gtx.Constraints.Max.X, int(focusBarY)+int(unit.Dp(50))),
			}.Push(&ops)
			paint.ColorOp{Color: color.NRGBA{R: 0xff, A: 0x66}}.Add(&ops)
			paint.PaintOp{}.Add(&ops)
			focusBar.Pop()

			// ---------- COLLECT INPUT ----------
			// Create an eventArea to collect events. It has the same size as the full windodw.
			// First we Push() it on the stack, then add code to catch keys and pointers
			// Finally we Pop() it
			eventArea := clip.Rect(
				image.Rectangle{
					// From top left
					Min: image.Point{0, 0},
					// To bottom right
					Max: image.Point{gtx.Constraints.Max.X, gtx.Constraints.Max.Y},
				},
			).Push(gtx.Ops)

			// Since Gio is stateless we must Tag events, to make sure we know where they came from.
			// Such a tag can anything really, so we simply use Tag: 0.
			// Later we retireve these events with gtx.Events(0)

			// 1) We first add a pointer.InputOp to catch scrolling:
			pointer.InputOp{
				Types: pointer.Scroll,
				Tag:   0, // Use Tag: 0 as the event routing tag, and retireve it through gtx.Events(0)
				// ScrollBounds sets bounds on scrolling, and we want it to be non-zero.
				// In practice it seldom reached 100, so [MinInt8,MaxInt8] or [-128,127] should be enough
				ScrollBounds: image.Rectangle{
					Min: image.Point{
						X: 0,
						Y: math.MinInt8, //-128
					},
					Max: image.Point{
						X: 0,
						Y: math.MaxInt8, //+127
					},
				},
			}.Add(gtx.Ops)

			// 2) Next we add key.FocusOp,
			// Needed for general keybaord input, except the ones defined explicitly in key.InputOp
			// These inputs are retrieved as key.EditEvent
			key.FocusOp{
				Tag: 0, // Use Tag: 0 as the event routing tag, and retireve it through gtx.Events(0)
			}.Add(gtx.Ops)

			// 3) Finally we add key.InputOp to catch specific keys
			// (Shift) means an optional Shift
			// These inputs are retrieved as key.Event
			key.InputOp{
				Keys: key.Set("(Shift)-F|(Shift)-S|(Shift)-U|(Shift)-D|(Shift)-J|(Shift)-K|(Shift)-W|(Shift)-N|Space"),
				Tag:  0, // Use Tag: 0 as the event routing tag, and retireve it through gtx.Events(0)
			}.Add(gtx.Ops)

			// Finally Pop() the eventArea from the stack
			eventArea.Pop()

			// ---------- FINALIZE ----------
			// Frame completes the FrameEvent by drawing the graphical operations from ops into the window.
			winE.Frame(&ops)

		// Should we shut down?
		case system.DestroyEvent:
			return winE.Err

		}
	}
	return nil
}
