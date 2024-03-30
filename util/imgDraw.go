package util

import (
	"image"
	"image/color"
	"image/draw"
	"os"
	"strconv"

	co "github.com/McaxDev/Back/config"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func Draw(x, y, fs int, text string, color string, img image.Image) (image.Image, error) {

	fontBytes, err := os.ReadFile(co.Config.McFont)
	if err != nil {
		return nil, err
	}

	f, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    float64(fs),
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	draw.Draw(newImg, bounds, img, bounds.Min, draw.Src)

	d := &font.Drawer{
		Dst:  newImg,
		Src:  image.NewUniform(RGBA(color)),
		Face: face,
	}

	d.Dot = fixed.P((bounds.Dx()*x)/100, (bounds.Dy()*y)/100)
	d.DrawString(text)

	return newImg, nil
}

func RGBA(hex string) color.Color {
	var err error
	if len(hex) != 6 {
		return color.RGBA{}
	}
	var r, g, b uint64
	r, err = strconv.ParseUint(hex[0:2], 16, 8)
	g, err = strconv.ParseUint(hex[2:4], 16, 8)
	b, err = strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return color.RGBA{}
	}
	return color.RGBA{uint8(r), uint8(g), uint8(b), 0xff}
}
