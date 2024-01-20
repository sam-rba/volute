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
		displacementChan = make(chan uint)
		output           = make(chan uint)
		focus            = NewFocus(2)

		displacement uint = 0
	)

	pad := 10
	r := image.Rect(pad, pad, pad+widget.TextWidth(6), pad+widget.TextHeight())
	go widget.Input(
		displacementChan,
		r,
		focus.widgets[0],
		mux.MakeEnv(),
	)
	r = image.Rect(
		r.Max.X+pad,
		r.Min.Y,
		r.Max.X+pad+widget.TextWidth(len("cc")),
		r.Max.Y,
	)
	go widget.Label("cc", r, mux.MakeEnv())

	r = image.Rect(
		r.Max.X+pad,
		r.Min.Y,
		r.Max.X+pad+widget.TextWidth(6),
		r.Max.Y,
	)
	go widget.Input(
		displacementChan,
		r,
		focus.widgets[1],
		mux.MakeEnv(),
	)

	r = image.Rect(
		pad,
		r.Max.Y+pad,
		pad+widget.TextWidth(6),
		r.Max.Y+pad+widget.TextHeight(),
	)
	go widget.Output(output, r, mux.MakeEnv())

	focus.widgets[focus.i] <- true

Loop:
	for {
		select {
		case displacement = <-displacementChan:
			output <- displacement
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
	close(output)
}

func main() {
	mainthread.Run(run)
}
