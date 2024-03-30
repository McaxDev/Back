package util

import (
	"image"
	"image/color"
	"image/draw"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// RenderTextOnImage 渲染文本到图片
func RenderTextOnImage(xPct, yPct, fontSize int, text string, textColor color.Color, img image.Image) (image.Image, error) {
	// 加载字体文件
	fontBytes, err := os.ReadFile("path/to/font/file") // 需要指定字体文件路径
	if err != nil {
		return nil, err
	}

	f, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    float64(fontSize), // 将整数字体大小转换为浮点数
		DPI:     72,                // 根据需要调整DPI
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
		Src:  image.NewUniform(textColor),
		Face: face,
	}

	// 根据百分比计算实际坐标
	x := (bounds.Dx() * xPct) / 100
	y := (bounds.Dy() * yPct) / 100
	d.Dot = fixed.P(x, y)

	d.DrawString(text)

	return newImg, nil
}
