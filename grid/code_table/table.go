package main

import (
	"flag"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"golang.org/x/exp/slices"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Datastructure for the simulation
var sectors = []string{
	"Technology",
	"Telecommunications",
	"Health Care",
	"Banks",
	"Financial Services",
	"Insurance",
	"Real Estate",
	"Automobiles and Parts",
	"Consumer Products and Services",
	"Media",
	"Retail",
	"Travel and Leisure",
	"Food, Beverage and Tobacco",
	"Personal Care, Drug and Grocery Stores",
	"Construction and Materials",
	"Industrial Goods and Services",
	"Basic Resources",
	"Chemicals",
	"Energy",
	"Utilities",
}

var markets = []string{
	"United Kingdom",
	"Germany",
	"France",
	"Switzerland",
	"Netherlands",
	"Spain",
	"Italy",
	"Sweden",
	"Belgium",
	"Denmark",
	"Finland",
	"Austria",
	"Poland",
}

// Command line input variables
var n_sim *int

func main() {
	// Step 1 - Read input from command line
	n_sim = flag.Int("N", 1000, "how many new stockprices per second?")
	*n_sim = *n_sim / 10 // Divide by 10 since we sim 10 times per second
	flag.Parse()

	// Step 2 - Simulate data
	data := initData(sectors, markets)

	var dataMutex sync.RWMutex
	go func() {
		for {
			// Sim baby sim
			data = simulateData(data, *n_sim, &dataMutex)
			// Chill out
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Step 3 - Start the GUI
	go func() {
		w := new(app.Window)
		w.Option(app.Title("GIO - Table"))
		w.Option(app.Size(unit.Dp(1000), unit.Dp(700)))

		if err := draw(w, data, &dataMutex); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func initData(rowNames []string, colNames []string) *map[string]map[string]float64 {

	rowNames = append(rowNames, "Total")
	colNames = append(colNames, "Total")

	data := make(map[string]map[string]float64, len(rowNames))

	// Preallocate inner maps
	for _, row := range rowNames {
		data[row] = make(map[string]float64, len(colNames))
	}
	// Preallocate values
	for _, row := range rowNames {
		for _, col := range colNames {
			data[row][col] = 0
		}
	}

	return &data
}

func simulateData(data *map[string]map[string]float64, n int, mu *sync.RWMutex) *map[string]map[string]float64 {

	for i := 0; i <= n; i++ {
		sec := sectors[rand.Intn(len(sectors))]
		mkt := markets[rand.Intn(len(markets))]
		pnl := (rand.NormFloat64())
		/*
			if (*data)[sec] == nil {
				(*data)[sec] = map[string]float64{}
			}
		*/
		mu.Lock()
		(*data)[sec][mkt] += pnl
		(*data)["Total"][mkt] += pnl
		(*data)[sec]["Total"] += pnl
		(*data)["Total"]["Total"] += pnl
		mu.Unlock()
	}

	return data
}

type (
	C = layout.Context
	D = layout.Dimensions
)

var colPos = color.NRGBA{0x1e, 0xb9, 0x80, 255} //rally green
var colNeg = color.NRGBA{0xff, 0x68, 0x59, 255} //rally orange
var colWhite = color.NRGBA{255, 255, 255, 255}
var colBlack = color.NRGBA{18, 18, 18, 255}

func draw(w *app.Window, data *map[string]map[string]float64, mu *sync.RWMutex) error {
	th := material.NewTheme()
	var (
		ops  op.Ops
		grid component.GridState
	)

	var rowNames = []string{}
	var colNames = []string{}

	// Find all row- and colnames
	mu.Lock()
	for row, inner := range *data {
		for col := range inner {
			// Collect names
			if !slices.Contains(rowNames, row) {
				rowNames = append(rowNames, row)
			}
			if !slices.Contains(colNames, col) {
				colNames = append(colNames, col)
			}
		}
	}
	mu.Unlock()
	// Sorting is nice
	// Custom sort to put "Total" last
	sort.SliceStable(rowNames, func(i, j int) bool {
		if rowNames[i] == "Total" {
			return false
		}
		if rowNames[j] == "Total" {
			return true
		}
		return rowNames[i] < rowNames[j]
	})

	sort.SliceStable(colNames, func(i, j int) bool {
		if colNames[i] == "Total" {
			return false
		}
		if colNames[j] == "Total" {
			return true
		}
		return colNames[i] < colNames[j]
	})

	// Used for thousand separator
	printer := message.NewPrinter(language.English)

	for {

		// -- PART 2 -- Visualize the grid
		windowevent := w.Event()
		switch e := windowevent.(type) {
		case app.DestroyEvent:
			return e.Err

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			paint.ColorOp{Color: colBlack}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			inset := layout.UniformInset(unit.Dp(2))

			// Configure a label styled to be a heading
			colHead := material.Body1(th, "")
			// colHead.Font.Weight = text.Bold
			colHead.Alignment = text.End
			colHead.MaxLines = 1
			colHead.Color = colWhite

			// Configure a label styled to be a cell
			cell := material.Body1(th, "")
			// cell.Font = "Mono"
			cell.Alignment = text.End
			cell.MaxLines = 1
			cell.Color = colWhite

			// Measure the height of a heading row.
			orig := gtx.Constraints
			gtx.Constraints.Min = image.Point{}
			macro := op.Record(gtx.Ops)
			dims := inset.Layout(gtx, colHead.Layout)
			_ = macro.Stop()
			gtx.Constraints = orig

			numRows := len(rowNames)
			numCols := len(colNames)

			component.Table(th, &grid).Layout(gtx, numRows, numCols+1,
				// Dimensioner func
				func(axis layout.Axis, index, constraint int) int {
					switch axis {
					case layout.Horizontal:
						minWidth := gtx.Dp(unit.Dp(50))
						return max(int(float32(constraint)/float32(numCols+1)), minWidth)
					default:
						return dims.Size.Y
					}
				},
				// Heading func
				func(gtx C, col int) D {
					return inset.Layout(gtx, func(gtx C) D {
						colHead.Text = ""
						if col > 0 {
							colName := colNames[col-1]
							colHead.Text = colName
							colHead.Font.Weight = font.Bold
						}
						return colHead.Layout(gtx)
					})
				},
				// Cell func
				func(gtx C, row, col int) D {
					rowName := rowNames[row]
					colName := ""
					if col > 0 {
						colName = colNames[col-1]
					}

					return inset.Layout(gtx, func(gtx C) D {
						// Ensure an empty cell
						cell.Text = ""
						// Default color
						cell.Color = colWhite
						// Normal non-bold font weight
						cell.Font.Weight = font.Normal

						// Zero the value
						var value float64
						// First col is rowName
						if col == 0 {
							cell.Alignment = text.Start
							cell.Text = rowName
							cell.Font.Weight = font.Bold
						}
						// Next columns are for data
						if col >= 1 {
							mu.Lock()
							value = (*data)[rowName][colName]
							mu.Unlock()
							cell.Text = printer.Sprintf("%.1f", value)
							cell.Alignment = text.End
							if value > 0 {
								cell.Color = colPos
							}
							if value < 0 {
								cell.Color = colNeg
							}
							if math.Abs(value) < 25 {
								cell.Color.A = 25
							}
						}

						if row == len(rowNames)-1 || col == len(colNames) {
							cell.Font.Weight = font.Bold
						}

						return cell.Layout(gtx)
					})
				},
			)

			// Request a redraw after 100ms
			gtx.Execute(op.InvalidateCmd{
				At: e.Now.Add(100 * time.Millisecond),
			})

			// Draw the frame
			e.Frame(gtx.Ops)
		}

	}
}
