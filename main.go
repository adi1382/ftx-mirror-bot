package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	ws "github.com/adi1382/ftx-mirror-bot/websocket"
	"github.com/gorilla/websocket"
	"net/url"
	"time"
)

func main() {
	conn := ws.NewSocketConnection("kqAyKxRHgQreYe4iNLB7qnpSp1zQsjQP2ePFUDjq", "PhqPf5qpoCp7aFjYC4Ua5ZJTAHuBP20P0TwyZvOX")
	go func() {
		for {
			_, msg, _ := conn.Conn.ReadMessage()
			fmt.Println("K")
			fmt.Println(string(msg))
		}
	}()

	conn.AuthenticateWebsocketConnection()
	conn.SubscribeToPrivateStreams()

	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//
	//ch := make(chan realtime.Response)
	//realtime.ConnectForPrivate(ctx, ch, "kqAyKxRHgQreYe4iNLB7qnpSp1zQsjQP2ePFUDjq", "PhqPf5qpoCp7aFjYC4Ua5ZJTAHuBP20P0TwyZvOX", []string{"orders", "fills"}, nil)
	//
	//for {
	//	message := <- ch
	//
	//	fmt.Println(message)
	//}

	//conn := ws.NewSocketConnection()
	//
	//go ws.ReadFromWSToChannel(conn)
	//
	//ws.AuthenticateWebsocketConnection(conn, "HdLU4D7OqAaZJy4Nl9fSJBTf1BO3V_yKLdkqHYIS", "YoZS9RJthKyhqgsdeDMED3y9wcwwQioc23zf5xYa")
	select {}
	fmt.Println("conn")

}

func Connect(host string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "wss", Host: host, Path: "/ws/"}
	fmt.Println(u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	//if err != nil{
	//	panic(err)
	//}

	return conn, err
}

//func ReadFromWSToChannel(c *websocket.Conn, chRead chan<- []byte, RestartCounter *int32) {
//
//	for {
//		_, message, err := c.ReadMessage()
//		if err != nil {
//			panic(err)
//		}
//
//		fmt.Println(string(message))
//
//		if atomic.LoadInt32(RestartCounter) > 0 {
//			fmt.Println("Read Socket Closed")
//			break
//		}
//
//		chRead <- message
//	}
//}

// GetAuthMessage func ReadFromWSToChannel(c *websocket.Conn, chRead chan<- []byte, RestartCounter *int32) {
//	var message []byte
//	var err error
//	var readWSLocal int32
//L:
//	for {
//		atomic.StoreInt32(&readWSLocal, 0)
//		go func() {
//			atomic.StoreInt32(&readWSLocal, 0)
//			_, message, err = c.ReadMessage()
//			atomic.StoreInt32(&readWSLocal, 1)
//		}()
//
//		for {
//			time.Sleep(time.Nanosecond)
//			if atomic.LoadInt32(&readWSLocal) == 1 {
//				break
//			} else {
//				if atomic.LoadInt32(RestartCounter) > 0 {
//					receiveLogger.Println("ReadFromWSToChannel Closed")
//					InfoLogger.Printf("ReadFromWSToChannel Closed")
//					break L
//				}
//			}
//		}
//
//		if atomic.LoadInt32(RestartCounter) > 0 {
//			receiveLogger.Println("ReadFromWSToChannel Closed")
//			InfoLogger.Printf("ReadFromWSToChannel Closed")
//			break L
//		}
//
//		tools.WebsocketErr(err, RestartCounter)
//		receiveLogger.Println("Length of channel:", len(chRead), "Message:", string(message))
//
//		chRead <- message
//	}
//}
//
//func WriteFromChannelToWS(c *websocket.Conn, chWrite <-chan interface{}, RestartCounter *int32, wg *sync.WaitGroup) {
//L:
//	for {
//		//a := chWrite.(string)
//		time.Sleep(time.Nanosecond)
//
//		select {
//		case message := <-chWrite:
//
//			if atomic.LoadInt32(RestartCounter) > 0 {
//
//				if atomic.LoadInt32(RestartCounter) >= 3 {
//					for len(chWrite) > 0 {
//						<-chWrite
//					}
//					sendLogger.Printf("\n\nWriteFromChannelToWS Closed\n\n")
//					_ = c.Close()
//					InfoLogger.Println(runtime.NumGoroutine())
//					InfoLogger.Printf("\n\nWriteFromChannelToWS Closed\n\n")
//					break L
//				}
//
//				for len(chWrite) > 0 {
//					<-chWrite
//				}
//				continue L
//			}
//
//			message, err := json.Marshal(message)
//			tools.WebsocketErr(err, RestartCounter)
//
//			sendLogger.Println("Length of channel:", len(chWrite), "Message:", string(message.([]byte)))
//			err = c.WriteMessage(websocket.TextMessage, message.([]byte))
//
//			tools.WebsocketErr(err, RestartCounter)
//
//		default:
//			//fmt.Println(atomic.LoadInt32(RestartCounter))
//			if atomic.LoadInt32(RestartCounter) >= 3 {
//				for len(chWrite) > 0 {
//					<-chWrite
//				}
//				sendLogger.Printf("\n\nWriteFromChannelToWS Closed\n\n")
//				_ = c.Close()
//				InfoLogger.Println(runtime.NumGoroutine())
//				InfoLogger.Printf("\n\nWriteFromChannelToWS Closed\n\n")
//				break L
//			}
//		}
//	}
//	wg.Done()
//}
//

type publicWebsocketMessage struct {
	op      string
	channel string
	market  string
}

type privateWebsocketMessage struct {
	op   string
	args map[string]interface{}
}

func getAuthMessage(key string, secret string) privateWebsocketMessage {
	timestamp := time.Now().UTC().UnixNano() / int64(time.Millisecond)

	sig := hmac.New(sha256.New, []byte(secret))
	sig.Write([]byte(fmt.Sprintf("%dwebsocket_login", timestamp)))
	args := map[string]interface{}{
		"key":  key,
		"sign": hex.EncodeToString(sig.Sum(nil)),
		"time": timestamp,
	}

	return privateWebsocketMessage{"login", args}
}
