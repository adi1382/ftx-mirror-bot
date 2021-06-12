package client

import (
	"math"

	"github.com/adi1382/ftx-mirror-bot/fpe"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/private/orders"
)

func (s *Sub) cancelOrderByIDs(orderIDs []int64) {
	for i := range orderIDs {
		s.client.postCancelOrderByID(int(orderIDs[i]))
	}
}

func (s *Sub) cancelOrderByClOrdID(ClOrdID []string) {
	for i := range ClOrdID {
		s.client.postCancelOrderByClOrdID(ClOrdID[i])
	}
}

func (s *Sub) placeOrders(orderRequests ...*orders.RequestForPlaceOrder) []*orders.ResponseForPlaceOrder {
	orderResponses := make([]*orders.ResponseForPlaceOrder, 0, 5)
	for i := range orderRequests {
		orderResponses = append(orderResponses, s.client.postPlaceOrder(orderRequests[i]))
	}
	return orderResponses
}

func (s *Sub) updateLeverage(isLeverageChangeRequired bool, newLeverage float64) {
	if isLeverageChangeRequired {
		s.client.postChangeLeverage(newLeverage)
	}
}

//generateMarketOrder generates requests for market order, it does not place any order.
func (s *Sub) generateMarketOrder(market, side string, size float64, reduceOnly bool) *orders.RequestForPlaceOrder {
	var orderSide string

	if side == "auto" {
		if size > 0 {
			orderSide = "buy"
		} else if size < 0 {
			orderSide = "sell"
		} else {
			panic("size zero order creation")
		}
	} else {
		orderSide = side
	}

	return &orders.RequestForPlaceOrder{
		ClientID:   fpe.GenerateRandomClOrdID(),
		Type:       "market",
		Market:     market,
		Side:       orderSide,
		ReduceOnly: reduceOnly,
		Size:       math.Abs(size),
	}
}
