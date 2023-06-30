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
