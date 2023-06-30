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
	"strconv"
)

type limitedReader struct {
	rc io.ReadCloser
	n  int64
}

// URLを取得して、リーダーを返す
func fetchImage(url string) (*bytes.Buffer, string, error) {
	core.MsgDebug(fmt.Sprintf("Download image: %s", url))

	// 現状では30MBに制限しているが変えられるようにするべきかも
	maxSize := int64(30 * 1024 * 1024)

	resp, err := downloadFile(url, maxSize)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch image: %v", err)
	}

	// これだと偽装できる
	//contentType := resp.Header.Get("Content-Type")

	body, _ := ioutil.ReadAll(resp.Body)
	contentType := http.DetectContentType(body)

	// SVGが平文扱いになるのをなんとかする
	if contentType == "text/plain; charset=utf-8" && resp.Header.Get("Content-Type") == "image/svg+xml" {
		contentType = "image/svg+xml"
	}

	core.MsgDebug("Detected MIME: " + contentType)

	if resp.StatusCode != http.StatusOK {
		core.MsgWarn("Failed to download image. URL: " + url + ", Status: " + resp.Status)
		resp.Body.Close()
		return nil, contentType, fmt.Errorf("failed to fetch image: error status code %d", resp.StatusCode)
	} else {
		core.MsgDebug("Request OK.")
	}

	//var bodyReader io.ReadCloser = resp.Body

	imageBuffer := bytes.NewBuffer(body)

	return imageBuffer, contentType, nil
}

func downloadFile(url string, maxSize int64) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	contentLength := resp.Header.Get("Content-Length")
	if contentLength != "" {
		length, err := strconv.ParseInt(contentLength, 10, 64)
		if err != nil {
			return nil, err
		}
		if length > maxSize {
			resp.Body.Close()
			return nil, fmt.Errorf("file size exceeds the limit")
		}
	}

	resp.Body = &limitedReader{rc: resp.Body, n: maxSize}
	return resp, nil
}

func (lr *limitedReader) Read(p []byte) (int, error) {
	if lr.n <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > lr.n {
		p = p[:lr.n]
	}
	n, err := lr.rc.Read(p)
	lr.n -= int64(n)
	return n, err
}

func (lr *limitedReader) Close() error {
	return lr.rc.Close()
}
