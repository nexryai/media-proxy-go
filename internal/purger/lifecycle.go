package purger

import (
	"git.sda1.net/media-proxy-go/internal/media"
	"time"
)

func StartLifecycleManager() {
	for {
		media.CleanCache()
		time.Sleep(15 * time.Minute)
	}
}
