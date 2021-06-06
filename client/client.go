package client

import (
	"github.com/adi1382/ftx-mirror-bot/constants"
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
	c.subRoutineCloser = make(chan int, 1)
	c.wsConnection = websocket.NewSocketConnection(apiKey, secret, isRestartRequired, c.subRoutineCloser)
	c.userStream = make(chan []byte, 100)
	c.running.Store(true)
	return &c
}

type Client struct {
	apiKey                        string
	rest                          *rest.Client
	symbolsInfo                   map[string]symbolInfo
	wsConnection                  *websocket.WSConnection
	userStream                    chan []byte
	subscriptionsToUserStream     []chan []byte //is subscribed by subAccounts to hostAccounts
	subscriptionsToUserStreamLock sync.Mutex
	isRestartRequired             *atomic.Bool
	leverage                      atomic.Float64
	totalCollateral               atomic.Float64
	running                       atomic.Bool
	isPositionCoolDownPeriod      atomic.Bool
	subRoutineCloser              chan int // Pass 1 to close all sub routines
	wg                            *sync.WaitGroup
	lastFillUnixTime              int64
	symbolTickers                 map[string]float64
	symbolTickerLock              sync.Mutex
	symbolTickerLastUpdated       atomic.Int64
	openOrders                    []*order
	openPositions                 []*position
	openOrdersLock                sync.Mutex
	openPositionsLock             sync.Mutex
	balanceUpdateRate             float64 // Applied till here
	isInitializationCompleted     atomic.Bool
	lastBalanceUpdateTimeUnix     atomic.Int64
	nextBalanceUpdateTimeUnix     atomic.Int64
	//fillsForPositionInitialization *fills.Response // not needed, only last fill unix could be used to remove unnecessary fills through stream

}

//func SubscribeToClientStream(c *Client, ch chan []byte) {
//	c.subscriptionsToUserStreamLock.Lock()
//	c.subscriptionsToUserStream = append(c.subscriptionsToUserStream, ch)
//	c.subscriptionsToUserStreamLock.Unlock()
//}

func (c *Client) Initialize() {
	c.wsConnection.Connect(c.userStream)

	c.wsConnection.AuthenticateWebsocketConnection()
	c.wsConnection.SubscribeToPrivateStreams()
	c.checkIfStreamsAreSuccessfullySubscribed([]string{"fills", "orders"}, constants.TimeoutToCheckForSubscriptions)
	if !c.runningStatus() {
		return
	}

	c.initializeAccountInfoAndPositions()
	c.initializeOrders()

	go c.receiveStreamingData()
}

func (c *Client) runningStatus() bool {
	return c.running.Load()
}

func (c *Client) restart() {
	c.subRoutineCloser <- 0
	c.isRestartRequired.Store(true)
	c.running.Store(false)
}

//func (c *Client) checkForRestart() {
//
//	for {
//		time.Sleep(time.Millisecond)
//		if c.isRestartRequired.Load() {
//			c.restart()
//			return
//		}
//
//		if !c.runningStatus() {
//			return
//		}
//	}
//}

func (c *Client) ActiveOrders() []order {
	c.openOrdersLock.Lock()
	defer c.openOrdersLock.Unlock()

	openOrders := make([]order, 0, 5)
	for i := range c.openOrders {
		openOrders = append(openOrders, *c.openOrders[i])
	}

	return openOrders

}

func (c *Client) ActivePositions() []position {
	c.openPositionsLock.Lock()
	defer c.openPositionsLock.Unlock()

	openPositions := make([]position, 0, 5)
	for i := range c.openPositions {
		openPositions = append(openPositions, *c.openPositions[i])
	}

	return openPositions

}
