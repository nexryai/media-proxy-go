package media

import (
	"bytes"
	"golang.org/x/image/webp"
	"testing"
)

func testProxyImageDecodingFromUrl(t *testing.T, url string, widthLimit int, expectedHeight int, isStatic bool) {

	options := &ProxyOpts{
		Url:          url,
		WidthLimit:   widthLimit,
		HeightLimit:  expectedHeight,
		IsStatic:     isStatic,
		TargetFormat: "webp",
	}

	// ProxyImage関数の呼び出し
	imageBytes, _, _ := ProxyImage(options)
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
	if imgHeight != expectedHeight {
		t.Errorf("Decoded image height not corrent: expectedHeight: %d, height: %d", expectedHeight, imgHeight)
	}
	if imgWidth > widthLimit {
		t.Errorf("Decoded image size exceeds the limits: width: %d, height: %d", imgWidth, imgHeight)
	}
}

func TestProxyImageDecoding(t *testing.T) {
	// 比率チェックも兼ねてる
	testProxyImageDecodingFromUrl(t, "https://s3.sda1.net/firefish/contents/5dbff670-9539-496e-b625-97c59ff7804b.png", 360, 203, false)
	testProxyImageDecodingFromUrl(t, "https://s3.sda1.net/firefish/contents/4525b647-e47f-4fe7-b6d8-77de1fdfb102.png", 1024, 576, false)
	testProxyImageDecodingFromUrl(t, "https://s3.sda1.net/smelt/contents/bf911d10-9faa-4e8c-b6b2-9d0021355f16.jpg", 800, 342, false)
	testProxyImageDecodingFromUrl(t, "https://s3.sda1.net/firefish/contents/7afe76e0-7d6f-4827-bc72-7c9d613aa7b9.jpg", 1280, 1280, false)
}

func TestAnimatedImageProxy(t *testing.T) {
	options := &ProxyOpts{
		Url:          "https://www.easygifanimator.net/images/samples/video-to-gif-sample.gif",
		WidthLimit:   80,
		HeightLimit:  80,
		IsStatic:     false,
		TargetFormat: "webp",
	}

	// ProxyImage関数の呼び出し
	imageBytes, _, _ := ProxyImage(options)
	if imageBytes == nil {
		t.Fatal("Failed to fetch and encode image")
	}

	if !isAnimatedWebP(imageBytes) {
		t.Errorf("Failed to encode as animated WebP")
	}
}
