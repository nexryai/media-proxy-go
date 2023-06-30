package media

import (
	"git.sda1.net/media-proxy-go/core"
	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/webp"
	// ref: https://github.com/strukturag/libheif/issues/824
	// _ "github.com/strukturag/libheif/go/heif"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

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
