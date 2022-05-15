package pressure

import (
	"errors"
	"fmt"
)

type Pressure float32

const (
	Pascal              Pressure = 1
	Kilopascal          Pressure = 1_000
	Bar                 Pressure = 100_000
	PoundsPerSquareInch Pressure = 6_894.757
)

// UnitStrings returns a slice of strings, each representing a
// unit.
// This is necessary because giu.Combo only works with strings.
func UnitStrings() []string {
	return []string{"Pa", "kPa", "bar", "psi"}
}

const (
	DefaultUnit Pressure = Kilopascal
	// DefaultUnitIndex is used to index UnitStrings().
	DefaultUnitIndex int32 = 1 // kPa
)

func UnitFromString(s string) (Pressure, error) {
	// Each case corresponds to a value in UnitStrings().
	switch s {
	case "Pa":
		return Pascal, nil
	case "kPa":
		return Kilopascal, nil
	case "bar":
		return Bar, nil
	case "psi":
		return PoundsPerSquareInch, nil
	default:
		return *new(Pressure), errors.New(fmt.Sprintf("invalid unit: '%s'", s))
	}
}

func Atmospheric() Pressure {
	return 1 * Bar
}
