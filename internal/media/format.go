package media

import (
	"bytes"
)

func isAnimatedPNG(data *[]byte) bool {
	// PNGヘッダーをチェック
	isPNG := bytes.HasPrefix(*data, []byte("\x89PNG\r\n\x1a\n"))
	if !isPNG {
		return false
	}

	// APNGのシグネチャをチェック
	return len(*data) > 41 && string((*data)[37:41]) == "acTL"
}

func isAnimatedWebP(data *[]byte) bool {
	// ヘッダーをチェック
	// Animated WebPの場合、ファイルの0x1Eから0x22がANIMになる
	return len(*data) > 0x22 && string((*data)[0x1E:0x22]) == "ANIM"
}

func isAVIF(data *[]byte) bool {
	return len(*data) > 12 && string((*data)[4:12]) == "ftypavif"
}

func isUsesColorPalette(data *[]byte) bool {
	// PNGの場合はPLTEチャンクが存在するかで判断
	return len(*data) > 64 && string((*data)[57:61]) == "PLTE"
}
