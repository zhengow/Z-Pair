package core

import (
	"github.com/zhengow/Z-Quant/model"
	"github.com/zhengow/Z-Quant/utils"
	"gonum.org/v1/gonum/stat"
)

func Rolling(prices []model.Price, window int, f model.Handler) model.Factors {
	factors := make(model.Factors, len(prices)-window+1)
	for i := window; i <= len(prices); i++ {
		factors[i-window] = f(prices[i-window : i])
	}
	return factors
}

func Pos(factors model.Factors) model.Positions {
	window := 60
	positions := make(model.Positions, 0)
	for i := window; i < len(factors); i++ {
		currentFactors := factors[i-window : i]
		mean, std := stat.MeanStdDev(currentFactors.Val(), nil)
		val := currentFactors.LastVal()
		pos := utils.Divide(val-mean, std)
		positions = append(positions, model.Base{
			Datetime: currentFactors.LastTime(),
			Val:      pos,
		})
	}
	return positions
}

func Rev(factors model.Factors, prices model.Prices) model.BaseSeries {
	positions := Pos(factors)
	priceReturns := prices.Return()
	rate := 0.0002
	posReturns := positions.Returns(priceReturns, rate)
	return posReturns
}

func CumSum(series model.BaseSeries) model.BaseSeries {
	for i := 1; i < len(series); i++ {
		series[i].Val += series[i-1].Val
	}
	return series
}
