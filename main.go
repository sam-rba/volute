package main

import (
	"fmt"
	"image"
	"os"
	"sync"

	"github.com/faiface/mainthread"
	"volute/gui"
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

	focus := widget.NewFocus([]int{1, POINTS, POINTS, POINTS, POINTS})
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
		&focus,
		mux,
		wg,
	)

	imChan := make(chan image.Image)
	defer close(imChan)
	wg.Add(1)
	go widget.Image(
		imChan,
		image.Rect(0, 200, 100, 300),
		mux.MakeEnv(),
		wg,
	)

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

func eventLoop(env gui.Env, focus *widget.Focus) {
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
