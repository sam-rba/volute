package widget

import "image"

// Focus keeps track of the currently selected widget.
// A widget receives true when it gains focus and false when it loses focus.
type Focus struct {
	Widgets [][]chan bool
	p       image.Point // coordinates of currently focused widget
}

func NewFocus(rows []int) Focus {
	f := Focus{
		make([][]chan bool, len(rows)),
		image.Point{},
	}
	for i := range f.Widgets {
		f.Widgets[i] = make([]chan bool, rows[i])
		for j := range f.Widgets[i] {
			f.Widgets[i][j] = make(chan bool)
		}
	}
	return f
}

func (f *Focus) Close() {
	for i := range f.Widgets {
		for j := range f.Widgets[i] {
			close(f.Widgets[i][j])
		}
	}
}

func (f *Focus) Focus(focus bool) {
	f.Widgets[f.p.Y][f.p.X] <- focus
}

func (f *Focus) Left() {
	f.Focus(false)
	if f.p.X <= 0 {
		f.p.X = len(f.Widgets[f.p.Y]) - 1
	} else {
		f.p.X--
	}
	f.Focus(true)
}

func (f *Focus) Right() {
	f.Focus(false)
	f.p.X = (f.p.X + 1) % len(f.Widgets[f.p.Y])
	f.Focus(true)
}

func (f *Focus) Up() {
	f.Focus(false)
	if f.p.Y <= 0 {
		f.p.Y = len(f.Widgets) - 1
	} else {
		f.p.Y--
	}
	f.p.X = min(f.p.X, len(f.Widgets[f.p.Y])-1)
	f.Focus(true)
}

func (f *Focus) Down() {
	f.Focus(false)
	f.p.Y = (f.p.Y + 1) % len(f.Widgets)
	f.p.X = min(f.p.X, len(f.Widgets[f.p.Y])-1)
	f.Focus(true)
}
