package manager

import (
	"fmt"
	"time"

	json2 "encoding/json"
	"github.com/go-redis/redis/v8"
)

var cacheClient *redis.Client

func initCache() *redis.Client {
	cacheAddr := fmt.Sprintf("%s:%s", conf.ServerHost, conf.CachePort)

	rdb := redis.NewClient(&redis.Options{
		Addr:     cacheAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb
}

func getCache(key string) (string, error) {
	ctx := initCtx()
	val, err := cacheClient.Get(ctx, key).Result()

	return val, err
}

func setCache(key string, data interface{}) {
	json, err := json2.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	ctx := initCtx()
	err = cacheClient.Set(ctx, key, json, 10*time.Minute).Err()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func delCache(key string) {
	ctx := initCtx()
	err := cacheClient.Del(ctx, key).Err()
	if err != nil {
		fmt.Println(err)
	}
}
