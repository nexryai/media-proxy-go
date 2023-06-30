package security

import "testing"

func testIsFileTypeBrowserSafeResult(t *testing.T, fileType string, expectedResult bool) {
	result := IsFileTypeBrowserSafe(fileType)
	if result != expectedResult {
		t.Errorf("isSafeUrl(%s) = %t, expected %t", fileType, result, expectedResult)
	}
}

func TestIsFileTypeBrowserSafe(t *testing.T) {
	testIsFileTypeBrowserSafeResult(t, "audio/flac", true)
	testIsFileTypeBrowserSafeResult(t, "image/avif", true)
	testIsFileTypeBrowserSafeResult(t, "application/octet-stream", false)
	testIsFileTypeBrowserSafeResult(t, "text/javascript", false)
	testIsFileTypeBrowserSafeResult(t, "dummy/dummy", false)
}
