package main

import (
	"fmt"
	g "github.com/AllenDang/giu"
	"image"
	"image/draw"
	_ "image/jpeg"
	"os"
	"time"

	"github.com/sam-anthony/volute/compressor"
	"github.com/sam-anthony/volute/mass"
	"github.com/sam-anthony/volute/pressure"
	"github.com/sam-anthony/volute/temperature"
	"github.com/sam-anthony/volute/util"
	"github.com/sam-anthony/volute/volume"
)

const (
	gasConstant  = 8.314472
	airMolarMass = 0.0289647 // kg/mol
)

// numPoints is the number of datapoints on the compressor map.
var numPoints = 1

var (
	displacement = volume.Volume{2000, volume.CubicCentimetre}
	// selectedVolumeUnit is used to index volume.UnitStrings().
	selectedVolumeUnit = volume.DefaultUnitIndex

	engineSpeed = []int32{2000}

	volumetricEfficiency = []int32{80}

	intakeAirTemperature = []temperature.Temperature{{25, temperature.Celcius}}
	// selectedTemperatureUnit is used to index temperature.UnitStrings().
	selectedTemperatureUnit = temperature.DefaultUnitIndex

	manifoldPressure = []pressure.Pressure{{100, pressure.DefaultUnit}}
	// selectedPressureUnit is used to index pressure.UnitStrings().
	selectedPressureUnit = pressure.DefaultUnitIndex
)

var pressureRatio []float32

func pressureRatioAt(point int) float32 {
	u := pressure.Pascal
	m := manifoldPressure[point].AsUnit(u)
	a := pressure.Atmospheric().AsUnit(u)
	return m / a
}
func init() {
	pressureRatio = append(pressureRatio, pressureRatioAt(0))
}

var (
	engineMassFlowRate []mass.FlowRate
	// selectedMassFlowRateUnit is used to index mass.FlowRateUnitStrings().
	selectedMassFlowRateUnit = mass.DefaultFlowRateUnitIndex
)

func massFlowRateAt(point int) mass.FlowRate {
	rpm := float32(engineSpeed[point])
	disp := displacement.AsUnit(volume.CubicMetre)
	ve := float32(volumetricEfficiency[point]) / 100.0
	cubicMetresPerMin := (rpm / 2.0) * disp * ve

	iat, err := intakeAirTemperature[point].AsUnit(temperature.Kelvin)
	util.Check(err)
	pres := manifoldPressure[point].AsUnit(pressure.Pascal)
	molsPerMin := (pres * cubicMetresPerMin) / (gasConstant * iat)

	kgPerMin := molsPerMin * airMolarMass

	massPerMin := mass.Mass{kgPerMin, mass.Kilogram}

	u, err := mass.FlowRateUnitFromString(mass.FlowRateUnitStrings()[selectedMassFlowRateUnit])
	util.Check(err)

	mfr, err := mass.NewFlowRate(massPerMin, time.Minute, u)
	util.Check(err)
	return mfr
}
func init() {
	engineMassFlowRate = append(engineMassFlowRate, massFlowRateAt(0))
}

func loop() {
	g.SingleWindow().Layout(
		engineDisplacementRow(),
		g.Table().
			Size(g.Auto, 190).
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
		selectCompressor(),
		g.Custom(compressorWidget),
	)
}

var (
	compressorImage    *image.RGBA
	compressorTexture  *g.Texture
	selectedCompressor compressor.Compressor
)

func setCompressor(c compressor.Compressor) {
	f, err := os.Open(c.FileName)
	util.Check(err)
	defer f.Close()

	j, _, err := image.Decode(f)
	util.Check(err)

	b := j.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), j, b.Min, draw.Src)

	selectedCompressor = c
	compressorImage = m

	go updateCompImg()
}

func init() {
	c, ok := compressor.Compressors()["Garrett"]["G"]["25-660"]
	if !ok {
		fmt.Println("Garrett G25-660 not in compressor.Compressors().")
		os.Exit(1)
	}

	setCompressor(c)
}

func main() {
	wnd := g.NewMasterWindow("volute", 400, 200, 0)

	go updateCompImg()
	m := <-updatedCompImg
	g.EnqueueNewTextureFromRgba(m, func(tex *g.Texture) {
		compressorTexture = tex
	})

	wnd.Run(loop)
}
