package client

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/adi1382/ftx-mirror-bot/websocket"
)

func (s *Sub) StartMirroring() {
	calibrationTicker := time.NewTicker(time.Duration(s.calibrationDuration) * time.Second)
	defer calibrationTicker.Stop()

	s.calibrate()

	for {
		select {
		case hostSocketMessage := <-s.hostMessageUpdates:
			if s.client.checkQuitStream(hostSocketMessage) {
				return
			}
			fmt.Println("Mirrorring message")
			s.processHostSocketUpdate(hostSocketMessage)
		case <-calibrationTicker.C:
			s.calibrate()
		}
	}
}

func (s *Sub) processHostSocketUpdate(hostSocketMessage []byte) {
	wsResponse := new(websocket.Response)
	err := json.Unmarshal(hostSocketMessage, wsResponse)
	s.client.unhandledError(err)

	if wsResponse.Channel != "order" {
		return
	}

}
