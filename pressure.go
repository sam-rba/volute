package main

import (
	"errors"
	"fmt"
)

type pressureUnit float32

const (
	pascal              pressureUnit = 1
	kilopascal          pressureUnit = 1_000
	bar                 pressureUnit = 100_000
	poundsPerSquareInch pressureUnit = 6_894.757
)

// pressureUnitStrings returns a slice of strings, each representing a
// pressureUnit.
// This is necessary because giu.Combo only works with strings.
func pressureUnitStrings() []string {
	return []string{"Pa", "kPa", "bar", "psi"}
}

const (
	defaultPressureUnit pressureUnit = kilopascal
	// Used to index pressureUnitStrings
	defaultPressureUnitIndex int32 = 1 // kPa
)

func pressureUnitFromString(s string) (pressureUnit, error) {
	// Each case corresponds to a value in pressureUnitStrings
	switch s {
	case "Pa":
		return pascal, nil
	case "kPa":
		return kilopascal, nil
	case "bar":
		return bar, nil
	case "psi":
		return poundsPerSquareInch, nil
	default:
		return *new(pressureUnit), errors.New(fmt.Sprintf("invalid pressureUnit: '%s'", s))
	}
}

type pressure struct {
	val  float32
	unit pressureUnit
}

func newPressure() pressure {
	return pressure{100, defaultPressureUnit}
}

func (p pressure) asUnit(u pressureUnit) float32 {
	pa := p.val * float32(p.unit) // Convert to pascals.
	return pa / float32(u)        // Convert to desired unit.
}
