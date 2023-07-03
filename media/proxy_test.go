package media

import (
	"bytes"
	"golang.org/x/image/webp"
	"testing"
)

func testProxyImageDecodingFromUrl(t *testing.T, url string, widthLimit int, heightLimit int, isStatic bool) {

	// ProxyImage関数の呼び出し
	imageBytes, _, _ := ProxyImage(url, widthLimit, heightLimit, isStatic, "webp")
	if imageBytes == nil {
		t.Fatal("Failed to fetch and encode image")
	}

	// WebP画像のデコード
	img, err := webp.Decode(bytes.NewReader(*imageBytes))
	if err != nil {
		t.Fatalf("Failed to decode returnd WebP image: %v", err)
	}

	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	// デコードされた画像のサイズチェック
	if imgWidth > widthLimit || imgHeight > heightLimit {
		t.Errorf("Decoded image size exceeds the limits: width: %d, height: %d", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

func TestProxyImageDecoding(t *testing.T) {
	// 比率チェックも兼ねてる
	testProxyImageDecodingFromUrl(t, "https://s3.sda1.net/misskey/contents/c45f5574-7bed-458e-b003-2014a13147ff.png", 360, 202, false)
}
