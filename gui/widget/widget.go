package widget

import (
	"cmp"
	"fmt"

	"image"
	"image/draw"

	"volute/gui"
	"volute/gui/win"
)

func Input(env gui.Env, r image.Rectangle, val chan<- float64) {
	redraw := func(text []byte) func(draw.Image) image.Rectangle {
		return func(drw draw.Image) image.Rectangle {
			drawText(text, drw, r)
			return r
		}
	}
	text := []byte{'0'}
	focus := false

	env.Draw() <- redraw(text)

	for event := range env.Events() {
		switch event := event.(type) {
		case win.WiFocus:
			if event.Focused {
				env.Draw() <- redraw(text)
			}
		case win.MoDown:
			if event.Point.In(r) {
				focus = true
			}
		case win.KbType:
			if !focus ||
				(!isDigit(event.Rune) && event.Rune != '.') ||
				(event.Rune == '.' && contains(text, '.')) {
				continue
			}
			text = fmt.Appendf(text, "%c", event.Rune)
			env.Draw() <- redraw(text)
		case win.KbDown:
			if event.Key == win.KeyBackspace && focus && len(text) > 0 {
				text = text[:len(text)-1]
				env.Draw() <- redraw(text)
			}
		}
	}
	close(env.Draw())
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func contains[T cmp.Ordered](slc []T, v T) bool {
	for i := range slc {
		if slc[i] == v {
			return true
		}
	}
	return false
}
