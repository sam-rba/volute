package mass

import (
	"errors"
	"fmt"
)

type unit float32

const (
	Gram     unit = 1
	Kilogram unit = 1_000
	Pound    unit = 453.5924
)

type Mass struct {
	val float32
}

func New(i float32, u unit) Mass {
	return Mass{i * float32(u)}
}

func (m Mass) AsUnit(u unit) float32 {
	return m.val / float32(u)
}

type flowRateUnit float32

const (
	KilogramsPerSecond flowRateUnit = 1
	PoundsPerMinute    flowRateUnit = 0.007_559_872_833
)

// FlowRateUnitStrings returns a slice of strings, each representing a
// flowRateUnit.
// This is necessary because giu.Combo only works with strings.
func FlowRateUnitStrings() []string {
	return []string{"kg/s", "lb/min"}
}

const (
	DefaultFlowRateUnit flowRateUnit = KilogramsPerSecond
	// DefaultFlowRateUnitIndex is used to index FlowRateUnitStrings()
	DefaultFlowRateUnitIndex int32 = 0 // kg/s
)

func FlowRateUnitFromString(s string) (flowRateUnit, error) {
	// Each case corresponds to a value in FlowRateUnitStrings().
	switch s {
	case "kg/s":
		return KilogramsPerSecond, nil
	case "lb/min":
		return PoundsPerMinute, nil
	default:
		return *new(flowRateUnit), errors.New(
			fmt.Sprintf("invalid mass flow rate unit: '%s'", s))
	}
}

type FlowRate struct {
	val float32
}

func NewFlowRate(i float32, u flowRateUnit) FlowRate {
	return FlowRate{i * float32(u)}
}

func (fr FlowRate) AsUnit(u flowRateUnit) float32 {
	return fr.val / float32(u)
}
