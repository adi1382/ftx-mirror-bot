package client

import (
	"github.com/adi1382/binance-copy-bot/rest/futures"
	"github.com/adi1382/binance-copy-bot/tools"
	"go.uber.org/atomic"
	"sync"
)

func NewClient() {

}

type Client struct {
	leverage                  float64
	UserStream                chan []byte
	subAccountSubscriptions   []chan []byte
	apiKey                    string
	testnet                   bool
	running                   atomic.Bool
	balanceUpdateRate         float64
	exchangeInfo              futures.ExchangeInfo
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
	logger                    *tools.LocalLogger
}

func (c *Client) Initialize() {

}
