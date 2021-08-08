package client

import (
	"time"

	"github.com/adi1382/ftx-mirror-bot/fpe"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/private/orders"
)

func (s *Sub) StartMirroring() {
	calibrationTicker := time.NewTicker(time.Duration(s.calibrationDuration) * time.Second)
	defer calibrationTicker.Stop()

	s.calibrate()

	for {
		select {
		case hostNewOrderUpdate := <-s.hostNewOrderUpdates:
			s.processNewOrderUpdate(hostNewOrderUpdate)
		case hostCancelOrderUpdate := <-s.hostExistingOrderUpdates:
			s.processExistingOrderUpdate(hostCancelOrderUpdate)
		case <-calibrationTicker.C:
			s.calibrate()
		case <-s.subRoutineCloser:
			s.subRoutineCloser <- 1
			return
		}
	}
}

func (s *Sub) processNewOrderUpdate(hostNewOrderUpdate *order) {
	s.client.activeOrdersLock.Lock()
	s.client.openPositionsLock.Lock()
	defer s.client.openPositionsLock.Unlock()
	defer s.client.activeOrdersLock.Unlock()

	orderRequest := new(orders.RequestForPlaceOrder)

	if hostNewOrderUpdate.Type == "market" {
		orderRequest = s.generateMarketOrder(hostNewOrderUpdate.Market,
			hostNewOrderUpdate.Side,
			s.adjustedSize(hostNewOrderUpdate.Size, hostNewOrderUpdate.Market),
			hostNewOrderUpdate.ReduceOnly)
		s.placeOrders(orderRequest)
		s.updatePositionsInLocalStateFromMarketOrderRequests(orderRequest)
	} else {
		orderRequest = &orders.RequestForPlaceOrder{
			ClientID:   fpe.GenerateClOrdIDFromOrdID(hostNewOrderUpdate.Id),
			Type:       hostNewOrderUpdate.Type,
			Market:     hostNewOrderUpdate.Market,
			Side:       hostNewOrderUpdate.Side,
			Price:      hostNewOrderUpdate.Price,
			Size:       s.adjustedSize(hostNewOrderUpdate.Size, hostNewOrderUpdate.Market),
			ReduceOnly: hostNewOrderUpdate.ReduceOnly,
			Ioc:        hostNewOrderUpdate.Ioc,
			PostOnly:   hostNewOrderUpdate.PostOnly,
		}
		orderResponse := s.placeOrders(orderRequest)
		s.updateOrdersInLocalStateFromOrderResponses(orderResponse...)
	}
}

func (s *Sub) processExistingOrderUpdate(hostExistingOrderUpdate *order) {
	s.client.activeOrdersLock.Lock()
	s.client.openPositionsLock.Lock()
	defer s.client.openPositionsLock.Unlock()
	defer s.client.activeOrdersLock.Unlock()
	//s.client.postCancelOrderByClOrdID()

	if hostExistingOrderUpdate.FilledSize != hostExistingOrderUpdate.Size {
		toCancelClOrdIDs := s.matchingClientIDs(hostExistingOrderUpdate.Id)
		s.cancelOrderByClOrdID(toCancelClOrdIDs)
		s.removeOrdersInLocalStateFromClOrdIDs(toCancelClOrdIDs)
	}
	//TODO: CHECK FOR ORDER AMENDMENT
}

func (s *Sub) matchingClientIDs(orderID int64) []string {
	toCancelClOrdIDs := make([]string, 0, 1)
	for i := range s.client.activeOrders {
		if orderID == fpe.GenerateOrdIDFromClOrdID(s.client.activeOrders[i].ClientId.Value) {
			toCancelClOrdIDs = append(toCancelClOrdIDs, s.client.activeOrders[i].ClientId.Value)
		}
	}
	return toCancelClOrdIDs
}
