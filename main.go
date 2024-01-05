package main

import (
	"fmt"
	"git.sda1.net/media-proxy-go/core"
	"git.sda1.net/media-proxy-go/server"
	"github.com/davidbyttow/govips/v2/vips"
	"net/http"
	"os"
	"strconv"
)

func createServer(port int) {
	listenPort := strconv.Itoa(port)
	fmt.Println("listens on port " + listenPort)

	http.HandleFunc("/", server.RequestHandler)

	err := http.ListenAndServe(":"+listenPort, nil)
	core.ExitOnError(err, "Failed to start server")
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

	// vipsの初期化
	vips.Startup(nil)
	defer vips.Shutdown()

	port := getPort()
	createServer(port)
}
