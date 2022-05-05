package pressure

import (
	"errors"
	"fmt"
)

type unit float32

const (
	Pascal              unit = 1
	Kilopascal          unit = 1_000
	Bar                 unit = 100_000
	PoundsPerSquareInch unit = 6_894.757
)

// UnitStrings returns a slice of strings, each representing a
// unit.
// This is necessary because giu.Combo only works with strings.
func UnitStrings() []string {
	return []string{"Pa", "kPa", "bar", "psi"}
}

const (
	DefaultUnit unit = Kilopascal
	// DefaultUnitIndex is used to index UnitStrings().
	DefaultUnitIndex int32 = 1 // kPa
)

func UnitFromString(s string) (unit, error) {
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
		return *new(unit), errors.New(fmt.Sprintf("invalid unit: '%s'", s))
	}
}

type Pressure struct {
	val float32
}

func New(i float32, u unit) Pressure {
	return Pressure{i * float32(u)}
}

func (p Pressure) AsUnit(u unit) float32 {
	return p.val / float32(u)
}

func Atmospheric() Pressure {
	return Pressure{float32(1 * Bar)}
}
