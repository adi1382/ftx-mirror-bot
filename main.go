package main

import (
	"github.com/adi1382/ftx-mirror-bot/client"
	"go.uber.org/atomic"
)

var isRestartRequired *atomic.Bool

func init() {
	isRestartRequired = atomic.NewBool(false)
}

func main() {
	hostClient := client.NewClient("kqAyKxRHgQreYe4iNLB7qnpSp1zQsjQP2ePFUDjq", "PhqPf5qpoCp7aFjYC4Ua5ZJTAHuBP20P0TwyZvOX", isRestartRequired)

	hostClient.Initialize()

	//n := 0

	//go func() {
	//	for {
	//		fmt.Printf("\n\nActive Positions: %v\n", hostClient.ActivePositions())
	//		time.Sleep(time.Second * 5)
	//	}
	//}()
	//
	//go func() {
	//	for {
	//		fmt.Printf("\n\nActive Orders: %v\n", hostClient.ActiveOrders())
	//		time.Sleep(time.Second * 5)
	//	}
	//}()

	//fmt.Println("$$$$$$$444")

	select {}
}
