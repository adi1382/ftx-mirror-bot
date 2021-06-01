package wallet

import (
	"net/http"

	jsonIter "github.com/json-iterator/go"
)

var json = jsonIter.ConfigCompatibleWithStandardLibrary

type RequestForCoins struct {
}

type ResponseForCoins []Coin

type Coin struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	CanDeposit  bool `json:"canDeposit"`
	CanWithdraw bool `json:"canWithdraw"`
	HasTag      bool `json:"hasTag"`
}

func (req *RequestForCoins) Path() string {
	return "/wallet/coins"
}

func (req *RequestForCoins) Method() string {
	return http.MethodGet
}

func (req *RequestForCoins) Query() string {
	return ""
}

func (req *RequestForCoins) Payload() []byte {
	return nil
}
