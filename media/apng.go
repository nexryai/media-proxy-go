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
