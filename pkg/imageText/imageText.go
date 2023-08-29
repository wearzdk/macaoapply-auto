package imageText

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// 插入白色背景和文字
func InsertTextToImage(imageBytes []byte, text string) ([]byte, error) {
	src, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}

	r := image.Rect(0, 0, src.Bounds().Dx(), src.Bounds().Dy())

	img := image.NewRGBA(r)
	draw.Draw(img, r, src, image.Point{}, draw.Src)

	col := color.RGBA{255, 255, 255, 0xff} // 白色背景
	point := fixed.P(10, 20)               // Positioned at the top left corner

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(text)

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
