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
					color := color.NRGBA{
						R: uint8(255 / sideLength * row),
						G: uint8(255 / sideLength * col),
						B: uint8(row * col),
						A: 255,
					}
					btn.Background = color
					return btn.Layout(gtx)
				})

			e.Frame(gtx.Ops)
		}
	}
}
