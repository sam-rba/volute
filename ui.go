package main

import g "github.com/AllenDang/giu"

func engineDisplacementRow() *g.RowWidget {
	return g.Row(
		g.Label("Engine Displacement"),
		g.InputFloat(&displacement.val).Format("%.2f"),
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
		g.InputInt(&engineSpeed[0]),
		g.InputInt(&engineSpeed[1]),
		g.InputInt(&engineSpeed[2]),
		g.InputInt(&engineSpeed[3]),
		g.InputInt(&engineSpeed[4]),
		g.InputInt(&engineSpeed[5]),
	)
}

func volumetricEfficiencyRow() *g.TableRowWidget {
	return g.TableRow(
		g.Label("Volumetric Efficiency"),
		g.Label("%"),
		g.InputInt(&volumetricEfficiency[0]),
		g.InputInt(&volumetricEfficiency[1]),
		g.InputInt(&volumetricEfficiency[2]),
		g.InputInt(&volumetricEfficiency[3]),
		g.InputInt(&volumetricEfficiency[4]),
		g.InputInt(&volumetricEfficiency[5]),
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
		g.InputFloat(&intakeAirTemperature[0].val).Format("%.2f"),
		g.InputFloat(&intakeAirTemperature[1].val).Format("%.2f"),
		g.InputFloat(&intakeAirTemperature[2].val).Format("%.2f"),
		g.InputFloat(&intakeAirTemperature[3].val).Format("%.2f"),
		g.InputFloat(&intakeAirTemperature[4].val).Format("%.2f"),
		g.InputFloat(&intakeAirTemperature[5].val).Format("%.2f"),
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

				for i := range manifoldPressure {
					manifoldPressure[i] = pressure{
						manifoldPressure[i].asUnit(u),
						u,
					}
				}
			}),
		g.InputFloat(&manifoldPressure[0].val).Format("%.2f"),
		g.InputFloat(&manifoldPressure[1].val).Format("%.2f"),
		g.InputFloat(&manifoldPressure[2].val).Format("%.2f"),
		g.InputFloat(&manifoldPressure[3].val).Format("%.2f"),
		g.InputFloat(&manifoldPressure[4].val).Format("%.2f"),
		g.InputFloat(&manifoldPressure[5].val).Format("%.2f"),
	)
}
