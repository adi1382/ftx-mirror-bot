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

func (c *client) initializeOrders() {
	c.activeOrdersLock.Lock()
	c.activeOrders = c.activeOrders[:0]
	c.activeOrdersLock.Unlock()

	openOrders := c.getAllOpenOrders()
	if openOrders == nil {
		return
	}
	c.generateNativeOrdersFromRestResponse(openOrders)
}

func (c *client) generateNativeOrdersFromRestResponse(openOrders *orders.ResponseForOpenOrder) {
	c.activeOrdersLock.Lock()
	defer c.activeOrdersLock.Unlock()

	for i := range *openOrders {
		c.activeOrders = append(c.activeOrders, c.generateNativeOrderFromRestOrder((*openOrders)[i]))
	}
}

func (c *client) generateNativeOrderFromRestOrder(restOrder orders.OpenOrder) *order {

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

func (c *client) handleOrderUpdateFromStream(newOrder *order) {
	c.activeOrdersLock.Lock()
	defer c.activeOrdersLock.Unlock()

	removalRequired := c.checkIfOrderNeedsToBeRemoved(newOrder)

	if index := c.checkIfOrderAlreadyExists(newOrder); index > -1 {
		c.sendExistingOrderUpdateToSubscriptions(*newOrder)
		c.updateExistingOrder(newOrder, index, removalRequired)
		return
	}

	c.sendNewOrderUpdateToSubscriptions(*newOrder)

	if !removalRequired {
		c.insertNewOrder(newOrder)
		return
	}

}

// This function should only be called from handleOrderUpdateFromStream because of mutex synchronizations
func (c *client) checkIfOrderAlreadyExists(newOrder *order) int {
	indexOfOrder := -1

	for i := range c.activeOrders {
		if c.activeOrders[i].Id == newOrder.Id {
			return i
		}
	}

	return indexOfOrder
}

// This function should only be called from handleOrderUpdateFromStream because of mutex synchronizations
func (c *client) updateExistingOrder(newOrder *order, existingOrderIndex int, isRemovalRequired bool) {
	c.activeOrders[existingOrderIndex] = newOrder

	if isRemovalRequired {
		c.activeOrders[existingOrderIndex] = c.activeOrders[len(c.activeOrders)-1]
		c.activeOrders = c.activeOrders[:len(c.activeOrders)-1]
	}

}

// This function should only be called from handleOrderUpdateFromStream because of mutex synchronizations
func (c *client) checkIfOrderNeedsToBeRemoved(newOrder *order) bool {
	if newOrder.RemainingSize == 0 || newOrder.Type == "market" || newOrder.Status == "closed" {
		return true
	}
	return false
}

// This function should only be called from handleOrderUpdateFromStream because of mutex synchronizations
func (c *client) insertNewOrder(newOrder *order) {
	c.activeOrders = append(c.activeOrders, newOrder)
}

///////////////////////// END --> STREAM ORDER FUNCTIONALITIES /////////////////////////
