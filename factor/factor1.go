package factor

import "github.com/zhengow/Z-Quant/model"

func factor1(prices []model.Price) model.Factor {
	var ret float64
	for i := 0; i < len(prices)-1; i++ {
		ret += prices[i+1].Open/prices[i].Open - 1
	}
	ret = ret / float64(len(prices)-1)
	return model.Factor{
		Datetime: prices[len(prices)-1].Datetime,
		Val:      ret,
	}
}

func Factor1() (int, model.Handler) {
	return 14400, factor1
}
