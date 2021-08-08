package mirror

import (
	"sync"

	"github.com/adi1382/ftx-mirror-bot/client"
)

func NewMirrorInstance(subRoutineCloser chan int, wg *sync.WaitGroup) *Mirror {
	m := &Mirror{
		wg:               wg,
		subRoutineCloser: subRoutineCloser,
	}
	m.subClients = make([]*client.Sub, 0, 5)
	return m
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

func (m *Mirror) Initialize() {
	m.hostClient.Initialize()
	for i := range m.subClients {
		m.subClients[i].Initialize()
	}
}

func (m *Mirror) StartMirroring() {
	for i := range m.subClients {
		go m.subClients[i].StartMirroring()
	}
}
