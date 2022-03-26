package main

import (
	"errors"
	"fmt"
	"time"
)

type massUnit float32

const (
	gram     massUnit = 1
	kilogram massUnit = 1_000
	pound    massUnit = 453.5924
)

type mass struct {
	val  float32
	unit massUnit
}

func (m mass) asUnit(u massUnit) float32 {
	g := m.val * float32(m.unit) // Convert to grams.
	return g / float32(u)        // Convert to desired unit.
}

type massFlowRateUnit float32

const (
	kilogramsPerSecond massFlowRateUnit = 1
	poundsPerMinute    massFlowRateUnit = 0.007_559_872_833
)

// massFlowRateUnitStrings returns a slice of strings, each representing a
// massFlowRateUnit.
// This is necessary because giu.Combo only works with strings.
func massFlowRateUnitStrings() []string {
	return []string{"kg/s", "lb/min"}
}

const (
	defaultMassFlowRateUnit massFlowRateUnit = kilogramsPerSecond
	//Used to index massFlowRateUnitStrings
	defaultMassFlowRateUnitIndex int32 = 0 // kg/s
)

func massFlowRateUnitFromString(s string) (massFlowRateUnit, error) {
	// Each case corresponds to a value in massFlowRateUnitStrings.
	switch s {
	case "kg/s":
		return kilogramsPerSecond, nil
	case "lb/min":
		return poundsPerMinute, nil
	default:
		return *new(massFlowRateUnit), errors.New(
			fmt.Sprintf("invalid massFlowRateUnit: '%s'", s))
	}
}

type massFlowRate struct {
	val  float32
	unit massFlowRateUnit
}

func newMassFlowRate(m mass, t time.Duration, u massFlowRateUnit) (massFlowRate, error) {
	switch u {
	case kilogramsPerSecond:
		return massFlowRate{
			m.asUnit(kilogram) / float32(t.Seconds()),
			u,
		}, nil
	case poundsPerMinute:
		return massFlowRate{
			m.asUnit(pound) / float32(t.Minutes()),
			u,
		}, nil
	default:
		return *new(massFlowRate), errors.New(
			fmt.Sprintf("invalid massFlowRateUnit: '%v'", u))
	}
}

func (fr massFlowRate) asUnit(u massFlowRateUnit) float32 {
	kgps := fr.val * float32(fr.unit) // Convert to kilogramsPerSecond.
	return kgps / float32(u)          // Convert to desired unit.
}
