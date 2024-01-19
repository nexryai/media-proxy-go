package media

import (
	"bytes"
	"fmt"
	"git.sda1.net/media-proxy-go/internal/core"
	"git.sda1.net/media-proxy-go/internal/logger"
	"github.com/davidbyttow/govips/v2/vips"
	"math"
	"os"
	"os/exec"
	"strconv"
)

func convertWithFfmpeg(opts *ffmpegOpts) (*[]byte, error) {
	log := logger.GetLogger("ffmpeg")

	ffmpegArgs := []string{"-i", "pipe:0"}

	// 奇数だとエラーになるので偶数にする
	if opts.height%2 != 0 {
		opts.height -= 1
		opts.shouldResize = true
	} else if opts.width%2 != 0 {
		// オプションに奇数を指定しなくても元画像の幅が奇数かつリサイズ無しでプロキシしようとするとエラーになるっぽい
		opts.shouldResize = true
	}

	if opts.shouldResize {
		ffmpegArgs = append(ffmpegArgs, "-vf", fmt.Sprintf("scale=-2:%d", opts.height))
	}

	ffmpegArgs = append(ffmpegArgs, "-loop", "0", "-pix_fmt", "yuva420p", "-crf", strconv.Itoa(int(opts.ffmpegCrf)), "-f", opts.targetFormat, "-")

	cmd := exec.Command("ffmpeg", ffmpegArgs...)
	log.Debug(fmt.Sprintf("ffmpeg args: %s", ffmpegArgs))

	// パイプ周り
	var stdoutBuffer bytes.Buffer
	cmd.Stdout = &stdoutBuffer

	if core.IsDebugMode() {
		cmd.Stderr = os.Stderr
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	defer stdin.Close()

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %v", err)
	}

	_, err = stdin.Write(*opts.imageBufferPtr)
	if err != nil {
		return nil, fmt.Errorf("error writing to stdin: %v", err)
	}
	stdin.Close()

	// 終了を待機
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("command execution error: %v", err)
	}

	result := stdoutBuffer.Bytes()

	return &result, nil

}

func convertAndResizeImage(opts *transcodeImageOpts) (*[]byte, error) {
	log := logger.GetLogger("MediaService")

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
	var width int
	var height int

	if opts.isAnimated {
		numFrames := image.Pages()
		log.Debug(fmt.Sprintf("frames: %d", numFrames))

		width = image.Width()
		// vipsにアニメーション画像を読ませると全部のフレームの合計が高さとして認識されるのでフレーム数で割る
		height = image.Height() / numFrames

	} else {
		width = image.Width()
		height = image.Height()
	}

	log.Debug(fmt.Sprintf("w: %d h: %d", width, height))

	if width > 5120 || height > 5120 {
		return nil, fmt.Errorf("too large image")
	}

	// リサイズ系処理
	var scale float64
	var shouldResize bool

	if width > opts.widthLimit || height > opts.heightLimit {
		shouldResize = true
	}

	if opts.isAnimated || opts.useLibsvtav1ForAvif {
		log.Debug("Encode as animated image!")

		// リサイズ系処理（animated）
		newWidth := width
		newHeight := height

		if shouldResize {

			// 縦横比率を計算
			aspectRatio := float64(width) / float64(height)

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

			log.Debug(fmt.Sprintf("newWidth: %d newHeight: %d aspectRatio: %v", newWidth, newHeight, aspectRatio))
		}

		err = image.ThumbnailWithSize(newWidth, newHeight, vips.InterestingAll, vips.SizeDown)
		if err != nil {
			return nil, err
		}

		// WebP形式に変換
		encodeOpts := vips.WebpExportParams{
			Quality:  70,
			Lossless: false, // Set to true for lossless compression
		}

		// 変換後の画像データを取得
		convertedDataBuffer, _, err := image.ExportWebp(&encodeOpts)
		if err != nil {
			return nil, err
		}

		return &convertedDataBuffer, nil

	} else {
		log.Debug("Encode as static image!")

		// 画像をリサイズ
		if shouldResize && !opts.isAnimated {
			// 超過が大きい方に合わせる
			widthExcess := width - opts.widthLimit
			heightExcess := height - opts.heightLimit

			if widthExcess < heightExcess {
				scale = float64(opts.heightLimit) / float64(height)
			} else {
				scale = float64(opts.widthLimit) / float64(width)
			}

			log.Debug(fmt.Sprintf("scale: %v ", scale))

			err = image.Resize(scale, vips.KernelAuto)
			if err != nil {
				return nil, err
			}
		}

		if opts.targetFormat == "avif" {
			// AVIF形式に変換
			encodeOpts := vips.AvifExportParams{
				Quality:  65,
				Effort:   1,
				Lossless: false,
			}

			// 変換後の画像データを取得
			convertedData, _, err := image.ExportAvif(&encodeOpts)
			if err != nil {
				return nil, err
			}

			return &convertedData, nil
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
