package media

import (
	"bytes"
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	"git.sda1.net/media-proxy-go/security"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"image"
)

func isStaticFormat(contentType string) bool {
	switch contentType {
	case "image/ico":
		// 現状icoをデコードできないのでfalseを返してデコードせずそのままプロキシするようにする
		return false
	case "image/jpeg":
		return true
	case "image/heif":
		return true
	default:
		// pngはapngがあるのでfalse
		return false
	}
}

func isStaticImage(contentType string, fetchedImage *[]byte) bool {
	if contentType == "image/png" && !isAnimatedPNG(fetchedImage) {
		return true
	}
	if contentType == "image/webp" && !isAnimatedWebP(fetchedImage) {
		return true
	}
	core.MsgDebug("Animated")
	return false
}

func ProxyImage(url string, widthLimit int, heightLimit int, isStatic bool) (*[]byte, string) {

	var img *image.Image

	imageBufferPtr, contentType, err := fetchImage(url)
	if err != nil {
		core.MsgErrWithDetail(err, "Failed to download image")
		return nil, contentType
	}

	core.MsgDebug("Content-Type: " + contentType)

	if contentType == "image/svg+xml" {
		// TODO: SVG対応
		return nil, contentType

	} else if isStaticFormat(contentType) || isStaticImage(contentType, imageBufferPtr) || isStatic {

		// 完全に静止画像のフォーマット or apngでない or static指定 ならdecodeStaticImageでデコードする
		img, err = decodeStaticImage(imageBufferPtr, contentType)

		if err != nil {
			core.MsgWarn(fmt.Sprintf("Failed to decode image: %v", err))
			return nil, contentType
		} else {
			core.MsgDebug("Decode ok.")
		}

		// widthLimitかheightLimitを超えている場合のみ処理する
		core.MsgDebug(fmt.Sprintf("widthLimit: %d  heightLimit: %d", widthLimit, heightLimit))

		if (*img).Bounds().Dx() > widthLimit || (*img).Bounds().Dy() > heightLimit {
			resizedImg := resizeImage(img, widthLimit, heightLimit)
			img = &resizedImg
		}

		var buf bytes.Buffer

		// TODO: エンコードオプション変えられるようにする
		options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
		if err != nil {
			return nil, contentType
		}

		errEncode := webp.Encode(&buf, *img, options)
		if errEncode != nil {
			return nil, contentType
		}

		imageBytes := buf.Bytes()
		buf.Reset()

		return &imageBytes, contentType

	} else if security.IsFileTypeBrowserSafe(contentType) {
		// どれにも当てはまらないかつブラウザセーフな形式ならそのままプロキシ
		// AVIFは敢えて無変換でプロキシする（サイズがwebpより小さくEdgeユーザーの存在を無視すれば変換する意義がほぼない）
		core.MsgDebug("Proxy image without transcode")
		return imageBufferPtr, contentType

	}

	//どれにも当てはまらないならnilを返してクライアントに400を返す
	return nil, contentType
}
