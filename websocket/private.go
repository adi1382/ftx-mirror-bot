package websocket

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"time"

	"github.com/adi1382/ftx-mirror-bot/constants"
	"github.com/gorilla/websocket"
)

type wsMessage struct {
	Op      string `json:"op"`
	Channel string `json:"channel"`
	Market  string `json:"market"`
}

type wsAuthorizationMessage struct {
	Op   string                 `json:"op"`
	Args map[string]interface{} `json:"args"`
}

func connect(host string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "wss", Host: host, Path: "/ws/"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	return conn, err
}

func (ws *WSConnection) readFromWSToChannel(chReadWS chan<- []byte) {
	for {
		_, message, err := ws.Conn.ReadMessage()

		if err != nil {
			ws.websocketError(err)
			chReadWS <- []byte("quit")
			return
		}
		chReadWS <- message
	}
}

func (ws *WSConnection) subscribeToPrivateChannels(channels []string) {
	for i := range channels {
		ws.wsWriteLock.Lock()
		err := ws.Conn.WriteJSON(&wsMessage{
			Op:      "subscribe",
			Channel: channels[i]})
		ws.wsWriteLock.Unlock()
		if err != nil {
			ws.websocketError(err)
			return
		}
	}
}

func (ws *WSConnection) getAuthMessage() *wsAuthorizationMessage {
	timestamp := time.Now().UTC().UnixNano() / int64(time.Millisecond)

	sig := hmac.New(sha256.New, []byte(ws.secret))
	sig.Write([]byte(fmt.Sprintf("%dwebsocket_login", timestamp)))
	args := map[string]interface{}{
		"key":  ws.key,
		"sign": hex.EncodeToString(sig.Sum(nil)),
		"time": timestamp,
	}

	return &wsAuthorizationMessage{"login", args}
}

func (ws *WSConnection) pingPong() {
	ticker := time.NewTicker(constants.PingPeriod)
	defer ticker.Stop()

	_ = ws.Conn.SetReadDeadline(time.Now().Add(constants.PongWait))
	ws.Conn.SetPongHandler(func(string) error { err := ws.Conn.SetReadDeadline(time.Now().Add(constants.PongWait)); return err })

	for {
		select {
		case <-ticker.C:
			ws.wsWriteLock.Lock()
			err := ws.Conn.WriteMessage(websocket.PingMessage, []byte(`{"op": "ping"}`))
			ws.wsWriteLock.Unlock()

			if err != nil {
				ws.websocketError(err)
				return
			}
		case c := <-ws.subRoutineCloser:
			ws.subRoutineCloser <- c + 1
			_ = ws.Conn.Close()
			return

		}
	}
}

func (ws *WSConnection) websocketError(err error) {
	fmt.Println(err)
}

//func (ws *WSConnection) closeOnRestart() {
//	for {
//		if ws.isRestartRequired.Load() {
//			err := ws.Conn.Close()
//
//			if err != nil {
//				return
//			}
//		}
//
//		time.Sleep(time.Millisecond)
//	}
//}
