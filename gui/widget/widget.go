package widget

import (
	"cmp"
	"fmt"
	"sync"

	xdraw "golang.org/x/image/draw"
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

func Input(val chan<- uint, r image.Rectangle, focus FocusSlave, env gui.Env, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(env.Draw())
	defer close(val)

	text := []byte{'0'}
	focused := false
	env.Draw() <- inputDraw(text, focused, r)
Loop:
	for {
		select {
		case _, ok := <-focus.gain:
			if !ok {
				break Loop
			}
			focused = true
			env.Draw() <- inputDraw(text, focused, r)
		case dir, ok := <-focus.lose:
			if !ok {
				break Loop
			}
			focus.yield <- dir
			focused = false
			env.Draw() <- inputDraw(text, focused, r)
		case event, ok := <-env.Events():
			if !ok {
				break Loop
			}
			switch event := event.(type) {
			case win.WiFocus:
				if event.Focused {
					env.Draw() <- inputDraw(text, focused, r)
				}
			case win.KbType:
				if focused && isDigit(event.Rune) {
					text = fmt.Appendf(text, "%c", event.Rune)
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

func inputDraw(text []byte, focused bool, r image.Rectangle) func(draw.Image) image.Rectangle {
	return func(drw draw.Image) image.Rectangle {
		if focused {
			drawText(text, drw, r, GREEN, FOCUS_COLOR)
		} else {
			drawText(text, drw, r, GREEN, WHITE)
		}
		return r
	}
}

func Output(val <-chan float64, r image.Rectangle, env gui.Env, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(env.Draw())

	var v float64 = 0.0
	env.Draw() <- outputDraw(v, r)
Loop:
	for {
		select {
		case v = <-val:
			env.Draw() <- outputDraw(v, r)
		case event, ok := <-env.Events():
			if !ok { // channel closed
				break Loop
			}
			if event, ok := event.(win.WiFocus); ok && event.Focused {
				env.Draw() <- outputDraw(v, r)
			}
		}
	}
}

func outputDraw(v float64, r image.Rectangle) func(draw.Image) image.Rectangle {
	return func(drw draw.Image) image.Rectangle {
		drawText([]byte(fmt.Sprintf("%.3f", v)), drw, r, BLACK, WHITE)
		return r
	}
}

func Image(imChan <-chan image.Image, r image.Rectangle, env gui.Env, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(env.Draw())

	interp := xdraw.ApproxBiLinear
	redraw := func(im image.Image) func(draw.Image) image.Rectangle {
		return func(drw draw.Image) image.Rectangle {
			interp.Scale(drw, r, im, im.Bounds(), draw.Src, nil)
			return r
		}
	}
	var im image.Image = image.NewGray(r)

	for {
		select {
		case im = <-imChan:
			env.Draw() <- redraw(im)
		case event, ok := <-env.Events():
			if !ok {
				return
			}
			if event, ok := event.(win.WiFocus); ok && event.Focused {
				env.Draw() <- redraw(im)
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
