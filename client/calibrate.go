package client

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/adi1382/ftx-mirror-bot/constants"
	"github.com/adi1382/ftx-mirror-bot/fpe"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/private/orders"
	"github.com/adi1382/ftx-mirror-bot/tools"
)

func (s *Sub) calibrate() {
	fmt.Println("Calibration started @", time.Now())
	s.client.activeOrdersLock.Lock()
	s.client.openPositionsLock.Lock()
	defer s.client.openPositionsLock.Unlock()
	defer s.client.activeOrdersLock.Unlock()

	hostOrders := s.hostClient.ActiveOrders()
	hostPositions := s.hostClient.OpenPositions()

	toCancelOrderIDs := s.findUnwantedOrderIDs(hostOrders)
	ordersForActiveOrdersOnHostAccount := s.generateNewOrdersForExistingOrdersOnHostAccount(hostOrders)
	isLeverageChangeRequired, newLeverage := s.isLeverageChangeRequired()
	ordersToCalibrateExistingPositions := s.generateMarketOrdersToCalibrateExistingPositions(hostPositions)
	ordersToCalibrateNewPositions := s.generateMarketOrdersToCalibrateNewPositions(hostPositions)

	//TODO: concurrent rest requests with could improve time
	s.cancelOrderIDs(toCancelOrderIDs)
	s.updateLeverage(isLeverageChangeRequired, newLeverage)
	s.placeOrders(ordersToCalibrateExistingPositions)                   // responses not needed as all are market orders
	s.placeOrders(ordersToCalibrateNewPositions)                        // responses not needed as all are market orders
	orderResponses := s.placeOrders(ordersForActiveOrdersOnHostAccount) // responses are needed for order ID

	s.removeOrdersInLocalStateFromOrderIDs(toCancelOrderIDs)
	s.changeLeverageInLocalState(isLeverageChangeRequired, newLeverage)
	s.updatePositionsInLocalStateFromMarketOrderRequests(ordersToCalibrateExistingPositions)
	s.updatePositionsInLocalStateFromMarketOrderRequests(ordersToCalibrateNewPositions)
	s.updateOrdersInLocalStateFromOrderResponses(orderResponses)
}

//findUnwantedOrderIDs modifies the map hostOrders by removing orders which have their sub orders already in place
//It returns the orderIDs on the sub account for the orders that needs to be canceled
//argument hostOrders is only left with orders that are not yet active on sub account
func (s *Sub) findUnwantedOrderIDs(hostOrders map[string]order) []int64 {
	toCancelOrderIDs := make([]int64, 0, 5)

	for i := range s.client.activeOrders {

		if !s.client.activeOrders[i].ClientId.Valid {
			// Sub order does not have ClientId
			toCancelOrderIDs = append(toCancelOrderIDs, s.client.activeOrders[i].Id)
			continue
		} else if !s.verifyClientID(s.client.activeOrders[i].ClientId.Value) {
			// Order not placed by mirror bot
			toCancelOrderIDs = append(toCancelOrderIDs, s.client.activeOrders[i].Id)
			continue
		}

		if v, ok := hostOrders[fpe.GenerateOrdIDFromClOrdID(s.client.activeOrders[i].ClientId.Value)]; !ok {
			// subsequent host order does not exists
			toCancelOrderIDs = append(toCancelOrderIDs, s.client.activeOrders[i].Id)
			continue
		} else if v.Price != s.client.activeOrders[i].Price ||
			s.adjustedSize(v.RemainingSize, v.Market) != s.client.activeOrders[i].RemainingSize {
			//TODO: Implement Order Amendment here
			toCancelOrderIDs = append(toCancelOrderIDs, s.client.activeOrders[i].Id)
			continue
		} else {
			// no change is required for this order
			delete(hostOrders, fpe.GenerateOrdIDFromClOrdID(s.client.activeOrders[i].ClientId.Value))
		}
	}

	return toCancelOrderIDs
}

//generateNewOrdersForExistingOrdersOnHostAccount must be called after calling findUnwantedOrderIDs
//It returns a slice containing requests for orders (limit) that calibrates the existing pending orders of host
//by creating new open orders
func (s *Sub) generateNewOrdersForExistingOrdersOnHostAccount(hostOrders map[string]order) []*orders.RequestForPlaceOrder {
	ordersForActiveOrdersOnHostAccount := make([]*orders.RequestForPlaceOrder, 0, 5)

	for i := range hostOrders {
		if s.adjustedSize(hostOrders[i].RemainingSize, hostOrders[i].Market) == 0 {
			continue
		}
		ordersForActiveOrdersOnHostAccount = append(ordersForActiveOrdersOnHostAccount,
			&orders.RequestForPlaceOrder{
				ClientID:   fpe.GenerateClOrdIDFromOrdID(strconv.Itoa(int(hostOrders[i].Id))),
				Type:       hostOrders[i].Type,
				Market:     hostOrders[i].Market,
				Side:       hostOrders[i].Side,
				Price:      hostOrders[i].Price,
				Size:       s.adjustedSize(hostOrders[i].RemainingSize, hostOrders[i].Market),
				ReduceOnly: hostOrders[i].ReduceOnly,
				Ioc:        hostOrders[i].Ioc,
				PostOnly:   hostOrders[i].PostOnly,
			})
	}
	return ordersForActiveOrdersOnHostAccount
}

//isLeverageChangeRequired checks if leverage is calibrated with host or not, and returns leverage as float64
//if not calibrated
func (s *Sub) isLeverageChangeRequired() (bool, float64) {
	if s.isCopyLeverage {
		hostLeverage := s.hostClient.AccountLeverage()
		if s.client.leverage.Load() != hostLeverage {
			return true, hostLeverage
		}
	}
	return false, s.client.leverage.Load()
}

//generateMarketOrdersToCalibrateExistingPositions modifies the map hostPositions by removing keys
//of Markets where no change is required
//It returns a slice containing requests for market orders that calibrates the existing positions with host
//argument hostPositions is only left with positions that are not yet open on sub account
func (s *Sub) generateMarketOrdersToCalibrateExistingPositions(hostPositions map[string]position) []*orders.RequestForPlaceOrder {

	requestsForMarketOrders := make([]*orders.RequestForPlaceOrder, 0, 5)

	for i := range s.client.openPositions {
		if v, ok := hostPositions[s.client.openPositions[i].Market]; !ok {
			// Sub position does not exists on host account
			requestsForMarketOrders = append(
				requestsForMarketOrders,
				s.generateMarketOrder(s.client.openPositions[i].Market, -s.client.openPositions[i].Size))
		} else if s.adjustedSize(v.Size, v.Market) != s.client.openPositions[i].Size {
			// sub positions size not as expected
			requestsForMarketOrders = append(
				requestsForMarketOrders,
				s.generateMarketOrder(
					s.client.openPositions[i].Market,
					tools.RoundFloat(s.adjustedSize(v.Size, v.Market)-s.client.openPositions[i].Size)))
			delete(hostPositions, s.client.openPositions[i].Market)
		} else {
			// no change is required for this Market
			delete(hostPositions, s.client.openPositions[i].Market)
		}
	}
	return requestsForMarketOrders
}

//generateMarketOrdersToCalibrateNewPositions must be called after calling generateMarketOrdersToCalibrateExistingPositions
//It returns a slice containing requests for market orders that calibrates the existing positions with host
//by creating new open positions
func (s *Sub) generateMarketOrdersToCalibrateNewPositions(hostPositions map[string]position) []*orders.RequestForPlaceOrder {
	requestsForMarketOrders := make([]*orders.RequestForPlaceOrder, 0, 5)
	for i := range hostPositions {
		if s.adjustedSize(hostPositions[i].Size, hostPositions[i].Market) == 0 {
			continue
		}
		requestsForMarketOrders = append(
			requestsForMarketOrders,
			s.generateMarketOrder(
				hostPositions[i].Market, s.adjustedSize(hostPositions[i].Size, hostPositions[i].Market)))
	}
	return requestsForMarketOrders
}

//generateMarketOrder generates requests for market order, it does not place any order.
func (s *Sub) generateMarketOrder(market string, size float64) *orders.RequestForPlaceOrder {
	fmt.Println("Generate Market Orders: ", market, size)
	var side string

	if size > 0 {
		side = "buy"
	} else if size < 0 {
		side = "sell"
	} else {
		panic("size zero order creation")
	}

	return &orders.RequestForPlaceOrder{
		ClientID: fpe.GenerateRandomClOrdID(),
		Type:     "market",
		Market:   market,
		Side:     side,
		Size:     math.Abs(size),
	}
}

//verifyClientID verifies if the clientID is placed by mirror bot or not
func (s *Sub) verifyClientID(clientID string) bool {
	if !strings.HasPrefix(clientID, constants.ClientOrderIDPrefix) {
		return false
	}

	clID := strings.TrimPrefix(clientID, constants.ClientOrderIDPrefix)
	if len(clientID) < 2+constants.ClientOrderIDSuffixLength {
		return false
	}

	clID = clID[:len(clID)-constants.ClientOrderIDSuffixLength]
	if _, err := strconv.Atoi(clID); err != nil {
		return false
	}

	return true
}

// These functions mutate the state of Sub type
// removeOrdersInLocalStateFromOrderIDs removes orders from s.client.activeOrders which got canceled by calibrator
func (s *Sub) removeOrdersInLocalStateFromOrderIDs(orderIDs []int64) {
	for i := range orderIDs {
		for j := range s.client.activeOrders {
			if orderIDs[i] == s.client.activeOrders[j].Id {
				s.client.activeOrders = append(s.client.activeOrders[:j], s.client.activeOrders[j+1:]...)
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

func (s *Sub) updatePositionsInLocalStateFromMarketOrderRequests(marketOrders []*orders.RequestForPlaceOrder) {
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

func (s *Sub) updateOrdersInLocalStateFromOrderResponses(orderResponses []*orders.ResponseForPlaceOrder) {
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
