package Binance

import (
	ws "github.com/gorilla/websocket"
)

type Spot_Websockets struct{}

type Spot_Websocket struct {
	Websocket *Websocket
	Conn      *ws.Conn
	// Host server's URL
	BaseURL string

	pendingRequests map[string]chan []byte
}

func (spot_ws *Spot_Websockets) CreateSocket(streams []string, isCombined bool) (*Spot_Websocket, error) {
	baseURL := SPOT_Constants.Websocket.URLs[0]

	socket, err := CreateSocket(baseURL, streams, isCombined)
	if err != nil {
		return nil, err
	}

	ws := &Spot_Websocket{
		Websocket:       socket,
		Conn:            socket.Conn,
		BaseURL:         baseURL,
		pendingRequests: make(map[string]chan []byte),
	}

	return ws, nil
}

// type Websocket_Error_Response struct {
// 	Code int    `json:"code"`
// 	Msg  string `json:"msg"`
// }

// type Websocket_Subscribe_Response struct {
// 	Id int64 `json:"id"`
// }

// type Websocket_Unsubscribe_Response struct {
// 	Id int64 `json:"id"`
// }

// type Websocket_ListingSubscriptions_Response struct {
// 	Id     int64    `json:"id"`
// 	Result []string `json:"result"`
// }

// type Websocket_SettingProperties_Response struct {
// 	Id int64 `json:"id"`
// }

// type Websocket_RetrievingProperties_Response struct {
// 	Id     int64 `json:"id"`
// 	Result bool  `json:"result"`
// }

// func generateMessageID() string {
// 	return uuid.New().String()
// }
