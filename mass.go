package main

import (
	"errors"
	"fmt"
)

type MassFlowRate float32

const (
	KilogramsPerSecond MassFlowRate = 1
	PoundsPerMinute    MassFlowRate = 0.007_559_872_833
)

var MassFlowRateUnits = []string{"kg/s", "lb/min"}

func ParseMassFlowRateUnit(s string) (MassFlowRate, error) {
	// Each case corresponds to a value in MassFlowRateUnits.
	switch s {
	case "kg/s":
		return KilogramsPerSecond, nil
	case "lb/min":
		return PoundsPerMinute, nil
	default:
		return *new(MassFlowRate), errors.New(
			fmt.Sprintf("invalid mass flow rate unit: '%s'", s))
	}
}
