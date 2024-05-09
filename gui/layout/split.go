package layout

import (
	"fmt"

	"volute/gui/text"
)

// SplitFunc represents a way to split a space among a number of elements.
// The length of the returned slice must be equal to the number of elements.
// The sum of all elements of the returned slice must be eqal to the space.
type SplitFunc func(elements int, space int) []int

var _ SplitFunc = EvenSplit

// EvenSplit is a SplitFunc used to split a space (almost) evenly among the elements.
// It is almost evenly because width may not be divisible by elements.
func EvenSplit(elements int, width int) []int {
	if elements <= 0 {
		panic(fmt.Errorf("EvenSplit: elements must be greater than 0"))
	}
	ret := make([]int, elements)
	for elements > 0 {
		v := width / elements
		width -= v
		elements -= 1
		ret[elements] = v
	}
	return ret
}

func TextRowSplit(elements int, space int) []int {
	bounds := make([]int, elements)
	height := text.Size("1").Y
	for i := 0; i < elements && space > 0; i++ {
		bounds[i] = min(height, space)
		space -= bounds[i]
	}
	return bounds
}
