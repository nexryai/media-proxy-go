package media

import (
	"bytes"
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"image"
	"io/ioutil"
)

func ProxyImage(url string, widthLimit int, heightLimit int, isStatic bool) []byte {

	var img image.Image

	imageBuffer, contentType, err := fetchImage(url)
	if err != nil {
		core.MsgErrWithDetail(err, "Failed to download image")
		return nil
	}

	// 何回も参照できるようにコピー
	fetchedImage, err := ioutil.ReadAll(imageBuffer)
	imageBuffer.Reset()

	core.MsgDebug("Content-Type: " + contentType)

	if contentType == "image/gif" {
		imgBytes, err := encodeAnimatedGifImage(bytes.NewReader(fetchedImage), contentType)

		if err != nil {
			core.MsgWarn(fmt.Sprintf("Failed to convert gif to webp: %v", err))
			return nil
		} else {
			return imgBytes
		}
	} else if isAnimatedPNG(bytes.NewReader(fetchedImage)) && !isStatic {
		// apngかつstatic出ない場合、apngをそのまま返す
		imgBytes, err := ioutil.ReadAll(bytes.NewReader(fetchedImage))

		if err != nil {
			core.MsgWarn(fmt.Sprintf("Failed to proxy apng: %v", err))
			return nil
		} else {
			return imgBytes
		}

	} else {
		img, err = decodeStaticImage(bytes.NewReader(fetchedImage), contentType)
		imageBuffer.Reset()

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

	return buf.Bytes()
}
