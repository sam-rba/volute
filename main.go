package main

import (
	"image"
	"image/color"

	"github.com/faiface/mainthread"
	"volute/gui"
	"volute/gui/layout"
	"volute/gui/widget"
	"volute/gui/win"
)

const (
	WIDTH  = 800
	HEIGHT = 600

	POINTS = 6
)

func run() {
	w, err := win.New(win.Title("volute"), win.Size(WIDTH, HEIGHT))
	if err != nil {
		panic(err)
	}
	mux, env := gui.NewMux(w)

	var (
		displacementChan = make(chan uint)
		rpmChan          = make([]chan uint, POINTS)
		veChan           = make([]chan uint, POINTS)
		focus            = NewFocus(1 + 2*POINTS)
	)
	for i := 0; i < POINTS; i++ {
		rpmChan[i] = make(chan uint)
		veChan[i] = make(chan uint)
	}

	bounds := layout.Grid{
		Rows:        []int{2, 7, 7},
		Background:  color.Gray{255},
		Gap:         1,
		Split:       split,
		SplitRows:   splitRows,
		Margin:      0,
		Border:      0,
		BorderColor: color.Gray{16},
		Flip:        false,
	}.Lay(image.Rect(0, 0, WIDTH, HEIGHT))

	go widget.Label("displacement (cc)", bounds[0], mux.MakeEnv())
	go widget.Input(
		displacementChan,
		bounds[1],
		focus.widgets[0],
		mux.MakeEnv(),
	)
	go widget.Label("speed (rpm)", bounds[2], mux.MakeEnv())
	go widget.Label("VE (%)", bounds[3+POINTS], mux.MakeEnv())
	for i := 0; i < POINTS; i++ {
		go widget.Input(
			rpmChan[i],
			bounds[3+i],
			focus.widgets[1+i],
			mux.MakeEnv(),
		)
		go widget.Input(
			veChan[i],
			bounds[3+POINTS+1+i],
			focus.widgets[1+POINTS+i],
			mux.MakeEnv(),
		)
	}

	focus.widgets[focus.i] <- true

Loop:
	for {
		select {
		case _ = <-displacementChan:
		case _ = <-rpmChan[0]:
		case event, ok := <-env.Events():
			if !ok { // channel closed
				break Loop
			}
			switch event := event.(type) {
			case win.WiClose:
				break Loop
			case win.KbType:
				switch event.Rune {
				case 'q':
					break Loop
				case 'j', 'l':
					focus.Next()
				case 'k', 'h':
					focus.Prev()
				}
			}
		}
	}
	close(env.Draw())
	close(displacementChan)
	for i := range rpmChan {
		close(rpmChan[i])
	}
}

func split(elements int, space int) []int {
	bounds := make([]int, elements)
	widths := []int{
		widget.TextSize("displacement (cc)").X,
		widget.TextSize("123456").X,
	}
	for i := 0; i < elements && space > 0; i++ {
		bounds[i] = min(widths[min(i, len(widths)-1)], space)
		space -= bounds[i]
	}
	return bounds
}

func splitRows(elements int, space int) []int {
	bounds := make([]int, elements)
	height := widget.TextSize("1").Y
	for i := 0; i < elements && space > 0; i++ {
		bounds[i] = min(height, space)
		space -= bounds[i]
	}
	return bounds
}

func main() {
	mainthread.Run(run)
}
