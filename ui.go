package main

import (
	"image"
	"image/color"
	"sync"

	"volute/gui"
	"volute/gui/layout"
	"volute/gui/widget"
)

func spawnWidgets(
	displacementChan chan uint,
	rpmChan, veChan, imapChan, actChan [POINTS]chan uint,
	flowChan [POINTS]chan float64,
	focus *widget.FocusMaster,
	mux *gui.Mux,
	wg *sync.WaitGroup,
) {
	bounds := layout.Grid{
		Rows:        []int{2, 7, 7, 7, 7, 7},
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
		focus.Slave(0, 0),
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
	wg.Add(1)
	go widget.Label("mass flow (kg/min)", bounds[6+4*POINTS], mux.MakeEnv(), wg)
	for i := 0; i < POINTS; i++ {
		wg.Add(1)
		go widget.Input( // speed
			rpmChan[i],
			bounds[3+i],
			focus.Slave(1, i),
			mux.MakeEnv(),
			wg,
		)
		wg.Add(1)
		go widget.Input( // VE
			veChan[i],
			bounds[4+POINTS+i],
			focus.Slave(2, i),
			mux.MakeEnv(),
			wg,
		)
		wg.Add(1)
		go widget.Input( // IMAP
			imapChan[i],
			bounds[5+2*POINTS+i],
			focus.Slave(3, i),
			mux.MakeEnv(),
			wg,
		)
		wg.Add(1)
		go widget.Input( // ACT
			actChan[i],
			bounds[6+3*POINTS+i],
			focus.Slave(4, i),
			mux.MakeEnv(),
			wg,
		)
		wg.Add(1)
		go widget.Output( // mass flow
			flowChan[i],
			bounds[7+4*POINTS+i],
			mux.MakeEnv(),
			wg,
		)
	}
}

func split(elements int, space int) []int {
	bounds := make([]int, elements)
	widths := []int{
		widget.TextSize(WIDEST_LABEL).X,
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
