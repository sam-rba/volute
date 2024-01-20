package widget

import (
	"log"
	"sync"

	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var (
	FONT      = goregular.TTF
	FONT_SIZE = 15
	DPI       = 72
	PAD       = 3
)

var face *concurrentFace

func init() {
	fnt, err := opentype.Parse(FONT)
	if err != nil {
		log.Fatal(err)
	}
	fce, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size: float64(FONT_SIZE),
		DPI:  float64(DPI),
	})
	if err != nil {
		log.Fatal(err)
	}
	face = &concurrentFace{sync.Mutex{}, fce}
}

func TextSize(text string) image.Point {
	bounds := textBounds([]byte(text), font.Drawer{Face: face})
	return image.Point{bounds.Max.X - bounds.Min.X + 2*PAD, bounds.Max.Y - bounds.Min.Y + 2*PAD}
}

func drawText(text []byte, dst draw.Image, r image.Rectangle, fg, bg color.Color) {
	drawer := font.Drawer{
		Src:  &image.Uniform{fg},
		Face: face,
		Dot:  fixed.P(0, 0),
	}

	// background
	draw.Draw(dst, r, &image.Uniform{bg}, image.ZP, draw.Src)

	// text image
	bounds := textBounds(text, drawer)
	textImg := image.NewRGBA(bounds)
	draw.Draw(textImg, bounds, &image.Uniform{bg}, image.ZP, draw.Src)
	drawer.Dst = textImg
	drawer.DrawBytes(text)

	// draw text image over background
	leftCentre := image.Pt(bounds.Min.X, (bounds.Min.Y+bounds.Max.Y)/2)
	target := image.Pt(r.Max.X-bounds.Max.X-PAD, (r.Min.Y+r.Max.Y)/2)
	delta := target.Sub(leftCentre)
	draw.Draw(dst, bounds.Add(delta).Intersect(r), drawer.Dst, bounds.Min, draw.Src)
}

func textBounds(text []byte, drawer font.Drawer) image.Rectangle {
	b, _ := drawer.BoundBytes(text)
	return image.Rect(
		b.Min.X.Floor(),
		b.Min.Y.Floor(),
		b.Max.X.Ceil(),
		b.Max.Y.Ceil(),
	)
}
