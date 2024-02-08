package main

import "image"

type Focus struct {
	widgets [][]chan bool
	p       image.Point // currently focused widget
}

func NewFocus(rows []int) Focus {
	f := Focus{
		make([][]chan bool, len(rows)),
		image.Point{},
	}
	for i := range f.widgets {
		f.widgets[i] = make([]chan bool, rows[i])
		for j := range f.widgets[i] {
			f.widgets[i][j] = make(chan bool)
		}
	}
	return f
}

func (f *Focus) Close() {
	for i := range f.widgets {
		for j := range f.widgets[i] {
			close(f.widgets[i][j])
		}
	}
}

func (f *Focus) Focus(focus bool) {
	f.widgets[f.p.Y][f.p.X] <- focus
}

func (f *Focus) Left() {
	f.Focus(false)
	if f.p.X <= 0 {
		f.p.X = len(f.widgets[f.p.Y]) - 1
	} else {
		f.p.X--
	}
	f.Focus(true)
}

func (f *Focus) Right() {
	f.Focus(false)
	f.p.X = (f.p.X + 1) % len(f.widgets[f.p.Y])
	f.Focus(true)
}

func (f *Focus) Up() {
	f.Focus(false)
	if f.p.Y <= 0 {
		f.p.Y = len(f.widgets) - 1
	} else {
		f.p.Y--
	}
	f.p.X = min(f.p.X, len(f.widgets[f.p.Y])-1)
	f.Focus(true)
}

func (f *Focus) Down() {
	f.Focus(false)
	f.p.Y = (f.p.Y + 1) % len(f.widgets)
	f.p.X = min(f.p.X, len(f.widgets[f.p.Y])-1)
	f.Focus(true)
}
