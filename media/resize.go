package media

import (
	"golang.org/x/image/draw"
	"image"
)

func resizeImage(img *image.Image, width, height int) image.Image {
	bounds := (*img).Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	// 縦横どちらかが0ならアスペクト比を保つよう適切な値を設定する
	if width == 0 {
		width = imgWidth * height / imgHeight
	} else if height == 0 {
		height = imgHeight * width / imgWidth
	}

	resizedImg := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(resizedImg, resizedImg.Bounds(), *img, (*img).Bounds(), draw.Over, nil)

	return resizedImg
}
