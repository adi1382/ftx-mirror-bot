package client

import (
	"github.com/adi1382/ftx-mirror-bot/constants"
	"github.com/adi1382/ftx-mirror-bot/fpe"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/private/orders"
	"strconv"
	"strings"
)

func (s *Sub) calibrate() {
	s.client.activeOrdersLock.Lock()
	s.client.openPositionsLock.Lock()
	defer s.client.activeOrdersLock.Unlock()
	defer s.client.openPositionsLock.Unlock()

	hostOrders := s.hostClient.ActiveOrders()
	hostPositions := s.hostClient.OpenPositions()

	toCancelOrderIDs := s.findUnwantedOrderIDs(hostOrders)
	//TODO: New orders that needs to be placed

	isLeverageChangeRequired, newLeverage := s.isLeverageChangeRequired()
	ordersToCalibrateExisitingPositions := s.generateMarketOrdersToCalibrateExistingPositions(hostPositions)
	ordersToCalibrateNewPositions := s.generateMarketOrdersToCalibrateNewPositions()
}

//generateMarketOrdersToCalibrateNewPositions must be called after calling generateMarketOrdersToCalibrateExistingPositions
//It returns a slice containing requests for market orders that calibrates the existing positions with host
//by creating new open positions
func (s *Sub) generateMarketOrdersToCalibrateNewPositions(hostPositions map[string]position) []*orders.RequestForPlaceOrder {
	requestsForMarketOrders := make([]*orders.RequestForPlaceOrder, 0, 5)
	//TODO
	return requestsForMarketOrders
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
					s.adjustedSize(v.Size, v.Market)-s.client.openPositions[i].Size))
		} else {
			// no change is required for this Market
			delete(hostPositions, s.client.openPositions[i].Market)
		}
	}
	return requestsForMarketOrders
}

func (s *Sub) generateMarketOrder(market string, size float64) *orders.RequestForPlaceOrder {
	var side string

	if size > 0 {
		side = "buy"
	} else if size < 0 {
		side = "sell"
	} else {
		return nil
	}

	return &orders.RequestForPlaceOrder{
		ClientID: fpe.GenerateRandomClOrdID(),
		Type:     "market",
		Market:   market,
		Side:     side,
		Size:     size,
	}
}

func (s *Sub) isLeverageChangeRequired() (bool, float64) {
	if s.isCopyLeverage {
		hostLeverage := s.hostClient.AccountLeverage()
		if s.client.leverage.Load() != hostLeverage {
			return true, hostLeverage
		}
	}
	return false, -1
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
