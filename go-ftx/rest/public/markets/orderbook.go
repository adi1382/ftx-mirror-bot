package markets

import (
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// RequestForOrderBook query
// ?depth={depth}
type RequestForOrderBook struct {
	ProductCode string `url:"-"`
	Depth       int    `url:"depth,omitempty"`
}

type ResponseForOrderBook OrderBook

type OrderBook struct {
	Asks [][]float64 `json:"asks"`
	Bids [][]float64 `json:"bids"`
}

func (req *RequestForOrderBook) Path() string {
	return fmt.Sprintf("/markets/%s/orderbook", req.ProductCode)
}

func (req *RequestForOrderBook) Method() string {
	return http.MethodGet
}

func (req *RequestForOrderBook) Query() string {
	values, _ := query.Values(req)
	return values.Encode()
}

func (req *RequestForOrderBook) Payload() []byte {
	return nil
}
