package main

import (
	"fmt"
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
	fmt.Println("$$$$$$$444")

	select {}
}
