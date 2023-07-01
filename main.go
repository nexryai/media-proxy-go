package main

import (
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	"git.sda1.net/media-proxy-go/server"
	"github.com/valyala/fasthttp"
	"net/http"
	"os"
	"strconv"
)

func createServer(port int) {
	listenPort := strconv.Itoa(port)
	core.MsgInfo("listens on port " + listenPort)
	if lowMemoryMode() {
		core.MsgInfo("Use low memory mode!")
		http.HandleFunc("/", server.RequestHandlerLowMemoryMode)
		err := http.ListenAndServe(fmt.Sprintf(":%s", listenPort), nil)
		core.ExitOnError(err, "Failed to create server")
	} else {
		err := fasthttp.ListenAndServe(fmt.Sprintf(":%s", listenPort), server.RequestHandler)
		core.ExitOnError(err, "Failed to create server")
	}
}

func getPort() int {
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

func lowMemoryMode() bool {
	if os.Getenv("LOW_MEMORY_MODE") == "1" {
		return true
	} else {
		return false
	}
}

func main() {
	core.MsgInfo("Starting media-proxy-go ...")
	port := getPort()
	createServer(port)
}
