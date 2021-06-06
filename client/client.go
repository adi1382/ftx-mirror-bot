package client

import (
	"sync"

	"github.com/adi1382/ftx-mirror-bot/go-ftx/auth"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest"
	"github.com/adi1382/ftx-mirror-bot/websocket"
	"go.uber.org/atomic"
)

func NewClient(apiKey, secret string, subRoutineCloser chan int, wg *sync.WaitGroup) *Client {
	c := Client{apiKey: apiKey}
	c.rest = rest.New(auth.New(apiKey, secret))
	c.subRoutineCloser = subRoutineCloser
	c.wg = wg
	c.wsConnection = websocket.NewSocketConnection(apiKey, secret, c.subRoutineCloser, c.wg)
	c.userStream = make(chan []byte, 100)
	c.running.Store(true)
	return &c
}

type Client struct {
	apiKey                        string
	rest                          *rest.Client
	symbolsInfo                   map[string]SymbolInfo
	symbolInfoLock                sync.Mutex
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
	openOrders                    []*Order
	openPositions                 []*Position
	openOrdersLock                sync.Mutex
	openPositionsLock             sync.Mutex
	balanceUpdateRate             float64 // Applied till here
	isInitializationCompleted     atomic.Bool
	lastBalanceUpdateTimeUnix     atomic.Int64
	nextBalanceUpdateTimeUnix     atomic.Int64
}

func (c *Client) runningStatus() bool {
	return c.running.Load()
}

func (c *Client) restart() {
	c.subRoutineCloser <- 0
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
