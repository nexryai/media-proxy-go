package media

import (
	"bytes"
	"image/jpeg"
	"image/png"
	"testing"
)

func testFetchImageFromUrl(t *testing.T, url string, expectedContentType string) {
	imageBufferPtr, contentType, err := fetchImage(url)
	if err != nil {
		t.Errorf("fetchImage returned an error: %v", err)
	}

	// 結果の検証
	if contentType != expectedContentType {
		t.Errorf("Unexpected content type. Expected: %s, Got: %s", expectedContentType, contentType)
	}

	imageReader := bytes.NewReader(*imageBufferPtr)

	// デコードできるか検証
	if expectedContentType == "image/png" {
		_, errDecode := png.Decode(imageReader)
		if errDecode != nil {
			t.Errorf("Failed to decode fetched png %v", errDecode)
		}
	} else if expectedContentType == "image/jpeg" {
		_, errDecode := jpeg.Decode(imageReader)
		if errDecode != nil {
			t.Errorf("Failed to decode fetched jpeg %v", errDecode)
		}
	}

}

func TestFetchImage(t *testing.T) {
	// pngデコードテスト
	testFetchImageFromUrl(t, "https://s3.sda1.net/firefish/contents/e853b975-6db7-4911-a7c2-0c77b8033201.png", "image/png")

	// jpeg
	testFetchImageFromUrl(t, "https://upload.wikimedia.org/wikipedia/en/a/a9/Example.jpg", "image/jpeg")

	// avif
	testFetchImageFromUrl(t, "https://s3.sda1.net/nyan/contents/bc0701f3-6a5e-471e-ae35-ddfae7d0b7f6.avif", "image/avif")
}

func TestDownloadFile(t *testing.T) {
	_, err := downloadFile("https://s3.sda1.net/firefish/contents/e853b975-6db7-4911-a7c2-0c77b8033201.png", 30)
	if err == nil {
		t.Errorf("File size limit not working!!!")
	}
}
