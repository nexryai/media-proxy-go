package security

import "testing"

func testIsSafeUrlResult(t *testing.T, url string, expectedResult bool) {
	result := IsSafeUrl(url)
	if result != expectedResult {
		t.Errorf("isSafeUrl(%s) = %t, expected %t", url, result, expectedResult)
	}
}

func TestIsSafeUrl(t *testing.T) {
	testIsSafeUrlResult(t, "https://google.com", true)
	testIsSafeUrlResult(t, "https://sda1.net:443", true)
	testIsSafeUrlResult(t, "http://192.168.1.1", false)
	testIsSafeUrlResult(t, "http://127.0.0.1", false)
	testIsSafeUrlResult(t, "http://0.0.0.0", false)
	testIsSafeUrlResult(t, "http://localhost", false)
	testIsSafeUrlResult(t, "http://192.168.1.1:8080", false)
	testIsSafeUrlResult(t, "https://test.sda1.net:3000", false)
}
