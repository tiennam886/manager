package manager

import (
	"context"
	"fmt"
	"time"

	json2 "encoding/json"
	"github.com/go-redis/redis/v8"
)

var (
	cacheAddr   = fmt.Sprintf("%s:%s", conf.ServerHost, conf.CachePort)
	cacheClient *redis.Client
)

func initCache() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cacheAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb
}

func getCache(key string) (string, error) {
	ctx := context.Background()
	val, err := cacheClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, err
}

func setCache(key string, data interface{}) {
	json, err := json2.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	ctx := context.Background()
	err = cacheClient.Set(ctx, key, json, 10*time.Minute).Err()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func delCache(key string) {
	ctx := context.Background()
	err := cacheClient.Del(ctx, key).Err()
	if err != nil {
		fmt.Println(err)
	}
}
