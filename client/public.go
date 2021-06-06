package client

import "github.com/adi1382/ftx-mirror-bot/constants"

func (c *Client) Initialize() {
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

func (c *Client) ActiveOrders() []Order {
	c.openOrdersLock.Lock()
	defer c.openOrdersLock.Unlock()

	openOrders := make([]Order, 0, 5)
	for i := range c.openOrders {
		openOrders = append(openOrders, *c.openOrders[i])
	}
	return openOrders
}

func (c *Client) ActivePositions() []Position {
	c.openPositionsLock.Lock()
	defer c.openPositionsLock.Unlock()

	openPositions := make([]Position, 0, 5)
	for i := range c.openPositions {
		openPositions = append(openPositions, *c.openPositions[i])
	}
	return openPositions
}

func (c *Client) FetchLeverage() float64 {
	return c.leverage.Load()
}

func (c *Client) FetchTotalCollateral() float64 {
	return c.totalCollateral.Load()
}

// SubscribeToClientStream is only called for host account
func (c *Client) SubscribeToClientStream(ch chan []byte) {
	c.subscriptionsToUserStreamLock.Lock()
	c.subscriptionsToUserStream = append(c.subscriptionsToUserStream, ch)
	c.subscriptionsToUserStreamLock.Unlock()
}

// UpdateSymbolInfoViaRest is only called for host account
func (c *Client) UpdateSymbolInfoViaRest() {
	c.symbolInfoLock.Lock()
	defer c.symbolInfoLock.Unlock()

	c.symbolsInfo = make(map[string]SymbolInfo, 1000)
	c.updateSymbolInfoForFutures()
	c.updateSymbolInfoForSpot()
}

// FetchSymbolInformation is only called for host account
func (c *Client) FetchSymbolInformation() map[string]SymbolInfo {
	c.symbolInfoLock.Lock()
	defer c.symbolInfoLock.Unlock()

	symbolInformation := make(map[string]SymbolInfo, 1000)
	for k, v := range c.symbolsInfo {
		symbolInformation[k] = v
	}
	return symbolInformation
}

// SetSymbolInformation is only called for sub account
func (c *Client) SetSymbolInformation(symbolInformation map[string]SymbolInfo) {
	c.symbolInfoLock.Lock()
	defer c.symbolInfoLock.Unlock()

	c.symbolsInfo = make(map[string]SymbolInfo, 1000)

	for k, v := range symbolInformation {
		c.symbolsInfo[k] = v
	}
}
