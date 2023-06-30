package media

import (
	"bytes"
	"git.sda1.net/media-proxy-go/core"
	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/sizeofint/webpanimation"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	// ref: https://github.com/strukturag/libheif/issues/824
	// _ "github.com/strukturag/libheif/go/heif"
	"image"
	"image/gif"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
)

// svgをwebpに変換する
func convertSvgToWebp(resp *http.Response) []byte {
	w, h := 512, 512

	icon, _ := oksvg.ReadIconStream(resp.Body)
	icon.SetTarget(0, 0, float64(w), float64(h))
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 1)

	var buf bytes.Buffer
	options, _ := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
	errEncode := webp.Encode(&buf, rgba, options)
	if errEncode != nil {
		return nil
	}

	return buf.Bytes()

}

// 画像をデコードし、image.Image型で返す
func decodeStaticImage(imageBuffer io.Reader, contentType string) (image.Image, error) {
	var img image.Image
	var errDecode error

	// 適切なデコーダーを使用して画像をデコード
	switch contentType {
	case "image/webp":
		core.MsgDebug("Decode as webp")
		img, errDecode = webp.Decode(imageBuffer, &decoder.Options{})
	default:
		core.MsgDebug("Decode as png/jpeg/heif")
		img, _, errDecode = image.Decode(imageBuffer)
	}

	if errDecode != nil {
		return nil, errDecode
	}

	return img, nil
}

// gif画像をAnimated webpにエンコードする
func encodeAnimatedGifImage(bodyReader io.Reader, contentType string) ([]byte, error) {

	gifImage, err := gif.DecodeAll(bodyReader)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer

	webpanim := webpanimation.NewWebpAnimation(gifImage.Config.Width, gifImage.Config.Height, gifImage.LoopCount)
	//webpanim.WebPAnimEncoderOptions.SetKmin(9999)
	//webpanim.WebPAnimEncoderOptions.SetKmax(9999)
	defer webpanim.ReleaseMemory() // これないとメモリリークする
	webpConfig := webpanimation.NewWebpConfig()
	webpConfig.SetLossless(1)

	timeline := 0

	for i, img := range gifImage.Image {
		err = webpanim.AddFrame(img, timeline, webpConfig)
		if err != nil {
			log.Fatal(err)
		}
		timeline += gifImage.Delay[i] * 10
	}

	err = webpanim.AddFrame(nil, timeline, webpConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = webpanim.Encode(&buf) // encode animation and write result bytes in buffer
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes(), nil

}
