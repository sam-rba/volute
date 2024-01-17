package widget

import (
	"image"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type concurrentFace struct {
	mu   sync.Mutex
	face font.Face
}

func (cf *concurrentFace) Close() error {
	cf.mu.Lock()
	defer cf.mu.Unlock()
	return cf.face.Close()
}

func (cf *concurrentFace) Glyph(dot fixed.Point26_6, r rune) (
	dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {
	cf.mu.Lock()
	defer cf.mu.Unlock()
	return cf.face.Glyph(dot, r)
}

func (cf *concurrentFace) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	cf.mu.Lock()
	defer cf.mu.Unlock()
	return cf.face.GlyphBounds(r)
}

func (cf *concurrentFace) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	cf.mu.Lock()
	defer cf.mu.Unlock()
	return cf.face.GlyphAdvance(r)
}

func (cf *concurrentFace) Kern(r0, r1 rune) fixed.Int26_6 {
	cf.mu.Lock()
	defer cf.mu.Unlock()
	return cf.face.Kern(r0, r1)
}

func (cf *concurrentFace) Metrics() font.Metrics {
	cf.mu.Lock()
	defer cf.mu.Unlock()
	return cf.face.Metrics()
}
