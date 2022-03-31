package mass

import (
	"errors"
	"fmt"
	"time"
)

type unit float32

const (
	Gram     unit = 1
	Kilogram unit = 1_000
	Pound    unit = 453.5924
)

type Mass struct {
	Val  float32
	Unit unit
}

func (m Mass) AsUnit(u unit) float32 {
	g := m.Val * float32(m.Unit) // Convert to grams.
	return g / float32(u)        // Convert to desired unit.
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
			fmt.Sprintf("invalid massFlowRateUnit: '%s'", s))
	}
}

type FlowRate struct {
	Val  float32
	Unit flowRateUnit
}

func NewFlowRate(m Mass, t time.Duration, u flowRateUnit) (FlowRate, error) {
	switch u {
	case KilogramsPerSecond:
		return FlowRate{
			m.AsUnit(Kilogram) / float32(t.Seconds()),
			u,
		}, nil
	case PoundsPerMinute:
		return FlowRate{
			m.AsUnit(Pound) / float32(t.Minutes()),
			u,
		}, nil
	default:
		return *new(FlowRate), errors.New(
			fmt.Sprintf("invalid massFlowRateUnit: '%v'", u))
	}
}

func (fr FlowRate) AsUnit(u flowRateUnit) float32 {
	kgps := fr.Val * float32(fr.Unit) // Convert to kilogramsPerSecond.
	return kgps / float32(u)          // Convert to desired unit.
}
