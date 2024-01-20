package server

import (
	"encoding/json"
	"fmt"
	"git.sda1.net/media-proxy-go/internal/logger"
	"git.sda1.net/media-proxy-go/internal/media"
	"git.sda1.net/media-proxy-go/internal/queue"
	"github.com/hibiken/asynq"
	"github.com/nexryai/archer"
	"io"
	"net/http"
	"os"
	"time"
)

const redisAddr = "127.0.0.1:6379"

func RequestHandler(w http.ResponseWriter, req *http.Request) {
	log := logger.GetLogger("Server")
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
		if !archer.IsSafeUrl(url) {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		targetFormat := "webp"

		var contentType string
		var err error

		// どこかでpanicになった場合の処理
		defer func() {
			if r := recover(); r != nil {
				// パニックが発生した場合、エラーレスポンスを返す
				log.Fatal("Panic occurred while proxying media")
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

		options := &media.ProxyRequest{
			Url:          url,
			WidthLimit:   widthLimit,
			HeightLimit:  heightLimit,
			IsStatic:     isStatic,
			TargetFormat: targetFormat,
			UseAVIF:      useAVIF || forceAVIF,
		}

		// キャッシュが存在しない場合はキューにタスクを投げてキャッシュを作成する
		if !media.CacheExists(options) {
			log.Debug("Cache not found. Waiting for cache to be created")
			client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
			defer client.Close()

			task, err := queue.NewProxyRequestTask(options)
			if err != nil {
				log.ErrorWithDetail("Failed to create task", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			_, err = client.Enqueue(task, asynq.MaxRetry(1), asynq.Timeout(15*time.Second))
			if err != nil {
				log.ErrorWithDetail("Failed to enqueue task", err)
			}

			// キャッシュが作成されるまで待つ
			i := 0
			for {
				i += 1
				time.Sleep(1 * time.Second)
				if media.CacheExists(options) {
					log.Debug("Cache created!")
					break
				} else if i > 15 {
					log.Error("Timeout")
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}
		}

		log.Debug("Streaming from cache")
		cachePath, err := media.GetCachePath(options)

		file, err := os.Open(cachePath)
		if err != nil {
			log.ErrorWithDetail("Failed to open cache file", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		if err != nil {
			log.ErrorWithDetail("Failed to get cache path", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if options.TargetFormat == "avif" {
			contentType = "image/avif"
		} else {
			contentType = "image/webp"
		}

		w.Header().Set("Content-Type", contentType)
		w.Header().Set("CDN-Cache-Control", "max-age=604800")
		w.Header().Set("Cache-Control", "max-age=432000")

		_, err = io.Copy(w, file)
		if err != nil {
			log.ErrorWithDetail("Failed to write response", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
