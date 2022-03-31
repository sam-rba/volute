package util

import (
	"fmt"
	"os"

	"github.com/sam-anthony/volute/mass"
	"github.com/sam-anthony/volute/pressure"
	"github.com/sam-anthony/volute/temperature"
)

func Check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Insert[T int32 | float32 | temperature.Temperature | pressure.Pressure | mass.FlowRate](slice []T, elem T, i int) []T {
	return append(
		slice[:i],
		append(
			[]T{elem},
			slice[i:]...,
		)...,
	)
}

func Remove[T int32 | float32 | temperature.Temperature | pressure.Pressure | mass.FlowRate](slice []T, i int) []T {
	return append(slice[:i], slice[i+1:]...)
}
