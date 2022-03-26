package main

import (
	"errors"
	"fmt"
)

type volumeUnit float32

const (
	cubicCentimetre volumeUnit = 1
	litre           volumeUnit = 1_000
	cubicMetre      volumeUnit = 1_000_000
	cubicInch       volumeUnit = 16.38706
)

// volumeUnitStrings returns a slice of strings, each representing a
// volumeUnit.
// This is necessary because giu.Combo only works with strings.
func volumeUnitStrings() []string {
	return []string{"cc", "L", "m³", "in³"}
}

const (
	defaultVolumeUnit volumeUnit = cubicCentimetre
	// Used to index volumeUnitStrings
	defaultVolumeUnitIndex int32 = 0 // cc
)

func volumeUnitFromString(s string) (volumeUnit, error) {
	// Each case corresponds to a value in volumeUnitStrings
	switch s {
	case "cc":
		return cubicCentimetre, nil
	case "L":
		return litre, nil
	case "m³":
		return cubicMetre, nil
	case "in³":
		return cubicInch, nil
	default:
		return *new(volumeUnit), errors.New(fmt.Sprintf("invalid volumeUnit: '%s'", s))
	}
}

type volume struct {
	val  float32
	unit volumeUnit
}

func (v volume) asUnit(u volumeUnit) float32 {
	cc := v.val * float32(v.unit) // Convert to cubic centimetres.
	return cc / float32(u)        // Convert to desired unit.
}
