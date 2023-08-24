package model

import (
	"math"
	"time"
)

type Base struct {
	Datetime time.Time
	Val      float64
}

type BaseSeries []Base

type Price struct {
	Datetime time.Time
	Open     float64
	Close    float64
	High     float64
	Low      float64
	Volume   float64
}

type Prices []Price

type Return Base

func (prices Prices) Return() []Return {
	val := make([]Return, len(prices)-1, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		val[i-1] = Return{
			Datetime: prices[i-1].Datetime,
			Val:      prices[i].Open/prices[i-1].Open - 1,
		}
	}
	return val
}

func (prices Prices) Open() []float64 {
	val := make([]float64, len(prices), len(prices))
	for i := 0; i < len(prices); i++ {
		val[i] = prices[i].Open
	}
	return val
}

type Factor Base

type Position Base

type Positions BaseSeries

func (positions Positions) Returns(priceReturns []Return, rate float64) BaseSeries {
	positionStart := positions[0].Datetime
	positionEnd := positions[len(positions)-1].Datetime
	priceReturnStart := priceReturns[0].Datetime
	priceReturnEnd := priceReturns[len(priceReturns)-1].Datetime

	positionStartIdx := 0
	priceReturnStartIdx := 0
	if positionStart.Before(priceReturnStart) {
		// position起始日期更早，position后移
		for positions[positionStartIdx].Datetime.Before(priceReturnStart) {
			positionStartIdx++
		}
	} else if positionStart.After(priceReturnStart) {
		// price起始日期更早，position后移
		for positionStart.After(priceReturns[priceReturnStartIdx].Datetime) {
			priceReturnStartIdx++
		}
	}

	positionEndIdx := len(positions) - 1
	priceReturnEndIdx := len(priceReturns) - 1
	if positionEnd.Before(priceReturnEnd) {
		// position结束的更早，price前移
		for positionEnd.Before(priceReturns[priceReturnEndIdx].Datetime) {
			priceReturnEndIdx--
		}
	} else if positionEnd.After(priceReturnEnd) {
		// price结束的更早，position前移
		for positions[positionEndIdx].Datetime.After(priceReturnEnd) {
			positionEndIdx--
		}
	}

	// 对齐index了
	positionReturns := make(BaseSeries, positionEndIdx-positionStartIdx+1)
	for i := positionStartIdx; i <= positionEndIdx; i++ {
		fee := float64(0)
		if i > 0 {
			fee = math.Abs(positions[i].Val-positions[i-1].Val) * rate
		}
		ret := float64(positions[i].Val)*priceReturns[priceReturnStartIdx-positionStartIdx+i].Val - fee
		positionReturns[i] = Base{
			Datetime: positions[i].Datetime,
			Val:      ret,
		}
	}

	return positionReturns
}

type Factors BaseSeries

func (factors Factors) Val() []float64 {
	val := make([]float64, len(factors), len(factors))
	for idx, factor := range factors {
		val[idx] = factor.Val
	}
	return val
}

func (factors Factors) LastTime() time.Time {
	return factors[len(factors)-1].Datetime
}

func (factors Factors) LastVal() float64 {
	return factors[len(factors)-1].Val
}

type Handler func(Prices) Base
