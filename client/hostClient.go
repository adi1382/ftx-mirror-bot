package client

import (
	"sync"
)

func NewHostClient(
	apiKey, apiSecret string,
	isFTXSubAccount bool, FTXSubAccountName string,
	leverageUpdateDuration, balanceUpdateDuration int64,
	subRoutineCloser chan int, wg *sync.WaitGroup) *Host {
	c := Host{
		client: newClient(apiKey, apiSecret, isFTXSubAccount, FTXSubAccountName, leverageUpdateDuration, balanceUpdateDuration, subRoutineCloser, wg),
	}
	return &c
}

type Host struct {
	client *client
}

func (h *Host) Initialize() {
	h.UpdateSymbolInfoViaRest()
	h.client.initialize()
}

func (h *Host) ActiveOrders() map[int64]order {
	h.client.activeOrdersLock.Lock()
	defer h.client.activeOrdersLock.Unlock()

	openOrders := make(map[int64]order, 5)
	for i := range h.client.activeOrders {
		openOrders[h.client.activeOrders[i].Id] = *h.client.activeOrders[i]
	}
	return openOrders
}

func (h *Host) OpenPositions() map[string]position {
	h.client.openPositionsLock.Lock()
	defer h.client.openPositionsLock.Unlock()

	openPositions := make(map[string]position, 5)
	for i := range h.client.openPositions {
		openPositions[h.client.openPositions[i].Market] = *h.client.openPositions[i]
	}
	return openPositions
}

func (h *Host) AccountLeverage() float64 {
	return h.client.leverage.Load()
}

func (h *Host) TotalCollateral() float64 {
	return h.client.totalCollateral.Load()
}

func (h *Host) UpdateSymbolInfoViaRest() {
	h.client.symbolInfoLock.Lock()
	defer h.client.symbolInfoLock.Unlock()

	h.client.symbolsInfo = make(map[string]symbolInfo, 1000)
	h.client.updateSymbolInfoForFutures()
	h.client.updateSymbolInfoForSpot()
}

func (h *Host) SymbolInformation() map[string]symbolInfo {
	h.client.symbolInfoLock.Lock()
	defer h.client.symbolInfoLock.Unlock()

	symbolInformation := make(map[string]symbolInfo, 1000)
	for k := range h.client.symbolsInfo {
		symbolInformation[k] = h.client.symbolsInfo[k]
	}
	return symbolInformation
}
