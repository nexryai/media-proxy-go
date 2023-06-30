package media

import (
	"io/ioutil"
	"testing"
)

func TestIsAPNG(t *testing.T) {
	imageBuffer, _, err := fetchImage("https://apng.onevcat.com/assets/elephant.png")
	if err != nil {
		t.Errorf("fetchImage returned an error: %v", err)
	}

	fetchedImage, err := ioutil.ReadAll(imageBuffer)
	if err != nil {
		t.Errorf("ioutil returned an error: %v", err)
	}

	imageBuffer.Reset()

	if !isAnimatedPNG(&fetchedImage) {
		t.Errorf("isAnimatedPng returned incorrect results")
	}
}

func TestIsAnimatedWebP(t *testing.T) {
	imageBuffer, _, err := fetchImage("https://mathiasbynens.be/demo/animated-webp-supported.webp")
	if err != nil {
		t.Errorf("fetchImage returned an error: %v", err)
	}

	fetchedImage, err := ioutil.ReadAll(imageBuffer)
	if err != nil {
		t.Errorf("ioutil returned an error: %v", err)
	}

	imageBuffer.Reset()

	if !isAnimatedWebP(&fetchedImage) {
		t.Errorf("isAnimatedPng returned incorrect results")
	}
}
