package main

import (
	"fmt"
	"os"
	"sync"

	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/faiface/mainthread"
	"volute/gui"
	"volute/gui/text"
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

	focus := widget.NewFocusMaster([]int{1, POINTS, POINTS, POINTS, POINTS, 1})
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
		focus,
		mux,
		wg,
	)

	compressors := []widget.Node[string]{
		{
			Label: "BorgWarner",
			Value: "bw",
			Children: []widget.Node[string]{
				{
					Label: "EFR",
					Value: "efr",
					Children: []widget.Node[string]{
						{
							Label:    "6258",
							Value:    "6258",
							Children: nil,
						}, {
							Label:    "7064",
							Value:    "7064",
							Children: nil,
						},
					},
				}, {
					Label: "K",
					Value: "k",
					Children: []widget.Node[string]{
						{
							Label:    "03",
							Value:    "03",
							Children: nil,
						}, {
							Label:    "04",
							Value:    "04",
							Children: nil,
						},
					},
				},
			},
		}, {
			Label: "Garrett",
			Value: "garrett",
			Children: []widget.Node[string]{
				{
					Label: "G",
					Value: "g",
					Children: []widget.Node[string]{
						{
							Label:    "25-550",
							Value:    "25-550",
							Children: nil,
						},
					},
				},
			},
		},
	}
	wg.Add(1)
	go widget.Tree(
		compressors,
		image.Rect(text.PAD, 125, 250, HEIGHT-text.PAD),
		focus.Slave(5, 0),
		mux,
		wg,
	)

	imChan := make(chan image.Image)
	defer close(imChan)
	wg.Add(1)
	go widget.Image(
		imChan,
		image.Rect(250+text.PAD, 125, WIDTH-text.PAD, HEIGHT-text.PAD),
		mux.MakeEnv(),
		wg,
	)
	f, err := os.Open("compressor_maps/borgwarner/efr/8374.jpg")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	im, _, err := image.Decode(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	f.Close()
	imChan <- im

	for i := 0; i < POINTS; i++ {
		wg.Add(1)
		go calculateFlow(
			flowChan[i],
			displacementBroadcast.AddDestination(),
			rpmChan[i], veChan[i], actChan[i], imapChan[i],
			wg,
		)
	}

	eventLoop(env, focus)
}

func eventLoop(env gui.Env, focus *widget.FocusMaster) {
	for event := range env.Events() {
		switch event := event.(type) {
		case win.WiClose:
			return
		case win.KbType:
			switch event.Rune {
			case 'q':
				return
			case 'h':
				focus.Shift(widget.LEFT)
			case 'j':
				focus.Shift(widget.DOWN)
			case 'k':
				focus.Shift(widget.UP)
			case 'l':
				focus.Shift(widget.RIGHT)
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
