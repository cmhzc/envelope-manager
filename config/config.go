package config

import (
	"envelope_manager/redis"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type snatchConfig struct {
	Probability float64
	min_amount  int64
	max_amount  int64
}

type rainConfig struct {
	budget_remain int64
	count_remain  int64
	Snatch_config *snatchConfig
}

var RainConfig *rainConfig
var lastBudget int64
var num_goroutines int

func InitRainConfig(name string) error {
	if err := initConfig(name); err != nil {
		return err
	}
	initSeed()
	RainConfig = &rainConfig{
		budget_remain: viper.GetInt64("rain.budget"),
		count_remain:  viper.GetInt64("rain.count"),
		Snatch_config: &snatchConfig{
			Probability: viper.GetFloat64("snatch.probability"),
			min_amount:  viper.GetInt64("snatch.min_amount"),
			max_amount:  viper.GetInt64("snatch.max_amount"),
		},
	}
	lastBudget = RainConfig.budget_remain
	num_goroutines = viper.GetInt("conc.threads")

	// produce envelopes
	start := time.Now()
	produce()
	elapsed := time.Since(start)
	log.Printf("[manager] produce took %s", elapsed)
	return nil
}

func initConfig(name string) error {
	// read config file
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName(name)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	return nil
}

func ReConfig(newBudget int64) {
	// retrieve envelopes
	start := time.Now()
	collect_amount, collect_count := consume()
	elapsed := time.Since(start)
	log.Printf("[manager] cosume took %s", elapsed)

	// calculate remain amount and count
	RainConfig.count_remain = collect_count + RainConfig.count_remain
	// amount_new - amount_old + amount_collect + amount_left = new amount
	RainConfig.budget_remain = newBudget - lastBudget + collect_amount + RainConfig.budget_remain
	lastBudget = RainConfig.budget_remain
	log.Printf("[manager] new config genreated: budget_remain %v, count_remain %v", RainConfig.budget_remain, RainConfig.count_remain)

	// regenerate envelopes
	start = time.Now()
	produce()
	elapsed = time.Since(start)
	log.Printf("[manager] produce took %s", elapsed)
}

func produce() {
	ch := make(chan []interface{}, 1000)
	wg := sync.WaitGroup{}
	go func() {
		for RainConfig.count_remain > 0 && RainConfig.budget_remain > 0 {
			s := make([]interface{}, 10000)
			count := 0
			for ; RainConfig.count_remain > 0 && RainConfig.budget_remain > 0 && count < 10000; count++ {
				s[count] = RainConfig.GetRandomMoney()
			}
			if count < 10000 {
				ch <- s[:count]
			} else {
				ch <- s
			}
		}
		close(ch)
	}()
	for i := 0; i < num_goroutines; i++ {
		wg.Add(1)
		go func() {
			for amounts := range ch {
				if _, err := redis.Rdb.LPush("envelope_list", amounts...).Result(); err != nil {
					log.Fatal("insert failed")
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func consume() (collect_amount int64, collect_count int64) {
	ch := make(chan []string, 1000)
	wg := sync.WaitGroup{}
	for i := 0; i < num_goroutines; i++ {
		wg.Add(1)
		go func() {
			for {
				pipe := redis.Rdb.TxPipeline()
				amounts := pipe.LRange("envelope_list", 0, 9999)
				pipe.LTrim("envelope_list", 10000, -1).Result()
				pipe.Exec()
				if len(amounts.Val()) == 0 {
					break
				}
				ch <- amounts.Val()
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	for amounts := range ch {
		for _, amount := range amounts {
			money, _ := strconv.ParseInt(amount, 10, 64)
			collect_amount += money
			collect_count++
		}
	}
	log.Printf("[manager] collected %v envelopes with value of %v", collect_count, collect_amount)
	return collect_amount, collect_count
}
