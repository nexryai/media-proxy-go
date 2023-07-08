package media

import (
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	"gopkg.in/gographics/imagick.v3/imagick"
	"math"
	"strconv"
)

func convertAndResizeImage(opts *transcodeImageOpts) (*[]byte, error) {

	if opts.isAnimated {
		core.MsgDebug("isAnimated: true")
	}

	// Imagickの初期化
	imagick.Initialize()
	defer imagick.Terminate()

	// MagickWandの作成
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// 画像データを読み込む
	err := mw.ReadImageBlob(*opts.imageBufferPtr)
	if err != nil {
		return nil, fmt.Errorf("画像の読み込みに失敗しました: %v", err)
	}

	// 画像サイズを取得
	width := int(mw.GetImageWidth())
	height := int(mw.GetImageHeight())
	delay := mw.GetImageDelay()

	core.MsgDebug(fmt.Sprintf("w: %d h: %d, delay: %d", width, height, delay))

	if width > 5120 || height > 5120 {
		return nil, fmt.Errorf("too large image")
	}

	// リサイズ系処理
	var newWidth int
	var newHeight int
	var shouldResize bool

	if width > opts.widthLimit || height > opts.heightLimit {
		shouldResize = true
	}

	if shouldResize {
		// 縦横比率を計算
		aspectRatio := float64(width) / float64(height)

		// リサイズ後のサイズを計算
		newWidth = width
		newHeight = height

		// 超過量を算出
		widthExcess := width - opts.widthLimit
		heightExcess := height - opts.heightLimit

		// widthLimitとheightLimitが両方超過してる場合、超過している部分が少ない方のlimitは0にして比率を維持する
		if opts.widthLimit != 0 && opts.heightLimit != 0 {
			if width > opts.widthLimit && height > opts.heightLimit {
				if widthExcess < heightExcess {
					opts.widthLimit = 0
				} else {
					opts.heightLimit = 0
				}
			}
		}

		// ChatGPTが考えてくれた
		if opts.widthLimit != 0 {
			if width > opts.widthLimit {
				newWidth = opts.widthLimit
				newHeight = int(math.Round(float64(newWidth) / aspectRatio))
			}
		} else if opts.heightLimit != 0 {
			if height > opts.heightLimit {
				newHeight = opts.heightLimit
				newWidth = int(math.Round(float64(newHeight) * aspectRatio))
			}
		}

		core.MsgDebug(fmt.Sprintf("newWidth: %d newHeight: %d aspectRatio: %v", newWidth, newHeight, aspectRatio))
	}

	if opts.isAnimated {
		core.MsgDebug("Encode as animated image!")

		aw := mw.CoalesceImages()
		mw.Destroy()
		defer aw.Destroy()

		// 新世界を創造する
		mw = imagick.NewMagickWand()
		mw.SetImageDelay(delay)
		defer mw.Destroy()

		for i := 0; i < int(aw.GetNumberImages()); i++ {
			core.MsgDebug("Encode animated image frame: " + strconv.Itoa(i))
			aw.SetIteratorIndex(i)
			img := aw.GetImage()
			if shouldResize {
				img.ResizeImage(uint(newWidth), uint(newHeight), imagick.FILTER_LANCZOS)
			}
			mw.AddImage(img)
			img.Destroy()
		}

		aw.Destroy()

		// WebP形式に変換
		mw.ResetIterator()

		// ToDo: この辺調整する
		mw.SetOption("webp:lossless", "false")
		//mw.SetOption("webp:method", "6")
		//mw.SetOption("webp:alpha-quality", "80")
		mw.SetFormat(opts.targetFormat)
		mw.SetIteratorIndex(0)

		// 変換後の画像データを取得
		convertedData := mw.GetImagesBlob()

		mw.Destroy()

		// なぜかこれ永遠に終わらん
		//imagick.Terminate()

		return &convertedData, nil

	} else {
		core.MsgDebug("Encode as static image!")

		// 画像をリサイズ
		if shouldResize {
			err = mw.ResizeImage(uint(newWidth), uint(newHeight), imagick.FILTER_LANCZOS)
			if err != nil {
				return nil, fmt.Errorf("画像のリサイズに失敗しました: %v", err)
			}
		}
	}

	// WebP形式に変換
	mw.ResetIterator()

	// ToDo: この辺調整する
	mw.SetOption("webp:lossless", "false")
	//mw.SetOption("webp:method", "6")
	//mw.SetOption("webp:alpha-quality", "80")
	mw.SetFormat(opts.targetFormat)

	// 変換後の画像データを取得
	convertedData := mw.GetImageBlob()

	mw.Destroy()
	imagick.Terminate()

	return &convertedData, nil
}
