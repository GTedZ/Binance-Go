package Binance

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
)

type Futures_Websockets struct {
	binance *Binance
}

type Futures_Websocket struct {
	Websocket *Websocket
	Conn      *ws.Conn
	// Host server's URL
	BaseURL string

	pendingRequests map[string]chan []byte
}

func (futures_ws *Futures_Websocket) Close() error {
	return futures_ws.Websocket.Close()
}

// Forcefully reconnects the socket
// Also makes it a reconnecting socket if it weren't before
// Useless, but there nonetheless...
func (futures_ws *Futures_Websocket) Reconnect() {
	futures_ws.Websocket.Reconnect()
}

func (futures_ws *Futures_Websocket) SetMessageListener(f func(messageType int, msg []byte)) {
	futures_ws.Websocket.OnMessage = f
}

func (futures_ws *Futures_Websocket) SetPingListener(f func(appData string)) {
	futures_ws.Websocket.OnPing = f
}

func (futures_ws *Futures_Websocket) SetPongListener(f func(appData string)) {
	futures_ws.Websocket.OnPong = f
}

// This is called when socket has been disconnected
// Called when the detected a disconnection and wants to reconnect afterwards
// Usually called right before the 'ReconnectingListener'
func (futures_ws *Futures_Websocket) SetDisconnectListener(f func(code int, text string)) {
	futures_ws.Websocket.OnDisconnect = f
}

// This is called when socket began reconnecting
func (futures_ws *Futures_Websocket) SetReconnectingListener(f func()) {
	futures_ws.Websocket.OnReconnecting = f
}

// This is called when the socket has successfully reconnected after a disconnection
func (futures_ws *Futures_Websocket) SetReconnectListener(f func()) {
	futures_ws.Websocket.OnReconnect = f
}

// This is called when the websocket closes indefinitely
// Meaning when you invoke the 'Close()' method
// Or any other way a websocket is set to never reconnect on a disconnection
func (futures_ws *Futures_Websocket) SetCloseListener(f func(code int, text string)) {
	futures_ws.Websocket.OnClose = f
}

type FuturesWS_PrivateMessage struct {
	Id string `json:"id"`
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type FuturesWS_AggTrade struct {

	// Event type
	Event string `json:"e"`

	// Event time
	EventTime int64 `json:"E"`

	// Symbol
	Symbol string `json:"s"`

	// Aggregate trade ID
	AggTradeId int64 `json:"a"`

	// Price
	Price string `json:"p"`

	// Quantity
	Quantity string `json:"q"`

	// First trade ID
	FirstTradeId int64 `json:"f"`

	// Last trade ID
	LastTradeId int64 `json:"l"`

	// Trade time
	Timestamp int64 `json:"T"`

	// Is the buyer the market maker?
	IsMaker bool `json:"m"`
}

type FuturesWS_AggTrade_Socket struct {
	Handler *Futures_Websocket
}

func (*FuturesWS_AggTrade_Socket) CreateStreamName(symbol string) string {
	return strings.ToLower(symbol) + "@aggTrade"
}

func (socket *FuturesWS_AggTrade_Socket) Subscribe(symbol ...string) (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	for i := range symbol {
		symbol[i] = socket.CreateStreamName(symbol[i])
	}

	return socket.Handler.Subscribe(symbol...)
}

func (socket *FuturesWS_AggTrade_Socket) Unsubscribe(symbol ...string) (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	for i := range symbol {
		symbol[i] = socket.CreateStreamName(symbol[i])
	}
	return socket.Handler.Unsubscribe(symbol...)
}

func (futures_ws *Futures_Websockets) AggTrade(publicOnMessage func(aggTrade *FuturesWS_AggTrade), symbol ...string) (*FuturesWS_AggTrade_Socket, *Error) {
	var newSocket FuturesWS_AggTrade_Socket
	for i := range symbol {
		symbol[i] = newSocket.CreateStreamName(symbol[i])
	}
	socket, err := futures_ws.CreateSocket(symbol, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var aggTrade FuturesWS_AggTrade
		err := json.Unmarshal(msg, &aggTrade)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(&aggTrade)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type FuturesWS_MarkPrice struct {

	// Event type
	Event string `json:"e"`

	// Event time
	EventTime int64 `json:"E"`

	// Symbol
	Symbol string `json:"s"`

	// Mark price
	MarkPrice string `json:"p"`

	// Index price
	IndexPrice string `json:"i"`

	// Estimated Settle Price, only useful in the last hour before the settlement starts
	EstimatedSettlePrice string `json:"P"`

	// Funding rate
	FundingRate string `json:"r"`

	// Next funding time
	NextFundingTime int64 `json:"T"`
}

type FuturesWS_MarkPrice_Params struct {
	Symbol string
	IsFast bool
}

type FuturesWS_MarkPrice_Socket struct {
	Handler *Futures_Websocket
}

func (*FuturesWS_MarkPrice_Socket) CreateStreamName(symbol string, IsFast bool) string {
	streamName := strings.ToLower(symbol) + "@markPrice"
	if IsFast {
		streamName += "@1s"
	}
	return streamName
}

func (socket *FuturesWS_MarkPrice_Socket) Subscribe(params ...FuturesWS_MarkPrice_Params) (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	streams := make([]string, len(params))
	for i := range params {
		streams[i] = socket.CreateStreamName(params[i].Symbol, params[i].IsFast)
	}

	return socket.Handler.Subscribe(streams...)
}

func (socket *FuturesWS_MarkPrice_Socket) Unsubscribe(params ...FuturesWS_MarkPrice_Params) (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	streams := make([]string, len(params))
	for i := range params {
		streams[i] = socket.CreateStreamName(params[i].Symbol, params[i].IsFast)
	}

	return socket.Handler.Unsubscribe(streams...)
}

func (futures_ws *Futures_Websockets) MarkPrice(publicOnMessage func(markPrice *FuturesWS_MarkPrice), params ...FuturesWS_MarkPrice_Params) (*FuturesWS_MarkPrice_Socket, *Error) {
	var newSocket FuturesWS_MarkPrice_Socket

	streams := make([]string, len(params))
	for i := range params {
		streams[i] = newSocket.CreateStreamName(params[i].Symbol, params[i].IsFast)
	}
	socket, err := futures_ws.CreateSocket(streams, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var markPrice FuturesWS_MarkPrice
		err := json.Unmarshal(msg, &markPrice)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(&markPrice)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type FuturesWS_AllMarkPrices_Socket struct {
	Handler *Futures_Websocket
}

func (*FuturesWS_AllMarkPrices_Socket) CreateStreamName(isFast bool) string {
	streamName := "!markPrice@arr"
	if isFast {
		streamName += "@1s"
	}
	return streamName
}

func (socket *FuturesWS_AllMarkPrices_Socket) Subscribe(isFast ...bool) (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	streams := make([]string, len(isFast))
	for i := range isFast {
		streams[i] = socket.CreateStreamName(isFast[i])
	}

	return socket.Handler.Subscribe(streams...)
}

func (socket *FuturesWS_AllMarkPrices_Socket) Unsubscribe(isFast ...bool) (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	streams := make([]string, len(isFast))
	for i := range isFast {
		streams[i] = socket.CreateStreamName(isFast[i])
	}

	return socket.Handler.Unsubscribe(streams...)
}

func (futures_ws *Futures_Websockets) AllMarkPrices(publicOnMessage func(markPrices []*FuturesWS_MarkPrice), isFast ...bool) (*FuturesWS_AllMarkPrices_Socket, *Error) {
	var newSocket FuturesWS_AllMarkPrices_Socket

	streams := make([]string, len(isFast))
	for i := range isFast {
		streams[i] = newSocket.CreateStreamName(isFast[i])
	}
	socket, err := futures_ws.CreateSocket(streams, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var markPrices []*FuturesWS_MarkPrice
		err := json.Unmarshal(msg, &markPrices)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(markPrices)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type FuturesWS_BookTicker struct {

	// event type
	Event string `json:"e"`

	// order book updateId
	UpdateId int64 `json:"u"`

	// event time
	EventTime int64 `json:"E"`

	// transaction time
	TransactTime int64 `json:"T"`

	// symbol
	Symbol string `json:"s"`

	// best bid price
	Bid string `json:"b"`

	// best bid qty
	BidQty string `json:"B"`

	// best ask price
	Ask string `json:"a"`

	// best ask qty
	AskQty string `json:"A"`
}

type FuturesWS_BookTicker_Socket struct {
	Handler *Futures_Websocket
}

func (*FuturesWS_BookTicker_Socket) CreateStreamName(symbol ...string) []string {
	for i := range symbol {
		symbol[i] = strings.ToLower(symbol[i]) + "@bookTicker"
	}
	return symbol
}

func (socket *FuturesWS_BookTicker_Socket) Subscribe(symbol ...string) (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	symbol = socket.CreateStreamName(symbol...)
	return socket.Handler.Subscribe(symbol...)
}

func (socket *FuturesWS_BookTicker_Socket) Unsubscribe(symbol ...string) (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	symbol = socket.CreateStreamName(symbol...)
	return socket.Handler.Unsubscribe(symbol...)
}

func (futures_ws *Futures_Websockets) BookTicker(publicOnMessage func(bookTicker *FuturesWS_BookTicker), symbol ...string) (*FuturesWS_BookTicker_Socket, *Error) {
	var newSocket FuturesWS_BookTicker_Socket

	symbol = newSocket.CreateStreamName(symbol...)
	socket, err := futures_ws.CreateSocket(symbol, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var bookTicker FuturesWS_BookTicker
		err := json.Unmarshal(msg, &bookTicker)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(&bookTicker)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type FuturesWS_AllBookTickers_Socket struct {
	Handler *Futures_Websocket
}

func (*FuturesWS_AllBookTickers_Socket) CreateStreamName() string {
	return "!bookTicker"
}

func (socket *FuturesWS_AllBookTickers_Socket) Subscribe() (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	streamName := socket.CreateStreamName()
	return socket.Handler.Subscribe(streamName)
}

func (socket *FuturesWS_AllBookTickers_Socket) Unsubscribe() (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	streamName := socket.CreateStreamName()
	return socket.Handler.Unsubscribe(streamName)
}

func (futures_ws *Futures_Websockets) AllBookTickers(publicOnMessage func(bookTickers []*FuturesWS_BookTicker)) (*FuturesWS_AllBookTickers_Socket, *Error) {
	var newSocket FuturesWS_AllBookTickers_Socket

	streamName := newSocket.CreateStreamName()
	socket, err := futures_ws.CreateSocket([]string{streamName}, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var bookTickers []*FuturesWS_BookTicker
		err := json.Unmarshal(msg, &bookTickers)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(bookTickers)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type FuturesWS_ContractInfo struct {

	// Event Type
	Event string `json:"e"`

	// Event Time
	EventTime int64 `json:"E"`

	// Symbol
	Symbol string `json:"s"`

	// Pair
	Pair string `json:"ps"`

	// Contract type
	ContractType string `json:"ct"`

	// Delivery date time
	DeliveryDate int64 `json:"dt"`

	// onboard date time
	OnboardDateTime int64 `json:"ot"`

	// Contract status
	ContractStatus string `json:"cs"`

	Bks []*FuturesWS_ContractInfo_Bracket `json:"bks"`
}
type FuturesWS_ContractInfo_Bracket struct {

	// Notional bracket
	NotionalBracket int64 `json:"bs"`

	// Floor notional of this bracket
	FloorNotional int64 `json:"bnf"`

	// Cap notional of this bracket
	MaxNotional int64 `json:"bnc"`

	// Maintenance ratio for this bracket
	MaintenanceRatio float64 `json:"mmr"`

	// Auxiliary number for quick calculation
	Auxiliary int64 `json:"cf"`

	// Min leverage for this bracket
	MinLeverage int64 `json:"mi"`

	// Max leverage for this bracket
	MaxLeverage int64 `json:"ma"`
}

type FuturesWS_ContractInfo_Socket struct {
	Handler *Futures_Websocket
}

func (*FuturesWS_ContractInfo_Socket) CreateStreamName() string {
	return "!contractInfo"
}

func (socket *FuturesWS_ContractInfo_Socket) Subscribe() (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	streamName := socket.CreateStreamName()
	return socket.Handler.Subscribe(streamName)
}

func (socket *FuturesWS_ContractInfo_Socket) Unsubscribe() (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	streamName := socket.CreateStreamName()
	return socket.Handler.Unsubscribe(streamName)
}

func (futures_ws *Futures_Websockets) ContractInfo(publicOnMessage func(contractInfo *FuturesWS_ContractInfo)) (*FuturesWS_ContractInfo_Socket, *Error) {
	var newSocket FuturesWS_ContractInfo_Socket
	streamName := newSocket.CreateStreamName()
	socket, err := futures_ws.CreateSocket([]string{streamName}, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var aggTrade FuturesWS_ContractInfo
		err := json.Unmarshal(msg, &aggTrade)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(&aggTrade)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (*Futures_Websockets) CreateSocket(streams []string, isCombined bool) (*Futures_Websocket, *Error) {
	baseURL := FUTURES_Constants.Websocket.URLs[0]

	socket, err := CreateSocket(baseURL, streams, isCombined)
	if err != nil {
		return nil, err
	}

	socket.privateMessageValidator = func(msg []byte) (isPrivate bool, Id string) {

		if len(msg) > 0 && msg[0] == '[' {
			return false, ""
		}

		var privateMessage FuturesWS_PrivateMessage
		err := json.Unmarshal(msg, &privateMessage)
		if err != nil {
			fmt.Println("[PRIVATEMESSAGEVALIDATOR ERR] WS Message is the following:", string(msg))
			LocalError(PARSING_ERROR, err.Error())
			return false, ""
		}

		if privateMessage.Id == "" {
			return false, ""
		}

		return true, privateMessage.Id
	}

	ws := &Futures_Websocket{
		Websocket:       socket,
		Conn:            socket.Conn,
		BaseURL:         baseURL,
		pendingRequests: make(map[string]chan []byte),
	}

	return ws, nil
}

func (futures *Futures_Websocket) createRequestObject() map[string]interface{} {
	requestObj := make(map[string]interface{})

	requestObj["id"] = uuid.New().String()
	return requestObj
}

type FuturesWS_ListSubscriptions_Response struct {
	Id     string   `json:"id"`
	Result []string `json:"result"`
}

func (futures_ws *Futures_Websocket) ListSubscriptions(timeout_sec ...int) (resp *FuturesWS_ListSubscriptions_Response, hasTimedOut bool, err *Error) {
	requestObj := futures_ws.createRequestObject()
	requestObj["method"] = "LIST_SUBSCRIPTIONS"
	data, timedOut, err := futures_ws.Websocket.SendRequest_sync(requestObj, timeout_sec[0])
	if err != nil || timedOut {
		return nil, timedOut, err
	}

	var response FuturesWS_ListSubscriptions_Response
	unmarshallErr := json.Unmarshal(data, &response)
	if unmarshallErr != nil {
		return nil, false, LocalError(ERROR_PROCESSING_ERR, unmarshallErr.Error())
	}

	return &response, false, nil
}

type FuturesWS_Subscribe_Response struct {
	Id string `json:"id"`
}

func (futures_ws *Futures_Websocket) Subscribe(stream ...string) (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	requestObj := futures_ws.createRequestObject()
	requestObj["method"] = "SUBSCRIBE"
	requestObj["params"] = stream
	data, timedOut, err := futures_ws.Websocket.SendRequest_sync(requestObj)
	if err != nil || timedOut {
		return nil, timedOut, err
	}

	var response FuturesWS_Subscribe_Response
	unmarshallErr := json.Unmarshal(data, &response)
	if unmarshallErr != nil {
		return nil, false, LocalError(ERROR_PROCESSING_ERR, unmarshallErr.Error())
	}

	futures_ws.Websocket.Streams = append(futures_ws.Websocket.Streams, stream...)

	LOG_WS_VERBOSE("Successfully Subscribed to", stream)

	return &response, false, nil
}

type FuturesWS_Unsubscribe_Response struct {
	Id string `json:"id"`
}

func (futures_ws *Futures_Websocket) Unsubscribe(stream ...string) (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	requestObj := futures_ws.createRequestObject()
	requestObj["method"] = "UNSUBSCRIBE"
	requestObj["params"] = stream
	data, timedOut, err := futures_ws.Websocket.SendRequest_sync(requestObj)
	if err != nil || timedOut {
		return nil, timedOut, err
	}

	var response FuturesWS_Unsubscribe_Response
	unmarshallErr := json.Unmarshal(data, &response)
	if unmarshallErr != nil {
		return nil, false, LocalError(ERROR_PROCESSING_ERR, unmarshallErr.Error())
	}

	// Filter out the streams to remove from futures_ws.Websocket.Streams
	streamMap := make(map[string]bool)
	for _, s := range stream {
		streamMap[s] = true
	}

	var updatedStreams []string
	for _, existingStream := range futures_ws.Websocket.Streams {
		if !streamMap[existingStream] {
			updatedStreams = append(updatedStreams, existingStream)
		}
	}
	futures_ws.Websocket.Streams = updatedStreams

	LOG_WS_VERBOSE("Successfully Unsubscribed from", stream)

	return &response, false, nil
}
