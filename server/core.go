package server

import (
	"encoding/json"
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	"git.sda1.net/media-proxy-go/media"
	"git.sda1.net/media-proxy-go/security"
	"github.com/valyala/fasthttp"
	"runtime"
)

func RequestHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())

	core.MsgInfo(fmt.Sprintf("Handled request: %s", path))

	switch path {
	case "/status":
		status := Status{Status: "OK"}
		jsonData, err := json.Marshal(status)
		if err != nil {
			ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
			return
		}
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Write(jsonData)

	default:
		queryArgs := ctx.QueryArgs()
		url := string(queryArgs.Peek("url"))
		isAvatar := string(queryArgs.Peek("avatar")) == "1"
		isEmoji := string(queryArgs.Peek("emoji")) == "1"
		isStatic := string(queryArgs.Peek("static")) == "1"
		isPreview := string(queryArgs.Peek("preview")) == "1"
		isBadge := string(queryArgs.Peek("badge")) == "1"

		// v13のプロキシ仕様にはないがv12はこれを使う?ため
		isThumbnail := string(queryArgs.Peek("thumbnail")) == "1"

		// 将来的にデフォルトになるので消す
		useAvif := string(queryArgs.Peek("avif")) == "1"

		if url == "" {
			ctx.Error("Bad request", fasthttp.StatusBadRequest)
			return
		}

		// ポートが指定されている、ホスト名がプライベートアドレスを示している場合はブロック
		if !security.IsSafeUrl(url) {
			ctx.Error("Access denied", fasthttp.StatusForbidden)
			return
		}

		targetFormat := "webp"
		if useAvif {
			// 試験的
			targetFormat = "avif"
		}

		var proxiedImage *[]byte
		var contentType string
		var err error

		// どこかでpanicになった場合の処理
		defer func() {
			if r := recover(); r != nil {
				// パニックが発生した場合、エラーレスポンスを返す
				core.MsgErr(fmt.Sprintf("Panic occurred while proxying media: %s", r.(error)))
				ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
				return
			}
		}()

		var widthLimit, heightLimit int

		switch {
		case isAvatar:
			widthLimit, heightLimit = 320, 320
		case isEmoji:
			widthLimit, heightLimit = 128, 128
		case isPreview:
			widthLimit, heightLimit = 200, 200
		case isBadge:
			widthLimit, heightLimit = 96, 96
		case isThumbnail:
			widthLimit, heightLimit = 500, 400
		default:
			widthLimit, heightLimit = 3200, 3200
		}

		options := &media.ProxyOpts{
			Url:          url,
			WidthLimit:   widthLimit,
			HeightLimit:  heightLimit,
			IsStatic:     isStatic,
			TargetFormat: targetFormat,
			IsEmoji:      isEmoji,
		}

		proxiedImage, contentType, err = media.ProxyImage(options)

		if err != nil {
			ctx.Error("Bad request", fasthttp.StatusBadRequest)
			return
		}

		ctx.Response.Header.SetContentType(contentType)
		ctx.Response.Header.Set("CDN-Cache-Control", "max-age=604800")
		ctx.Response.Header.Set("Cache-Control", "max-age=432000")
		ctx.Response.SetBody(*proxiedImage)

		// これ効果なさそう？
		runtime.GC()

	}
}
