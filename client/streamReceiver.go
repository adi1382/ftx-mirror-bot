package client

func (c *Client) receiveStreamingData() {
	for {
		msg := <-c.userStream
		c.sendMessageToSubscriptions(msg)
		if c.checkQuitStream(msg) {
			return
		}

		//eventName := c.getEventName(msg)
		//
		//switch eventName {
		//case "ORDER_TRADE_UPDATE":
		//	c.processOrderTradeUpdate(msg)
		//case "ACCOUNT_UPDATE":
		//	c.processAccountUpdate(msg)
		//case "ACCOUNT_CONFIG_UPDATE":
		//	c.processAccountConfigUpdate(msg)
		//case "listenKeyExpired":
		//	c.Restart()
		//default:
		//	continue
		//}
	}
}

func (c *Client) sendMessageToSubscriptions(msg []byte) {
	if len(c.subscriptionsToUserStream) > 0 {
		c.subscriptionsToUserStreamLock.Lock()
		for i := range c.subscriptionsToUserStream {
			c.subscriptionsToUserStream[i] <- msg
		}
		c.subscriptionsToUserStreamLock.Unlock()
	}
}

func (c *Client) checkQuitStream(msg []byte) bool {
	if string(msg) == "quit" {
		c.isRestartRequired.Store(true)
		return true
	}
	return false
}
