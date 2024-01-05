package server

import (
	"encoding/json"
	"fmt"
	"git.sda1.net/media-proxy-go/media"
	"git.sda1.net/media-proxy-go/security"
	"net/http"
)

func RequestHandler(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	fmt.Printf("Handled request: %s\n", path)

	switch path {
	case "/status":
		status := Status{Status: "OK"}
		jsonData, err := json.Marshal(status)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)

	default:
		queryArgs := req.URL.Query()
		url := queryArgs.Get("url")
		isAvatar := queryArgs.Get("avatar") == "1"
		isEmoji := queryArgs.Get("emoji") == "1"
		isStatic := queryArgs.Get("static") == "1"
		isPreview := queryArgs.Get("preview") == "1"
		isBadge := queryArgs.Get("badge") == "1"
		isTicker := queryArgs.Get("ticker") == "1"

		// v13のプロキシ仕様にはないがv12はこれを使う?ため
		isThumbnail := queryArgs.Get("thumbnail") == "1"

		forceAVIF := queryArgs.Get("avif") == "1"

		if url == "" {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// ポートが指定されている、ホスト名がプライベートアドレスを示している場合はブロック
		if !security.IsSafeUrl(url) {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		targetFormat := "webp"

		var proxiedImage *[]byte
		var contentType string
		var err error

		// どこかでpanicになった場合の処理
		defer func() {
			if r := recover(); r != nil {
				// パニックが発生した場合、エラーレスポンスを返す
				fmt.Printf("Panic occurred while proxying media: %s\n", r)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}()

		var widthLimit, heightLimit int
		// 巨大画像をつかうとたまにAVIFのエンコードに時間がすごいかかるので大きめの画像がありそうならWebPを使う
		var useAVIF bool

		switch {
		case isAvatar:
			widthLimit, heightLimit = 320, 320
			useAVIF = true
		case isEmoji:
			widthLimit, heightLimit = 700, 128
			useAVIF = true
		case isPreview:
			widthLimit, heightLimit = 200, 200
			useAVIF = false
		case isBadge:
			widthLimit, heightLimit = 96, 96
			useAVIF = true
		case isThumbnail:
			widthLimit, heightLimit = 500, 400
			useAVIF = false
		case isTicker:
			widthLimit, heightLimit = 64, 64
			useAVIF = true
		default:
			widthLimit, heightLimit = 3200, 3200
			useAVIF = false
		}

		if forceAVIF || useAVIF {
			targetFormat = "avif"
		}

		options := &media.ProxyOpts{
			Url:          url,
			WidthLimit:   widthLimit,
			HeightLimit:  heightLimit,
			IsStatic:     isStatic,
			TargetFormat: targetFormat,
			UseAVIF:      useAVIF || forceAVIF,
		}

		proxiedImage, contentType, err = media.ProxyImage(options)

		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("CDN-Cache-Control", "max-age=604800")
		w.Header().Set("Cache-Control", "max-age=432000")
		w.Write(*proxiedImage)
	}
}
