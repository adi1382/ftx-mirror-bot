package client

import (
	"sync"

	"github.com/adi1382/ftx-mirror-bot/constants"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/auth"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest"
	"github.com/adi1382/ftx-mirror-bot/websocket"
	"go.uber.org/atomic"
)

func newClient(
	apiKey, secret string,
	leverageUpdateDuration, balanceUpdateDuration int64,
	subRoutineCloser chan int, wg *sync.WaitGroup) *client {

	c := client{apiKey: apiKey}
	c.rest = rest.New(auth.New(apiKey, secret))
	c.leverageUpdateDuration = leverageUpdateDuration
	c.balanceUpdateDuration = balanceUpdateDuration
	c.subRoutineCloser = subRoutineCloser
	c.wg = wg
	c.wsConnection = websocket.NewSocketConnection(apiKey, secret, c.subRoutineCloser, c.wg)
	c.userStream = make(chan []byte, 100)
	c.running.Store(true)
	return &c
}

type client struct {
	apiKey                             string
	rest                               *rest.Client
	symbolsInfo                        map[string]symbolInfo
	symbolInfoLock                     sync.Mutex
	wsConnection                       *websocket.WSConnection
	userStream                         chan []byte
	subscriptionsToUserStream          []chan []byte //is subscribed by subAccounts to hostAccounts
	subscriptionsToUserStreamLock      sync.Mutex
	isRestartRequired                  *atomic.Bool
	leverage                           atomic.Float64
	totalCollateral                    atomic.Float64
	running                            atomic.Bool
	isPositionCoolDownPeriod           atomic.Bool
	subRoutineCloser                   chan int //Pass 1 to close all sub routines
	wg                                 *sync.WaitGroup
	lastFillUnixTime                   int64
	symbolTickers                      map[string]float64
	symbolTickerLock                   sync.Mutex
	symbolTickerLastUpdated            atomic.Int64
	activeOrders                       []*order
	openPositions                      []*position
	activeOrdersLock                   sync.Mutex
	openPositionsLock                  sync.Mutex //Applied till here
	fillsAdjuster                      map[string]float64
	fillsAdjusterLock                  sync.Mutex
	leverageUpdateDuration             int64
	balanceUpdateDuration              int64
	lastAccountInformationCallTimeUnix int64
}

func (c *client) initialize() {
	c.wsConnection.Connect(c.userStream)

	c.wsConnection.AuthenticateWebsocketConnection()
	c.wsConnection.SubscribeToPrivateStreams()
	c.checkIfStreamsAreSuccessfullySubscribed([]string{"fills", "orders"}, constants.TimeoutToCheckForSubscriptions)
	if !c.runningStatus() {
		return
	}

	c.initializeAccountInfoAndPositions()
	c.initializeOrders()
	c.setFillsAdjuster()

	c.wg.Add(1)
	go c.receiveStreamingData()
}

func (c *client) runningStatus() bool {
	return c.running.Load()
}

func (c *client) restart() {
	c.subRoutineCloser <- 0
	c.running.Store(false)
}

//func (c *client) checkForRestart() {
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
