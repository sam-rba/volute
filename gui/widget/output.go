package widget

import (
	"fmt"
	"image"
	"image/draw"
	"sync"

	"volute/gui"
	"volute/gui/text"
	"volute/gui/win"
)

func Output(val <-chan float64, r image.Rectangle, env gui.Env, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(env.Draw())

	var v float64 = 0.0
	env.Draw() <- outputDraw(v, r)
	for {
		select {
		case v = <-val:
			env.Draw() <- outputDraw(v, r)
		case event, ok := <-env.Events():
			if !ok {
				return
			}
			if event, ok := event.(win.WiFocus); ok && event.Focused {
				env.Draw() <- outputDraw(v, r)
			}
		}
	}
}

func outputDraw(v float64, r image.Rectangle) func(draw.Image) image.Rectangle {
	return func(drw draw.Image) image.Rectangle {
		text.Draw(fmt.Sprintf("%.3f", v), drw, r, BLACK, WHITE, text.ALIGN_RIGHT)
		return r
	}
}
