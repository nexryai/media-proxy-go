package media

import (
	"bytes"
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// URLを取得して、リーダーを返す
func fetchImage(url string) (io.ReadCloser, string, error) {
	core.MsgDebug(fmt.Sprintf("Download image: %s", url))

	resp, err := http.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch image: %v", err)
	}

	contentType := resp.Header.Get("Content-Type")

	if resp.StatusCode != http.StatusOK {
		core.MsgWarn("Failed to download image. URL: " + url + ", Status: " + resp.Status)
		resp.Body.Close()
		return nil, contentType, fmt.Errorf("failed to fetch image: error status code %d", resp.StatusCode)
	} else {
		core.MsgDebug("Request OK.")
	}

	var bodyReader io.ReadCloser = resp.Body

	// SVGなら一旦webpにする
	if contentType == "image/svg+xml" {
		body := convertSvgToWebp(resp)
		resp.Body.Close() // 元のレスポンスを閉じる
		if err != nil {
			return nil, contentType, fmt.Errorf("failed to convert SVG to WebP: %v", err)
		}

		bodyReader = ioutil.NopCloser(bytes.NewReader(body))
	}

	return bodyReader, contentType, nil
}

func saveResponseToFile(resp *http.Response, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save response to file: %v", err)
	}

	return nil
}
