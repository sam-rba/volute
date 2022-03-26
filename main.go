package main

import (
	"fmt"
	g "github.com/AllenDang/giu"
	"os"
	"time"
)

const (
	// numPoints is the number of datapoints on the compressor map.
	numPoints = 6

	gasConstant  = 8.314472
	airMolarMass = 0.0289647 // kg/mol
)

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

var pressureRatio [numPoints]float32

func pressureRatioAt(point int) float32 {
	u := pascal
	m := manifoldPressure[point].asUnit(u)
	a := atmosphericPressure().asUnit(u)
	return m / a
}

func init() {
	for i := 0; i < numPoints; i++ {
		pressureRatio[i] = pressureRatioAt(i)
	}
}

var (
	engineMassFlowRate [numPoints]massFlowRate

	// selectedMassFlowRateUnit is used to index massFlowRateUnitStrings.
	selectedMassFlowRateUnit = defaultMassFlowRateUnitIndex
)

func massFlowRateAt(point int) massFlowRate {
	rpm := float32(engineSpeed[point])
	disp := displacement.asUnit(cubicMetre)
	ve := float32(volumetricEfficiency[point]) / 100.0
	cubicMetresPerMin := (rpm / 2.0) * disp * ve

	iat, err := intakeAirTemperature[point].asUnit(kelvin)
	check(err)
	pres := manifoldPressure[point].asUnit(pascal)
	molsPerMin := (pres * cubicMetresPerMin) / (gasConstant * iat)

	kgPerMin := molsPerMin * airMolarMass

	massPerMin := mass{kgPerMin, kilogram}

	u, err := massFlowRateUnitFromString(massFlowRateUnitStrings()[selectedMassFlowRateUnit])
	check(err)

	mfr, err := newMassFlowRate(massPerMin, time.Minute, u)
	check(err)
	return mfr
}

func init() {
	for i := 0; i < numPoints; i++ {
		engineMassFlowRate[i] = massFlowRateAt(i)
	}
}

func loop() {
	g.SingleWindow().Layout(
		engineDisplacementRow(),
		g.Table().
			Rows(
				engineSpeedRow(),
				volumetricEfficiencyRow(),
				intakeAirTemperatureRow(),
				manifoldPressureRow(),
				pressureRatioRow(),
				massFlowRateRow(),
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
