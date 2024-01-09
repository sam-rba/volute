package main

import (
	"fmt"
	g "github.com/AllenDang/giu"
	"image"
	"image/color"
	"image/draw"
	"strconv"

	"volute/compressor"
	"volute/mass"
	"volute/pressure"
	"volute/temperature"
	"volute/util"
	"volute/volume"
)

func red() color.RGBA {
	return color.RGBA{255, 0, 0, 255}
}

func engineDisplacementRow() *g.RowWidget {
	s := volume.UnitStrings()[volumeUnitIndex]
	unit, err := volume.UnitFromString(s)
	util.Check(err)
	engDisp := float32(displacement / unit)
	valWid, _ := g.CalcTextSize("12345.67")
	unitWid, _ := g.CalcTextSize(volume.UnitStrings()[volumeUnitIndex])
	return g.Row(
		g.Label("Engine Displacement"),
		g.InputFloat(&engDisp).
			Format("%.2f").
			OnChange(func() {
				displacement = volume.Volume(engDisp) * unit
				for i := 0; i < numPoints; i++ {
					engineMassFlowRate[i] = massFlowRateAt(i)
					go updateCompImg()
				}
			}).
			Size(valWid),
		g.Combo(
			"",
			volume.UnitStrings()[volumeUnitIndex],
			volume.UnitStrings(),
			&volumeUnitIndex,
		).Size(unitWid*2),
	)
}

func engineSpeedRow() *g.TableRowWidget {
	widgets := []g.Widget{
		g.Label("Engine Speed"),
		g.Label("rpm"),
	}
	for i := 0; i < numPoints; i++ {
		i := i
		widgets = append(
			widgets,
			g.InputInt(&engineSpeed[i]).OnChange(func() {
				engineMassFlowRate[i] = massFlowRateAt(i)
				go updateCompImg()
			}),
		)
	}
	return g.TableRow(widgets...)
}

func volumetricEfficiencyRow() *g.TableRowWidget {
	widgets := []g.Widget{
		g.Label("Volumetric Efficiency"),
		g.Label("%"),
	}
	for i := 0; i < numPoints; i++ {
		i := i
		widgets = append(
			widgets,
			g.InputInt(&volumetricEfficiency[i]).OnChange(func() {
				engineMassFlowRate[i] = massFlowRateAt(i)
				go updateCompImg()
			}),
		)
	}
	return g.TableRow(widgets...)
}

func intakeAirTemperatureRow() *g.TableRowWidget {
	wid, _ := g.CalcTextSize(temperature.UnitStrings()[temperatureUnitIndex])
	widgets := []g.Widget{
		g.Label("Intake Air Temperature"),
		g.Combo(
			"",
			temperature.UnitStrings()[temperatureUnitIndex],
			temperature.UnitStrings(),
			&temperatureUnitIndex,
		).OnChange(func() {
			s := temperature.UnitStrings()[temperatureUnitIndex]
			u, err := temperature.UnitFromString(s)
			util.Check(err)

			for i := range intakeAirTemperature {
				t, err := intakeAirTemperature[i].AsUnit(u)
				util.Check(err)
				intakeAirTemperature[i] = temperature.Temperature{t, u}
			}
		}).Size(wid * 2),
	}
	for i := 0; i < numPoints; i++ {
		i := i
		widgets = append(
			widgets,
			g.InputFloat(&intakeAirTemperature[i].Val).
				Format("%.2f").
				OnChange(func() {
					engineMassFlowRate[i] = massFlowRateAt(i)
					go updateCompImg()
				}),
		)
	}
	return g.TableRow(widgets...)
}

func manifoldPressureRow() *g.TableRowWidget {
	s := pressure.UnitStrings()[pressureUnitIndex]
	unit, err := pressure.UnitFromString(s)
	util.Check(err)
	wid, _ := g.CalcTextSize(pressure.UnitStrings()[pressureUnitIndex])
	widgets := []g.Widget{
		g.Label("Manifold Absolute Pressure"),
		g.Combo(
			"",
			pressure.UnitStrings()[pressureUnitIndex],
			pressure.UnitStrings(),
			&pressureUnitIndex,
		).Size(wid * 2),
	}
	for i := 0; i < numPoints; i++ {
		i := i
		manPres := float32(manifoldPressure[i] / unit)
		widgets = append(
			widgets,
			g.InputFloat(&manPres).Format("%.2f").
				OnChange(func() {
					manifoldPressure[i] = pressure.Pressure(manPres * float32(unit))
					pressureRatio[i] = pressureRatioAt(i)
					engineMassFlowRate[i] = massFlowRateAt(i)
					go updateCompImg()
				}),
		)
	}
	return g.TableRow(widgets...)
}

func pressureRatioRow() *g.TableRowWidget {
	widgets := []g.Widget{
		g.Label("Pressure Ratio"),
		g.Label(""),
	}
	for i := 0; i < numPoints; i++ {
		pr := strconv.FormatFloat(float64(pressureRatio[i]), 'f', 1, 32)
		widgets = append(
			widgets,
			g.Label(pr),
		)
	}
	return g.TableRow(widgets...)
}

func massFlowRateRow() *g.TableRowWidget {
	s := mass.FlowRateUnitStrings()[selectedMassFlowRateUnit]
	mfrUnit, err := mass.FlowRateUnitFromString(s)
	util.Check(err)

	wid, _ := g.CalcTextSize(mass.FlowRateUnitStrings()[selectedMassFlowRateUnit])
	widgets := []g.Widget{
		g.Label("Mass Flow Rate"),
		g.Combo(
			"",
			mass.FlowRateUnitStrings()[selectedMassFlowRateUnit],
			mass.FlowRateUnitStrings(),
			&selectedMassFlowRateUnit,
		).Size(wid * 2),
	}
	for i := 0; i < numPoints; i++ {
		mfr := strconv.FormatFloat(
			float64(engineMassFlowRate[i]/mfrUnit),
			'f',
			3,
			32,
		)
		widgets = append(
			widgets,
			g.Label(mfr),
		)
	}
	return g.TableRow(widgets...)
}

func duplicateDeleteRow() *g.TableRowWidget {
	widgets := []g.Widget{g.Label(""), g.Label("")}
	for i := 0; i < numPoints; i++ {
		i := i
		widgets = append(widgets, g.Row(
			g.Button("Duplicate").OnClick(func() {
				numPoints++
				engineSpeed = util.Insert(
					engineSpeed,
					engineSpeed[i],
					i,
				)
				volumetricEfficiency = util.Insert(
					volumetricEfficiency,
					volumetricEfficiency[i],
					i,
				)
				intakeAirTemperature = util.Insert(
					intakeAirTemperature,
					intakeAirTemperature[i],
					i,
				)
				manifoldPressure = util.Insert(
					manifoldPressure,
					manifoldPressure[i],
					i,
				)
				pressureRatio = util.Insert(
					pressureRatio,
					pressureRatio[i],
					i,
				)
				engineMassFlowRate = util.Insert(
					engineMassFlowRate,
					engineMassFlowRate[i],
					i,
				)
				go updateCompImg()
			}),
			g.Button("Delete").OnClick(func() {
				if numPoints < 2 {
					return
				}
				numPoints--
				engineSpeed = util.Remove(engineSpeed, i)
				volumetricEfficiency = util.Remove(volumetricEfficiency, i)
				intakeAirTemperature = util.Remove(intakeAirTemperature, i)
				manifoldPressure = util.Remove(manifoldPressure, i)
				pressureRatio = util.Remove(pressureRatio, i)
				engineMassFlowRate = util.Remove(engineMassFlowRate, i)
				go updateCompImg()
			}),
		))
	}
	return g.TableRow(widgets...)
}

func columns() []*g.TableColumnWidget {
	widgets := []*g.TableColumnWidget{
		g.TableColumn("Parameter"),
		g.TableColumn("Unit"),
	}
	for i := 0; i < numPoints; i++ {
		widgets = append(
			widgets,
			g.TableColumn(fmt.Sprintf("Point %d", i+1)),
		)
	}
	return widgets
}

var compressorTree []g.Widget

func init() {
	compressors := compressor.Compressors()
	for man := range compressors {
		man := man // Manufacturer
		var serNodes []g.Widget
		for ser := range compressors[man] {
			ser := ser // Series
			var modNodes []g.Widget
			for mod, c := range compressors[man][ser] {
				mod := mod // Model
				c := c     // Compressor
				modNodes = append(
					modNodes,
					g.Selectable(mod).OnClick(func() {
						go setCompressor(c)
					}),
				)
			}
			serNodes = append(
				serNodes,
				g.TreeNode(ser).Layout(modNodes...),
			)
		}
		manNode := g.TreeNode(man).Layout(serNodes...)
		compressorTree = append(compressorTree, manNode)
	}
}

func selectCompressor() g.Widget {
	return g.ComboCustom("Compressor", selectedCompressor.Name).
		Layout(compressorTree...)
}

var updatedCompImg = make(chan image.Image)

func updateCompImg() {
	img := copyImage(compressorImage)
	for i := 0; i < numPoints; i++ {
		pos := pointPos(i)
		ps := img.Bounds().Dx() / 100 // Point size
		draw.Draw(img,
			image.Rect(pos.X-ps/2, pos.Y-ps/2, pos.X+ps/2, pos.Y+ps/2),
			&image.Uniform{red()},
			image.ZP,
			draw.Src,
		)
	}
	updatedCompImg <- img
}

func compressorWidget() {
	select {
	case m := <-updatedCompImg:
		g.EnqueueNewTextureFromRgba(m, func(tex *g.Texture) {
			compressorTexture = tex
		})
	default:
	}

	canvas := g.GetCanvas()
	if compressorTexture != nil {
		winWidth, winHeight := g.GetAvailableRegion()

		bounds := compressorImage.Bounds()
		imWidth := float32(bounds.Dx())
		imHeight := float32(bounds.Dy())

		var ratio, xratio, yratio float32
		xratio = winWidth / imWidth
		yratio = winHeight / imHeight
		if xratio < yratio {
			ratio = xratio
		} else {
			ratio = yratio
		}

		x := int(imWidth * ratio)
		y := int(imHeight * ratio)

		canvas.AddImage(
			compressorTexture,
			image.Pt(0, 250),
			image.Pt(x, y),
		)
	}
}

func copyImage(old *image.RGBA) *image.RGBA {
	b := old.Bounds()
	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(img, img.Bounds(), old, b.Min, draw.Src)
	return img
}

// The position on the compressor map of an operating point.
func pointPos(i int) (pos image.Point) {
	const unit = mass.KilogramsPerSecond
	mfr := engineMassFlowRate[i] / unit
	maxMfr := selectedCompressor.MaxFlow / unit
	min := selectedCompressor.MinX
	max := selectedCompressor.MaxX
	pos.X = min + int(float32(max-min)*float32(mfr/maxMfr))

	min = selectedCompressor.MinY
	max = selectedCompressor.MaxY
	pr := pressureRatio[i]
	maxPr := selectedCompressor.MaxPR
	pos.Y = min - int(float32((min-max))*((pr-1.0)/(maxPr-1.0)))
	return pos
}
