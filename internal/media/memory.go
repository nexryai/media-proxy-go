package media

import (
	"context"
	"fmt"
	"git.sda1.net/media-proxy-go/internal/logger"
	"github.com/redis/go-redis/v9"
	"os"
	"strings"
	"time"
)

var ctx = context.Background()

// ProxyRequest型のリクエストをRedisに格納するときにつかうキーに変換する
func proxyRequestToKvKey(request *ProxyRequest) string {
	return fmt.Sprintf("mediaProxyApi/%s:%s:%d:%d:%v:%v", request.Url, request.TargetFormat, request.WidthLimit, request.HeightLimit, request.IsStatic, request.IsEmoji)
}

func generateLifecycleKey(request *ProxyRequest) string {
	return fmt.Sprintf("lifecycleManagerStore/%s:%s:%d:%d:%v:%v", request.Url, request.TargetFormat, request.WidthLimit, request.HeightLimit, request.IsStatic, request.IsEmoji)
}

func timeToString(t time.Time) string {
	return t.Format(time.RFC3339)
}

func stringToTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		// 2000年1月1日
		return time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	return t
}

func connectToRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
	})

	return client
}

// キャッシュIDから実際のパスに変換する
func GetPathFromCacheId(id string) string {
	if id == "FAILED" {
		return "FAILED"
	}

	return os.Getenv("CACHE_DIR") + "/" + id
}

func StoreCachePath(request *ProxyRequest, cacheId string) error {
	kvStore := connectToRedis()
	defer kvStore.Close()

	kvStore.Set(ctx, proxyRequestToKvKey(request), GetPathFromCacheId(cacheId), 0)
	kvStore.Set(ctx, generateLifecycleKey(request), timeToString(time.Now()), 0)
	return nil
}

func CacheExists(request *ProxyRequest) bool {
	kvStore := connectToRedis()
	defer kvStore.Close()

	result, err := kvStore.Exists(ctx, proxyRequestToKvKey(request)).Result()
	if err != nil {
		return false
	}

	return result == 1
}

func GetCachePath(request *ProxyRequest) (string, error) {
	kvStore := connectToRedis()
	defer kvStore.Close()

	result, err := kvStore.Get(ctx, proxyRequestToKvKey(request)).Result()
	if err != nil {
		return "", err
	}

	// Last usedを更新
	if result != "FAILED" {
		kvStore.Set(ctx, generateLifecycleKey(request), timeToString(time.Now()), 0)
	}
	return result, nil
}

// 使われていないキャッシュを削除する
func CleanCache() {
	log := logger.GetLogger("LifecycleManager")
	log.Info("Cleaning cache ...")

	kvStore := connectToRedis()
	defer kvStore.Close()

	// 3日以上使われていないキャッシュを削除する
	result, err := kvStore.Keys(ctx, "lifecycleManagerStore/*").Result()
	if err != nil {
		return
	}

	for _, key := range result {
		lastUsed, err := kvStore.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		lastUsedTime := stringToTime(lastUsed)
		cacheStoreKey := fmt.Sprintf("mediaProxyApi/%s", strings.TrimPrefix(key, "lifecycleManagerStore/"))
		cachePath := kvStore.Get(ctx, cacheStoreKey).Val()

		if time.Since(lastUsedTime).Hours() > 72 {
			// 3日以上使われていないキャッシュを削除する
			log.Info(fmt.Sprintf("Removing cache: %s", cachePath))

			err = os.Remove(cachePath)
			if err != nil {
				log.Warn(fmt.Sprintf("Failed to remove cache: %v", err))
			}

			kvStore.Del(ctx, cacheStoreKey)
			kvStore.Del(ctx, key)
		} else if time.Since(lastUsedTime).Minutes() > 20 && cachePath == "FAILED" {
			// 20分以上使われていないFAILEDキャッシュを削除する
			log.Info(fmt.Sprintf("Removing FAILED cache: %s", cacheStoreKey))

			kvStore.Del(ctx, cacheStoreKey)
			kvStore.Del(ctx, key)
		}
	}
}
