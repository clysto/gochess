package ui

import (
	"github.com/clysto/gochess/resources"
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"image/color"
)

func MeasureString(face font.Face, s string) (int, int) {
	d := &font.Drawer{
		Face: face,
	}
	a := d.MeasureString(s)
	return int(a >> 6), int(face.Metrics().Height >> 6)
}

func Button(s string) *ebiten.Image {
	w, h := MeasureString(resources.MaShanZhengRegularFont, s)
	w += 40
	h += 10
	dc := gg.NewContext(w, h)
	dc.DrawRoundedRectangle(0, 0, float64(w), float64(h), float64(h/2))
	dc.SetColor(color.Black)
	dc.Fill()
	dc.SetColor(color.White)
	dc.SetFontFace(resources.MaShanZhengRegularFont)
	dc.DrawString(s, 20, float64(h)-10)
	img := ebiten.NewImageFromImage(dc.Image())
	return img
}
