package processor

import (
	"git.sda1.net/media-proxy-go/internal/queue"
	"github.com/hibiken/asynq"
	"log"
	"os"
)

func ProxyQueueProcessor() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: os.Getenv("REDIS_ADDR")},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(queue.TypeImageProxy, queue.HandleImageProxyTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
