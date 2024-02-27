package widget

import (
	"image"
	"sync"
)

type FocusMaster struct {
	slaves [][]chan bool
	mu     sync.Mutex
	p      image.Point // coordinates of currently focused slave
}

type FocusSlave struct {
	Focus <-chan bool
	Mu    *sync.Mutex
}

func NewFocusMaster(rows []int) FocusMaster {
	f := FocusMaster{
		make([][]chan bool, len(rows)),
		sync.Mutex{},
		image.Point{},
	}
	for i := range f.slaves {
		f.slaves[i] = make([]chan bool, rows[i])
		for j := range f.slaves[i] {
			f.slaves[i][j] = make(chan bool)
		}
	}
	return f
}

func (f *FocusMaster) Slave(y, x int) FocusSlave {
	return FocusSlave{f.slaves[y][x], &f.mu}
}

func (f *FocusMaster) Close() {
	for i := range f.slaves {
		for j := range f.slaves[i] {
			close(f.slaves[i][j])
		}
	}
}

func (f *FocusMaster) Focus(focus bool) {
	f.slaves[f.p.Y][f.p.X] <- focus
}

func (f *FocusMaster) TryLeft() {
	if !f.mu.TryLock() {
		return
	}
	defer f.mu.Unlock()
	f.Focus(false)
	if f.p.X <= 0 {
		f.p.X = len(f.slaves[f.p.Y]) - 1
	} else {
		f.p.X--
	}
	f.Focus(true)
}

func (f *FocusMaster) TryRight() {
	if !f.mu.TryLock() {
		return
	}
	defer f.mu.Unlock()
	f.Focus(false)
	f.p.X = (f.p.X + 1) % len(f.slaves[f.p.Y])
	f.Focus(true)
}

func (f *FocusMaster) TryUp() {
	if !f.mu.TryLock() {
		return
	}
	defer f.mu.Unlock()
	f.Focus(false)
	if f.p.Y <= 0 {
		f.p.Y = len(f.slaves) - 1
	} else {
		f.p.Y--
	}
	f.p.X = min(f.p.X, len(f.slaves[f.p.Y])-1)
	f.Focus(true)
}

func (f *FocusMaster) TryDown() {
	if !f.mu.TryLock() {
		return
	}
	defer f.mu.Unlock()
	f.Focus(false)
	f.p.Y = (f.p.Y + 1) % len(f.slaves)
	f.p.X = min(f.p.X, len(f.slaves[f.p.Y])-1)
	f.Focus(true)
}
