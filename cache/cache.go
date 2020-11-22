package cache

import "github.com/go-redis/redis"

var mainRedis *redis.Client

func InitCache(addr string,db int) {
	mainRedis = redis.NewClient(&redis.Options{
		Addr: addr,
		DB:db,
	})
}
