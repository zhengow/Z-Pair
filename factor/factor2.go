package factor

import (
	"github.com/zhengow/Z-Quant/model"
	"gonum.org/v1/gonum/stat"
)

func factor2(prices model.Prices) model.Base {
	ret := stat.Mean(prices.Open(), nil)
	return model.Base{
		Datetime: prices[len(prices)-1].Datetime,
		Val:      ret,
	}
}

func Factor2() (int, model.Handler, string) {
	return 1440, factor2, "factor2"
}
