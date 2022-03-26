package main

import (
	"fmt"
	g "github.com/AllenDang/giu"
	"os"
)

// numPoints is the number of datapoints on the compressor map.
const numPoints = 6

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var engineSpeed = [numPoints]int32{2000, 3000, 4000, 5000, 6000, 7000}

var volumetricEfficiency = [numPoints]int32{100, 100, 100, 100, 100, 100}

var (
	manifoldPressure [numPoints]pressure

	// selectedPressureUnit is used to index pressureUnits
	selectedPressureUnit int32
)

func init() {
	manifoldPressure = [numPoints]pressure{
		newPressure(),
		newPressure(),
		newPressure(),
		newPressure(),
		newPressure(),
		newPressure(),
	}

	// selectedPressureUnit is used to index pressureUnitStrings
	selectedPressureUnit = defaultPressureUnitIndex
}

func loop() {
	g.SingleWindow().Layout(
		g.Table().
			Rows(
				engineSpeedRow(),
				volumetricEfficiencyRow(),
				manifoldPressureRow(),
			).
			Columns(
				g.TableColumn("Parameter"),
				g.TableColumn("Unit"),
				g.TableColumn("Point 1"),
				g.TableColumn("Point 2"),
				g.TableColumn("Point 3"),
				g.TableColumn("Point 4"),
				g.TableColumn("Point 5"),
				g.TableColumn("Point 6"),
			),
	)
}

func main() {
	wnd := g.NewMasterWindow("volute", 400, 200, 0)
	wnd.Run(loop)
}
