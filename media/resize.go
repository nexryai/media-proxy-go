package media

import (
	"bytes"
	"git.sda1.net/media-proxy-go/core"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"golang.org/x/image/draw"
	"image"
)

func resizeImage(img image.Image, width, height int) image.Image {
	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	// 縦横どちらかが0ならアスペクト比を保つよう適切な値を設定する
	if width == 0 {
		width = imgWidth * height / imgHeight
	} else if height == 0 {
		height = imgHeight * width / imgWidth
	}

	resizedImg := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(resizedImg, resizedImg.Bounds(), img, img.Bounds(), draw.Over, nil)
	return resizedImg
}

func ProcessImage(url string, heightLimit int) []byte {
	img, err := fetchImage(url)
	if err != nil {
		core.MsgErrWithDetail(err, "Faild to download image")
		return nil
	}

	// heightLimitより小さい画像であるなら変換する必要はないため処理を飛ばす
	if img.Bounds().Dy() > heightLimit {
		resizedImg := resizeImage(img, 0, heightLimit)
		img = resizedImg
	}

	var buf bytes.Buffer

	// TODO: エンコードオプション変えられるようにする
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
	if err != nil {
		return nil
	}

	errEncode := webp.Encode(&buf, img, options)
	if errEncode != nil {
		return nil
	}

	return buf.Bytes()
}
