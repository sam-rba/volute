package mass

import (
	"errors"
	"fmt"
)

type Mass float32

const (
	Gram     Mass = 1
	Kilogram Mass = 1_000
	Pound    Mass = 453.5924
)

type FlowRate float32

const (
	KilogramsPerSecond FlowRate = 1
	PoundsPerMinute    FlowRate = 0.007_559_872_833
)

// FlowRateUnitStrings returns a slice of strings, each representing a
// flowRate unit.
// This is necessary because giu.Combo only works with strings.
func FlowRateUnitStrings() []string {
	return []string{"kg/s", "lb/min"}
}

const (
	DefaultFlowRateUnit FlowRate = KilogramsPerSecond
	// DefaultFlowRateUnitIndex is used to index FlowRateUnitStrings()
	DefaultFlowRateUnitIndex int32 = 0 // kg/s
)

func FlowRateUnitFromString(s string) (FlowRate, error) {
	// Each case corresponds to a value in FlowRateUnitStrings().
	switch s {
	case "kg/s":
		return KilogramsPerSecond, nil
	case "lb/min":
		return PoundsPerMinute, nil
	default:
		return *new(FlowRate), errors.New(
			fmt.Sprintf("invalid mass flow rate unit: '%s'", s))
	}
}
