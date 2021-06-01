package futures

import (
	"net/http"
	"time"
)

type RequestForFutures struct {
}

type ResponseForFutures []FutureForList

type FutureForList struct {
	Type       string `json:"type"`
	Name       string `json:"name"`
	Underlying string `json:"underlying"`

	Index        float64 `json:"index"`
	Mark         float64 `json:"mark"`
	Last         float64 `json:"last"`
	Ask          float64 `json:"ask"`
	Bid          float64 `json:"bid"`
	Change1H     float64 `json:"change1h"`
	Change24H    float64 `json:"change24h"`
	ChangeBod    float64 `json:"changeBod"`
	Volume       float64 `json:"volume"`
	VolumeUsd24H float64 `json:"volumeUsd24h"`

	PriceIncrement float64 `json:"priceIncrement"`
	SizeIncrement  float64 `json:"sizeIncrement"`

	UpperBound float64 `json:"upperBound"`
	LowerBound float64 `json:"lowerBound"`

	Description string    `json:"description"`
	Expiry      time.Time `json:"expiry"`
	Enabled     bool      `json:"enabled"`
	Expired     bool      `json:"expired"`
	Perpetual   bool      `json:"perpetual"`
	PostOnly    bool      `json:"postOnly"`
}

func (req *RequestForFutures) Path() string {
	return "/futures"
}

func (req *RequestForFutures) Method() string {
	return http.MethodGet
}

func (req *RequestForFutures) Query() string {
	return ""
}

func (req *RequestForFutures) Payload() []byte {
	return nil
}

func (futures ResponseForFutures) Products() []string {
	list := make([]string, len(futures))
	for i := range futures {
		list[i] = futures[i].Name
	}
	return list
}

// Len Sort by alphabetical order (by Name)
func (futures ResponseForFutures) Len() int           { return len(futures) }
func (futures ResponseForFutures) Swap(i, j int)      { futures[i], futures[j] = futures[j], futures[i] }
func (futures ResponseForFutures) Less(i, j int) bool { return futures[i].Name < futures[j].Name }
