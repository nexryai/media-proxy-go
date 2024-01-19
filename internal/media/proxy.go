package media

import (
	"fmt"
	"git.sda1.net/media-proxy-go/internal/logger"
)

// trueであればImageMagickによって処理される。そうでなければそのままプロキシされる
func isConvertible(contentType string, fetchedImage *[]byte) bool {
	switch contentType {
	case "image/avif":
		return true
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

func ProxyImage(opts *ProxyRequest) (string, string, error) {
	log := logger.GetLogger("MediaService")

	imageBufferPtr, contentType, err := fetchImage(opts.Url)
	if err != nil {
		log.ErrorWithDetail("Failed to download image", err)
		return "", contentType, fmt.Errorf("failed to download image")
	}

	log.Debug("Content-Type: " + contentType)

	var isAnimated bool

	if isAnimatedFormat(contentType, imageBufferPtr) && !opts.IsStatic {
		log.Debug("Animated image!")
		isAnimated = true
		opts.TargetFormat = "webp"
	} else {
		isAnimated = false
	}

	if isConvertible(contentType, imageBufferPtr) {
		log.Debug("Use vips")

		encodeOpts := &transcodeImageOpts{
			imageBufferPtr: imageBufferPtr,
			widthLimit:     opts.WidthLimit,
			heightLimit:    opts.HeightLimit,
			originalFormat: contentType,
			targetFormat:   opts.TargetFormat,
			isAnimated:     isAnimated,
		}

		if opts.UseAVIF && !isAnimated {
			encodeOpts.targetFormat = "avif"
		}

		// isAnimatedがTrueなら1フレームずつ処理する。!isAnimatedでアニメーション画像をプロキシすると最初の1フレームだけ返ってくる
		cacheId, err := resizeWithFfmpeg(encodeOpts)

		if err != nil {
			log.Warn(fmt.Sprintf("Failed to decode image: %v", err))
			return "", contentType, fmt.Errorf("failed to decode image")
		} else {
			log.Debug("Decode ok.")
		}

		contentType = fmt.Sprintf("image/%s", encodeOpts.targetFormat)

		err = StoreCachePath(opts, cacheId)
		if err != nil {
			log.ErrorWithDetail("Failed to store cache path", err)
			return "", contentType, fmt.Errorf("failed to store cache path")
		}

		return cacheId, contentType, nil

	}

	//どれにも当てはまらない
	log.Debug(contentType + " is not supported")
	return "", contentType, fmt.Errorf("invalid file format")
}
