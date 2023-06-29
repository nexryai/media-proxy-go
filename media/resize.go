package media

import (
	"bytes"
	"git.sda1.net/media-proxy-go/core"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/sizeofint/webpanimation"
	"golang.org/x/image/draw"
	"image"
	"log"
	"path/filepath"
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

func ProxyImage(url string, widthLimit int, heightLimit int, isStatic bool) []byte {

	var img image.Image

	extension := filepath.Ext(url)

	if extension == ".gif" && !isStatic {
		var buf bytes.Buffer

		gif, err := fetchGifImage(url)
		if err != nil {
			core.MsgErrWithDetail(err, "Faild to download image")
			return nil
		}

		webpanim := webpanimation.NewWebpAnimation(gif.Config.Width, gif.Config.Height, gif.LoopCount)
		//webpanim.WebPAnimEncoderOptions.SetKmin(9999)
		//webpanim.WebPAnimEncoderOptions.SetKmax(9999)
		defer webpanim.ReleaseMemory() // これないとメモリリークする
		webpConfig := webpanimation.NewWebpConfig()
		webpConfig.SetLossless(1)

		timeline := 0

		for i, img := range gif.Image {

			err = webpanim.AddFrame(img, timeline, webpConfig)
			if err != nil {
				log.Fatal(err)
			}
			timeline += gif.Delay[i] * 10
		}
		err = webpanim.AddFrame(nil, timeline, webpConfig)
		if err != nil {
			log.Fatal(err)
		}

		err = webpanim.Encode(&buf) // encode animation and write result bytes in buffer
		if err != nil {
			log.Fatal(err)
		}

		return buf.Bytes()

	} else {
		decodedImage, _, err := fetchImage(url)
		if err != nil {
			core.MsgErrWithDetail(err, "Faild to download image")
			return nil
		}
		img = decodedImage
	}

	// widthLimitかheightLimitを超えている場合のみ処理する
	if img.Bounds().Dx() > widthLimit || img.Bounds().Dy() > heightLimit {
		resizedImg := resizeImage(img, widthLimit, heightLimit)
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
