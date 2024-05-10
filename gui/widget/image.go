package widget

import (
	"sync"

	"image"
	"image/draw"

	"volute/gui"
	"volute/gui/win"
)

func Image(imChan <-chan image.Image, r image.Rectangle, env gui.Env, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(env.Draw())

	var im image.Image = image.NewGray(r)
	for {
		select {
		case im = <-imChan:
			env.Draw() <- imageDraw(im, r)
		case event, ok := <-env.Events():
			if !ok {
				return
			}
			if event, ok := event.(win.WiFocus); ok && event.Focused {
				env.Draw() <- imageDraw(im, r)
			}
		}
	}
}

func imageDraw(im image.Image, r image.Rectangle) func(draw.Image) image.Rectangle {
	return func(drw draw.Image) image.Rectangle {
		interpolator.Scale(drw, r, im, im.Bounds(), draw.Src, nil)
		return r
	}
}
