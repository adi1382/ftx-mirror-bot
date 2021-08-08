//These functions verify if a channel is subscribed on the websocket connection

package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/adi1382/ftx-mirror-bot/websocket"
)

func (c *client) pushChannelToTheTopOfUserStream(tempChannel chan []byte) error {
	c.flushUserStreamToChannel(tempChannel)

	n := len(tempChannel)

	for len(tempChannel) > 0 {
		c.userStream <- <-tempChannel
	}

	if len(c.userStream) != n {
		return fmt.Errorf("new message got into userstream while flushing tempChannel")
	}

	return nil
}

func (c *client) flushUserStreamToChannel(tempChannel chan []byte) {
	for len(c.userStream) > 0 {
		tempChannel <- <-c.userStream
	}
}

func (c *client) fetchMessagesFromChannel(tempChannel chan []byte) [][]byte {
	byteMessages := make([][]byte, 0, 2)

	for len(tempChannel) > 0 {
		byteMessages = append(byteMessages, <-tempChannel)
	}

	for i := range byteMessages {
		tempChannel <- byteMessages[i]
	}

	return byteMessages
}

func (c *client) fetchMessagesFromUserStreamWithoutModifyingUserStream(
	minChannelLength int, timeout time.Duration) ([][]byte, error) {
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
			return nil, fmt.Errorf("didn't received %d message(s) in timeout", minChannelLength)
		}
	}

	c.flushUserStreamToChannel(tempChannel)
	messages := c.fetchMessagesFromChannel(tempChannel)
	if err := c.pushChannelToTheTopOfUserStream(tempChannel); err != nil {
		return nil, err
	}

	return messages, nil
}

func (c *client) subscribedChannels(messages [][]byte) []string {
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

func (c *client) fetchSubscribedChannels(minChannelLength int, timeout time.Duration) ([]string, error) {
	if messages, err := c.fetchMessagesFromUserStreamWithoutModifyingUserStream(minChannelLength, timeout); err != nil {
		return nil, err
	} else {
		return c.subscribedChannels(messages), nil
	}
}

func (c *client) numberOfMatchingChannels(channelsToSubscribe, subscribedChannels []string) int {
	var noOfChannelsSubscribed int
	for i := range channelsToSubscribe {
		for j := range subscribedChannels {
			if channelsToSubscribe[i] == subscribedChannels[j] {
				noOfChannelsSubscribed++
			}
		}
	}
	return noOfChannelsSubscribed
}

func (c *client) checkIfStreamsAreSuccessfullySubscribed(channelsToSubscribe []string, timeout time.Duration) error {
	noOfChannelsToCheck := len(channelsToSubscribe)
	startTime := time.Now().Unix()

	for {

		subscribedChannels, err := c.fetchSubscribedChannels(noOfChannelsToCheck, timeout)
		if err != nil {
			return err
		}
		noOfChannelsSubscribed := c.numberOfMatchingChannels(channelsToSubscribe, subscribedChannels)

		if noOfChannelsSubscribed == noOfChannelsToCheck {
			return nil
		} else {
			noOfChannelsSubscribed = 0
		}

		if time.Now().Unix()-startTime > int64(timeout/time.Second) {
			return fmt.Errorf("Unable to verify subscriptions for channels %v\n", channelsToSubscribe)
		}
		time.Sleep(time.Millisecond)
	}
}
