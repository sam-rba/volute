package main

import (
	"fmt"
	g "github.com/AllenDang/giu"
	"os"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	manifoldPressure pressure

	// selectedPressureUnit is used to index pressureUnits
	selectedPressureUnit int32
)

func init() {
	manifoldPressure = pressure{100, defaultPressureUnit}

	// selectedPressureUnit is used to index pressureUnitStrings
	selectedPressureUnit = defaultPressureUnitIndex
}

func loop() {
	g.SingleWindow().Layout(
		g.Table().
			Rows(
				g.TableRow(
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

							manifoldPressure = pressure{
								manifoldPressure.asUnit(u),
								u,
							}
						}),
					g.InputFloat(&manifoldPressure.val).Format("%.2f"),
				),
			).
			Columns(
				g.TableColumn("Parameter"),
				g.TableColumn("Unit"),
				g.TableColumn("Point 1"),
			),
	)
}

func main() {
	wnd := g.NewMasterWindow("volute", 400, 200, 0)
	wnd.Run(loop)
}
