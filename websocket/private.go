package websocket

import (
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

// This function is called as a go routine
func (ws *WSConnection) readFromWSToChannel(chReadWS chan<- []byte) {
	defer func() {
		ws.wg.Done()
	}()
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

	args := map[string]interface{}{
		"key":  ws.config.Key,
		"sign": ws.config.Signature(fmt.Sprintf("%dwebsocket_login", timestamp)),
		"time": timestamp,
	}

	if ws.config.SubAccount.IsSubAccount {
		args["subaccount"] = ws.config.SubAccount.Name
	}

	return &wsAuthorizationMessage{"login", args}
}

// This function is called as a go routine
func (ws *WSConnection) pingPong() {
	defer func() {
		ws.wg.Done()
	}()
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
