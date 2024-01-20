package main

import (
	"fmt"
	g "github.com/AllenDang/giu"
	"image"
	"image/draw"
	_ "image/jpeg"
	"os"
)

const (
	gasConstant  = 8.314472
	airMolarMass = 0.0289647 // kg/mol
)

var (
	defaultDisplacement       = 2 * Litre
	defaultSpeed        int32 = 2000
	defaultVE           int32 = 80
	defaultTemperature        = Temperature{25, Celcius}
)

var (
	defaultManufacturer = "borgwarner"
	defaultSeries       = "efr"
	defaultModel        = "6258"
)

// Number of data points on the compressor map.
var numPoints = 1

var (
	displacement    = defaultDisplacement
	volumeUnitIndex int32

	// Angular crankshaft speed in RPM.
	speed = []int32{defaultSpeed}

	volumetricEfficiency = []int32{defaultVE}

	intakeAirTemperature = []Temperature{defaultTemperature}
	temperatureUnitIndex int32

	manifoldPressure  = []Pressure{AtmosphericPressure()}
	pressureUnitIndex int32
)

var pressureRatio []float32

func pressureRatioAt(point int) float32 {
	u := Pascal
	m := manifoldPressure[point] / u
	a := AtmosphericPressure() / u
	return float32(m / a)
}
func init() {
	pressureRatio = append(pressureRatio, pressureRatioAt(0))
}

var (
	massFlowRateAir       []MassFlowRate
	massFlowRateUnitIndex int32
)

func massFlowRateAt(point int) MassFlowRate {
	rpm := float32(speed[point])
	disp := float32(displacement / CubicMetre)
	ve := float32(volumetricEfficiency[point]) / 100.0
	cubicMetresPerMin := (rpm / 2.0) * disp * ve

	iat, err := intakeAirTemperature[point].AsUnit(Kelvin)
	Check(err)
	pres := manifoldPressure[point] / Pascal
	molsPerMin := (float32(pres) * cubicMetresPerMin) / (gasConstant * iat)

	kgPerMin := molsPerMin * airMolarMass

	mfr := MassFlowRate(kgPerMin/60.0) * KilogramsPerSecond
	return mfr
}
func init() {
	massFlowRateAir = append(massFlowRateAir, massFlowRateAt(0))
}

var (
	compressorImage    *image.RGBA
	compressorTexture  *g.Texture
	selectedCompressor Compressor
)

func init() {
	manufacturer := defaultManufacturer
	series := defaultSeries
	model := defaultModel
	c, ok := Compressors[manufacturer][series][model]
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

func setCompressor(c Compressor) {
	f, err := os.Open(c.FileName)
	Check(err)
	defer f.Close()

	j, _, err := image.Decode(f)
	Check(err)

	b := j.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), j, b.Min, draw.Src)

	selectedCompressor = c
	compressorImage = m

	go updateCompImg()
}

func loop() {
	g.SingleWindow().Layout(
		displacementRow(),
		g.Table().
			Size(g.Auto, 190).
			Rows(
				speedRow(),
				volumetricEfficiencyRow(),
				intakeAirTemperatureRow(),
				manifoldPressureRow(),
				pressureRatioRow(),
				massFlowRateRow(),
				duplicateDeleteRow(),
			).
			Columns(
				columns()...,
			).
			Flags(g.TableFlagsSizingFixedFit),
		selectCompressor(),
		g.Custom(compressorWidget),
	)
}
