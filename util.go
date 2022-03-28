package main

import (
	"fmt"
	"os"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func insert[T int32 | float32 | temperature | pressure | massFlowRate](slice []T, elem T, i int) []T {
	return append(
		slice[:i],
		append(
			[]T{elem},
			slice[i:]...,
		)...,
	)
}

func remove[T int32 | float32 | temperature | pressure | massFlowRate](slice []T, i int) []T {
	return append(slice[:i], slice[i+1:]...)
}
