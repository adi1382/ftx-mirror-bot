package client

import "github.com/adi1382/ftx-mirror-bot/constants"

func (c *client) Initialize() {
	c.wsConnection.Connect(c.userStream)

	c.wsConnection.AuthenticateWebsocketConnection()
	c.wsConnection.SubscribeToPrivateStreams()
	c.checkIfStreamsAreSuccessfullySubscribed([]string{"fills", "orders"}, constants.TimeoutToCheckForSubscriptions)
	if !c.runningStatus() {
		return
	}

	c.initializeAccountInfoAndPositions()
	c.initializeOrders()

	c.wg.Add(1)
	go c.receiveStreamingData()
}

func (c *client) ActiveOrders() []order {
	c.openOrdersLock.Lock()
	defer c.openOrdersLock.Unlock()

	openOrders := make([]order, 0, 5)
	for i := range c.openOrders {
		openOrders = append(openOrders, *c.openOrders[i])
	}
	return openOrders
}

func (c *client) ActivePositions() []position {
	c.openPositionsLock.Lock()
	defer c.openPositionsLock.Unlock()

	openPositions := make([]position, 0, 5)
	for i := range c.openPositions {
		openPositions = append(openPositions, *c.openPositions[i])
	}
	return openPositions
}

func (c *client) FetchLeverage() float64 {
	return c.leverage.Load()
}

func (c *client) FetchTotalCollateral() float64 {
	return c.totalCollateral.Load()
}

// SubscribeToClientStream is only called for host account
func (c *client) SubscribeToClientStream(ch chan []byte) {
	c.subscriptionsToUserStreamLock.Lock()
	c.subscriptionsToUserStream = append(c.subscriptionsToUserStream, ch)
	c.subscriptionsToUserStreamLock.Unlock()
}

// UpdateSymbolInfoViaRest is only called for host account
func (c *client) UpdateSymbolInfoViaRest() {
	c.symbolInfoLock.Lock()
	defer c.symbolInfoLock.Unlock()

	c.symbolsInfo = make(map[string]symbolInfo, 1000)
	c.updateSymbolInfoForFutures()
	c.updateSymbolInfoForSpot()
}

// FetchSymbolInformation is only called for host account
func (c *client) FetchSymbolInformation() map[string]symbolInfo {
	c.symbolInfoLock.Lock()
	defer c.symbolInfoLock.Unlock()

	symbolInformation := make(map[string]symbolInfo, 1000)
	for k, v := range c.symbolsInfo {
		symbolInformation[k] = v
	}
	return symbolInformation
}

// SetSymbolInformation is only called for sub account
func (c *client) SetSymbolInformation(symbolInformation map[string]symbolInfo) {
	c.symbolInfoLock.Lock()
	defer c.symbolInfoLock.Unlock()

	c.symbolsInfo = make(map[string]symbolInfo, 1000)

	for k, v := range symbolInformation {
		c.symbolsInfo[k] = v
	}
}
