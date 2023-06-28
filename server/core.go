package server

import (
	"encoding/json"
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
	default:
		ctx.Error("Not found", fasthttp.StatusNotFound)
	}
}
