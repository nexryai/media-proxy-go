package queue

import (
	"context"
	"encoding/json"
	"git.sda1.net/media-proxy-go/internal/logger"
	"git.sda1.net/media-proxy-go/internal/media"
	"github.com/hibiken/asynq"
)

const (
	TypeImageProxy = "image:proxy"
)

type ImageProxyPayload struct {
	Request *media.ProxyRequest
}

func NewProxyRequestTask(req *media.ProxyRequest) (*asynq.Task, error) {
	payload, err := json.Marshal(ImageProxyPayload{Request: req})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeImageProxy, payload), nil
}

func HandleImageProxyTask(ctx context.Context, t *asynq.Task) error {
	log := logger.GetLogger("QueueProcessor")

	var payload ImageProxyPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		log.ErrorWithDetail("Failed to unmarshal payload", err)
		return err
	}

	log.Info("Processing image proxy task: " + payload.Request.Url)
	_, _, err := media.ProxyImage(payload.Request)
	return err
}
