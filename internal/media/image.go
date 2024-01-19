package media

import (
	"bytes"
	"fmt"
	"git.sda1.net/media-proxy-go/internal/core"
	"git.sda1.net/media-proxy-go/internal/logger"
	"github.com/davidbyttow/govips/v2/vips"
	"github.com/google/uuid"
	"math"
	"os"
	"os/exec"
	"strconv"
)

func runFfmpeg(opts *ffmpegOpts, cacheId string) error {
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

	ffmpegArgs = append(ffmpegArgs, "-loop", "0", "-pix_fmt", "yuva420p", "-crf", strconv.Itoa(int(opts.ffmpegCrf)),
		"-f", opts.targetFormat, GetPathFromCacheId(cacheId))

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
		return err
	}
	defer stdin.Close()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}

	_, err = stdin.Write(*opts.imageBufferPtr)
	if err != nil {
		return fmt.Errorf("error writing to stdin: %v", err)
	}
	stdin.Close()

	// 終了を待機
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command execution error: %v", err)
	}

	return nil
}

func resizeWithFfmpeg(opts *transcodeImageOpts) (string, error) {
	log := logger.GetLogger("MediaService")
	cacheId := uuid.NewString()

	var image *vips.ImageRef
	var err error

	params := vips.NewImportParams()
	params.NumPages.Set(-1)
	image, err = vips.LoadImageFromBuffer(*opts.imageBufferPtr, params)

	// バッファーから読み込み
	if err != nil {
		return "", fmt.Errorf("failed to load image: %v", err)
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
		return "", fmt.Errorf("too large image")
	}

	// リサイズ系処理
	var shouldResize bool
	if width > opts.widthLimit || height > opts.heightLimit {
		shouldResize = true
	}

	// リサイズ系処理
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

	// 変換後の画像データを取得
	ffmpegOption := &ffmpegOpts{
		imageBufferPtr: opts.imageBufferPtr,
		shouldResize:   shouldResize,
		width:          newWidth,
		height:         newHeight,
	}

	if opts.targetFormat == "avif" {
		ffmpegOption.targetFormat = "avif"
		ffmpegOption.encoder = "libaom-av1"
		ffmpegOption.ffmpegCrf = 40
	} else {
		ffmpegOption.targetFormat = "webp"
		ffmpegOption.encoder = "libwebp"
		ffmpegOption.ffmpegCrf = 70
	}

	err = runFfmpeg(ffmpegOption, cacheId)
	if err != nil {
		return "", err
	}

	return cacheId, nil
}
