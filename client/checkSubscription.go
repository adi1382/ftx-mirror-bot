package client

import (
	"encoding/json"
	"fmt"
	"github.com/adi1382/ftx-mirror-bot/websocket"
	"time"
)

func (c *Client) pushChannelToTheTopOfUserStream(tempChannel chan []byte) {
	c.dumpUserStreamToChannel(tempChannel)

	n := len(tempChannel)

	for len(tempChannel) > 0 {
		c.userStream <- <-tempChannel
	}

	if len(c.userStream) != n {
		c.restart()
		return
	}
}

func (c *Client) dumpUserStreamToChannel(tempChannel chan []byte) {
	for len(c.userStream) > 0 {
		tempChannel <- <-c.userStream
	}
}

func (c *Client) fetchMessagesFromChannel(tempChannel chan []byte) [][]byte {
	byteMessages := make([][]byte, 0, 2)

	for len(tempChannel) > 0 {
		byteMessages = append(byteMessages, <-tempChannel)
	}

	for i := range byteMessages {
		tempChannel <- byteMessages[i]
	}

	return byteMessages
}

func (c *Client) fetchMessagesFromUserStreamWithoutModifyingUserStream(minChannelLength int, timeout time.Duration) [][]byte {
	tempChannel := make(chan []byte, 100)
	timer := time.NewTimer(timeout)
	defer timer.Stop()

L:
	for {
		select {
		case msg := <-c.userStream:
			tempChannel <- msg
			if len(tempChannel) < minChannelLength {
				continue
			} else {
				break L
			}
		case <-timer.C:
			fmt.Printf("Didn't received %d message(s) in 15 seconds. Trying to restart...\n", minChannelLength)
			c.restart()
			return nil
		}
	}

	c.dumpUserStreamToChannel(tempChannel)
	messages := c.fetchMessagesFromChannel(tempChannel)
	c.pushChannelToTheTopOfUserStream(tempChannel)

	return messages
}

func (c *Client) fetchSubscribedChannels(messages [][]byte) []string {
	wsMessage := new(websocket.Response)
	subscribed := make([]string, 0, 2)

	for i := range messages {
		err := json.Unmarshal(messages[i], wsMessage)
		c.unhandledError(err)
		if wsMessage.TypeOfResponse == "subscribed" {
			subscribed = append(subscribed, wsMessage.Channel)
		}
	}

	return subscribed
}

func (c *Client) checkIfStreamsAreSuccessfullySubscribed(channelsToSubscribe []string, timeout time.Duration) {
	noOfChannelsToCheck := len(channelsToSubscribe)
	noOfChannelsSubscribed := 0
	startTime := time.Now().Unix()

	for {
		byteMessages := c.fetchMessagesFromUserStreamWithoutModifyingUserStream(noOfChannelsToCheck, timeout)
		if byteMessages == nil {
			return
		}

		subscribedChannels := c.fetchSubscribedChannels(byteMessages)

		for i := range channelsToSubscribe {
			for j := range subscribedChannels {
				if channelsToSubscribe[i] == subscribedChannels[j] {
					noOfChannelsSubscribed++
				}
			}
		}

		if noOfChannelsSubscribed == noOfChannelsToCheck {
			return
		} else {
			noOfChannelsSubscribed = 0
		}

		if time.Now().Unix()-startTime > int64(timeout/time.Second) {
			fmt.Printf("Unable to verify subscriptions for channels %v\n", channelsToSubscribe)
			c.restart()
			return
		}
		time.Sleep(time.Millisecond)
	}
}
