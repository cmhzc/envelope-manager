package redis

import (
	"os"

	"github.com/go-redis/redis"
)

var Rdb *redis.Client

func InitRedis() error {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
		PoolSize: 100,
	})
	_, err := Rdb.Ping().Result()
	// DEL `node_id` key
	Rdb.Del("node_id")
	return err
}
