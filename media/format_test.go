package media

import (
	"testing"
)

func TestIsAPNG(t *testing.T) {
	imageBufferPtr, _, err := fetchImage("https://apng.onevcat.com/assets/elephant.png")
	if err != nil {
		t.Errorf("fetchImage returned an error: %v", err)
	}

	if !isAnimatedPNG(imageBufferPtr) {
		t.Errorf("isAnimatedPng returned incorrect results")
	}
}

func TestIsAnimatedWebP(t *testing.T) {
	imageBufferPtr, _, err := fetchImage("https://mathiasbynens.be/demo/animated-webp-supported.webp")
	if err != nil {
		t.Errorf("fetchImage returned an error: %v", err)
	}

	if !isAnimatedWebP(imageBufferPtr) {
		t.Errorf("isAnimatedPng returned incorrect results")
	}
}
func TestIsAVIF(t *testing.T) {
	imageBufferPtr, downloaderDetectedType, err := fetchImage("https://s3.sda1.net/nyan/contents/bc0701f3-6a5e-471e-ae35-ddfae7d0b7f6.avif")
	if err != nil {
		t.Errorf("fetchImage returned an error: %v", err)
	}

	if !isAVIF(imageBufferPtr) {
		t.Errorf("isAVIF returned incorrect results")
	}

	if downloaderDetectedType != "image/avif" {
		t.Errorf("fetchImage returned incorrect results")
	}
}
