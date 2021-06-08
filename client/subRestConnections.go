package client

import "github.com/adi1382/ftx-mirror-bot/go-ftx/rest/private/orders"

func (s *Sub) cancelOrderIDs(orderIDs []int64) {
	for i := range orderIDs {
		s.client.postCancelOrderByID(int(orderIDs[i]))
	}
}

func (s *Sub) placeOrders(orders []*orders.RequestForPlaceOrder) {
	for i := range orders {
		s.client.postPlaceOrder(orders[i])
	}
}

func (s *Sub) updateLeverage(isLeverageChangeRequired bool, newLeverage float64) {
	if isLeverageChangeRequired {
		s.client.postChangeLeverage(newLeverage)
	}
}
