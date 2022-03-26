package main

import g "github.com/AllenDang/giu"

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
