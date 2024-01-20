package main

type Focus struct {
	widgets [][]chan bool
	p       Point // currently focused widget
}

func NewFocus(rows []int) Focus {
	f := Focus{
		make([][]chan bool, len(rows)),
		Point{},
	}
	for i := range f.widgets {
		f.widgets[i] = make([]chan bool, rows[i])
		for j := range f.widgets[i] {
			f.widgets[i][j] = make(chan bool)
		}
	}
	return f
}

func (f *Focus) Left() {
	f.widgets[f.p.Y][f.p.X] <- false
	if f.p.X <= 0 {
		f.p.X = len(f.widgets[f.p.Y]) - 1
	} else {
		f.p.X--
	}
	f.widgets[f.p.Y][f.p.X] <- true
}

func (f *Focus) Right() {
	f.widgets[f.p.Y][f.p.X] <- false
	f.p.X = (f.p.X + 1) % len(f.widgets[f.p.Y])
	f.widgets[f.p.Y][f.p.X] <- true
}

func (f *Focus) Up() {
	f.widgets[f.p.Y][f.p.X] <- false
	if f.p.Y <= 0 {
		f.p.Y = len(f.widgets) - 1
	} else {
		f.p.Y--
	}
	f.p.X = min(f.p.X, len(f.widgets[f.p.Y])-1)
	f.widgets[f.p.Y][f.p.X] <- true
}

func (f *Focus) Down() {
	f.widgets[f.p.Y][f.p.X] <- false
	f.p.Y = (f.p.Y + 1) % len(f.widgets)
	f.p.X = min(f.p.X, len(f.widgets[f.p.Y])-1)
	f.widgets[f.p.Y][f.p.X] <- true
}

type Point struct {
	X, Y int
}
