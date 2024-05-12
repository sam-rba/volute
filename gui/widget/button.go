package widget

import (
	"image"
	"image/draw"
	"sync"

	"volute/gui"
	"volute/gui/text"
	"volute/gui/win"
)

// Sends signal over tx when Enter is pressed while focused.
func Button[T any](
	tx chan<- T,
	signal T,
	label string,
	r image.Rectangle,
	focus FocusSlave,
	env gui.Env,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	defer close(tx)

	focused := false
	env.Draw() <- buttonDraw(label, focused, r)
	for {
		select {
		case _, ok := <-focus.gain:
			if !ok {
				return
			}
			focused = true
			env.Draw() <- buttonDraw(label, focused, r)
		case dir, ok := <-focus.lose:
			if !ok {
				return
			}
			focus.yield <- dir
			focused = false
			env.Draw() <- buttonDraw(label, focused, r)
		case event, ok := <-env.Events():
			if !ok {
				return
			}
			if event, ok := event.(win.KbDown); ok && focused && event.Key == win.KeyEnter {
				tx <- signal
			}
		}
	}
}

func buttonDraw(label string, focused bool, r image.Rectangle) func(drw draw.Image) image.Rectangle {
	return func(drw draw.Image) image.Rectangle {
		if focused {
			text.Draw(label, drw, r, BLACK, FOCUS_COLOR, text.ALIGN_LEFT)
		} else {
			text.Draw(label, drw, r, BLACK, WHITE, text.ALIGN_LEFT)
		}
		return r
	}
}
