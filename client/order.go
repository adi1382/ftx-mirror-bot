package client

import (
	"github.com/adi1382/ftx-mirror-bot/go-ftx/rest/private/orders"
	"github.com/adi1382/ftx-mirror-bot/optional"
)

type order struct {
	Id            int64           `json:"id"`
	Market        string          `json:"market"`
	Type          string          `json:"type"`
	Side          string          `json:"side"`
	Price         float64         `json:"price"`
	Size          float64         `json:"size"`
	FilledSize    float64         `json:"filledSize"`
	RemainingSize float64         `json:"remainingSize"`
	AvgFillPrice  float64         `json:"avgFillPrice"`
	Status        string          `json:"status"`
	ReduceOnly    bool            `json:"reduceOnly"`
	Ioc           bool            `json:"ioc"`
	PostOnly      bool            `json:"postOnly"`
	ClientId      optional.String `json:"clientId"`
}

func (c *Client) initializeOrders() {
	c.openOrdersLock.Lock()
	c.openOrders = c.openOrders[:0]
	c.openOrdersLock.Unlock()

	openOrders := c.getAllOpenOrders()

	c.generateNativeOrdersFromRestResponse(openOrders)
}

func (c *Client) generateNativeOrdersFromRestResponse(openOrders *orders.ResponseForOpenOrder) {
	c.openOrdersLock.Lock()
	defer c.openOrdersLock.Unlock()

	for i := range *openOrders {
		c.openOrders = append(c.openOrders, c.generateNativeOrderFromRestOrder((*openOrders)[i]))
	}
}

func (c *Client) generateNativeOrderFromRestOrder(restOrder orders.OpenOrder) *order {

	nativeOrder := new(order)

	nativeOrder.Id = restOrder.ID
	nativeOrder.Market = restOrder.Market
	nativeOrder.Type = restOrder.Type
	nativeOrder.Side = restOrder.Side
	nativeOrder.Price = restOrder.Price
	nativeOrder.Size = restOrder.Size
	nativeOrder.FilledSize = restOrder.FilledSize
	nativeOrder.RemainingSize = restOrder.RemainingSize
	nativeOrder.AvgFillPrice = restOrder.AvgFillPrice
	nativeOrder.Status = restOrder.Status
	nativeOrder.ReduceOnly = restOrder.ReduceOnly
	nativeOrder.Ioc = restOrder.Ioc
	nativeOrder.PostOnly = restOrder.PostOnly

	if restOrder.ClientID != "" {
		nativeOrder.ClientId.SetValue(restOrder.ClientID)
	}

	return nativeOrder
}

///////////////////////// BEGIN --> STREAM ORDER FUNCTIONALITIES /////////////////////////

func (c *Client) handleOrderUpdateFromStream(newOrder *order) {
	c.openOrdersLock.Lock()
	defer c.openOrdersLock.Unlock()

	removalRequired := c.checkIfOrderNeedsToBeRemoved(newOrder)

	if index := c.checkIfOrderAlreadyExists(newOrder); index > -1 {
		c.updateExistingOrder(newOrder, index, removalRequired)
		return
	}

	if !removalRequired {
		c.insertNewOrder(newOrder)
		return
	}

}

// This function should only be called from handleOrderUpdateFromStream because of mutex synchronizations
func (c *Client) checkIfOrderAlreadyExists(newOrder *order) int {
	indexOfOrder := -1

	for i := range c.openOrders {
		if c.openOrders[i].Id == newOrder.Id {
			return i
		}
	}

	return indexOfOrder
}

// This function should only be called from handleOrderUpdateFromStream because of mutex synchronizations
func (c *Client) updateExistingOrder(newOrder *order, existingOrderIndex int, isRemovalRequired bool) {
	c.openOrders[existingOrderIndex] = newOrder

	if isRemovalRequired {
		c.openOrders = append(c.openOrders[:existingOrderIndex], c.openOrders[existingOrderIndex+1:]...)
	}

}

// This function should only be called from handleOrderUpdateFromStream because of mutex synchronizations
func (c *Client) checkIfOrderNeedsToBeRemoved(newOrder *order) bool {
	if newOrder.RemainingSize == 0 || newOrder.Type == "market" || newOrder.Status == "closed" {
		return true
	}
	return false
}

// This function should only be called from handleOrderUpdateFromStream because of mutex synchronizations
func (c *Client) insertNewOrder(newOrder *order) {
	c.openOrders = append(c.openOrders, newOrder)
}

///////////////////////// END --> STREAM ORDER FUNCTIONALITIES /////////////////////////
