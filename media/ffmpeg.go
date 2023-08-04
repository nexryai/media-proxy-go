package media

import (
	"fmt"
	"io"
	"os/exec"
)

func convertToAnimatedWebP(opts *ffmpegOpts) (*[]byte, error) {

	ffmpegArgs := []string{}

	ffmpegArgs = append(ffmpegArgs, "-i", "pipe:0", "-vcodec", opts.encoder)

	if opts.shouldResize {
		ffmpegArgs = append(ffmpegArgs, "-vf", fmt.Sprintf("scale=%d*dar:%d", opts.height, opts.height))
	}

	ffmpegArgs = append(ffmpegArgs, "-loop", "0", "-pix_fmt", "yuva420p", "-f", opts.targetFormat, "pipe:1")

	cmd := exec.Command("ffmpeg", ffmpegArgs...)

	// パイプ周り
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %v", err)
	}
	// 入力データを標準入力に書き込む
	_, err = stdin.Write(*opts.imageBufferPtr)
	if err != nil {
		return nil, fmt.Errorf("Error writing to stdin: %v", err)
	}
	stdin.Close()

	// stdoutを取得
	ffmpegOut, err := io.ReadAll(stdout)
	if err != nil {
		return nil, fmt.Errorf("Error reading from stdout: %v", err)
	}

	// 終了を待機
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("Command execution error: %v", err)
	}

	return &ffmpegOut, nil

}
