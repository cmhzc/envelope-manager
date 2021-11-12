package config

import (
	"math/rand"
	"time"
)

func initSeed() {
	rand.Seed(time.Now().UnixNano())
}

func (RainConfig *rainConfig) GetRandomMoney() int64 {
	money := int64(0)
	if RainConfig.count_remain <= 0 || RainConfig.budget_remain <= 0 {
		RainConfig.count_remain--
		return 0
	} else if RainConfig.count_remain == 1 {
		if RainConfig.budget_remain > 0 {
			money = RainConfig.budget_remain
			RainConfig.budget_remain = 0
		} else {
			money = 0
		}
		RainConfig.count_remain = 0
		return money
	}

	// 最大可调度金额
	max := RainConfig.budget_remain - RainConfig.snatch_config.min_amount*RainConfig.count_remain
	if max <= 0 {
		// todo: return 0 or min_amount?
		RainConfig.budget_remain -= RainConfig.snatch_config.min_amount
		RainConfig.count_remain--
		return RainConfig.snatch_config.min_amount
	}
	// 每个红包平均调度金额
	avgMax := max / RainConfig.count_remain

	// 根据平均调度金额来生成每个红包金额
	randNum := rand.Float64() - 0.5
	avgMax += int64(randNum * float64(avgMax))

	money = RainConfig.snatch_config.min_amount + avgMax

	// border clip
	if money < RainConfig.snatch_config.min_amount {
		money = RainConfig.snatch_config.min_amount
	} else if money > RainConfig.snatch_config.max_amount {
		money = RainConfig.snatch_config.max_amount
	}

	RainConfig.budget_remain -= money
	RainConfig.count_remain--
	return money
}
