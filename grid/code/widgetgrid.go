// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
)

func main() {
	go func() {
		w := app.NewWindow()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

const sideLength int = 8
const cellSize int = 100

type (
	C = layout.Context
	D = layout.Dimensions
)

func loop(w *app.Window) error {
	th := material.NewTheme(gofont.Collection())
	var (
		ops  op.Ops
		grid component.GridState
	)

	clickers := []widget.Clickable{}
	for i := 0; i < sideLength*sideLength; i++ {
		clickers = append(clickers, widget.Clickable{})
	}

	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			component.Grid(th, &grid).Layout(gtx, sideLength, sideLength,
				func(axis layout.Axis, index, constraint int) int {
					return gtx.Dp(unit.Dp(cellSize))
				},
				func(gtx C, row, col int) D {
					clk := &clickers[row*sideLength+col]
					btn := material.Button(th, clk, fmt.Sprintf("R%d C%d", row, col))
					color := color.NRGBA{R: uint8(float32(255) / float32(sideLength) * float32(row)), G: uint8(255 / sideLength * col), B: uint8(row * col), A: 255}
					btn.Background = color
					return btn.Layout(gtx)
				})

			e.Frame(gtx.Ops)
		}
	}
}

/*
type Square struct {
	pressed bool
	hovered bool
	color   color.NRGBA
}

func (s *Square) Layout(gtx C) D {
	// here we loop through all the events associated with this button.
	for _, e := range gtx.Events(s) {
		if e, ok := e.(pointer.Event); ok {
			switch e.Type {
			case pointer.Press:
				s.pressed = true
			case pointer.Release:
				s.pressed = false
			case pointer.Enter:
				s.hovered = true
			case pointer.Leave:
				s.hovered = false
			}
		}
	}
	// Draw the button.
	//s.color = color.NRGBA{R: 0x80, A: 0xFF}
	if s.pressed {
		s.color.G = 0x80
	}
	if s.hovered {
		s.color.B = 0x80
	}

	return s.draw(gtx.Ops)
}

func (s *Square) draw(ops *op.Ops) D {
	/*defer clip.Rect{Max: image.Pt(cellSize, cellSize)}.Push(ops).Pop()
	paint.ColorOp{Color: s.color}.Add(ops)
	paint.PaintOp{}.Add(ops)

	// Define the area for pointer events.
	area := clip.Rect(image.Rect(0, 0, cellSize, cellSize)).Push(ops)
	fmt.Println(s)
	pointer.InputOp{
		Tag:   s,
		Types: pointer.Press | pointer.Release | pointer.Enter | pointer.Leave,
	}.Add(ops)
	area.Pop()


	return D{Size: image.Pt(cellSize, cellSize)}
}
*/
