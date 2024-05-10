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
	"volute/gui/text"
	"volute/gui/win"
)

var (
	FOCUS_COLOR = color.RGBA{179, 217, 255, 255}
	GREEN       = color.RGBA{51, 102, 0, 255}
	BLACK       = color.Gray{0}
	WHITE       = color.Gray{255}

	interpolator = xdraw.ApproxBiLinear
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
		text.Draw([]byte(fmt.Sprintf("%.3f", v)), drw, r, BLACK, WHITE)
		return r
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
