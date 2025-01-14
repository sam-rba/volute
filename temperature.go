package main

import "fmt"

type TemperatureUnit int

const (
	Celcius TemperatureUnit = iota
	Kelvin
	Fahrenheit
)

type Temperature struct {
	Val  float64
	Unit TemperatureUnit
}

func (t Temperature) AsUnit(u TemperatureUnit) float64 {
	// Convert to celcius
	var c float64
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
		return c
	case Kelvin:
		return c + 272.15
	case Fahrenheit:
		return c*(9.0/5.0) + 32.0
	default:
		panic(fmt.Sprintf("invalid unit: '%v'", u))
	}
}
