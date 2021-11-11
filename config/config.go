package config

import (
	"envelope_manager/redis"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type snatchConfig struct {
	probability float64
	min_amount  int64
	max_amount  int64
}

type rainConfig struct {
	budget_remain int64
	count_remain  int64
	snatch_config *snatchConfig
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
		snatch_config: &snatchConfig{
			probability: viper.GetFloat64("snatch.probability"),
			min_amount:  viper.GetInt64("snatch.min_amount"),
			max_amount:  viper.GetInt64("snatch.max_amount"),
		},
	}
	lastBudget = RainConfig.budget_remain
	num_goroutines = viper.GetInt("conc.threads")
	// start watching config file
	watchConfig()

	// produce envelopes
	produce()
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

func watchConfig() {
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("[viper] Config file content changed:", e.Name)
		// on config file content change, reconfig
		reConfig()
	})
	viper.WatchConfig()
}

func reConfig() {
	// retrieve envelopes
	collect_amount, collect_count := consume()

	// calculate remain amount and count
	RainConfig.count_remain = collect_count
	// amount_new - amount_old + amount_collect + amount_left = new amount
	RainConfig.budget_remain += viper.GetInt64("rain.budget") - lastBudget + collect_amount
	lastBudget = RainConfig.budget_remain

	// regenerate envelopes
	produce()
}

func produce() {
	ch := make(chan int64, 1000)
	go func() {
		for RainConfig.count_remain > 0 {
			ch <- GetRandomMoney()
		}
		close(ch)
	}()
	for i := 0; i < num_goroutines; i++ {
		go func() {
			for amount := range ch {
				if n, err := redis.Rdb.LPush("envelope_list", amount).Result(); err != nil {
					log.Fatal("insert failed")
				} else {
					log.Printf("insert success %v,the value is %v", n, amount)
				}
			}
		}()
	}
}

func consume() (int64, int64) {
	ch := make(chan int64, 1000)
	wg := sync.WaitGroup{}
	for i := 0; i < num_goroutines; i++ {
		wg.Add(1)
		go func() {
			for money, err := redis.Rdb.LPop("envelope_list").Result(); err == nil; {
				amount, _ := strconv.ParseInt(money, 10, 64)
				ch <- amount
				money, err = redis.Rdb.LPop("envelope_list").Result()
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	var collect_amount int64 = 0
	var collect_count int64 = 0
	for amount := range ch {
		collect_amount += amount
		collect_count++
	}
	return collect_amount, collect_count
}
