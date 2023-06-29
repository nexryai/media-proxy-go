package server

import (
	"encoding/json"
	"git.sda1.net/media-proxy-go/media"
	"github.com/valyala/fasthttp"
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
	case "/proxy/avatar.webp":
		queryArgs := ctx.QueryArgs()
		url := string(queryArgs.Peek("url"))

		if url == "" {
			ctx.Error("Bad request", fasthttp.StatusBadRequest)
		}

		convertedImage := media.ProcessImage(url, 320)

		if convertedImage == nil {
			ctx.Error("Internal server error", fasthttp.StatusInternalServerError)
		}

		ctx.Response.Header.SetContentType("image/webp")
		ctx.Response.SetBody(convertedImage)

	default:
		ctx.Error("Not found", fasthttp.StatusNotFound)
	}
}
