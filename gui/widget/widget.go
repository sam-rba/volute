package widget

import (
	"cmp"
	"fmt"

	"image"
	"image/draw"

	"volute/gui"
	"volute/gui/win"
)

func Label(text string, r image.Rectangle, env gui.Env) {
	redraw := func(drw draw.Image) image.Rectangle {
		drawText([]byte(text), drw, r)
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
	close(env.Draw())
}

func Input(val chan<- uint, r image.Rectangle, env gui.Env) {
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
			if focus && isDigit(event.Rune) {
				text = fmt.Appendf(text, "%c", event.Rune)
				env.Draw() <- redraw(text)
				val <- atoi(text)
			}
		case win.KbDown:
			if focus && event.Key == win.KeyBackspace && len(text) > 0 {
				text = text[:len(text)-1]
				env.Draw() <- redraw(text)
				val <- atoi(text)
			}
		}
	}
	close(env.Draw())
}

func Output(val <-chan uint, r image.Rectangle, env gui.Env) {
	redraw := func(n uint) func(draw.Image) image.Rectangle {
		return func(drw draw.Image) image.Rectangle {
			drawText([]byte(fmt.Sprint(n)), drw, r)
			return r
		}
	}

	var n uint = 0
	env.Draw() <- redraw(n)

Loop:
	for {
		select {
		case n = <-val:
			env.Draw() <- redraw(n)
		case event, ok := <-env.Events():
			if !ok { // channel closed
				break Loop
			}
			if event, ok := event.(win.WiFocus); ok && event.Focused {
				env.Draw() <- redraw(n)
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

func atoi(s []byte) uint {
	var n uint = 0
	for _, d := range s {
		n = n*10 + uint(d-'0')
	}
	return n
}
