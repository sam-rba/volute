package main

import (
	"fmt"
	"os"
)

func Check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Insert[T any](slice []T, elem T, i int) []T {
	return append(slice[:i], append([]T{elem}, slice[i:]...)...)
}

func Remove[T any](slice []T, i int) []T {
	return append(slice[:i], slice[i+1:]...)
}
