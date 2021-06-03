package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/atomic"
)

type WSConnection struct {
	Conn              *websocket.Conn
	isRestartRequired *atomic.Bool
	key               string
	secret            string
}

func NewSocketConnection(key, secret string, restartCounter *atomic.Bool) WSConnection {
	ws := WSConnection{key: key, secret: secret, isRestartRequired: restartCounter}
	return ws
}

func (ws *WSConnection) Connect() {
	conn, err := connect("ftx.com")
	if err != nil {
		ws.websocketError(err)
	}
	go ws.pingPong()
	ws.Conn = conn
}

func (ws *WSConnection) AuthenticateWebsocketConnection() {
	err := ws.Conn.WriteJSON(ws.getAuthMessage())
	if err != nil {
		panic(err)
	}
}

func (ws *WSConnection) SubscribeToPrivateStreams() {
	ws.subscribeToPrivateChannels([]string{"fills", "orders"})
}

func ReadFromWSToChannel(c *websocket.Conn) {

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			panic(err)
		}

		fmt.Println("Message Received: ", string(message))
	}
}
