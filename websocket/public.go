package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WSConnection struct {
	Conn             *websocket.Conn
	wsWriteLock      sync.Mutex
	key              string
	secret           string
	subRoutineCloser chan int
	wg               *sync.WaitGroup
}

func NewSocketConnection(key, secret string, subRoutineCloser chan int, wg *sync.WaitGroup) *WSConnection {
	ws := WSConnection{key: key, secret: secret, subRoutineCloser: subRoutineCloser, wg: wg}
	return &ws
}

func (ws *WSConnection) Connect(chReadWS chan<- []byte) {
	conn, err := connect("ftx.com")
	if err != nil {
		ws.websocketError(err)
	}
	ws.Conn = conn

	ws.wg.Add(2)
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
