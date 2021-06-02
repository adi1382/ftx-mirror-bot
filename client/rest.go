package client

import (
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/public/futures"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/public/markets"
)

func (c *Client) getAllFuturesList() *futures.ResponseForFutures {
	resp, err := c.rest.Futures(&futures.RequestForFutures{})
	if err != nil {
		c.restError(err)
	}

	return resp
}

func (c *Client) getAllMarkets() *markets.ResponseForMarkets {
	resp, err := c.rest.Markets(&markets.RequestForMarkets{})
	if err != nil {
		c.restError(err)
	}

	return resp
}

func (c *Client) restError(err error) {
	panic(err)
}
