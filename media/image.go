package media

import (
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	"github.com/davidbyttow/govips/v2/vips"
	"math"
)

func isTooBigFile(d *[]byte) bool {
	// 2MB以上ならTrue
	return len(*d) >= 2*1024*1024
}

func convertAndResizeImage(opts *transcodeImageOpts) (*[]byte, error) {

	var image *vips.ImageRef
	var err error

	params := vips.NewImportParams()
	params.NumPages.Set(-1)
	image, err = vips.LoadImageFromBuffer(*opts.imageBufferPtr, params)

	// メモリ使用量が97.5%以上なら処理を中断
	core.RaisePanicOnHighMemoryUsage(97.5)

	// バッファーから読み込み

	if err != nil {
		return nil, fmt.Errorf("failed to load image: %v", err)
	}

	defer image.Close()

	// 画像サイズを取得
	width := image.Width()
	height := image.Height()
	core.MsgDebug(fmt.Sprintf("w: %d h: %d", width, height))

	if width > 5120 || height > 5120 {
		return nil, fmt.Errorf("too large image")
	}

	// リサイズ系処理
	var scale float64
	var shouldResize bool

	if width > opts.widthLimit || height > opts.heightLimit {
		shouldResize = true
	}

	if opts.isAnimated {
		core.MsgDebug("Encode as animated image!")

		// リサイズ系処理（animated）
		var newWidth int
		var newHeight int

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

		err := image.ThumbnailWithSize(newWidth, newHeight, vips.InterestingAll, vips.SizeDown)
		if err != nil {
			return nil, err
		}

		// WebP形式に変換
		encodeOpts := vips.WebpExportParams{
			Quality:  70,
			Lossless: false, // Set to true for lossless compression
		}

		// 変換後の画像データを取得
		convertedData, _, err := image.ExportWebp(&encodeOpts)
		if err != nil {
			return nil, err
		}

		return &convertedData, nil

	} else {
		core.MsgDebug("Encode as static image!")

		// 画像をリサイズ
		if shouldResize && !opts.isAnimated {
			scale = float64(opts.widthLimit) / float64(width)
			core.MsgDebug(fmt.Sprintf("scale: %v ", scale))

			err = image.Resize(scale, vips.KernelAuto)
			if err != nil {
				return nil, err
			}
		}

		// WebP形式に変換
		encodeOpts := vips.WebpExportParams{
			Quality:  70,
			Lossless: false,
		}

		// 変換後の画像データを取得
		convertedData, _, err := image.ExportWebp(&encodeOpts)
		if err != nil {
			return nil, err
		}

		return &convertedData, nil

	}

}
