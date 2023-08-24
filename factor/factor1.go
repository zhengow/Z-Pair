package factor

import "github.com/zhengow/Z-Quant/model"

func factor1(prices model.Prices) model.Base {
	ret := (prices[len(prices)-1].Open/prices[0].Open - 1) / float64(len(prices))
	return model.Base{
		Datetime: prices[len(prices)-1].Datetime,
		Val:      ret,
	}
}

func Factor1() (int, model.Handler, string) {
	return 1440, factor1, "factor1"
}
