package volume

import (
	"errors"
	"fmt"
)

type unit float32

const (
	CubicCentimetre unit = 1
	Litre           unit = 1_000
	CubicMetre      unit = 1_000_000
	CubicInch       unit = 16.38706
)

// UnitStrings returns a slice of strings, each representing a
// unit.
// This is necessary because giu.Combo only works with strings.
func UnitStrings() []string {
	return []string{"cc", "L", "m³", "in³"}
}

const (
	DefaultUnit unit = CubicCentimetre
	// DefaulUnitIndex is used to index UnitStrings().
	DefaultUnitIndex int32 = 0 // cc
)

func UnitFromString(s string) (unit, error) {
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
		return *new(unit), errors.New(fmt.Sprintf("invalid volume unit: '%s'", s))
	}
}

type Volume struct {
	val float32
}

func New(i float32, u unit) Volume {
	return Volume{i * float32(u)}
}

func (v Volume) AsUnit(u unit) float32 {
	return v.val / float32(u)
}
