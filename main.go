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
	conn.Connect()

	go func() {
		for {
			_, msg, _ := conn.Conn.ReadMessage()
			fmt.Println("K")
			fmt.Println(string(msg))
		}
	}()

	conn.AuthenticateWebsocketConnection()
	conn.SubscribeToPrivateStreams()
	select {}

}
