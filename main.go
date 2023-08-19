package main

import (
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	"git.sda1.net/media-proxy-go/server"
	"github.com/davidbyttow/govips/v2/vips"
	"github.com/valyala/fasthttp"
	"os"
	"strconv"
)

func createServer(port int) {
	listenPort := strconv.Itoa(port)
	core.MsgInfo("listens on port " + listenPort)

	err := fasthttp.ListenAndServe(fmt.Sprintf(":%s", listenPort), server.RequestHandler)
	core.ExitOnError(err, "Failed to create server")

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

func main() {
	core.MsgInfo("Starting media-proxy-go ...")
	if core.IsDebugMode() {
		fmt.Println("\u001B[31m@@>>>>> Debug mode is enabled!!! NEVER use this in a production environment!! Debugging endpoints can leak sensitive information!!!!! <<<<<@@\u001B[0m")
	}

	mallocTrimThreshold := os.Getenv("MALLOC_TRIM_THRESHOLD_")
	if mallocTrimThreshold == "" {
		core.MsgWarn("MALLOC_TRIM_THRESHOLD_ environment variable is not set. It is highly recommended to set MALLOC_TRIM_THRESHOLD_ to 16 or a number close to it.")
	}

	mallocNmapThreshold := os.Getenv("MALLOC_MMAP_THRESHOLD_")
	if mallocNmapThreshold == "" {
		core.MsgWarn("MALLOC_MMAP_THRESHOLD_ environment variable is not set. It is highly recommended to set MALLOC_MMAP_THRESHOLD_ to 16 or a number close to it.")
	}

	// vipsの初期化
	vips.Startup(nil)
	defer vips.Shutdown()

	port := getPort()
	createServer(port)
}
