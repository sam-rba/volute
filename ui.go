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
	return g.Row(
		g.Label("Engine Displacement"),
		g.InputFloat(&displacement.Val).Format("%.2f").OnChange(func() {
			for i := 0; i < numPoints; i++ {
				engineMassFlowRate[i] = massFlowRateAt(i)
				go updateCompImg()
			}
		}),
		g.Combo(
			"",
			volume.UnitStrings()[selectedVolumeUnit],
			volume.UnitStrings(),
			&selectedVolumeUnit,
		).OnChange(func() {
			s := volume.UnitStrings()[selectedVolumeUnit]
			u, err := volume.UnitFromString(s)
			util.Check(err)
			displacement = volume.Volume{
				displacement.AsUnit(u),
				u,
			}
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
			temperature.UnitStrings()[selectedTemperatureUnit],
			temperature.UnitStrings(),
			&selectedTemperatureUnit,
		).OnChange(func() {
			s := temperature.UnitStrings()[selectedTemperatureUnit]
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
	widgets := []g.Widget{
		g.Label("Manifold Absolute Pressure"),
		g.Combo(
			"",
			pressure.UnitStrings()[selectedPressureUnit],
			pressure.UnitStrings(),
			&selectedPressureUnit,
		).OnChange(func() {
			s := pressure.UnitStrings()[selectedPressureUnit]
			u, err := pressure.UnitFromString(s)
			util.Check(err)

			for i := 0; i < numPoints; i++ {
				manifoldPressure[i] = pressure.Pressure{
					manifoldPressure[i].AsUnit(u),
					u,
				}
			}
		}),
	}
	for i := 0; i < numPoints; i++ {
		i := i
		widgets = append(
			widgets,
			g.InputFloat(&manifoldPressure[i].Val).Format("%.2f").
				OnChange(func() {
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
	widgets := []g.Widget{
		g.Label("Mass Flow Rate"),
		g.Combo(
			"",
			mass.FlowRateUnitStrings()[selectedMassFlowRateUnit],
			mass.FlowRateUnitStrings(),
			&selectedMassFlowRateUnit,
		).OnChange(func() {
			s := mass.FlowRateUnitStrings()[selectedMassFlowRateUnit]
			u, err := mass.FlowRateUnitFromString(s)
			util.Check(err)

			for i := 0; i < numPoints; i++ {
				engineMassFlowRate[i] = mass.FlowRate{
					engineMassFlowRate[i].AsUnit(u),
					u,
				}
			}
		}),
	}
	for i := 0; i < numPoints; i++ {
		mfr := strconv.FormatFloat(
			float64(engineMassFlowRate[i].Val),
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
	for manufacturer := range compressors {
		manufacturerNode := g.TreeNode(manufacturer)
		for series := range compressors[manufacturer] {
			seriesNode := g.TreeNode(series)
			for model, c := range compressors[manufacturer][series] {
				seriesNode = seriesNode.Layout(
					g.Selectable(model).OnClick(func() {
						go setCompressor(c)
					}),
				)
			}
			manufacturerNode = manufacturerNode.Layout(seriesNode)
		}
		compressorTree = append(compressorTree, manufacturerNode)
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
		unit := selectedCompressor.MaxFlow.Unit
		mfr := engineMassFlowRate[i].AsUnit(unit)
		maxMfr := selectedCompressor.MaxFlow.AsUnit(unit)
		x := min + int(float32(max-min)*(mfr/maxMfr))

		min = selectedCompressor.MinY
		max = selectedCompressor.MaxY
		pr := pressureRatio[i]
		maxPr := selectedCompressor.MaxPressureRatio
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
		canvas.AddImage(
			compressorTexture,
			image.Pt(0, 250),
			image.Pt(int(winWidth), int(winHeight)),
		)
	}
}
