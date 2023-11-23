package media

import (
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	"git.sda1.net/media-proxy-go/security"
	"os"
)

// trueであればImageMagickによって処理される。そうでなければそのままプロキシされる
func isConvertible(contentType string, fetchedImage *[]byte) bool {
	switch contentType {
	case "image/ico":
		// 現状icoをデコードできないのでfalseを返してデコードせずそのままプロキシするようにする
		return false
	case "image/jpeg":
		return true
	case "image/heif":
		return true
	case "image/png":
		// apng無理。滅びろ。お前に人権はない。
		return !isAnimatedPNG(fetchedImage)
	case "image/webp":
		return true
	case "image/gif":
		return true
	case "image/svg+xml":
		return true
	default:
		// AVIFは敢えて無変換でプロキシする（サイズがwebpより小さくEdgeユーザーの存在を無視すれば変換する意義がほぼない）
		return false
	}
}

// アニメーション画像かどうか（convertAndResizeImage()がアニメーション画像としてレンダリングするかどうか決める時に使う）
func isAnimatedFormat(contentType string, fetchedImage *[]byte) bool {
	if contentType == "image/webp" {
		return isAnimatedWebP(fetchedImage)
	}
	if contentType == "image/gif" {
		return true
	}
	return false
}

func ProxyImage(opts *ProxyOpts) (*[]byte, string, error) {

	imageBufferPtr, contentType, err := fetchImage(opts.Url)
	if err != nil {
		core.MsgErrWithDetail(err, "Failed to download image")
		return nil, contentType, fmt.Errorf("failed to download image")
	}

	core.MsgDebug("Content-Type: " + contentType)

	var isAnimated bool

	if isAnimatedFormat(contentType, imageBufferPtr) && !opts.IsStatic {
		core.MsgDebug("Animated image!")
		isAnimated = true
	} else {
		isAnimated = false
	}

	if isConvertible(contentType, imageBufferPtr) {
		core.MsgDebug("Use vips (or ffmpeg)")

		encodeOpts := &transcodeImageOpts{
			imageBufferPtr: imageBufferPtr,
			widthLimit:     opts.WidthLimit,
			heightLimit:    opts.HeightLimit,
			originalFormat: contentType,
			targetFormat:   opts.TargetFormat,
			isAnimated:     isAnimated,
		}

		if os.Getenv("USE_FFMPEG") == "1" {
			// なんかffmpegだと透過画像を正しく扱えないことが多い
			if !opts.IsEmoji {
				encodeOpts.useLibsvtav1ForAvif = true
			}
			encodeOpts.targetFormat = "avif"
		}

		// isAnimatedがTrueなら1フレームずつ処理する。!isAnimatedでアニメーション画像をプロキシすると最初の1フレームだけ返ってくる
		img, err := convertAndResizeImage(encodeOpts)

		if err != nil {
			core.MsgWarn(fmt.Sprintf("Failed to decode image: %v", err))
			return nil, contentType, fmt.Errorf("failed to decode image")
		} else {
			core.MsgDebug("Decode ok.")
		}

		contentType = fmt.Sprintf("image/%s", encodeOpts.targetFormat)

		return img, contentType, nil

	} else if security.IsFileTypeBrowserSafe(contentType) {
		// どれにも当てはまらないかつブラウザセーフな形式ならそのままプロキシ
		core.MsgDebug("Proxy image without transcode")
		return imageBufferPtr, contentType, nil
	}

	//どれにも当てはまらないならnilを返してクライアントに400を返す
	return nil, contentType, fmt.Errorf("invalid file format")
}
