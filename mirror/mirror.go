package mirror

import (
	"sync"

	"github.com/adi1382/ftx-mirror-bot/host"
	"github.com/adi1382/ftx-mirror-bot/sub"
)

func NewMirrorInstance(wg *sync.WaitGroup, subRoutineCloser chan int) *Mirror {
	return &Mirror{
		wg:               wg,
		subRoutineCloser: subRoutineCloser,
	}
}

type Mirror struct {
	hostClient       *host.Host
	subClients       []*sub.Sub
	wg               *sync.WaitGroup
	subRoutineCloser chan int
}

func (m *Mirror) SetHostClient(host *host.Host) {
	m.hostClient = host
}

func (m *Mirror) AddSubClient(sub *sub.Sub) {
	m.subClients = append(m.subClients, sub)
}
