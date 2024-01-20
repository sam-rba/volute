package main

import (
	"errors"
	"fmt"
)

type Volume float32

const (
	CubicCentimetre Volume = 1
	Litre           Volume = 1_000
	CubicMetre      Volume = 1_000_000
	CubicInch       Volume = 16.38706
)

var VolumeUnits = []string{"cc", "L", "m³", "in³"}

func ParseVolumeUnit(s string) (Volume, error) {
	// Each case corresponds to a value in VolumeUnits.
	switch s {
	case "cc":
		return CubicCentimetre, nil
	case "L":
		return Litre, nil
	case "m³":
		return CubicMetre, nil
	case "in³":
		return CubicInch, nil
	default:
		return *new(Volume), errors.New(fmt.Sprintf("invalid volume unit: '%s'", s))
	}
}
