package main

import (
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	"git.sda1.net/media-proxy-go/server"
	"github.com/pkg/profile"
	"github.com/valyala/fasthttp"
)

func createServer(listenPort string) {
	core.MsgInfo("listens on port " + listenPort)
	err := fasthttp.ListenAndServe(fmt.Sprintf(":%s", listenPort), server.RequestHandler)
	core.ExitOnError(err, "Failed to create server")
}

func main() {
	defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()
	core.MsgInfo("Starting media-proxy-go ...")
	createServer("8080")
}
