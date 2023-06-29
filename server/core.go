package server

import (
	"encoding/json"
	"git.sda1.net/media-proxy-go/media"
	"git.sda1.net/media-proxy-go/security"
	"github.com/valyala/fasthttp"
	"strings"
)

func RequestHandler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/":
		status := Status{Status: "OK"}
		jsonData, err := json.Marshal(status)
		if err != nil {
			ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
			return
		}
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Write(jsonData)

	default:
		if strings.HasSuffix(string(ctx.Path()), ".webp") {
			queryArgs := ctx.QueryArgs()
			url := string(queryArgs.Peek("url"))

			if url == "" {
				ctx.Error("Bad request", fasthttp.StatusBadRequest)
				return
			}

			// ポートが指定されている、ホスト名がプライベートアドレスを示している場合はブロック
			if !security.IsSafeUrl(url) {
				ctx.Error("Access denied", fasthttp.StatusForbidden)
				return
			}

			convertedImage := media.ProcessImage(url, 320)

			if convertedImage == nil {
				ctx.Error("Internal server error", fasthttp.StatusInternalServerError)
				return
			}

			ctx.Response.Header.SetContentType("image/webp")
			ctx.Response.SetBody(convertedImage)
		} else {
			ctx.Error("Not found", fasthttp.StatusNotFound)
		}

	}
}
