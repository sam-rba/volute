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
		).
			OnChange(func() {
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
	return g.TableRow(
		g.Label("Engine Speed"),
		g.Label("rpm"),
		g.InputInt(&engineSpeed[0]).OnChange(func() {
			engineMassFlowRate[0] = massFlowRateAt(0)
		}),
		g.InputInt(&engineSpeed[1]).OnChange(func() {
			engineMassFlowRate[1] = massFlowRateAt(1)
		}),
		g.InputInt(&engineSpeed[2]).OnChange(func() {
			engineMassFlowRate[2] = massFlowRateAt(2)
		}),
		g.InputInt(&engineSpeed[3]).OnChange(func() {
			engineMassFlowRate[3] = massFlowRateAt(3)
		}),
		g.InputInt(&engineSpeed[4]).OnChange(func() {
			engineMassFlowRate[4] = massFlowRateAt(4)
		}),
		g.InputInt(&engineSpeed[5]).OnChange(func() {
			engineMassFlowRate[5] = massFlowRateAt(5)
		}),
	)
}

func volumetricEfficiencyRow() *g.TableRowWidget {
	return g.TableRow(
		g.Label("Volumetric Efficiency"),
		g.Label("%"),
		g.InputInt(&volumetricEfficiency[0]).OnChange(func() {
			engineMassFlowRate[0] = massFlowRateAt(0)
		}),
		g.InputInt(&volumetricEfficiency[1]).OnChange(func() {
			engineMassFlowRate[1] = massFlowRateAt(1)
		}),
		g.InputInt(&volumetricEfficiency[2]).OnChange(func() {
			engineMassFlowRate[2] = massFlowRateAt(2)
		}),
		g.InputInt(&volumetricEfficiency[3]).OnChange(func() {
			engineMassFlowRate[3] = massFlowRateAt(3)
		}),
		g.InputInt(&volumetricEfficiency[4]).OnChange(func() {
			engineMassFlowRate[4] = massFlowRateAt(4)
		}),
		g.InputInt(&volumetricEfficiency[5]).OnChange(func() {
			engineMassFlowRate[5] = massFlowRateAt(5)
		}),
	)
}

func intakeAirTemperatureRow() *g.TableRowWidget {
	return g.TableRow(
		g.Label("Intake Air Temperature"),
		g.Combo(
			"",
			temperatureUnitStrings()[selectedTemperatureUnit],
			temperatureUnitStrings(),
			&selectedTemperatureUnit,
		).
			OnChange(func() {
				s := temperatureUnitStrings()[selectedTemperatureUnit]
				u, err := temperatureUnitFromString(s)
				check(err)

				for i := range intakeAirTemperature {
					t, err := intakeAirTemperature[i].asUnit(u)
					check(err)
					intakeAirTemperature[i] = temperature{t, u}
				}
			}),
		g.InputFloat(&intakeAirTemperature[0].val).Format("%.2f").
			OnChange(func() {
				engineMassFlowRate[0] = massFlowRateAt(0)
			}),
		g.InputFloat(&intakeAirTemperature[1].val).Format("%.2f").
			OnChange(func() {
				engineMassFlowRate[1] = massFlowRateAt(1)
			}),
		g.InputFloat(&intakeAirTemperature[2].val).Format("%.2f").
			OnChange(func() {
				engineMassFlowRate[2] = massFlowRateAt(2)
			}),
		g.InputFloat(&intakeAirTemperature[3].val).Format("%.2f").
			OnChange(func() {
				engineMassFlowRate[3] = massFlowRateAt(3)
			}),
		g.InputFloat(&intakeAirTemperature[4].val).Format("%.2f").
			OnChange(func() {
				engineMassFlowRate[4] = massFlowRateAt(4)
			}),
		g.InputFloat(&intakeAirTemperature[5].val).Format("%.2f").
			OnChange(func() {
				engineMassFlowRate[5] = massFlowRateAt(5)
			}),
	)
}

func manifoldPressureRow() *g.TableRowWidget {
	return g.TableRow(
		g.Label("Manifold Absolute Pressure"),
		g.Combo(
			"",
			pressureUnitStrings()[selectedPressureUnit],
			pressureUnitStrings(),
			&selectedPressureUnit,
		).
			OnChange(func() {
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
		g.InputFloat(&manifoldPressure[0].val).Format("%.2f").
			OnChange(func() {
				pressureRatio[0] = pressureRatioAt(0)
				engineMassFlowRate[0] = massFlowRateAt(0)
			}),
		g.InputFloat(&manifoldPressure[1].val).Format("%.2f").
			OnChange(func() {
				pressureRatio[1] = pressureRatioAt(1)
				engineMassFlowRate[1] = massFlowRateAt(1)
			}),
		g.InputFloat(&manifoldPressure[2].val).Format("%.2f").
			OnChange(func() {
				pressureRatio[2] = pressureRatioAt(2)
				engineMassFlowRate[2] = massFlowRateAt(2)
			}),
		g.InputFloat(&manifoldPressure[3].val).Format("%.2f").
			OnChange(func() {
				pressureRatio[3] = pressureRatioAt(3)
				engineMassFlowRate[3] = massFlowRateAt(3)
			}),
		g.InputFloat(&manifoldPressure[4].val).Format("%.2f").
			OnChange(func() {
				pressureRatio[4] = pressureRatioAt(4)
				engineMassFlowRate[4] = massFlowRateAt(4)
			}),
		g.InputFloat(&manifoldPressure[5].val).Format("%.2f").
			OnChange(func() {
				pressureRatio[5] = pressureRatioAt(5)
				engineMassFlowRate[5] = massFlowRateAt(5)
			}),
	)
}

func pressureRatioRow() *g.TableRowWidget {
	return g.TableRow(
		g.Label("Pressure Ratio"),
		g.Label(""),
		g.Label(strconv.FormatFloat(float64(pressureRatio[0]), 'f', 1, 32)),
		g.Label(strconv.FormatFloat(float64(pressureRatio[1]), 'f', 1, 32)),
		g.Label(strconv.FormatFloat(float64(pressureRatio[2]), 'f', 1, 32)),
		g.Label(strconv.FormatFloat(float64(pressureRatio[3]), 'f', 1, 32)),
		g.Label(strconv.FormatFloat(float64(pressureRatio[4]), 'f', 1, 32)),
		g.Label(strconv.FormatFloat(float64(pressureRatio[5]), 'f', 1, 32)),
	)
}

func massFlowRateRow() *g.TableRowWidget {
	return g.TableRow(
		g.Label("Mass Flow Rate"),
		g.Combo(
			"",
			massFlowRateUnitStrings()[selectedMassFlowRateUnit],
			massFlowRateUnitStrings(),
			&selectedMassFlowRateUnit,
		).
			OnChange(func() {
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
		g.Label(strconv.FormatFloat(float64(engineMassFlowRate[0].val), 'f', 3, 32)),
		g.Label(strconv.FormatFloat(float64(engineMassFlowRate[1].val), 'f', 3, 32)),
		g.Label(strconv.FormatFloat(float64(engineMassFlowRate[2].val), 'f', 3, 32)),
		g.Label(strconv.FormatFloat(float64(engineMassFlowRate[3].val), 'f', 3, 32)),
		g.Label(strconv.FormatFloat(float64(engineMassFlowRate[4].val), 'f', 3, 32)),
		g.Label(strconv.FormatFloat(float64(engineMassFlowRate[5].val), 'f', 3, 32)),
	)
}
