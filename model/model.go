package model

import "time"

type Price struct {
	Datetime time.Time
	Open     float64
	Close    float64
	High     float64
	Low      float64
	Volume   float64
}

type Factor struct {
	Datetime time.Time
	Val      float64
}

type Handler func([]Price) Factor
