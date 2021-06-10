package client

import (
	"fmt"
	"time"
)

func (s *Sub) StartMirroring() {
	calibrationTicker := time.NewTicker(time.Duration(s.calibrationDuration) * time.Second)
	defer calibrationTicker.Stop()

	s.calibrate()

	for {
		select {
		case hostNewOrderUpdate := <-s.hostNewOrderUpdates:
			s.processNewOrderUpdate(hostNewOrderUpdate)
		case hostCancelOrderUpdate := <-s.hostExistingOrderUpdates:
			s.processExistingOrderUpdate(hostCancelOrderUpdate)
		case <-calibrationTicker.C:
			s.calibrate()
		case <-s.subRoutineCloser:
			s.subRoutineCloser <- 1
			return
		}
	}
}

func (s *Sub) processNewOrderUpdate(hostNewOrderUpdate *order) {
	fmt.Println("$$$$$$$$$$$ NEW ORDER")
	fmt.Println(*hostNewOrderUpdate)
	fmt.Println("$$$$$$$$$$$ NEW ORDER")
}

func (s *Sub) processExistingOrderUpdate(hostExistingOrderUpdate *order) {
	fmt.Println("$$$$$$$$$$$ Existing ORDER")
	fmt.Println(*hostExistingOrderUpdate)
	fmt.Println("$$$$$$$$$$$ Existing ORDER")
}
