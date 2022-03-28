package main

import (
	g "github.com/AllenDang/giu"
	"time"
)

const (
	gasConstant  = 8.314472
	airMolarMass = 0.0289647 // kg/mol
)

// numPoints is the number of datapoints on the compressor map.
var numPoints = 1

var (
	displacement = volume{2000, cubicCentimetre}
	// selectedVolumeUnit is used to index volumeUnitStrings.
	selectedVolumeUnit = defaultVolumeUnitIndex

	engineSpeed = []int32{2000}

	volumetricEfficiency = []int32{80}

	intakeAirTemperature = []temperature{{25, celcius}}
	// selectedTemperatureUnit is used to index temperatureUnitStrings.
	selectedTemperatureUnit = defaultTemperatureUnitIndex

	manifoldPressure = []pressure{{100, defaultPressureUnit}}
	// selectedPressureUnit is used to index pressureUnitStrings.
	selectedPressureUnit = defaultPressureUnitIndex
)

var pressureRatio []float32

func pressureRatioAt(point int) float32 {
	u := pascal
	m := manifoldPressure[point].asUnit(u)
	a := atmosphericPressure().asUnit(u)
	return m / a
}
func init() {
	pressureRatio = append(pressureRatio, pressureRatioAt(0))
}

var (
	engineMassFlowRate []massFlowRate
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
	engineMassFlowRate = append(engineMassFlowRate, massFlowRateAt(0))
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
				duplicateDeleteRow(),
			).
			Columns(
				columns()...,
			),
	)
}

func main() {
	wnd := g.NewMasterWindow("volute", 400, 200, 0)
	wnd.Run(loop)
}
