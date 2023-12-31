package media

import (
	"bytes"
)

func isAnimatedPNG(imageData *[]byte) bool {
	// PNGヘッダーをチェック
	isPNG := bytes.HasPrefix(*imageData, []byte("\x89PNG\r\n\x1a\n"))
	if !isPNG {
		return false
	}

	// APNGのシグネチャをチェック
	isAPNG := bytes.Contains(*imageData, []byte("acTL"))
	return isAPNG
}

func isAnimatedWebP(imageData *[]byte) bool {
	// ヘッダーをチェック
	// Animated WebPの場合、ファイルの0x1Eから0x22がANIMになる
	return string((*imageData)[0x1E:0x22]) == "ANIM"
}

func isAVIF(data *[]byte) bool {
	return len(*data) > 12 && string((*data)[4:12]) == "ftypavif"
}
