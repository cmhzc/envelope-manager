package config

import (
	"fmt"
	"testing"
)

func TestGetRandomMoney(t *testing.T) {

	RainConfig = &rainConfig{
		budget_remain: 100000000000,
		count_remain:  100000000,
		snatch_config: &snatchConfig{
			probability: 0.6,
			min_amount:  100,
			max_amount:  10000,
		},
	}

	moneyArr := make([]int64, 0)
	for RainConfig.count_remain > 0 {
		x := GetRandomMoney()
		moneyArr = append(moneyArr, x)
	}

	fmt.Println("分配的红包金额:")
	fmt.Println(moneyArr)

	fmt.Println("分配的红包金额总和:")
	var sum int64 = 0
	for _, num := range moneyArr {
		sum += num
	}
	fmt.Println(sum)
}
