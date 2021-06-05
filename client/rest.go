package client

import (
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/private/account"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/private/fills"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/private/orders"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/public/futures"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/public/markets"
	"time"
)

func (c *Client) getAllFuturesList() *futures.ResponseForFutures {
	resp, err := c.rest.Futures(&futures.RequestForFutures{})
	c.restError(err)

	return resp
}

func (c *Client) getAllMarkets() *markets.ResponseForMarkets {
	resp, err := c.rest.Markets(&markets.RequestForMarkets{})
	c.restError(err)

	return resp
}

func (c *Client) getAccountInformation() *account.ResponseForInformation {
	accountInformation, err := c.rest.Information(&account.RequestForInformation{})
	c.restError(err)
	return accountInformation
}

func (c *Client) getFills(seconds int64) *fills.Response {
	//resp, err := c.rest.Fills(&fills.Request{})
	resp, err := c.rest.Fills(&fills.Request{
		Start: time.Now().Unix() - seconds,
	})
	c.restError(err)

	return resp
}

// This Function is not used by the mirror bot, as same results could be obtained from getAccountInformation()
//func (c *Client) getAllPositions(showAvgPrice bool) *account.ResponseForPositions {
//	resp, err := c.rest.Positions(&account.RequestForPositions{ShowAvgPrice: showAvgPrice})
//	c.restError(err)
//
//	return resp
//}

func (c *Client) getAllOpenOrders() *orders.ResponseForOpenOrder {
	resp, err := c.rest.OpenOrder(&orders.RequestForOpenOrder{})
	c.restError(err)

	return resp
}

func (c *Client) restError(err error) {
	if err != nil {
		panic(err)
	}
}
