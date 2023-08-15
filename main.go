package main

import (
	"fmt"
	"github.com/zhengow/Z-Quant/factor"
	"github.com/zhengow/Z-Quant/model"
	"github.com/zhengow/Z-Quant/utils"
	"time"
)

func rolling(prices []model.Price, window int, f model.Handler) {
	var factors []model.Factor
	for i := window; i <= len(prices); i++ {
		factors = append(factors, f(prices[i-window:i]))
	}
	for i := 0; i < 5; i++ {
		fmt.Println(factors[i].Datetime, factors[i].Val)
	}
}

func main() {
	start := time.Now()
	prices := utils.LoadData("DOGEUSDT", []string{"open", "high"})
	window, handler := factor.Factor1()
	fmt.Println("start rolling", time.Since(start))
	rolling(prices, window, handler)
	fmt.Println(time.Since(start))
}
