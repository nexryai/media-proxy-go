package main

import (
	"fmt"
	"git.sda1.net/media-proxy-go/internal/core"
	"git.sda1.net/media-proxy-go/internal/logger"
	"git.sda1.net/media-proxy-go/internal/server"
	"github.com/davidbyttow/govips/v2/vips"
	"github.com/nexryai/visualog"
	"net/http"
	"os"
	"strconv"
)

func createServer(port int, log visualog.Logger) {
	log.ProgressInfo("Starting server ...")

	listenPort := strconv.Itoa(port)

	http.HandleFunc("/", server.RequestHandler)

	log.ProgressOk()
	fmt.Print("\n")
	log.Info("Listening on port " + listenPort)

	err := http.ListenAndServe(":"+listenPort, nil)
	if err != nil {
		log.FatalWithDetail("Failed to start server", err)
		os.Exit(1)
	}
}

func getPort(log visualog.Logger) int {
	port := os.Getenv("PORT")
	const defaultPort = 8080

	if port == "" {
		return defaultPort
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		log.Error("PORT " + port + " is not valid")
		return defaultPort
	}
	return portNum

}

func main() {
	log := logger.GetLogger("Boot")

	log.Info("Starting media-proxy-go ...")
	if core.IsDebugMode() {
		log.Warn("@@>>>>> Debug mode is enabled!!! NEVER use this in a production environment!! Debugging endpoints can leak sensitive information!!!!! <<<<<@@")
	}

	// vipsの初期化
	log.ProgressInfo("Initializing vips ...")
	vips.Startup(&vips.Config{
		ConcurrencyLevel: 1,
		MaxCacheMem:      8 * 1024 * 1024,
		MaxCacheSize:     32,
		MaxCacheFiles:    32,
	})
	defer vips.Shutdown()
	log.ProgressOk()

	port := getPort(*log)
	createServer(port, *log)
}
