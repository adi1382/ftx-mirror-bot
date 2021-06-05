package client

import (
	"encoding/json"
	"fmt"
	"github.com/adi1382/ftx-mirror-bot/websocket"
	"time"
)

func (c *Client) receiveStreamingData() {
	for {
		msg := <-c.userStream
		fmt.Println(string(msg))
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

func (c *Client) pushChannelToTheTopOfUserStream(tempChannel chan []byte) {
	if len(c.userStream) > 0 {
		for {
			tempChannel <- <-c.userStream
			if len(c.userStream) == 0 {
				break
			}
		}
	}

	n := len(tempChannel)
	for {
		c.userStream <- <-tempChannel
		if len(tempChannel) == 0 {
			if len(c.userStream) == n {
				return
			} else {
				c.restart()
				return
			}
		}
	}
}

func (c *Client) getMessagesFromChannelWithoutModifyingUserStream(minChannelLength int) []string {
	startTimeUnix := time.Now().Unix()
	rawMessages := make([][]byte, 0, 2)
	messages := make([]string, 0, 2)
	tempChannel := make(chan []byte, 100)

	for {
		if len(c.userStream) >= minChannelLength {
			for {
				tempChannel <- <-c.userStream
				if len(c.userStream) == 0 {
					break
				}
			}

			for {
				rawMessages = append(rawMessages, <-tempChannel)
				if len(tempChannel) == 0 {
					break
				}
			}

			for i := range rawMessages {
				tempChannel <- rawMessages[i]
			}

			for i := range rawMessages {
				messages = append(messages, string(rawMessages[i]))
			}

			c.pushChannelToTheTopOfUserStream(tempChannel)

			return messages
		}

		if time.Now().Unix()-startTimeUnix > 15 {
			fmt.Println("Didn't received subscribed message in 15 seconds. Trying to restart...")
			c.restart()
			return nil
		}
	}
}

func (c *Client) checkIfStreamsAreSuccessfullySubscribed(channels ...string) {
	noOfChannelsToCheck := len(channels)
	noOfChannelsSubscribed := 0
	messagesToFind := make([]string, 0, 2)
	startTime := time.Now().Unix()

	for i := range channels {
		messagesToFind = append(messagesToFind, fmt.Sprintf(`{"type": "subscribed", "channel": "%s"}`, channels[i]))
	}

	for {
		messages := c.getMessagesFromChannelWithoutModifyingUserStream(noOfChannelsToCheck)

		for i := range messagesToFind {
			for j := range messages {
				if messagesToFind[i] == messages[j] {
					noOfChannelsSubscribed++
					break
				}
			}
		}

		if noOfChannelsSubscribed == noOfChannelsToCheck {
			return
		} else {
			noOfChannelsSubscribed = 0
		}

		if time.Now().Unix()-startTime > 15 {
			c.restart()
			return
		}
		time.Sleep(time.Millisecond)
	}
}

func (c *Client) checkIfWSDataNil(data interface{}) bool {
	if data == nil {
		fmt.Println("Data Nil")
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
		fmt.Println("New Order Detected!")
		c.handleOrderUpdateFromStream(newOrderUpdate)
	case "fills":
		newFillUpdate := new(websocket.FillsData)
		err := json.Unmarshal(data, newFillUpdate)
		c.unhandledError(err)
		fmt.Println("New Fill Detected!")
		c.handleFillUpdateFromStream(newFillUpdate)
	}
}
