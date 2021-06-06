package client

import (
	"sync"
)

func NewHostClient(apiKey, apiSecret string, subRoutineCloser chan int, wg *sync.WaitGroup) *Host {
	c := Host{
		client: newClient(apiKey, apiSecret, subRoutineCloser, wg),
	}
	return &c
}

type Host struct {
	client *client
}

func (h *Host) Initialize() {
	h.client.UpdateSymbolInfoViaRest()
	h.client.Initialize()
}

func (h *Host) FetchOpenOrders() []order {
	return h.client.ActiveOrders()
}

func (h *Host) FetchOpenPositions() []position {
	return h.client.ActivePositions()
}

func (h *Host) FetchAccountLeverage() float64 {
	return h.client.FetchLeverage()
}

func (h *Host) FetchTotalCollateral() float64 {
	return h.FetchTotalCollateral()
}

func (h *Host) FetchSymbolInformation() map[string]symbolInfo {
	return h.client.FetchSymbolInformation()
}

func (h *Host) SubscribeToHostUpdates(ch chan []byte) {
	h.client.SubscribeToClientStream(ch)
}
