package widget

import (
	"image"
	"image/color"
	"sync"

	"volute/gui"
	"volute/gui/layout"
)

type Node[T any] struct {
	Label    string
	Value    T
	Children []Node[T]

	expanded bool
}

func Tree[T any](trees []Node[T], r image.Rectangle, focus FocusSlave, mux *gui.Mux, wg *sync.WaitGroup) {
	defer wg.Done()

	var nodes []string
	for _, root := range trees {
		nodes = append(nodes, flatten(root, 0)...)
	}

	bounds := layout.Grid{
		Rows:        populate(make([]int, len(nodes)), 1),
		Background:  color.Gray{255},
		Gap:         1,
		Split:       layout.EvenSplit,
		SplitRows:   layout.TextRowSplit,
		Margin:      0,
		Border:      0,
		BorderColor: color.Gray{16},
		Flip:        false,
	}.Lay(r)
	for i := range nodes {
		wg.Add(1)
		go Label(nodes[i], bounds[i], mux.MakeEnv(), wg)
	}

	/*
		globalFocus := focus;
		localFocus := NewFocusMaster([]int{1, 1, 1});
		defer localFocus.Close()
	*/
	// TODO
}

func flatten[T any](root Node[T], depth int) []string {
	indent := string(populate(make([]byte, 2*depth), ' '))
	nodes := []string{indent + root.Label}
	root.expanded = true // TODO: remove me
	if root.expanded {
		for _, c := range root.Children {
			nodes = append(nodes, flatten(c, depth+1)...)
		}
	}
	return nodes
}

func populate[T any](arr []T, v T) []T {
	for i := range arr {
		arr[i] = v
	}
	return arr
}
