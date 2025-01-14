package widget

import (
	"image"
	"image/draw"
	"sync"

	"volute/gui"
	"volute/gui/text"
	"volute/gui/win"
)

func Input(val chan<- uint, r image.Rectangle, focus FocusSlave, env gui.Env, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(env.Draw())
	defer close(val)

	text := "0"
	focused := false
	env.Draw() <- inputDraw(text, focused, r)
	for {
		select {
		case _, ok := <-focus.gain:
			if !ok {
				return
			}
			focused = true
			env.Draw() <- inputDraw(text, focused, r)
		case dir, ok := <-focus.lose:
			if !ok {
				return
			}
			focus.yield <- dir
			focused = false
			env.Draw() <- inputDraw(text, focused, r)
		case event, ok := <-env.Events():
			if !ok {
				return
			}
			switch event := event.(type) {
			case win.WiFocus:
				if event.Focused {
					env.Draw() <- inputDraw(text, focused, r)
				}
			case win.KbType:
				if focused && isDigit(event.Rune) {
					text = text + string(event.Rune)
					env.Draw() <- inputDraw(text, focused, r)
					val <- atoi(text)
				}
			case win.KbDown:
				if focused && event.Key == win.KeyBackspace && len(text) > 0 {
					text = text[:len(text)-1]
					env.Draw() <- inputDraw(text, focused, r)
					val <- atoi(text)
				}
			}
		}
	}
}

func inputDraw(str string, focused bool, r image.Rectangle) func(draw.Image) image.Rectangle {
	return func(drw draw.Image) image.Rectangle {
		if focused {
			text.Draw(str, drw, r, GREEN, FOCUS_COLOR, text.ALIGN_RIGHT)
		} else {
			text.Draw(str, drw, r, GREEN, WHITE, text.ALIGN_RIGHT)
		}
		return r
	}
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func atoi(s string) uint {
	var n uint = 0
	for _, d := range s {
		n = n*10 + uint(d-'0')
	}
	return n
}
