package main

import (
	g "github.com/AllenDang/giu"
	"strconv"
)

func engineDisplacementRow() *g.RowWidget {
	return g.Row(
		g.Label("Engine Displacement"),
		g.InputFloat(&displacement.val).Format("%.2f").OnChange(func() {
			for i := 0; i < numPoints; i++ {
				engineMassFlowRate[i] = massFlowRateAt(i)
			}
		}),
		g.Combo(
			"",
			volumeUnitStrings()[selectedVolumeUnit],
			volumeUnitStrings(),
			&selectedVolumeUnit,
		).OnChange(func() {
			s := volumeUnitStrings()[selectedVolumeUnit]
			u, err := volumeUnitFromString(s)
			check(err)
			displacement = volume{
				displacement.asUnit(u),
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
			temperatureUnitStrings()[selectedTemperatureUnit],
			temperatureUnitStrings(),
			&selectedTemperatureUnit,
		).OnChange(func() {
			s := temperatureUnitStrings()[selectedTemperatureUnit]
			u, err := temperatureUnitFromString(s)
			check(err)

			for i := range intakeAirTemperature {
				t, err := intakeAirTemperature[i].asUnit(u)
				check(err)
				intakeAirTemperature[i] = temperature{t, u}
			}
		}),
	}
	for i := 0; i < numPoints; i++ {
		i := i
		widgets = append(
			widgets,
			g.InputFloat(&intakeAirTemperature[i].val).
				Format("%.2f").
				OnChange(func() {
					engineMassFlowRate[i] = massFlowRateAt(i)
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
			pressureUnitStrings()[selectedPressureUnit],
			pressureUnitStrings(),
			&selectedPressureUnit,
		).OnChange(func() {
			s := pressureUnitStrings()[selectedPressureUnit]
			u, err := pressureUnitFromString(s)
			check(err)

			for i := 0; i < numPoints; i++ {
				manifoldPressure[i] = pressure{
					manifoldPressure[i].asUnit(u),
					u,
				}
			}
		}),
	}
	for i := 0; i < numPoints; i++ {
		i := i
		widgets = append(
			widgets,
			g.InputFloat(&manifoldPressure[i].val).Format("%.2f").
				OnChange(func() {
					pressureRatio[i] = pressureRatioAt(i)
					engineMassFlowRate[i] = massFlowRateAt(i)
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
			massFlowRateUnitStrings()[selectedMassFlowRateUnit],
			massFlowRateUnitStrings(),
			&selectedMassFlowRateUnit,
		).OnChange(func() {
			s := massFlowRateUnitStrings()[selectedMassFlowRateUnit]
			u, err := massFlowRateUnitFromString(s)
			check(err)

			for i := 0; i < numPoints; i++ {
				engineMassFlowRate[i] = massFlowRate{
					engineMassFlowRate[i].asUnit(u),
					u,
				}
			}
		}),
	}
	for i := 0; i < numPoints; i++ {
		mfr := strconv.FormatFloat(
			float64(engineMassFlowRate[i].val),
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
