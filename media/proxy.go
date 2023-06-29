package media

import (
	"bytes"
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"image"
)

func ProxyImage(url string, widthLimit int, heightLimit int, isStatic bool) []byte {

	var img image.Image

	fetchedImage, contentType, err := fetchImage(url)
	if err != nil {
		core.MsgErrWithDetail(err, "Failed to download image")
		return nil
	}

	if contentType == "image/gif" {
		imgBuffer, err := encodeAnimatedGifImage(fetchedImage, contentType)
		if err != nil {
			core.MsgWarn(fmt.Sprintf("Failed to convert gif to webp: %v", err))
			return nil
		} else {
			return imgBuffer
		}
	} else {
		img, err = decodeStaticImage(fetchedImage, contentType)
		if err != nil {
			core.MsgWarn(fmt.Sprintf("Failed to decode image: %v", err))
			return nil
		} else {
			core.MsgDebug("Decode ok.")
		}
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

	// buf.Bytes()を直接返すとメモリリークの原因になる
	encodedImg := make([]byte, buf.Len())
	copy(encodedImg, buf.Bytes())

	return encodedImg
}
