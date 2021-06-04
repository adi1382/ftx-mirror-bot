package websocket

import (
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

func (ws *WSConnection) Connect() chan []byte {
	conn, err := connect("ftx.com")
	if err != nil {
		ws.websocketError(err)
	}
	go ws.pingPong()
	ws.Conn = conn

	chReadWS := make(chan []byte, 100)
	go ws.readFromWSToChannel(chReadWS)
	go ws.closeOnRestart()

	return chReadWS
}

func (ws *WSConnection) AuthenticateWebsocketConnection() {
	err := ws.Conn.WriteJSON(ws.getAuthMessage())
	if err != nil {
		ws.websocketError(err)
	}
}

func (ws *WSConnection) SubscribeToPrivateStreams() {
	ws.subscribeToPrivateChannels([]string{"fills", "orders"})
}
