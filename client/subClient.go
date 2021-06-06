package client

import (
	"sync"
)

func NewSubClient(
	apiKey, apiSecret string,
	subRoutineCloser chan int,
	wg *sync.WaitGroup,
	host *Host) *Sub {
	c := Sub{
		client: newClient(apiKey, apiSecret, subRoutineCloser, wg), hostClient: host,
	}
	c.hostMessageUpdates = make(chan []byte, 100)
	return &c
}

type Sub struct {
	client             *client
	hostClient         *Host
	hostMessageUpdates chan []byte
	subRoutineCloser   chan int
	wg                 *sync.WaitGroup
}

func (s *Sub) Initialize() {
	s.client.SetSymbolInformation(s.hostClient.FetchSymbolInformation())
	s.client.Initialize()
	s.hostClient.SubscribeToHostUpdates(s.hostMessageUpdates)
}
