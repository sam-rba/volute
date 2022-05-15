package volume

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

// UnitStrings returns a slice of strings, each representing a
// unit.
// This is necessary because giu.Combo only works with strings.
func UnitStrings() []string {
	return []string{"cc", "L", "m³", "in³"}
}

const (
	DefaultUnit Volume = CubicCentimetre
	// DefaulUnitIndex is used to index UnitStrings().
	DefaultUnitIndex int32 = 0 // cc
)

func UnitFromString(s string) (Volume, error) {
	// Each case corresponds to a value in UnitStrings().
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
