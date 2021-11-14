package main

import (
	"envelope_manager/config"
	"envelope_manager/dao"
	"envelope_manager/redis"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	// Flush MySQL DB
	dao.FlushDB()

	// Redis connection init
	if err := redis.InitRedis(); err != nil {
		log.Fatal("failed to connect to Redis")
	}
	defer redis.Rdb.Close()

	// config init, produce envelopes
	config.InitRainConfig(os.Getenv("CONFIG_NAME"))
	if err := redis.WriteProb(config.RainConfig.Snatch_config.Probability); err != nil {
		panic("failed writing prob")
	}
	log.Printf("[manager] wrote to redis: prob %v", config.RainConfig.Snatch_config.Probability)

	// config secret
	secret := gin.Accounts{
		os.Getenv("AUTH_USERNAME"): os.Getenv("AUTH_PASSWORD"),
	}
	// receive requests for reconfig
	r := gin.Default()
	r.POST("/reconfig", gin.BasicAuth(secret), func(c *gin.Context) {
		newBudget, err := strconv.ParseInt(c.PostForm("newBudget"), 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "failed due to wrong input format",
			})
			return
		}
		config.ReConfig(newBudget)
		c.JSON(http.StatusOK, gin.H{
			"status": "succeeded",
		})
	})
	r.Run()
}
