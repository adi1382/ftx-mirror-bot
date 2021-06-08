package client

import (
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/private/account"
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/private/orders"
)

func (c *client) postCancelOrderByID(orderId int) {
	_, err := c.rest.CancelByID(&orders.RequestForCancelByID{
		OrderID: orderId,
	})
	c.restError(err)
}

func (c *client) postPlaceOrder(order *orders.RequestForPlaceOrder) *orders.ResponseForPlaceOrder {
	resp, err := c.rest.PlaceOrder(order)
	c.restError(err)
	return resp
}

func (c *client) postChangeLeverage(leverage float64) {
	_, err := c.rest.Leverage(&account.RequestForLeverage{Leverage: int(leverage)})
	c.restError(err)
}
