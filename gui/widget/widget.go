package widget

import (
	"cmp"
	"fmt"
	"sync"

	"image"
	"image/color"
	"image/draw"

	"volute/gui"
	"volute/gui/win"
)

var (
	FOCUS_COLOR = color.RGBA{179, 217, 255, 255}
	GREEN       = color.RGBA{51, 102, 0, 255}
	BLACK       = color.Gray{0}
	WHITE       = color.Gray{255}
)

func Label(text string, r image.Rectangle, env gui.Env, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(env.Draw())

	redraw := func(drw draw.Image) image.Rectangle {
		drawText([]byte(text), drw, r, BLACK, WHITE)
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

func Input(val chan<- uint, r image.Rectangle, focusChan <-chan bool, env gui.Env, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(env.Draw())

	redraw := func(text []byte, focus bool) func(draw.Image) image.Rectangle {
		return func(drw draw.Image) image.Rectangle {
			if focus {
				drawText(text, drw, r, GREEN, FOCUS_COLOR)
			} else {
				drawText(text, drw, r, GREEN, WHITE)
			}
			return r
		}
	}
	text := []byte{'0'}
	focus := false

	env.Draw() <- redraw(text, focus)
Loop:
	for {
		select {
		case focus = <-focusChan:
			env.Draw() <- redraw(text, focus)
		case event, ok := <-env.Events():
			if !ok { // channel closed
				break Loop
			}
			switch event := event.(type) {
			case win.WiFocus:
				if event.Focused {
					env.Draw() <- redraw(text, focus)
				}
			case win.KbType:
				if focus && isDigit(event.Rune) {
					text = fmt.Appendf(text, "%c", event.Rune)
					env.Draw() <- redraw(text, focus)
					val <- atoi(text)
				}
			case win.KbDown:
				if focus && event.Key == win.KeyBackspace && len(text) > 0 {
					text = text[:len(text)-1]
					env.Draw() <- redraw(text, focus)
					val <- atoi(text)
				}
			}
		}
	}
}

func Output(val <-chan float64, r image.Rectangle, env gui.Env, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(env.Draw())

	redraw := func(v float64) func(draw.Image) image.Rectangle {
		return func(drw draw.Image) image.Rectangle {
			drawText([]byte(fmt.Sprintf("%.3f", v)), drw, r, BLACK, WHITE)
			return r
		}
	}
	var v float64 = 0.0

	env.Draw() <- redraw(v)
Loop:
	for {
		select {
		case v = <-val:
			env.Draw() <- redraw(v)
		case event, ok := <-env.Events():
			if !ok { // channel closed
				break Loop
			}
			if event, ok := event.(win.WiFocus); ok && event.Focused {
				env.Draw() <- redraw(v)
			}
		}
	}
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
