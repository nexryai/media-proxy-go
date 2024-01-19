package media

import (
	"bytes"
	"fmt"
	"git.sda1.net/media-proxy-go/internal/logger"
	"github.com/nexryai/archer"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
)

// URLを取得して、バッファーのポインタを返す
func fetchImage(url string) (*[]byte, string, error) {
	log := logger.GetLogger("MediaService")
	log.Debug(fmt.Sprintf("Download image: %s", url))

	// 現状では6MBに制限しているが変えられるようにするべきかも
	maxSize := int64(6 * 1024 * 1024)

	resp, err := downloadFile(url, maxSize)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch image: %v", err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		// エラーハンドリング
		return nil, "", fmt.Errorf("failed to fetch image: %v", err)
	}

	body := buf.Bytes()

	resp.Body.Close()
	buf.Reset()

	contentType := http.DetectContentType(body)

	// SVGが平文扱いになるのをなんとかする
	if contentType == "text/plain; charset=utf-8" && resp.Header.Get("Content-Type") == "image/svg+xml" {
		contentType = "image/svg+xml"
	} else if contentType == "application/octet-stream" {
		// http.DetectContentTypeがAVIFを正しく認識しないのでバッファーの先頭を読んで気合で判別する
		if isAVIF(&body) {
			contentType = "image/avif"
		}
	}

	log.Debug("Detected MIME: " + contentType)

	if resp.StatusCode != http.StatusOK {
		log.Warn("Failed to download image. URL: " + url + ", Status: " + resp.Status)
		return nil, contentType, fmt.Errorf("failed to fetch image: error status code %d", resp.StatusCode)
	} else {
		log.Debug("Request OK.")
	}

	return &body, contentType, nil
}

// サイズ制限付きダウンローダー
func downloadFile(targetUrl string, maxSize int64) (*http.Response, error) {
	// リクエストを作成
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		return nil, err
	}

	// ユーザーエージェントを設定
	req.Header.Set("User-Agent", "Misskey-Media-Proxy-Go v0.10")

	secureRequestService := archer.SecureRequest{
		Request: req,
		MaxSize: maxSize,
	}

	// リクエストを送信
	resp, err := secureRequestService.Send()
	if err != nil {
		return nil, err
	}

	return resp, nil
}
