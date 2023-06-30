package media

import (
	"io"
	"io/ioutil"
	"strings"
)

func isAnimatedPNG(bodyReader io.Reader) bool {

	// レスポンスのデータを取得
	data, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return false
	}

	// PNGヘッダーをチェック
	isPNG := strings.HasPrefix(string(data), "\x89PNG\r\n\x1a\n")
	if !isPNG {
		return false
	}

	// APNGのシグネチャをチェック
	isAPNG := strings.Contains(string(data), "acTL")
	return isAPNG
}

func isAnimatedWebP(bodyReader io.Reader) bool {
	header := make([]byte, 34)
	_, err := bodyReader.Read(header)
	if err != nil {
		return false
	}

	// ファイルのヘッダーをチェック
	// Animated WebPの場合、ファイルの0x1Eから0x22がANIMになる
	isAnimated := string(header[0x1E:0x22]) == "ANIM"
	return isAnimated
}
