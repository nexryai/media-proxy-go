package media

import (
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	"git.sda1.net/media-proxy-go/security"
	"gopkg.in/gographics/imagick.v3/imagick"
)

func isConvertible(contentType string) bool {
	switch contentType {
	case "image/ico":
		// 現状icoをデコードできないのでfalseを返してデコードせずそのままプロキシするようにする
		return false
	case "image/jpeg":
		return true
	case "image/heif":
		return true
	case "image/png":
		return true
	case "image/webp":
		return true
	case "image/avif":
		return true
	case "image/gif":
		return true
	case "image/svg+xml":
		return true
	default:
		return false
	}
}

func isStaticFormat(contentType string) bool {
	switch contentType {
	case "image/ico":
		return true
	case "image/jpeg":
		return true
	case "image/heif":
		return true
	case "image/svg+xml":
		return true
	default:
		// pngはapngがあるのでfalse
		return false
	}
}

func isConvertibleAnimatedFormat(contentType string, fetchedImage *[]byte) bool {
	if contentType == "image/png" && !isAnimatedPNG(fetchedImage) {
		return true
	}
	if contentType == "image/webp" {
		return true
	}
	if contentType == "image/gif" {
		return true
	}
	return false
}

func ProxyImage(url string, widthLimit int, heightLimit int, isStatic bool, targetFormat string) (*[]byte, string, error) {

	imageBufferPtr, contentType, err := fetchImage(url)
	if err != nil {
		core.MsgErrWithDetail(err, "Failed to download image")
		return nil, contentType, fmt.Errorf("failed to download image")
	}

	core.MsgDebug("Content-Type: " + contentType)

	var isAnimated bool

	if !isStaticFormat(contentType) && isConvertibleAnimatedFormat(contentType, imageBufferPtr) && !isStatic {
		core.MsgDebug("Animated image!")
		isAnimated = true
	} else {
		isAnimated = false
	}

	if isConvertible(contentType) {
		core.MsgDebug("Use ImageMagick")
		img, err := convertAndResizeImage(imageBufferPtr, widthLimit, heightLimit, targetFormat, isAnimated)
		imagick.Terminate()

		if err != nil {
			core.MsgWarn(fmt.Sprintf("Failed to decode image: %v", err))
			return nil, contentType, fmt.Errorf("failed to decode image")
		} else {
			core.MsgDebug("Decode ok.")
		}

		contentType = fmt.Sprintf("image/%s", targetFormat)

		return img, contentType, nil

	} else if security.IsFileTypeBrowserSafe(contentType) {
		// どれにも当てはまらないかつブラウザセーフな形式ならそのままプロキシ
		// AVIFは敢えて無変換でプロキシする（サイズがwebpより小さくEdgeユーザーの存在を無視すれば変換する意義がほぼない）
		core.MsgDebug("Proxy image without transcode")
		return imageBufferPtr, contentType, nil

	}

	//どれにも当てはまらないならnilを返してクライアントに400を返す
	return nil, contentType, nil
}
