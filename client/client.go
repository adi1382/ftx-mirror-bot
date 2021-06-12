package client

import (
	"strconv"
	"strings"
	"sync"

	"github.com/adi1382/ftx-mirror-bot/constants"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/auth"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest"
	"github.com/adi1382/ftx-mirror-bot/websocket"
	"go.uber.org/atomic"
)

func newClient(
	apiKey, secret string,
	isFTXSubAccount bool, FTXSubAccountName string,
	leverageUpdateDuration, balanceUpdateDuration int64,
	subRoutineCloser chan int, wg *sync.WaitGroup) *client {

	c := client{}
	c.config = auth.New(apiKey, secret, isFTXSubAccount, FTXSubAccountName)
	c.rest = rest.New(c.config)
	c.leverageUpdateDuration = leverageUpdateDuration
	c.balanceUpdateDuration = balanceUpdateDuration
	c.subRoutineCloser = subRoutineCloser
	c.wg = wg
	c.wsConnection = websocket.NewSocketConnection(c.config, c.subRoutineCloser, c.wg)
	c.userStream = make(chan []byte, 100)
	c.running.Store(true)
	return &c
}

type client struct {
	rest                               *rest.Client
	config                             *auth.Config
	symbolsInfo                        map[string]symbolInfo
	symbolInfoLock                     sync.Mutex
	wsConnection                       *websocket.WSConnection
	userStream                         chan []byte
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

	subscriptionsToNewOrderUpdates      []chan *order //is subscribed by subAccounts to hostAccounts
	subscriptionsToNewOrderUpdatesLock  sync.Mutex
	subscriptionsToExistingOrderUpdates []chan *order //is subscribed by subAccounts to hostAccounts
	subscriptionsToExistingUpdatesLock  sync.Mutex
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

//verifyClientID verifies if the clientID is placed by mirror bot or not
func (s *Sub) verifyClientID(clientID string) bool {
	if !strings.HasPrefix(clientID, constants.ClientOrderIDPrefix) {
		return false
	}

	clID := strings.TrimPrefix(clientID, constants.ClientOrderIDPrefix)
	if len(clientID) < 2+constants.ClientOrderIDSuffixLength {
		return false
	}

	clID = clID[:len(clID)-constants.ClientOrderIDSuffixLength]
	if _, err := strconv.Atoi(clID); err != nil {
		return false
	}

	return true
}
