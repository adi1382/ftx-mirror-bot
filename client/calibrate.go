//This file contains the calibration routine
//Calibration routine is ran by the StartMirroring fn at a user-defined interval
//This routine fetches all orders from host and sub account, matches these orders and positions
//If they're are any discrepancies it adjusts them, so the sub account is in sync with the host account
//This routine is only called for subClient
//All the functional available in this file excluding calibrate should only be called from the calibrate function

package client

import (
	"fmt"
	"github.com/adi1382/ftx-mirror-bot/fpe"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/private/orders"
	"github.com/adi1382/ftx-mirror-bot/tools"
)

func (s *Sub) calibrate() {
	s.client.activeOrdersLock.Lock()
	s.client.openPositionsLock.Lock()
	defer s.client.openPositionsLock.Unlock()
	defer s.client.activeOrdersLock.Unlock()

	hostOrders := s.hostClient.ActiveOrders()
	hostPositions := s.hostClient.OpenPositions()

	fmt.Println("Host:", s.hostClient.TotalCollateral())
	fmt.Println("Sub:", s.client.totalCollateral.Load())

	toCancelOrderIDs := s.findUnwantedOrderIDs(hostOrders)
	ordersForActiveOrdersOnHostAccount := s.generateNewOrdersForExistingOrdersOnHostAccount(hostOrders)
	isLeverageChangeRequired, newLeverage := s.isLeverageChangeRequired()
	ordersToCalibrateExistingPositions := s.generateMarketOrdersToCalibrateExistingPositions(hostPositions)
	ordersToCalibrateNewPositions := s.generateMarketOrdersToCalibrateNewPositions(hostPositions)

	//TODO: concurrent rest requests with could improve time
	s.cancelOrderByIDs(toCancelOrderIDs)
	if newLeverage > s.client.leverage.Load() {
		s.updateLeverage(isLeverageChangeRequired, newLeverage)
	}

	s.placeOrders(ordersToCalibrateExistingPositions...) //responses not needed as all are market orders
	//s.updateLeverage(isLeverageChangeRequired, newLeverage)
	s.placeOrders(ordersToCalibrateNewPositions...)                        //responses not needed as all are market orders
	orderResponses := s.placeOrders(ordersForActiveOrdersOnHostAccount...) //responses are needed for order ID

	if newLeverage < s.client.leverage.Load() {
		s.updateLeverage(isLeverageChangeRequired, newLeverage)
	}

	s.removeOrdersInLocalStateFromOrderIDs(toCancelOrderIDs)
	s.changeLeverageInLocalState(isLeverageChangeRequired, newLeverage)
	s.updatePositionsInLocalStateFromMarketOrderRequests(ordersToCalibrateExistingPositions...)
	s.updatePositionsInLocalStateFromMarketOrderRequests(ordersToCalibrateNewPositions...)
	s.updateOrdersInLocalStateFromOrderResponses(orderResponses...)
}

//findUnwantedOrderIDs modifies the map hostOrders by removing orders which have their sub orders already in place
//It returns the orderIDs on the sub account for the orders that needs to be canceled
//argument hostOrders is only left with orders that are not yet active on sub account
func (s *Sub) findUnwantedOrderIDs(hostOrders map[int64]order) []int64 {
	toCancelOrderIDs := make([]int64, 0, 5)

	for i := range s.client.activeOrders {

		if !s.client.activeOrders[i].ClientId.Valid {
			// Sub order does not have ClientId
			toCancelOrderIDs = append(toCancelOrderIDs, s.client.activeOrders[i].Id)
			continue
		} else if tools.VerifyClientID(s.client.activeOrders[i].ClientId.Value) {
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
func (s *Sub) generateNewOrdersForExistingOrdersOnHostAccount(hostOrders map[int64]order) []*orders.RequestForPlaceOrder {
	ordersForActiveOrdersOnHostAccount := make([]*orders.RequestForPlaceOrder, 0, 5)

	for i := range hostOrders {
		if s.adjustedSize(hostOrders[i].RemainingSize, hostOrders[i].Market) == 0 {
			continue
		}
		ordersForActiveOrdersOnHostAccount = append(ordersForActiveOrdersOnHostAccount,
			&orders.RequestForPlaceOrder{
				ClientID:   fpe.GenerateClOrdIDFromOrdID(hostOrders[i].Id),
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
	fmt.Println("Leverage")
	fmt.Println(s.hostClient.AccountLeverage(), s.client.leverage.Load())
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
				s.generateMarketOrder(s.client.openPositions[i].Market,
					"auto", -s.client.openPositions[i].Size, true))
		} else if s.adjustedSize(v.Size, v.Market) != s.client.openPositions[i].Size {
			// sub positions size not as expected
			requestsForMarketOrders = append(
				requestsForMarketOrders,
				s.generateMarketOrder(
					s.client.openPositions[i].Market,
					"auto",
					tools.RoundFloat(s.adjustedSize(v.Size, v.Market)-s.client.openPositions[i].Size),
					false))
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

		fmt.Println(hostPositions[i].Side, hostPositions[i].Size, hostPositions[i].Market)
		fmt.Println(s.adjustedSize(hostPositions[i].Size, hostPositions[i].Market))

		requestsForMarketOrders = append(
			requestsForMarketOrders,
			s.generateMarketOrder(
				hostPositions[i].Market,
				"auto",
				s.adjustedSize(hostPositions[i].Size, hostPositions[i].Market),
				false))
	}
	return requestsForMarketOrders
}
