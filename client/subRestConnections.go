package client

import "github.com/adi1382/ftx-mirror-bot/go-ftx/rest/private/orders"

func (s *Sub) cancelOrderIDs(orderIDs []int64) {
	for i := range orderIDs {
		s.client.postCancelOrderByID(int(orderIDs[i]))
	}
}

func (s *Sub) placeOrders(orderRequests []*orders.RequestForPlaceOrder) []*orders.ResponseForPlaceOrder {
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
