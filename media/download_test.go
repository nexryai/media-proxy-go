package media

import (
	"bytes"
	"image/jpeg"
	"image/png"
	"testing"
)

func testFetchImageFromUrl(t *testing.T, url string, expectedContentType string) {
	// テスト対象の関数を呼び出し
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
	testFetchImageFromUrl(t, "https://s3.sda1.net/firefish/contents/741748a6-3185-4b19-99d7-b704e16aecbc.png", "image/png")

	// jpeg
	testFetchImageFromUrl(t, "https://upload.wikimedia.org/wikipedia/en/a/a9/Example.jpg", "image/jpeg")
}

func TestDownloadFile(t *testing.T) {
	_, err := downloadFile("https://s3.sda1.net/firefish/contents/741748a6-3185-4b19-99d7-b704e16aecbc.png", 30)
	if err == nil {
		t.Errorf("File size limit not working!!!")
	}
}
