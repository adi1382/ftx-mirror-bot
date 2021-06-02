package client

import (
	"github.com/adi1382/ftx-mirror-bot/go-ftx/auth"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest"
	"go.uber.org/atomic"
	"sync"
)

func NewClient(apiKey, secret string) *Client {

	c := Client{apiKey: apiKey}
	c.rest = rest.New(auth.New(apiKey, secret))

	return &c
}

type Client struct {
	apiKey                    string
	rest                      *rest.Client
	leverage                  float64
	UserStream                chan []byte
	subAccountSubscriptions   []chan []byte
	testnet                   bool
	running                   atomic.Bool
	balanceUpdateRate         float64
	symbolTickers             map[string]float64
	symbolTickerLock          sync.Mutex
	symbolTickerLastUpdated   atomic.Int64
	openOrders                []*order
	allPositions              []*position
	openOrdersLock            sync.Mutex
	positionLock              sync.Mutex
	isRestartRequired         *atomic.Bool
	isInitializationCompleted atomic.Bool
	lastBalanceUpdateTimeUnix atomic.Int64
	nextBalanceUpdateTimeUnix atomic.Int64
	listenKeyLastUpdated      atomic.Int64 // This is still required to be fully implemented
	wg                        *sync.WaitGroup
}

func (c *Client) Initialize() {

}
