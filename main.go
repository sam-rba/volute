package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"sync"

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

	R = 8314.3 // gas constant
	M = 28.962 // molar mass of air
)

func run() {
	var (
		wg = new(sync.WaitGroup)

		focus = NewFocus([]int{1, POINTS, POINTS, POINTS, POINTS})

		displacementChan = make(chan uint)

		rpmChan  [POINTS]chan uint
		veChan   [POINTS]chan uint
		imapChan [POINTS]chan uint
		actChan  [POINTS]chan uint
	)
	defer wg.Wait()
	defer focus.Close()
	defer close(displacementChan)
	for i := 0; i < POINTS; i++ {
		rpmChan[i] = make(chan uint)
		veChan[i] = make(chan uint)
		imapChan[i] = make(chan uint)
		actChan[i] = make(chan uint)

		defer close(rpmChan[i])
		defer close(veChan[i])
		defer close(imapChan[i])
		defer close(actChan[i])
	}

	w, err := win.New(win.Title("volute"), win.Size(WIDTH, HEIGHT))
	if err != nil {
		fmt.Println("error creating window:", err)
		os.Exit(1)
	}
	mux, env := gui.NewMux(w)
	defer close(env.Draw())

	bounds := layout.Grid{
		Rows:        []int{2, 7, 7, 7, 7},
		Background:  color.Gray{255},
		Gap:         1,
		Split:       split,
		SplitRows:   splitRows,
		Margin:      0,
		Border:      0,
		BorderColor: color.Gray{16},
		Flip:        false,
	}.Lay(image.Rect(0, 0, WIDTH, HEIGHT))

	wg.Add(1)
	go widget.Label("displacement (cc)", bounds[0], mux.MakeEnv(), wg)
	wg.Add(1)
	go widget.Input(
		displacementChan,
		bounds[1],
		focus.widgets[0][0],
		mux.MakeEnv(),
		wg,
	)
	wg.Add(1)
	go widget.Label("speed (rpm)", bounds[2], mux.MakeEnv(), wg)
	wg.Add(1)
	go widget.Label("VE (%)", bounds[3+POINTS], mux.MakeEnv(), wg)
	wg.Add(1)
	go widget.Label("IMAP (mbar)", bounds[4+2*POINTS], mux.MakeEnv(), wg)
	wg.Add(1)
	go widget.Label("ACT (°C)", bounds[5+3*POINTS], mux.MakeEnv(), wg)
	for i := 0; i < POINTS; i++ {
		wg.Add(1)
		go widget.Input( // speed
			rpmChan[i],
			bounds[3+i],
			focus.widgets[1][i],
			mux.MakeEnv(),
			wg,
		)
		wg.Add(1)
		go widget.Input( // VE
			veChan[i],
			bounds[4+POINTS+i],
			focus.widgets[2][i],
			mux.MakeEnv(),
			wg,
		)
		wg.Add(1)
		go widget.Input( // IMAP
			imapChan[i],
			bounds[5+2*POINTS+i],
			focus.widgets[3][i],
			mux.MakeEnv(),
			wg,
		)
		wg.Add(1)
		go widget.Input( // ACT
			actChan[i],
			bounds[6+3*POINTS+i],
			focus.widgets[4][i],
			mux.MakeEnv(),
			wg,
		)
	}

	focus.widgets[focus.p.Y][focus.p.X] <- true

Loop:
	for {
		select {
		case _ = <-displacementChan:
		case _ = <-rpmChan[0]:
		case _ = <-veChan[0]:
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
				case 'h':
					focus.Left()
				case 'j':
					focus.Down()
				case 'k':
					focus.Up()
				case 'l':
					focus.Right()
				}
			}
		}
	}
	fmt.Println("Shutting down...")
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
