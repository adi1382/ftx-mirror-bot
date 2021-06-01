package websocket

import (
	"github.com/gorilla/websocket"
)

type WSConnection struct {
	Conn   *websocket.Conn
	key    string
	secret string
}

func NewSocketConnection(key, secret string) WSConnection {
	conn, err := connect("ftx.com")
	if err != nil {
		panic(err)
	}
	ws := WSConnection{conn, key, secret}
	go ws.pingPong()
	return ws
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
