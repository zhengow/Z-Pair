package main

import (
	"fmt"
	"github.com/zhengow/Z-Quant/core"
	"github.com/zhengow/Z-Quant/factor"
	"github.com/zhengow/Z-Quant/utils"
	"time"
)

func main() {
	start := time.Now()
	prices := utils.LoadData("DOGEUSDT", []string{"open", "high"})
	window, handler, name := factor.Factor1()
	fmt.Println("start rolling", time.Since(start))
	factors := core.Rolling(prices, window, handler)
	fmt.Println(time.Since(start))
	rev := core.Rev(factors, prices)
	cumRev := core.CumSum(rev)
	utils.DrawSimple(cumRev, name)
}
