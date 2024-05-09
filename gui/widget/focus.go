package widget

import (
	"fmt"
	"image"
)

type Direction int

const (
	UP Direction = iota
	DOWN
	LEFT
	RIGHT
)

type FocusSlave struct {
	gain  <-chan bool
	lose  <-chan Direction // forward to yield if widget accepts focus loss
	yield chan<- Direction
}

type focusSlave struct {
	gain chan bool
	lose chan Direction
}

func (fs focusSlave) close() {
	close(fs.gain)
	close(fs.lose)
}

type FocusMaster struct {
	slaves  [][]focusSlave
	yield   chan Direction
	focused image.Point // coordinates of focused slave
}

func NewFocusMaster(rows []int) *FocusMaster {
	fm := &FocusMaster{
		slaves:  make([][]focusSlave, len(rows)),
		yield:   make(chan Direction),
		focused: image.Point{0, 0},
	}
	for y := range fm.slaves {
		fm.slaves[y] = make([]focusSlave, rows[y])
		for x := range fm.slaves[y] {
			fm.slaves[y][x] = focusSlave{make(chan bool), make(chan Direction)}
		}
	}

	go func() {
		fm.slaves[0][0].gain <- true

		for dir := range fm.yield {
			fm.focused = fm.neighborPos(fm.focused, dir)
			fm.slaves[fm.focused.Y][fm.focused.X].gain <- true
		}
	}()

	return fm
}

func (fm FocusMaster) Slave(y, x int) FocusSlave {
	return FocusSlave{
		gain:  fm.slaves[y][x].gain,
		lose:  fm.slaves[y][x].lose,
		yield: fm.yield,
	}
}

func (fm FocusMaster) Close() {
	for y := range fm.slaves {
		for x := range fm.slaves[y] {
			fm.slaves[y][x].close()
		}
	}
	close(fm.yield)
}

func (fm FocusMaster) Shift(dir Direction) {
	fm.slaves[fm.focused.Y][fm.focused.X].lose <- dir
}

func (fm FocusMaster) neighborPos(pos image.Point, dir Direction) image.Point {
	switch dir {
	case UP:
		return fm.upNeighborPos(pos)
	case DOWN:
		return fm.downNeighborPos(pos)
	case LEFT:
		return fm.leftNeighborPos(pos)
	case RIGHT:
		return fm.rightNeighborPos(pos)
	default:
		panic(fmt.Sprintf("invalid Direction: %v", dir))
	}
}

func (fm FocusMaster) upNeighborPos(pos image.Point) image.Point {
	if pos.Y <= 0 {
		pos.Y = len(fm.slaves) - 1
	} else {
		pos.Y--
	}
	pos.X = min(pos.X, len(fm.slaves[pos.Y])-1)
	return pos
}

func (fm FocusMaster) downNeighborPos(pos image.Point) image.Point {
	if pos.Y >= len(fm.slaves)-1 {
		pos.Y = 0
	} else {
		pos.Y++
	}
	pos.X = min(pos.X, len(fm.slaves[pos.Y])-1)
	return pos
}

func (fm FocusMaster) leftNeighborPos(pos image.Point) image.Point {
	if pos.X <= 0 {
		pos.X = len(fm.slaves[pos.Y]) - 1
	} else {
		pos.X--
	}
	return pos
}

func (fm FocusMaster) rightNeighborPos(pos image.Point) image.Point {
	if pos.X >= len(fm.slaves[pos.Y])-1 {
		pos.X = 0
	} else {
		pos.X++
	}
	return pos
}
