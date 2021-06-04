package main

import (
	"fmt"
	ws "github.com/adi1382/ftx-mirror-bot/websocket"
	"go.uber.org/atomic"
)

var isRestartRequired *atomic.Bool

func init() {
	isRestartRequired = atomic.NewBool(false)
}

func main() {
	conn := ws.NewSocketConnection("kqAyKxRHgQreYe4iNLB7qnpSp1zQsjQP2ePFUDjq", "PhqPf5qpoCp7aFjYC4Ua5ZJTAHuBP20P0TwyZvOX", isRestartRequired)
	ch := conn.Connect()
	conn.AuthenticateWebsocketConnection()
	conn.SubscribeToPrivateStreams()

	go func() {
		for {
			msg := <-ch
			fmt.Println(string(msg))
		}
	}()

	select {}

}
