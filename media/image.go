package media

import (
	"bytes"
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	"github.com/kolesa-team/go-webp/decoder"
	"github.com/kolesa-team/go-webp/webp"
	"gopkg.in/gographics/imagick.v3/imagick"
	"math"

	// ref: https://github.com/strukturag/libheif/issues/824
	// _ "github.com/strukturag/libheif/go/heif"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

// 画像をデコードし、image.Image型で返す
func decodeStaticImage(imageBufferPtr *[]byte, contentType string) (*image.Image, error) {
	var img image.Image
	var errDecode error

	imageReader := bytes.NewReader(*imageBufferPtr)

	// 適切なデコーダーを使用して画像をデコード
	switch contentType {
	case "image/webp":
		core.MsgDebug("Decode as webp")
		img, errDecode = webp.Decode(imageReader, &decoder.Options{})
	default:
		core.MsgDebug("Decode as png/jpeg/heif")
		img, _, errDecode = image.Decode(imageReader)
	}

	if errDecode != nil {
		return nil, errDecode
	}

	return &img, nil
}

func convertAndResizeImage(imageBufferPtr *[]byte, widthLimit int, heightLimit int, targetFormat string) (*[]byte, error) {
	// Imagickの初期化
	imagick.Initialize()
	defer imagick.Terminate()

	// MagickWandの作成
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// 画像データを読み込む
	err := mw.ReadImageBlob(*imageBufferPtr)
	if err != nil {
		return nil, fmt.Errorf("画像の読み込みに失敗しました: %v", err)
	}

	// 画像サイズを取得
	width := int(mw.GetImageWidth())
	height := int(mw.GetImageHeight())

	if width > 5120 || height > 5120 {
		return nil, fmt.Errorf("too large image")
	}

	if width > widthLimit || height > heightLimit {
		// 縦横比率を計算
		aspectRatio := float64(width) / float64(height)

		// リサイズ後のサイズを計算
		newWidth := width
		newHeight := height

		// 超過量を算出
		widthExcess := width - widthLimit
		heightExcess := height - heightLimit

		// widthLimitとheightLimitが両方超過してる場合、超過している部分が少ない方のlimitは0にして比率を維持する
		if widthLimit != 0 && heightLimit != 0 {
			if width > widthLimit && height > heightLimit {
				if widthExcess < heightExcess {
					widthLimit = 0
				} else {
					heightLimit = 0
				}
			}
		}

		if widthLimit != 0 {
			if width > widthLimit {
				newWidth = widthLimit
				newHeight = int(math.Round(float64(newWidth) / aspectRatio))
			}
		} else if heightLimit != 0 {
			if height > heightLimit {
				newHeight = heightLimit
				newWidth = int(math.Round(float64(newHeight) * aspectRatio))
			}
		}

		// 画像をリサイズ
		err = mw.ResizeImage(uint(newWidth), uint(newHeight), imagick.FILTER_LANCZOS)
		if err != nil {
			return nil, fmt.Errorf("画像のリサイズに失敗しました: %v", err)
		}
	}

	// WebP形式に変換
	mw.SetImageIterations(0)

	// ToDo: この辺調整する
	//mw.SetOption("webp:lossless", "false")
	//mw.SetOption("webp:method", "6")
	//mw.SetOption("webp:alpha-quality", "80")
	mw.SetFormat(targetFormat)

	// 変換後の画像データを取得
	convertedData := mw.GetImageBlob()

	mw.Destroy()
	imagick.Terminate()

	return &convertedData, nil
}
