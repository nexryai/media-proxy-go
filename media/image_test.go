package media

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func testDecodeStaticImageFromUrl(t *testing.T, url string, contentType string, expectedImageWidth int, expectedImageHeight int) {
	imageBuffer, _, err := fetchImage(url)
	if err != nil {
		t.Errorf("fetchImage returned an error: %v", err)
	}

	fetchedImage, err := ioutil.ReadAll(imageBuffer)
	if err != nil {
		t.Errorf("ioutil returned an error: %v", err)
	}

	imageBuffer.Reset()

	img, errDecode := decodeStaticImage(bytes.NewReader(fetchedImage), contentType)
	if errDecode != nil {
		t.Errorf("decodeStaticImage returned an error: %v", errDecode)
	}

	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	if imgWidth != expectedImageWidth || imgHeight != expectedImageHeight {
		t.Errorf("Decoded image different than expectedImageWidth(%dx%d): %dx%d", expectedImageWidth, expectedImageHeight, imgWidth, imgHeight)
	}
}

func TestDecodeStaticImage(t *testing.T) {
	testDecodeStaticImageFromUrl(t, "https://s3.sda1.net/misskey/contents/94f005bc-1a77-4c57-a72f-43f50cc144ea.png", "image/png", 2048, 877)
}
