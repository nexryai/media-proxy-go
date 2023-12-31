package media

import (
	"bytes"
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type limitedReader struct {
	rc io.ReadCloser
	n  int64
}

// URLを取得して、バッファーのポインタを返す
func fetchImage(url string) (*[]byte, string, error) {
	core.MsgDebug(fmt.Sprintf("Download image: %s", url))

	// 現状では6MBに制限しているが変えられるようにするべきかも
	maxSize := int64(6 * 1024 * 1024)

	resp, err := downloadFile(url, maxSize)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch image: %v", err)
	}

	// これだと偽装できる
	//contentType := resp.Header.Get("Content-Type")

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

	core.MsgDebug("Detected MIME: " + contentType)

	if resp.StatusCode != http.StatusOK {
		core.MsgWarn("Failed to download image. URL: " + url + ", Status: " + resp.Status)
		return nil, contentType, fmt.Errorf("failed to fetch image: error status code %d", resp.StatusCode)
	} else {
		core.MsgDebug("Request OK.")
	}

	return &body, contentType, nil
}

// サイズ制限付きダウンローダー
func downloadFile(targetUrl string, maxSize int64) (*http.Response, error) {
	// リクエストを作成
	client := mkHttpClient()
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		return nil, err
	}

	// ユーザーエージェントを設定
	req.Header.Set("User-Agent", "Misskey-Media-Proxy-Go v0.10")

	// リクエストを送信
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// ファイルサイズが制限を超えているかチェック
	contentLength := resp.Header.Get("Content-Length")
	if contentLength != "" {
		length, err := strconv.ParseInt(contentLength, 10, 64)
		if err != nil {
			return nil, err
		}
		if length > maxSize {
			return nil, fmt.Errorf("file size exceeds the limit")
		}
	}

	resp.Body = &limitedReader{rc: resp.Body, n: maxSize}
	return resp, nil
}

// プロキシ設定に応じていい感じのhttp.Clientを生成する
func mkHttpClient() *http.Client {
	if core.GetProxyConfig() != "" {
		core.MsgDebug("Use proxy")

		proxyUrl, err := url.Parse(core.GetProxyConfig())
		if err != nil {
			core.MsgWarn("Invalid proxy config. Ignore.")
		}

		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}

		return &http.Client{
			Transport: transport,
		}

	} else {
		return &http.Client{}
	}
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
