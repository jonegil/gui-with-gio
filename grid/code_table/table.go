package main

import (
	"encoding/csv"
	"flag"
	"image"
	"image/color"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/io/system"
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

// Command line input variables
var filename *string
var simulate *bool

func main() {
	// Step 1 - Read input from command line
	filename = flag.String("file", "example.csv", "Which .csv file shall I present? ")
	simulate = flag.Bool("simulate", false, "or should I simulate 1 million random rows")
	flag.Parse()

	// Step 2 - Read from file
	dataset := []data{}
	if *simulate {
		dataset = simulateData(1e6)
	} else {
		dataset = readCSV(filename)
	}

	// Step 3 - Start the GUI
	go func() {
		w := app.NewWindow(
			app.Title("Gio example - Table"),
			app.Size(unit.Dp(800), unit.Dp(600)),
		)
		if err := draw(w, dataset); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

type data struct {
	rowName string
	colName string
	value   float64
}

func readCSV(filename *string) []data {
	// open file
	f, err := os.Open(*filename)
	if err != nil {
		log.Fatal("Error when reading file:\n  ", err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// Create a slice of data
	dataset := []data{}

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error when parsing csv:\n  ", err)
		}
		if len(line) == 3 {
			val, err := strconv.ParseFloat(line[2], 64)
			if err != nil {
				log.Fatal("Error when converting data to float:\n  ", err)
			}
			d := data{rowName: line[0],
				colName: line[1],
				value:   val}
			dataset = append(dataset, d)
		}
	}
	return dataset
}

func simulateData(n int) []data {

	sectors := []string{
		"Basic Materials",
		"Consumer Discretionary",
		"Consumer Staples",
		"Energy",
		"Financials",
		"Health Care",
		"Industrials",
		"Real Estate",
		"Technology",
		"Telecommunications",
		"Utilities",
	}

	regions := []string{
		"Americas",
		"Europe",
		"Asia Pacific",
		"Middle East & Africa",
	}

	dataset := []data{}

	for i := 0; i <= n; i++ {
		sector := sectors[rand.Intn(len(sectors))]
		region := regions[rand.Intn(len(regions))]
		earnings := 10 * (rand.Float64() - 0.5)
		d := data{
			rowName: sector,
			colName: region,
			value:   earnings,
		}
		dataset = append(dataset, d)
	}

	return (dataset)

}

type (
	C = layout.Context
	D = layout.Dimensions
)

var colPos = color.NRGBA{0x1e, 0xb9, 0x80, 255} //rally green
var colNeg = color.NRGBA{0xff, 0x68, 0x59, 255} //rally orange
var colWhite = color.NRGBA{255, 255, 255, 255}
var colBlack = color.NRGBA{18, 18, 18, 255}

func draw(w *app.Window, dataset []data) error {
	th := material.NewTheme()
	var (
		ops  op.Ops
		grid component.GridState
	)

	// -- PART 1 -- Convert the dataset to maps for the grid

	// Convert dataset to a grid of cells, and also add sums
	cells := map[string]map[string]float64{}
	rowNames := []string{}
	colNames := []string{}
	rowSums := map[string]float64{}
	colSums := map[string]float64{}
	totSum := 0.0

	// Iterate through the whole dataset
	for _, d := range dataset {
		// Collect names
		if !slices.Contains(rowNames, d.rowName) {
			rowNames = append(rowNames, d.rowName)
		}
		if !slices.Contains(colNames, d.colName) {
			colNames = append(colNames, d.colName)
		}
		// build and populate the 2d grid
		var ok bool
		if _, ok = cells[d.rowName]; !ok {
			// Create the first map
			cells[d.rowName] = map[string]float64{}
		}
		if _, ok = cells[d.rowName][d.colName]; !ok {
			// Create the second map
			cells[d.rowName][d.colName] = 0
		}
		// Calcualte the cell value
		cells[d.rowName][d.colName] += d.value
		// Callculate rowSums and colSums
		rowSums[d.rowName] += d.value
		colSums[d.colName] += d.value
		totSum += d.value
	}

	// Add rowSums and colSums to the datasets
	for _, v := range rowNames {
		cells[v]["Total"] = rowSums[v]
	}
	for _, v := range colNames {
		if _, ok := cells["Total"]; !ok {
			cells["Total"] = map[string]float64{}
		}
		cells["Total"][v] = colSums[v]
	}
	cells["Total"]["Total"] = totSum

	// Append Total to rowNames and colNames
	rowNames = append(rowNames, "Total")
	colNames = append(colNames, "Total")

	// -- PART 2 -- Visualize the grid

	// Used for thousand separator
	printer := message.NewPrinter(language.English)

	for {

		windowevent := <-w.Events()
		switch e := windowevent.(type) {
		case system.DestroyEvent:
			return e.Err

		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
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
							value = cells[rowName][colName]
							cell.Text = printer.Sprintf("%.1f", value)
							cell.Alignment = text.End
							if value > 0 {
								cell.Color = colPos
							}
							if value < 0 {
								cell.Color = colNeg
							}
						}

						if row == len(rowNames)-1 || col == len(colNames) {
							cell.Font.Weight = font.Bold
						}

						return cell.Layout(gtx)
					})
				},
			)

			e.Frame(gtx.Ops)
		}
	}
}
