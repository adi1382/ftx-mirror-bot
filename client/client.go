package client

import (
	"github.com/adi1382/ftx-mirror-bot/go-ftx/auth"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest"
	"github.com/adi1382/ftx-mirror-bot/websocket"
	"go.uber.org/atomic"
	"sync"
)

func NewClient(apiKey, secret string, isRestartRequired *atomic.Bool) *Client {

	c := Client{apiKey: apiKey}
	c.rest = rest.New(auth.New(apiKey, secret))
	c.isRestartRequired = isRestartRequired
	c.updateSymbolInfo()
	c.wsConnection = websocket.NewSocketConnection(apiKey, secret, isRestartRequired)
	return &c
}

type Client struct {
	apiKey                        string
	rest                          *rest.Client
	symbolsInfo                   []*symbolInfo
	wsConnection                  websocket.WSConnection
	userStream                    chan []byte
	subscriptionsToUserStream     []chan []byte //is subscribed by subAccounts to hostAccounts
	subscriptionsToUserStreamLock sync.Mutex
	isRestartRequired             *atomic.Bool
	leverage                      float64
	subAccountSubscriptions       []chan []byte
	testnet                       bool
	running                       atomic.Bool
	balanceUpdateRate             float64
	symbolTickers                 map[string]float64
	symbolTickerLock              sync.Mutex
	symbolTickerLastUpdated       atomic.Int64
	openOrders                    []*order
	allPositions                  []*position
	openOrdersLock                sync.Mutex
	positionLock                  sync.Mutex
	isInitializationCompleted     atomic.Bool
	lastBalanceUpdateTimeUnix     atomic.Int64
	nextBalanceUpdateTimeUnix     atomic.Int64
	listenKeyLastUpdated          atomic.Int64 // This is still required to be fully implemented
	wg                            *sync.WaitGroup
}

func SubscribeToClientStream(c *Client, ch chan []byte) {
	c.subscriptionsToUserStreamLock.Lock()
	c.subscriptionsToUserStream = append(c.subscriptionsToUserStream, ch)
	c.subscriptionsToUserStreamLock.Unlock()
}

func (c *Client) Initialize() {

}
