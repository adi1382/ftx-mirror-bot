package sub

import (
	"sync"

	"github.com/adi1382/ftx-mirror-bot/client"
	"github.com/adi1382/ftx-mirror-bot/host"
)

func NewSubClient(
	apiKey, apiSecret string,
	subRoutineCloser chan int,
	wg *sync.WaitGroup,
	host *host.Host) *Sub {
	c := Sub{
		client: client.NewClient(apiKey, apiSecret, subRoutineCloser, wg), hostClient: host,
	}
	c.hostMessageUpdates = make(chan []byte, 100)
	return &c
}

type Sub struct {
	client             *client.Client
	hostClient         *host.Host
	hostMessageUpdates chan []byte
	subRoutineCloser   chan int
	wg                 *sync.WaitGroup
}

func (s *Sub) Initialize() {
	s.client.SetSymbolInformation(s.hostClient.FetchSymbolInformation())
	s.client.Initialize()
	s.hostClient.SubscribeToHostUpdates(s.hostMessageUpdates)
}
