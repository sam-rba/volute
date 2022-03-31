package temperature

import (
	"errors"
	"fmt"
)

type unit int

const (
	Celcius unit = iota
	Kelvin
	Fahrenheit
)

// UnitStrings returns a slice of strings, each representing a
// unit.
// This is necessary because giu.Combo only works with strings.
func UnitStrings() []string {
	return []string{"°C", "°K", "°F"}
}

const (
	DefaultUnit unit = Celcius
	// DefaultUnitIndex is used to index UnitStrings().
	DefaultUnitIndex int32 = 0 // celcius
)

func UnitFromString(s string) (unit, error) {
	// Each case corresponds to a value in UnitStrings().
	switch s {
	case "°C":
		return Celcius, nil
	case "°K":
		return Kelvin, nil
	case "°F":
		return Fahrenheit, nil
	default:
		return *new(unit), errors.New(fmt.Sprintf("invalid unit: '%s'", s))
	}
}

type Temperature struct {
	Val  float32
	Unit unit
}

func (t Temperature) AsUnit(u unit) (float32, error) {
	// Convert to celcius
	var c float32
	switch t.Unit {
	case Celcius:
		c = t.Val
	case Kelvin:
		c = t.Val - 272.15
	case Fahrenheit:
		c = (t.Val - 32.0) * (5.0 / 9.0)
	}

	// Convert to desired unit
	switch u {
	case Celcius:
		return c, nil
	case Kelvin:
		return c + 272.15, nil
	case Fahrenheit:
		return c*(9.0/5.0) + 32.0, nil
	default:
		return 0, errors.New(fmt.Sprintf("invalid unit: '%v'", u))
	}
}
