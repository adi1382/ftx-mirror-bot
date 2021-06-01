package client

import (
	"github.com/adi1382/ftx-mirror-bot/optional"
)

type order struct {
	Id            optional.Int64   `json:"id"`
	ClientId      optional.String  `json:"clientId"`
	Market        optional.String  `json:"market"`
	Type          optional.String  `json:"type"`
	Side          optional.String  `json:"side"`
	Size          optional.Float64 `json:"size"`
	Price         optional.Float64 `json:"price"`
	ReduceOnly    optional.Bool    `json:"reduceOnly"`
	Ioc           optional.Bool    `json:"ioc"`
	PostOnly      optional.Bool    `json:"postOnly"`
	Status        optional.String  `json:"status"`
	FilledSize    optional.Float64 `json:"filledSize"`
	RemainingSize optional.Float64 `json:"remainingSize"`
	AvgFillPrice  optional.Float64 `json:"avgFillPrice"`
}
