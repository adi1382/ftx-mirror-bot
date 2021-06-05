package client

import (
	"github.com/adi1382/ftx-mirror-bot/go-ftx/auth"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/private/fills"
	"github.com/adi1382/ftx-mirror-bot/websocket"
	"go.uber.org/atomic"
	"sync"
	"time"
)

func NewClient(apiKey, secret string, isRestartRequired *atomic.Bool) *Client {

	c := Client{apiKey: apiKey}
	c.rest = rest.New(auth.New(apiKey, secret))
	c.isRestartRequired = isRestartRequired
	c.updateSymbolInfo()
	c.wsConnection = websocket.NewSocketConnection(apiKey, secret, isRestartRequired)
	c.userStream = make(chan []byte, 100)
	c.running.Store(true)
	return &c
}

type Client struct {
	typeOfAccount                  string
	apiKey                         string
	rest                           *rest.Client
	symbolsInfo                    map[string]symbolInfo
	wsConnection                   websocket.WSConnection
	userStream                     chan []byte
	subscriptionsToUserStream      []chan []byte //is subscribed by subAccounts to hostAccounts
	subscriptionsToUserStreamLock  sync.Mutex
	isRestartRequired              *atomic.Bool
	leverage                       atomic.Float64
	totalCollateral                atomic.Float64
	running                        atomic.Bool
	isPositionCoolDownPeriod       atomic.Bool
	fillsForPositionInitialization *fills.Response
	lastFillUnixTime               int64 // Applied till here
	balanceUpdateRate              float64
	symbolTickers                  map[string]float64
	symbolTickerLock               sync.Mutex
	symbolTickerLastUpdated        atomic.Int64
	openOrders                     []*order
	openPositions                  []*position
	openOrdersLock                 sync.Mutex
	openPositionsLock              sync.Mutex
	isInitializationCompleted      atomic.Bool
	lastBalanceUpdateTimeUnix      atomic.Int64
	nextBalanceUpdateTimeUnix      atomic.Int64
	listenKeyLastUpdated           atomic.Int64 // This is still required to be fully implemented
	wg                             *sync.WaitGroup
}

func SubscribeToClientStream(c *Client, ch chan []byte) {
	c.subscriptionsToUserStreamLock.Lock()
	c.subscriptionsToUserStream = append(c.subscriptionsToUserStream, ch)
	c.subscriptionsToUserStreamLock.Unlock()
}

func (c *Client) Initialize() {
	c.wsConnection.Connect(c.userStream)

	c.wsConnection.AuthenticateWebsocketConnection()
	c.wsConnection.SubscribeToPrivateStreams()
	c.checkIfStreamsAreSuccessfullySubscribed("fills", "orders")
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
	c.isRestartRequired.Store(true)
	c.running.Store(false)
}

func (c *Client) checkForRestart() {

	for {
		time.Sleep(time.Millisecond)
		if c.isRestartRequired.Load() {
			c.restart()
			return
		}

		if !c.runningStatus() {
			return
		}
	}
}
