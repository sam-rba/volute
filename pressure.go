package main

import (
	"errors"
	"fmt"
)

type Pressure float32

const (
	Millibar            Pressure = 100
	Pascal              Pressure = 1
	Kilopascal          Pressure = 1_000
	Bar                 Pressure = 100_000
	PoundsPerSquareInch Pressure = 6_894.757
)

var PressureUnits = []string{"mbar", "Pa", "kPa", "bar", "psi"}

func ParsePressureUnit(s string) (Pressure, error) {
	// Each case corresponds to a value in PressureUnits.
	switch s {
	case "mbar":
		return Millibar, nil
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

func AtmosphericPressure() Pressure {
	return 101.325 * Pascal
}
