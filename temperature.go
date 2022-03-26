package main

import (
	"errors"
	"fmt"
)

type temperatureUnit int

const (
	celcius temperatureUnit = iota
	kelvin
	fahrenheit
)

// temperatureUnitStrings returns a slice of strings, each representing a
// temperatureUnit.
// This is necessary because giu.Combo only works with strings.
func temperatureUnitStrings() []string {
	return []string{"°C", "°K", "°F"}
}

const (
	defaultTemperatureUnit temperatureUnit = celcius
	// Used to index temperatureUnitStrings
	defaultTemperatureUnitIndex int32 = 0 // celcius
)

func temperatureUnitFromString(s string) (temperatureUnit, error) {
	// Each case corresponds to a value in volumeUnitStrings.
	switch s {
	case "°C":
		return celcius, nil
	case "°K":
		return kelvin, nil
	case "°F":
		return fahrenheit, nil
	default:
		return *new(temperatureUnit), errors.New(fmt.Sprintf("invalid temperatureUnit: '%s'", s))
	}
}

type temperature struct {
	val  float32
	unit temperatureUnit
}

func (t temperature) asUnit(u temperatureUnit) (float32, error) {
	// Convert to celcius
	var c float32
	switch t.unit {
	case celcius:
		c = t.val
	case kelvin:
		c = t.val - 272.15
	case fahrenheit:
		c = (t.val - 32.0) * (5.0 / 9.0)
	}

	// Convert to desired unit
	switch u {
	case celcius:
		return c, nil
	case kelvin:
		return c + 272.15, nil
	case fahrenheit:
		return c*(9.0/5.0) + 32.0, nil
	default:
		return 0, errors.New(fmt.Sprintf("invalid temperatureUnit: '%v'", u))
	}
}
