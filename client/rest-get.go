package client

import (
	"time"

	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/private/account"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/private/fills"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/private/orders"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/public/futures"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/public/markets"
)

func (c *client) getAllFuturesList() *futures.ResponseForFutures {
	resp, err := c.rest.Futures(&futures.RequestForFutures{})
	c.restError(err)

	return resp
}

func (c *client) getAllMarkets() *markets.ResponseForMarkets {
	resp, err := c.rest.Markets(&markets.RequestForMarkets{})
	c.restError(err)

	return resp
}

func (c *client) getAccountInformation() *account.ResponseForInformation {
	accountInformation, err := c.rest.Information(&account.RequestForInformation{})
	c.restError(err)
	return accountInformation
}

func (c *client) getFills(seconds time.Duration) *fills.Response {
	//resp, err := c.rest.Fills(&fills.Request{})
	resp, err := c.rest.Fills(&fills.Request{
		Start: time.Now().Unix() - int64(seconds/time.Second),
	})
	c.restError(err)

	return resp
}

// This Function is not used by the mirror bot, as same results could be obtained from getAccountInformation()
//func (c *client) getAllPositions(showAvgPrice bool) *account.ResponseForPositions {
//	resp, err := c.rest.Positions(&account.RequestForPositions{ShowAvgPrice: showAvgPrice})
//	c.restError(err)
//
//	return resp
//}

func (c *client) getAllOpenOrders() *orders.ResponseForOpenOrder {
	resp, err := c.rest.OpenOrder(&orders.RequestForOpenOrder{})
	c.restError(err)

	return resp
}

func (c *client) restError(err error) {
	if err != nil {
		c.restart()
		panic(err)
	}
}
