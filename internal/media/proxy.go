package media

import (
	"fmt"
	"git.sda1.net/media-proxy-go/internal/logger"
	"github.com/google/uuid"
	"os"
)

// trueであればImageMagickによって処理される。そうでなければそのままプロキシされる
func isConvertible(contentType string, fetchedImage *[]byte) bool {
	switch contentType {
	case "image/avif":
		return true
	case "image/ico":
		return true
	case "image/jpeg":
		return true
	case "image/heif":
		return true
	case "image/png":
		return true
	case "image/webp":
		return true
	case "image/gif":
		return true
	case "image/svg+xml":
		return true
	case "image/x-icon":
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

	// どこかでoptsが変わるとキャッシュキーが狂って困る
	var originalOpts ProxyRequest
	originalOpts = *opts

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
		convertedImageBuffer, err := convertAndResizeImage(encodeOpts)
		cacheId := uuid.NewString()

		if err != nil {
			log.Warn(fmt.Sprintf("Failed to resize image: %v", err))

			// 失敗したならFAILEDをキャッシュする
			err = StoreCachePath(&originalOpts, "FAILED")
			if err != nil {
				log.ErrorWithDetail("Failed to store error status", err)
				return "", contentType, fmt.Errorf("failed to store error status")
			}

			return "", contentType, fmt.Errorf("failed to resize image")
		} else {
			log.Debug("ok.")
		}

		// ディスクに保存
		file, err := os.Create(GetPathFromCacheId(cacheId))
		if err != nil {
			log.ErrorWithDetail("Failed to create file", err)
			_ = StoreCachePath(&originalOpts, "FAILED")
			return "", contentType, fmt.Errorf("failed to create file")
		}

		defer file.Close()

		_, err = file.Write(*convertedImageBuffer)
		if err != nil {
			log.ErrorWithDetail("Failed to write image", err)
			_ = StoreCachePath(&originalOpts, "FAILED")
			return "", contentType, fmt.Errorf("failed to write image")
		}

		contentType = fmt.Sprintf("image/%s", encodeOpts.targetFormat)
		err = StoreCachePath(&originalOpts, cacheId)
		if err != nil {
			log.ErrorWithDetail("Failed to store cache path", err)
			return "", contentType, fmt.Errorf("failed to store cache path")
		}

		return cacheId, contentType, nil

	}

	// 画像ではない
	log.Debug(contentType + " is not supported")

	// FAILEDをキャッシュする
	err = StoreCachePath(&originalOpts, "FAILED")
	if err != nil {
		log.ErrorWithDetail("Failed to store error status", err)
		return "", contentType, fmt.Errorf("failed to store error status")
	}

	return "", contentType, fmt.Errorf("invalid file format")
}
