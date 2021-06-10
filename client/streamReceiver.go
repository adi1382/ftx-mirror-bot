package client

import (
	"encoding/json"

	"github.com/adi1382/ftx-mirror-bot/websocket"
)

// This function is called as a sub routine
func (c *client) receiveStreamingData() {
	defer func() {
		c.wg.Done()
	}()
	for {
		msg := <-c.userStream
		//TODO: REMOVE IT FROM HERE
		//c.sendNewOrderUpdateToSubscriptions(msg) //This is used only for host accounts
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

func (c *client) handleWebSocketData(data []byte, channel string) {
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

func (c *client) checkIfWSDataNil(data interface{}) bool {
	if data == nil {
		return true
	}
	return false
}

func (c *client) checkQuitStream(msg []byte) bool {
	if string(msg) == "quit" {
		return true
	}
	return false
}

func (c *client) unhandledError(err error) {
	if err != nil {
		panic(err)
	}
}
