package websocket

import "github.com/adi1382/ftx-mirror-bot/optional"

//Response format

//channel
//market
//type: The type of message
////error: Occurs when there is an error. When type is error, there will also be a code and msg field. code takes on the values of standard HTTP error codes.
////subscribed: Indicates a successful subscription to the channel and market.
////unsubscribed: Indicates a successful unsubscription to the channel and market.
////info: Used to convey information to the user. Is accompanied by a code and msg field.
////When our servers restart, you may see an info message with code 20001. If you do, please reconnect.
////partial: Contains a snapshot of current market data. The data snapshot can be found in the accompanying data field.
////update: Contains an update about current market data. The update can be found in the accompanying data field.
//code (optional)
//msg (optional)
//data(optional)

type Response struct {
	Channel        string          `json:"channel"`
	Market         optional.String `json:"market"`
	TypeOfResponse string          `json:"type"`
	Code           optional.Int64  `json:"code"`
	Msg            optional.String `json:"msg"`
	Data           interface{}     `json:"data"`
}
