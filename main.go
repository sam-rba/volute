package main

import (
	"image"

	"github.com/faiface/mainthread"
	"volute/gui"
	"volute/gui/widget"
	"volute/gui/win"
)

func run() {
	w, err := win.New(win.Title("volute"), win.Size(800, 600))
	if err != nil {
		panic(err)
	}

	mux, env := gui.NewMux(w)

	var (
		displacementChan = make(chan float64)
	)

	go widget.Input(mux.MakeEnv(), image.Rect(20, 20, 100, 40), displacementChan)

Loop:
	for event := range env.Events() {
		switch event := event.(type) {
		case win.WiClose:
			break Loop
		case win.KbType:
			if event.Rune == 'q' {
				break Loop
			}
		}
		select {
		case _ = <-displacementChan:
		default:
		}
	}
	close(env.Draw())
	close(displacementChan)
}

func main() {
	mainthread.Run(run)
}
