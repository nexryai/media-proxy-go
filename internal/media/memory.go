package media

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
)

var ctx = context.Background()

// ProxyRequest型のリクエストをRedisに格納するときにつかうキーに変換する
func proxyRequestToKvKey(request *ProxyRequest) string {
	return fmt.Sprintf("mediaProxyApi/%s:%s:%d:%d:%v:%v", request.Url, request.TargetFormat, request.WidthLimit, request.HeightLimit, request.IsStatic, request.IsEmoji)
}

func connectToRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
	})

	return client
}

// キャッシュIDから実際のパスに変換する
func GetPathFromCacheId(id string) string {
	return os.Getenv("CACHE_DIR") + "/" + id
}

func StoreCachePath(request *ProxyRequest, cacheId string) error {
	kvStore := connectToRedis()
	defer kvStore.Close()

	kvStore.Set(ctx, proxyRequestToKvKey(request), GetPathFromCacheId(cacheId), 0)
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

	return result, nil
}
