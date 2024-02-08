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

	WIDEST_LABEL = "mass flow (kg/min)"
)

func run() {
	wg := new(sync.WaitGroup)
	defer wg.Wait()

	focus := NewFocus([]int{1, POINTS, POINTS, POINTS, POINTS})
	defer focus.Close()

	displacementChan := make(chan uint)
	displacementBroadcast := NewBroadcast(displacementChan)
	defer displacementBroadcast.Wait()

	var (
		rpmChan  [POINTS]chan uint
		veChan   [POINTS]chan uint
		imapChan [POINTS]chan uint
		actChan  [POINTS]chan uint

		flowChan [POINTS]chan float64
	)
	makeChans(rpmChan[:], veChan[:], imapChan[:], actChan[:])
	makeChans(flowChan[:])

	w, err := win.New(win.Title("volute"), win.Size(WIDTH, HEIGHT))
	if err != nil {
		fmt.Println("error creating window:", err)
		os.Exit(1)
	}
	mux, env := gui.NewMux(w)
	defer close(env.Draw())

	spawnWidgets(
		displacementChan,
		rpmChan, veChan, imapChan, actChan,
		flowChan,
		&focus, mux, wg,
	)

	// TODO: make these output properly on screen.
	for i := 0; i < POINTS; i++ {
		wg.Add(1)
		go calculateFlow(
			flowChan[i],
			displacementBroadcast.AddDestination(),
			rpmChan[i], veChan[i], actChan[i], imapChan[i],
			wg,
		)
	}

	focus.Focus(true)
	eventLoop(env, &focus)
}

func spawnWidgets(
	displacementChan chan uint,
	rpmChan, veChan, imapChan, actChan [POINTS]chan uint,
	flowChan [POINTS]chan float64,
	focus *Focus, mux *gui.Mux, wg *sync.WaitGroup,
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
	wg.Add(1)
	go widget.Label("mass flow (kg/min)", bounds[6+4*POINTS], mux.MakeEnv(), wg)
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
		wg.Add(1)
		go widget.Output( // mass flow
			flowChan[i],
			bounds[7+4*POINTS+i],
			mux.MakeEnv(),
			wg,
		)
	}
}

func eventLoop(env gui.Env, focus *Focus) {
	for event := range env.Events() {
		switch event := event.(type) {
		case win.WiClose:
			return
		case win.KbType:
			switch event.Rune {
			case 'q':
				return
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

func makeChans[T any](chanss ...[]chan T) {
	for i := range chanss {
		for j := range chanss[i] {
			chanss[i][j] = make(chan T)
		}
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

func calculateFlow(
	flow chan<- float64,
	displacementChan, rpmChan, veChan, actChan, imapChan <-chan uint,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	defer close(flow)

	var (
		displacement Volume
		rpm          uint
		ve           uint
		act          Temperature
		imap         Pressure

		v  uint
		ok bool
	)

	for {
		select {
		case v, ok = <-displacementChan:
			displacement = Volume(v) * CubicCentimetre
		case rpm, ok = <-rpmChan:
		case ve, ok = <-veChan:
		case v, ok = <-actChan:
			act = Temperature{float64(v), Celcius}
		case v, ok = <-imapChan:
			imap = Pressure(v) * Millibar
		}
		if !ok {
			return
		}
		flow <- massFlow(displacement, rpm, ve, act, imap)
	}
}

func massFlow(displacement Volume, rpm, ve uint, act Temperature, imap Pressure) float64 {
	density := (M / R) * float64(imap/Pascal) / act.AsUnit(Kelvin)                          // kg/m3
	volumeFlow := float64(displacement/CubicMetre) * float64(rpm/2) * (float64(ve) / 100.0) // m3/min
	return density * volumeFlow
}

func main() {
	mainthread.Run(run)
}
