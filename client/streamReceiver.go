package client

import (
	"encoding/json"
	"github.com/adi1382/ftx-mirror-bot/websocket"
)

func (c *Client) receiveStreamingData() {
	for {
		msg := <-c.userStream
		c.sendMessageToSubscriptions(msg) //This is used only for host accounts
		if c.checkQuitStream(msg) {
			return
		}

		wsResponse := new(websocket.Response)
		err := json.Unmarshal(msg, wsResponse)
		c.unhandledError(err)

		if c.checkIfWSDataNil(wsResponse.Data) {
			continue
		}

		rawData, err := json.Marshal(wsResponse.Data)
		c.unhandledError(err)

		c.handleWebSocketData(rawData, wsResponse.Channel)
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

func (c *Client) checkIfWSDataNil(data interface{}) bool {
	if data == nil {
		return true
	}
	return false
}

func (c *Client) unhandledError(err error) {
	if err != nil {
		panic(err)
	}
}

func (c *Client) handleWebSocketData(data []byte, channel string) {
	switch channel {
	case "orders":
		newOrderUpdate := new(order)
		err := json.Unmarshal(data, newOrderUpdate)
		c.unhandledError(err)
		c.handleOrderUpdateFromStream(newOrderUpdate)
	case "fills":
		newFillUpdate := new(websocket.FillsData)
		err := json.Unmarshal(data, newFillUpdate)
		c.unhandledError(err)
		c.handleFillUpdateFromStream(newFillUpdate)
	}
}
