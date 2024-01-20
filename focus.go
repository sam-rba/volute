package main

type Focus struct {
	widgets []chan bool
	i       int // index of focused widget
}

func NewFocus(nWidgets int) Focus {
	f := Focus{make([]chan bool, nWidgets), 0}
	for i := range f.widgets {
		f.widgets[i] = make(chan bool)
	}
	return f
}

func (f *Focus) Next() {
	f.widgets[f.i] <- false
	f.i = (f.i + 1) % len(f.widgets)
	f.widgets[f.i] <- true
}

func (f *Focus) Prev() {
	f.widgets[f.i] <- false
	f.i = abs(f.i-1) % len(f.widgets)
	f.widgets[f.i] <- true
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
