package server

import (
	"encoding/json"
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	"git.sda1.net/media-proxy-go/media"
	"git.sda1.net/media-proxy-go/security"
	"github.com/valyala/fasthttp"
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

		if url == "" {
			ctx.Error("Bad request", fasthttp.StatusBadRequest)
			return
		}

		// ポートが指定されている、ホスト名がプライベートアドレスを示している場合はブロック
		if !security.IsSafeUrl(url) {
			ctx.Error("Access denied", fasthttp.StatusForbidden)
			return
		}

		var proxiedImage *[]byte
		var contentType string

		if isAvatar {
			// アバター用
			proxiedImage, contentType = media.ProxyImage(url, 0, 320, isStatic)
		} else if isEmoji {
			// 絵文字用
			proxiedImage, contentType = media.ProxyImage(url, 0, 128, isStatic)
		} else if isPreview {
			proxiedImage, contentType = media.ProxyImage(url, 200, 200, isStatic)
		} else if isBadge {
			proxiedImage, contentType = media.ProxyImage(url, 96, 96, true)
		} else {
			// TODO: Misskeyの仕様的にはsvgでない場合、無変換でプロキシするのが望ましいらしい (ref: https://github.com/misskey-dev/media-proxy/blob/master/SPECIFICATION.md#%E5%A4%89%E6%8F%9B%E3%82%AF%E3%82%A8%E3%83%AA%E3%81%8C%E5%AD%98%E5%9C%A8%E3%81%97%E3%81%AA%E3%81%84%E5%A0%B4%E5%90%88%E3%81%AE%E6%8C%99%E5%8B%95)
			proxiedImage, contentType = media.ProxyImage(url, 3200, 3200, isStatic)
		}

		// media.ProxyImage()のどこかでpanicになった場合の処理
		defer func() {
			if r := recover(); r != nil {
				// パニックが発生した場合、エラーレスポンスを返します
				core.MsgErr(fmt.Sprintf("Panic occurred while proxying media: %s", r.(error)))
				ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
			}
		}()

		if *proxiedImage == nil {
			ctx.Error("Bad request", fasthttp.StatusBadRequest)
			return
		}

		ctx.Response.Header.SetContentType(contentType)
		ctx.Response.SetBody(*proxiedImage)

	}
}
