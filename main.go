package main

import (
	"envelope_manager/config"
	"envelope_manager/redis"
	"log"
	"os"
)

func main() {
	// Redis connection init
	if err := redis.InitRedis(); err != nil {
		log.Fatal("failed to connect to Redis")
	}
	defer redis.Rdb.Close()

	// config init, produce envelopes
	config.InitRainConfig(os.Getenv("CONFIG_NAME"))

	// block the main goroutine
	select {}
}
