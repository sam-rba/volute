package widget

import (
	"image"
	"image/draw"
	"sync"

	"volute/gui"
	"volute/gui/text"
	"volute/gui/win"
)

func Label(str string, r image.Rectangle, env gui.Env, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(env.Draw())

	redraw := func(drw draw.Image) image.Rectangle {
		text.Draw([]byte(str), drw, r, BLACK, WHITE)
		return r
	}

	env.Draw() <- redraw
	for event := range env.Events() {
		switch event := event.(type) {
		case win.WiFocus:
			if event.Focused {
				env.Draw() <- redraw
			}
		}
	}
}
