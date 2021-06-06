package host

import (
	"sync"

	"github.com/adi1382/ftx-mirror-bot/client"
)

func NewHostClient(apiKey, apiSecret string, subRoutineCloser chan int, wg *sync.WaitGroup) *Host {
	c := Host{
		client: client.NewClient(apiKey, apiSecret, subRoutineCloser, wg),
	}
	return &c
}

type Host struct {
	client *client.Client
}

func (h *Host) Initialize() {
	h.client.UpdateSymbolInfoViaRest()
	h.client.Initialize()
}

func (h *Host) FetchOpenOrders() []client.Order {
	return h.client.ActiveOrders()
}

func (h *Host) FetchOpenPositions() []client.Position {
	return h.client.ActivePositions()
}

func (h *Host) FetchAccountLeverage() float64 {
	return h.client.FetchLeverage()
}

func (h *Host) FetchTotalCollateral() float64 {
	return h.FetchTotalCollateral()
}

func (h *Host) FetchSymbolInformation() map[string]client.SymbolInfo {
	return h.client.FetchSymbolInformation()
}

func (h *Host) SubscribeToHostUpdates(ch chan []byte) {
	h.client.SubscribeToClientStream(ch)
}
