package media

import (
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	"github.com/google/uuid"
	"os"
	"os/exec"
	"strconv"
)

func convertWithFfmpeg(opts *ffmpegOpts) (*[]byte, error) {

	// libsvtav1はパイプ出力をサポートしていないため一時ファイルを使う
	tmpFileId := uuid.New()
	tmpFilePath := fmt.Sprintf("/var/mediaproxy-tmp/%s.avif", tmpFileId.String())

	ffmpegArgs := []string{}

	ffmpegArgs = append(ffmpegArgs, "-i", "pipe:0", "-vcodec", opts.encoder)

	if opts.shouldResize {
		ffmpegArgs = append(ffmpegArgs, "-vf", fmt.Sprintf("scale=%d*dar:%d", opts.height, opts.height))
	}

	ffmpegArgs = append(ffmpegArgs, "-loop", "0", "-crf", strconv.Itoa(int(opts.ffmpegCrf)), "-f", opts.targetFormat, tmpFilePath)

	cmd := exec.Command("ffmpeg", ffmpegArgs...)

	// パイプ周り
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

	// 変換後のデータを一時ファイルから読んで消す
	ffmpegOut, err := os.ReadFile(tmpFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading from stdout: %v", err)
	}

	if !core.IsDebugMode() {
		err = os.Remove(tmpFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to remove tmp file: %v", err)
		}
	} else {
		core.MsgDebug(fmt.Sprintf("tmpFilePath: %s", tmpFilePath))
	}

	return &ffmpegOut, nil

}
