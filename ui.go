package main

import (
	"fmt"
	g "github.com/AllenDang/giu"
	"image"
	"image/color"
	"image/draw"
	"strconv"

	"github.com/sam-anthony/volute/compressor"
	"github.com/sam-anthony/volute/mass"
	"github.com/sam-anthony/volute/pressure"
	"github.com/sam-anthony/volute/temperature"
	"github.com/sam-anthony/volute/util"
	"github.com/sam-anthony/volute/volume"
)

func red() color.RGBA {
	return color.RGBA{255, 0, 0, 255}
}

func engineDisplacementRow() *g.RowWidget {
	s := volume.UnitStrings()[volumeUnitIndex]
	unit, err := volume.UnitFromString(s)
	util.Check(err)
	engDisp := displacement.AsUnit(unit)
	return g.Row(
		g.Label("Engine Displacement"),
		g.InputFloat(&engDisp).Format("%.2f").OnChange(func() {
			displacement = volume.New(engDisp, unit)
			for i := 0; i < numPoints; i++ {
				engineMassFlowRate[i] = massFlowRateAt(i)
				go updateCompImg()
			}
		}),
		g.Combo(
			"",
			volume.UnitStrings()[volumeUnitIndex],
			volume.UnitStrings(),
			&volumeUnitIndex,
		).OnChange(func() {
			displacement = volume.New(
				displacement.AsUnit(unit),
				unit,
			)
		}),
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
		}),
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

	widgets := []g.Widget{
		g.Label("Manifold Absolute Pressure"),
		g.Combo(
			"",
			pressure.UnitStrings()[pressureUnitIndex],
			pressure.UnitStrings(),
			&pressureUnitIndex,
		).OnChange(func() {
			for i := 0; i < numPoints; i++ {
				manifoldPressure[i] = pressure.New(
					manifoldPressure[i].AsUnit(unit),
					unit,
				)
			}
		}),
	}
	for i := 0; i < numPoints; i++ {
		i := i
		manPres := manifoldPressure[i].AsUnit(unit)
		widgets = append(
			widgets,
			g.InputFloat(&manPres).Format("%.2f").
				OnChange(func() {
					manifoldPressure[i] = pressure.New(manPres, unit)
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

	widgets := []g.Widget{
		g.Label("Mass Flow Rate"),
		g.Combo(
			"",
			mass.FlowRateUnitStrings()[selectedMassFlowRateUnit],
			mass.FlowRateUnitStrings(),
			&selectedMassFlowRateUnit,
		).OnChange(func() {
			for i := 0; i < numPoints; i++ {
				engineMassFlowRate[i] = mass.NewFlowRate(
					engineMassFlowRate[i].AsUnit(mfrUnit),
					mfrUnit,
				)
			}
		}),
	}
	for i := 0; i < numPoints; i++ {
		mfr := strconv.FormatFloat(
			float64(engineMassFlowRate[i].AsUnit(mfrUnit)),
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
	// Copy compressorImage
	b := compressorImage.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), compressorImage, b.Min, draw.Src)

	for i := 0; i < numPoints; i++ {
		min := selectedCompressor.MinX
		max := selectedCompressor.MaxX

		unit := mass.KilogramsPerSecond
		mfr := engineMassFlowRate[i].AsUnit(unit)
		maxMfr := selectedCompressor.MaxFlow.AsUnit(unit)

		x := min + int(float32(max-min)*(mfr/maxMfr))

		min = selectedCompressor.MinY
		max = selectedCompressor.MaxY

		pr := pressureRatio[i]
		maxPr := selectedCompressor.MaxPR

		y := min - int(float32((min-max))*((pr-1.0)/(maxPr-1.0)))

		ps := m.Bounds().Dx() / 100 // Point size

		draw.Draw(m,
			image.Rect(x-ps/2, y-ps/2, x+ps/2, y+ps/2),
			&image.Uniform{red()},
			image.ZP,
			draw.Src,
		)
	}

	updatedCompImg <- m
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
