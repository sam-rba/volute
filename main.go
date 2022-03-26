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

var (
	displacement = volume{2000, cubicCentimetre}
	// selectedVolumeUnit is used to index volumeUnitStrings.
	selectedVolumeUnit = defaultVolumeUnitIndex
)

var engineSpeed = [numPoints]int32{2000, 3000, 4000, 5000, 6000, 7000}

var volumetricEfficiency = [numPoints]int32{100, 100, 100, 100, 100, 100}

var (
	intakeAirTemperature = [numPoints]temperature{
		{35, celcius},
		{35, celcius},
		{35, celcius},
		{35, celcius},
		{35, celcius},
		{35, celcius},
	}

	// selectedTemperatureUnit is used to index temperatureUnitStrings.
	selectedTemperatureUnit = defaultTemperatureUnitIndex
)

var (
	manifoldPressure = [numPoints]pressure{
		{100, defaultPressureUnit},
		{100, defaultPressureUnit},
		{100, defaultPressureUnit},
		{100, defaultPressureUnit},
		{100, defaultPressureUnit},
		{100, defaultPressureUnit},
	}

	// selectedPressureUnit is used to index pressureUnitStrings.
	selectedPressureUnit = defaultPressureUnitIndex
)

func loop() {
	g.SingleWindow().Layout(
		engineDisplacementRow(),
		g.Table().
			Rows(
				engineSpeedRow(),
				volumetricEfficiencyRow(),
				intakeAirTemperatureRow(),
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
