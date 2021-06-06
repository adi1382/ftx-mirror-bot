package host

import (
	"sync"

	"github.com/adi1382/ftx-mirror-bot/client"
)

func NewHostClient() {

}

type Host struct {
	client            *client.Client
	openOrders        []*client.Order
	openPositions     []*client.Position
	openOrdersLock    sync.Mutex
	openPositionsLock sync.Mutex
}
