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
	// FlushDB
	Rdb.FlushDB()
	return err
}

func WriteProb(prob float64) error {
	if _, err := Rdb.IncrByFloat("prob", prob).Result(); err != nil {
		return err
	}
	return nil
}
