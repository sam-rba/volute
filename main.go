package main

import (
	"fmt"
	g "github.com/AllenDang/giu"
	"image"
	"image/draw"
	_ "image/jpeg"
	"os"

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
	displacement = 2000 * volume.CubicCentimetre
	// volumeUnitIndex is used to index volume.UnitStrings().
	volumeUnitIndex = volume.DefaultUnitIndex

	engineSpeed = []int32{2000}

	volumetricEfficiency = []int32{80}

	intakeAirTemperature = []temperature.Temperature{{25, temperature.Celcius}}
	// temperatureUnitIndex is used to index temperature.UnitStrings().
	temperatureUnitIndex = temperature.DefaultUnitIndex

	manifoldPressure = []pressure.Pressure{pressure.Atmospheric()}
	// pressureUnitIndex is used to index pressure.UnitStrings().
	pressureUnitIndex = pressure.DefaultUnitIndex
)

var pressureRatio []float32

func pressureRatioAt(point int) float32 {
	u := pressure.Pascal
	m := manifoldPressure[point] / u
	a := pressure.Atmospheric() / u
	return float32(m / a)
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
	disp := float32(displacement / volume.CubicMetre)
	ve := float32(volumetricEfficiency[point]) / 100.0
	cubicMetresPerMin := (rpm / 2.0) * disp * ve

	iat, err := intakeAirTemperature[point].AsUnit(temperature.Kelvin)
	util.Check(err)
	pres := manifoldPressure[point] / pressure.Pascal
	molsPerMin := (float32(pres) * cubicMetresPerMin) / (gasConstant * iat)

	kgPerMin := molsPerMin * airMolarMass

	mfr := mass.FlowRate(kgPerMin/60.0) * mass.KilogramsPerSecond
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
	manufacturer := "garrett"
	series := "g"
	model := "25-660"
	c, ok := compressor.Compressors()[manufacturer][series][model]
	if !ok {
		fmt.Printf("compressor.Compressors()[\"%s\"][\"%s\"][\"%s\"] does not exist.\n",
			manufacturer, series, model,
		)
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
