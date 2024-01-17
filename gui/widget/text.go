package widget

import (
	"log"
	"sync"

	"image"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var (
	FONT               = goregular.TTF
	FONT_SIZE  float64 = 15
	DPI        float64 = 72
	BG_COLOR           = image.White
	TEXT_COLOR         = image.Black
)

var face *concurrentFace

func init() {
	fnt, err := opentype.Parse(FONT)
	if err != nil {
		log.Fatal(err)
	}
	fce, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size: FONT_SIZE,
		DPI:  DPI,
	})
	if err != nil {
		log.Fatal(err)
	}
	face = &concurrentFace{sync.Mutex{}, fce}
}

func drawText(text []byte, dst draw.Image, r image.Rectangle) {
	drawer := font.Drawer{
		Src:  TEXT_COLOR,
		Face: face,
		Dot:  fixed.P(0, 0),
	}

	// background
	draw.Draw(dst, r, BG_COLOR, image.ZP, draw.Src)

	// text image
	bounds := textBounds(text, drawer)
	textImg := image.NewRGBA(bounds)
	draw.Draw(textImg, bounds, BG_COLOR, image.ZP, draw.Src)
	drawer.Dst = textImg
	drawer.DrawBytes(text)

	// draw text image over background
	left := image.Pt(bounds.Min.X, (bounds.Min.Y+bounds.Max.Y)/2)
	target := image.Pt(r.Min.X, (r.Min.Y+r.Max.Y)/2)
	delta := target.Sub(left)
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
