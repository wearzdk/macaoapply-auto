package imageText

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/golang/freetype/truetype"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// 插入白色背景和文字
func InsertTextToImage(imageBytes []byte, text string) ([]byte, error) {
	src, err := jpeg.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}

	r := image.Rect(0, 0, src.Bounds().Dx(), src.Bounds().Dy())

	img := image.NewRGBA(r)
	draw.Draw(img, r, src, image.Point{}, draw.Src)

	col := color.RGBA{255, 255, 255, 0xff} // 白色背景

	// 绘制文字背景 宽度100% 高度10% 颜色黑色
	draw.Draw(img, image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()/10), &image.Uniform{color.RGBA{0, 0, 0, 0xff}}, image.Point{}, draw.Src)

	// 绘制文字
	point := fixed.P(15, 25) // Positioned at the top left corner
	fontBytes, err := os.ReadFile("msyh.ttf")
	if err != nil {
		return nil, err
	}
	PingFangFont, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	d := &font.Drawer{
		Dst: img,
		Src: image.NewUniform(col),
		Face: truetype.NewFace(PingFangFont, &truetype.Options{
			Size: 20,
		}),
		Dot: point,
	}
	d.DrawString(text)

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
