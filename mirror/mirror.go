package mirror

import (
	"github.com/adi1382/ftx-mirror-bot/client"
	"sync"
)

func NewMirrorInstance(wg *sync.WaitGroup, subRoutineCloser chan int) *Mirror {
	return &Mirror{
		wg:               wg,
		subRoutineCloser: subRoutineCloser,
	}
}

type Mirror struct {
	hostClient       *client.Host
	subClients       []*client.Sub
	wg               *sync.WaitGroup
	subRoutineCloser chan int
}

func (m *Mirror) SetHostClient(host *client.Host) {
	m.hostClient = host
}

func (m *Mirror) AddSubClient(sub *client.Sub) {
	m.subClients = append(m.subClients, sub)
}
