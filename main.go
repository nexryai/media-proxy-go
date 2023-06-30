package main

import (
	"fmt"
	"os"
	"strconv"

	"git.sda1.net/media-proxy-go/core"
	"git.sda1.net/media-proxy-go/server"
	"github.com/pkg/profile"
	"github.com/valyala/fasthttp"
)

func createServer(port int) {
	listenPort := strconv.Itoa(port)
	core.MsgInfo("listens on port " + listenPort)
	err := fasthttp.ListenAndServe(fmt.Sprintf(":%s", listenPort), server.RequestHandler)
	core.ExitOnError(err, "Failed to create server")
}

func getPort() (int) {
	port := os.Getenv("PORT")
	const defaultPort = 8080

	if port == "" {
		return defaultPort
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		messeage := "PORT" + port + "is not valid"
		core.MsgErr(messeage)
		return defaultPort
	}
	return portNum

}

func main() {
	if core.IsDebugMode() {
		defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()
	}
	core.MsgInfo("Starting media-proxy-go ...")
	port := getPort()
	createServer(port)
}
