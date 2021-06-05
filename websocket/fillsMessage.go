package websocket

import (
	"time"
)

type FillsWSMessage struct {
	Channel        string    `json:"channel"`
	TypeOfResponse string    `json:"type"`
	Data           FillsData `json:"data"`
}

type FillsData struct {
	Fee       float64   `json:"fee"`
	FeeRate   float64   `json:"feeRate"`
	Future    string    `json:"future"`
	Id        int       `json:"id"`
	Liquidity string    `json:"liquidity"`
	Market    string    `json:"market"`
	OrderId   int       `json:"orderId"`
	TradeId   int       `json:"tradeId"`
	Price     float64   `json:"price"`
	Side      string    `json:"side"`
	Size      float64   `json:"size"`
	Time      time.Time `json:"time"`
	Type      string    `json:"type"`
}
