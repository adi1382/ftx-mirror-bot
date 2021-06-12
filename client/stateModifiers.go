package client

import (
	"math"

	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/private/orders"
	"github.com/adi1382/ftx-mirror-bot/tools"
)

// These functions mutate the state of Sub type
// removeOrdersInLocalStateFromOrderIDs removes orders from s.client.activeOrders which got canceled by calibrator
func (s *Sub) removeOrdersInLocalStateFromClOrdIDs(ClOrdIDs []string) {
	for i := range ClOrdIDs {
		for j := range s.client.activeOrders {
			if ClOrdIDs[i] == s.client.activeOrders[j].ClientId.Value {
				s.client.activeOrders[j] = s.client.activeOrders[len(s.client.activeOrders)-1]
				s.client.activeOrders = s.client.activeOrders[:len(s.client.activeOrders)-1]
				break
			}
		}
	}
}

// These functions mutate the state of Sub type
// removeOrdersInLocalStateFromOrderIDs removes orders from s.client.activeOrders which got canceled by calibrator
func (s *Sub) removeOrdersInLocalStateFromOrderIDs(orderIDs []int64) {
	for i := range orderIDs {
		for j := range s.client.activeOrders {
			if orderIDs[i] == s.client.activeOrders[j].Id {
				s.client.activeOrders[j] = s.client.activeOrders[len(s.client.activeOrders)-1]
				s.client.activeOrders = s.client.activeOrders[:len(s.client.activeOrders)-1]
				break
			}
		}
	}
}

func (s *Sub) changeLeverageInLocalState(isLeverageChangeRequired bool, newLeverage float64) {
	if isLeverageChangeRequired {
		s.client.leverage.Store(newLeverage)
	}
}

func (s *Sub) updatePositionsInLocalStateFromMarketOrderRequests(marketOrders ...*orders.RequestForPlaceOrder) {
	for i := range marketOrders {
		if marketOrders[i].Type != "market" {
			continue
		}

		for j := range s.client.openPositions {
			if s.client.openPositions[j].Market == marketOrders[i].Market {
				if marketOrders[i].Side == "buy" {
					s.client.openPositions[j].Size += math.Abs(marketOrders[i].Size)
					tools.RoundFloatPointer(&s.client.openPositions[j].Size)

					// remove positions if size is 0
					if s.client.openPositions[j].Size == 0 {
						s.client.openPositions[j] = s.client.openPositions[len(s.client.openPositions)-1]
						s.client.openPositions = s.client.openPositions[:len(s.client.openPositions)-1]
					}

					s.client.addToFillsAdjuster(marketOrders[i].Market, math.Abs(marketOrders[i].Size))
				} else {
					s.client.openPositions[j].Size += -math.Abs(marketOrders[i].Size)
					tools.RoundFloatPointer(&s.client.openPositions[j].Size)

					// remove positions if size is 0
					if s.client.openPositions[j].Size == 0 {
						s.client.openPositions[j] = s.client.openPositions[len(s.client.openPositions)-1]
						s.client.openPositions = s.client.openPositions[:len(s.client.openPositions)-1]
					}

					s.client.addToFillsAdjuster(marketOrders[i].Market, -math.Abs(marketOrders[i].Size))
				}
				break
			}
		}
	}
}

func (s *Sub) updateOrdersInLocalStateFromOrderResponses(orderResponses ...*orders.ResponseForPlaceOrder) {
	for i := range orderResponses {
		if orderResponses[i].Type != "limit" {
			continue
		}

		newOrder := order{
			Id:            int64(orderResponses[i].ID),
			Market:        orderResponses[i].Market,
			Type:          orderResponses[i].Type,
			Side:          orderResponses[i].Side,
			Price:         orderResponses[i].Price,
			Size:          orderResponses[i].Size,
			FilledSize:    orderResponses[i].FilledSize,
			RemainingSize: orderResponses[i].RemainingSize,
			AvgFillPrice:  0,
			Status:        orderResponses[i].Status,
			ReduceOnly:    orderResponses[i].ReduceOnly,
			Ioc:           orderResponses[i].Ioc,
			PostOnly:      orderResponses[i].PostOnly,
		}

		if newOrder.RemainingSize == 0 {
			continue
		}

		newOrder.ClientId.SetValue(orderResponses[i].ClientID)

		s.client.activeOrders = append(s.client.activeOrders, &newOrder)
	}
}
