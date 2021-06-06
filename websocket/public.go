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
	subRoutineCloser  chan int
}

func NewSocketConnection(key, secret string, restartCounter *atomic.Bool, subRoutineCloser chan int) WSConnection {
	ws := WSConnection{key: key, secret: secret, isRestartRequired: restartCounter, subRoutineCloser: subRoutineCloser}
	return ws
}

func (ws *WSConnection) Connect(chReadWS chan<- []byte) {
	conn, err := connect("ftx.com")
	if err != nil {
		ws.websocketError(err)
	}
	ws.Conn = conn

	go ws.pingPong()
	go ws.readFromWSToChannel(chReadWS)
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
