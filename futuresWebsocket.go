package Binance

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"slices"

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

func (*FuturesWS_AggTrade_Socket) CreateStreamName(symbol ...string) []string {
	streamNames := make([]string, len(symbol))
	for i := range symbol {
		streamNames[i] = strings.ToLower(symbol[i]) + "@aggTrade"
	}
	return streamNames
}

func (socket *FuturesWS_AggTrade_Socket) Subscribe(symbol ...string) (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	streamNames := socket.CreateStreamName(symbol...)
	return socket.Handler.Subscribe(streamNames...)
}

func (socket *FuturesWS_AggTrade_Socket) Unsubscribe(symbol ...string) (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	streamNames := socket.CreateStreamName(symbol...)
	return socket.Handler.Unsubscribe(streamNames...)
}

func (futures_ws *Futures_Websockets) AggTrade(publicOnMessage func(aggTrade *FuturesWS_AggTrade), symbol ...string) (*FuturesWS_AggTrade_Socket, *Error) {
	var newSocket FuturesWS_AggTrade_Socket

	streamNames := newSocket.CreateStreamName(symbol...)

	socket, err := futures_ws.CreateSocket(streamNames, false)
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

func (*FuturesWS_MarkPrice_Socket) CreateStreamName(params ...FuturesWS_MarkPrice_Params) []string {
	streamNames := make([]string, len(params))
	for i := range params {
		streamNames[i] = strings.ToLower(params[i].Symbol) + "@markPrice"
		if params[i].IsFast {
			streamNames[i] += "@1s"
		}
	}
	return streamNames
}

func (socket *FuturesWS_MarkPrice_Socket) Subscribe(params ...FuturesWS_MarkPrice_Params) (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	streams := socket.CreateStreamName(params...)
	return socket.Handler.Subscribe(streams...)
}

func (socket *FuturesWS_MarkPrice_Socket) Unsubscribe(params ...FuturesWS_MarkPrice_Params) (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	streams := socket.CreateStreamName(params...)
	return socket.Handler.Unsubscribe(streams...)
}

func (futures_ws *Futures_Websockets) MarkPrice(publicOnMessage func(markPrice *FuturesWS_MarkPrice), params ...FuturesWS_MarkPrice_Params) (*FuturesWS_MarkPrice_Socket, *Error) {
	var newSocket FuturesWS_MarkPrice_Socket

	streams := newSocket.CreateStreamName(params...)

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

type FuturesWS_Candlestick struct {

	// Event type
	Event string `json:"e"`

	// Event time
	EventTime int64 `json:"E"`

	// Symbol
	Symbol string `json:"s"`

	Kline *FuturesWS_Candlestick_Kline `json:"k"`
}
type FuturesWS_Candlestick_Kline struct {

	// Kline start time
	OpenTime int64 `json:"t"`

	// Kline close time
	CloseTime int64 `json:"T"`

	// Symbol
	Symbol string `json:"s"`

	// Interval
	Interval string `json:"i"`

	// First trade ID
	FirstTradeId int64 `json:"f"`

	// Last trade ID
	LastTradeId int64 `json:"L"`

	// Open price
	Open string `json:"o"`

	// Close price
	Close string `json:"c"`

	// High price
	High string `json:"h"`

	// Low price
	Low string `json:"l"`

	// Base asset volume
	BaseAssetVolume string `json:"v"`

	// Number of trades
	TradeCount int64 `json:"n"`

	// Is this kline closed?
	IsClosed bool `json:"x"`

	// Quote asset volume
	QuoteAssetVolume string `json:"q"`

	// Taker buy base asset volume
	TakerBuyBaseAssetVolume string `json:"V"`

	// Taker buy quote asset volume
	TakerBuyQuoteAssetVolume string `json:"Q"`

	// Ignore
	Ignore string `json:"B"`
}

type FuturesWS_Candlestick_Socket struct {
	Handler *Futures_Websocket
}

type FuturesWS_Candlestick_Params struct {
	Symbol   string
	Interval string
}

func (*FuturesWS_Candlestick_Socket) CreateStreamName(params ...FuturesWS_Candlestick_Params) []string {
	streamNames := make([]string, len(params))
	for i := range params {
		streamNames[i] = strings.ToLower(params[i].Symbol) + "@kline_" + params[i].Interval
	}
	return streamNames
}

func (socket *FuturesWS_Candlestick_Socket) Subscribe(params ...FuturesWS_Candlestick_Params) (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	streamNames := socket.CreateStreamName(params...)
	return socket.Handler.Subscribe(streamNames...)
}

func (socket *FuturesWS_Candlestick_Socket) Unsubscribe(params ...FuturesWS_Candlestick_Params) (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	streamNames := socket.CreateStreamName(params...)
	return socket.Handler.Unsubscribe(streamNames...)
}

func (futures_ws *Futures_Websockets) Candlesticks(publicOnMessage func(candlestick *FuturesWS_Candlestick), params ...FuturesWS_Candlestick_Params) (*FuturesWS_Candlestick_Socket, *Error) {
	var newSocket FuturesWS_Candlestick_Socket

	streamNames := newSocket.CreateStreamName(params...)

	socket, err := futures_ws.CreateSocket(streamNames, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var kline *FuturesWS_Candlestick
		err := json.Unmarshal(msg, &kline)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(kline)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type FuturesWS_ContinuousCandlestick struct {

	// Event type
	Event string `json:"e"`

	// Event time
	EventTime int64 `json:"E"`

	// Pair
	Pair string `json:"ps"`

	// Contract type
	ContractType string `json:"ct"`

	Kline *FuturesWS_ContinuousCandlestick_Kline `json:"k"`
}
type FuturesWS_ContinuousCandlestick_Kline struct {

	// Kline start time
	OpenTime int64 `json:"t"`

	// Kline close time
	CloseTime int64 `json:"T"`

	// Interval
	Interval string `json:"i"`

	// First updateId
	FirstUpdateId int64 `json:"f"`

	// Last updateId
	LastUpdateId int64 `json:"L"`

	// Open price
	Open string `json:"o"`

	// Close price
	Close string `json:"c"`

	// High price
	High string `json:"h"`

	// Low price
	Low string `json:"l"`

	// volume
	Volume string `json:"v"`

	// Number of trades
	TradeCount int64 `json:"n"`

	// Is this kline closed?
	IsClosed bool `json:"x"`

	// Quote asset volume
	QuoteAssetVolume string `json:"q"`

	// Taker buy volume
	TakerBuyVolume string `json:"V"`

	// Taker buy quote asset volume
	TakerBuyQuoteAssetVolume string `json:"Q"`

	// Ignore
	Ignore string `json:"B"`
}

type FuturesWS_ContinuousCandlestick_Socket struct {
	Handler *Futures_Websocket
}

type FuturesWS_ContinuousCandlestick_Params struct {
	Symbol       string
	ContractType string
	Interval     string
}

func (*FuturesWS_ContinuousCandlestick_Socket) CreateStreamName(params ...FuturesWS_ContinuousCandlestick_Params) []string {
	streamNames := make([]string, len(params))
	for i := range params {
		streamNames[i] = strings.ToLower(params[i].Symbol) + "@kline_" + params[i].Interval
	}
	return streamNames
}

func (socket *FuturesWS_ContinuousCandlestick_Socket) Subscribe(params ...FuturesWS_ContinuousCandlestick_Params) (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	streamNames := socket.CreateStreamName(params...)
	return socket.Handler.Subscribe(streamNames...)
}

func (socket *FuturesWS_ContinuousCandlestick_Socket) Unsubscribe(params ...FuturesWS_ContinuousCandlestick_Params) (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	streamNames := socket.CreateStreamName(params...)
	return socket.Handler.Unsubscribe(streamNames...)
}

func (futures_ws *Futures_Websockets) ContinuousCandlesticks(publicOnMessage func(candlestick *FuturesWS_ContinuousCandlestick), params ...FuturesWS_ContinuousCandlestick_Params) (*FuturesWS_ContinuousCandlestick_Socket, *Error) {
	var newSocket FuturesWS_ContinuousCandlestick_Socket

	streamNames := newSocket.CreateStreamName(params...)

	socket, err := futures_ws.CreateSocket(streamNames, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var kline *FuturesWS_ContinuousCandlestick
		err := json.Unmarshal(msg, &kline)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(kline)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type FuturesWS_MiniTicker struct {

	// Event type
	Event string `json:"e"`

	// Event time
	EventTime int64 `json:"E"`

	// Symbol
	Symbol string `json:"s"`

	// Close price
	Close string `json:"c"`

	// Open price
	Open string `json:"o"`

	// High price
	High string `json:"h"`

	// Low price
	Low string `json:"l"`

	// Total traded base asset volume
	BaseAssetVolume string `json:"v"`

	// Total traded quote asset volume
	QuoteAssetVolume string `json:"q"`
}

type FuturesWS_MiniTicker_Socket struct {
	Handler *Futures_Websocket
}

func (*FuturesWS_MiniTicker_Socket) CreateStreamName(symbol ...string) []string {
	streamNames := make([]string, len(symbol))
	for i := range symbol {
		streamNames[i] = strings.ToLower(symbol[i]) + "@miniTicker"
	}
	return streamNames
}

func (socket *FuturesWS_MiniTicker_Socket) Subscribe(symbol ...string) (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	streamNames := socket.CreateStreamName(symbol...)
	return socket.Handler.Subscribe(streamNames...)
}

func (socket *FuturesWS_MiniTicker_Socket) Unsubscribe(symbol ...string) (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	streamNames := socket.CreateStreamName(symbol...)
	return socket.Handler.Unsubscribe(streamNames...)
}

func (futures_ws *Futures_Websockets) MiniTicker(publicOnMessage func(miniTicker *FuturesWS_MiniTicker), symbol ...string) (*FuturesWS_MiniTicker_Socket, *Error) {
	var newSocket FuturesWS_MiniTicker_Socket

	streamNames := newSocket.CreateStreamName(symbol...)

	socket, err := futures_ws.CreateSocket(streamNames, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var miniTicker *FuturesWS_MiniTicker
		err := json.Unmarshal(msg, &miniTicker)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(miniTicker)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type FuturesWS_AllMiniTickers_Socket struct {
	Handler *Futures_Websocket
}

func (*FuturesWS_AllMiniTickers_Socket) CreateStreamName() string {
	return "!miniTicker@arr"
}

func (socket *FuturesWS_AllMiniTickers_Socket) Subscribe() (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	return socket.Handler.Subscribe(socket.CreateStreamName())
}

func (socket *FuturesWS_AllMiniTickers_Socket) Unsubscribe() (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	return socket.Handler.Unsubscribe(socket.CreateStreamName())
}

func (futures_ws *Futures_Websockets) AllMiniTickers(publicOnMessage func(miniTickers []*FuturesWS_MiniTicker)) (*FuturesWS_AllMiniTickers_Socket, *Error) {
	var newSocket FuturesWS_AllMiniTickers_Socket

	streamName := newSocket.CreateStreamName()

	socket, err := futures_ws.CreateSocket([]string{streamName}, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var miniTickers []*FuturesWS_MiniTicker
		err := json.Unmarshal(msg, &miniTickers)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(miniTickers)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type FuturesWS_Ticker struct {

	// Event type
	Event string `json:"e"`

	// Event time
	EventTime int64 `json:"E"`

	// Symbol
	Symbol string `json:"s"`

	// Price change
	PriceChange string `json:"p"`

	// Price change percent
	PriceChangePercent string `json:"P"`

	// Weighted average price
	WeightedAveragePrice string `json:"w"`

	// Last price
	LastPrice string `json:"c"`

	// Last quantity
	LastQty string `json:"Q"`

	// Open price
	Open string `json:"o"`

	// High price
	High string `json:"h"`

	// Low price
	Low string `json:"l"`

	// Total traded base asset volume
	BaseAssetVolume string `json:"v"`

	// Total traded quote asset volume
	QuoteAssetVolume string `json:"q"`

	// Statistics open time
	OpenTime int64 `json:"O"`

	// Statistics close time
	CloseTime int64 `json:"C"`

	// First trade ID
	FirstTradeId int64 `json:"F"`

	// Last trade Id
	LastTradeId int64 `json:"L"`

	// Total number of trades
	TradeCount int64 `json:"n"`
}

type FuturesWS_Ticker_Socket struct {
	Handler *Futures_Websocket
}

func (*FuturesWS_Ticker_Socket) CreateStreamName(symbol ...string) []string {
	streamNames := make([]string, len(symbol))
	for i := range symbol {
		streamNames[i] = strings.ToLower(symbol[i]) + "@ticker"
	}
	return streamNames
}

func (socket *FuturesWS_Ticker_Socket) Subscribe(symbol ...string) (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	streamNames := socket.CreateStreamName(symbol...)
	return socket.Handler.Subscribe(streamNames...)
}

func (socket *FuturesWS_Ticker_Socket) Unsubscribe(symbol ...string) (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	streamNames := socket.CreateStreamName(symbol...)
	return socket.Handler.Unsubscribe(streamNames...)
}

func (futures_ws *Futures_Websockets) Ticker(publicOnMessage func(ticker *FuturesWS_Ticker), symbol ...string) (*FuturesWS_Ticker_Socket, *Error) {
	var newSocket FuturesWS_Ticker_Socket

	streamNames := newSocket.CreateStreamName(symbol...)

	socket, err := futures_ws.CreateSocket(streamNames, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var ticker *FuturesWS_Ticker
		err := json.Unmarshal(msg, &ticker)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(ticker)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type FuturesWS_AllTickers_Socket struct {
	Handler *Futures_Websocket
}

func (*FuturesWS_AllTickers_Socket) CreateStreamName() string {
	return "!ticker@arr"
}

func (socket *FuturesWS_AllTickers_Socket) Subscribe() (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	return socket.Handler.Subscribe(socket.CreateStreamName())
}

func (socket *FuturesWS_AllTickers_Socket) Unsubscribe() (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	return socket.Handler.Unsubscribe(socket.CreateStreamName())
}

func (futures_ws *Futures_Websockets) AllTickers(publicOnMessage func(tickers []*FuturesWS_Ticker)) (*FuturesWS_AllTickers_Socket, *Error) {
	var newSocket FuturesWS_AllTickers_Socket

	streamName := newSocket.CreateStreamName()

	socket, err := futures_ws.CreateSocket([]string{streamName}, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var tickers []*FuturesWS_Ticker
		err := json.Unmarshal(msg, &tickers)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(tickers)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

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

type FuturesWS_LiquidationOrder struct {

	// Event Type
	Event string `json:"e"`

	// Event Time
	EventTime int64 `json:"E"`

	Type *FuturesWS_LiquidationOrder_Order `json:"o"`
}
type FuturesWS_LiquidationOrder_Order struct {

	// Symbol
	Symbol string `json:"s"`

	// Side
	Side string `json:"S"`

	// Order Type
	Type string `json:"o"`

	// Time in Force
	TimeInForce string `json:"f"`

	// Original Quantity
	OrigQty string `json:"q"`

	// Price
	Price string `json:"p"`

	// Average Price
	AvgPrice string `json:"ap"`

	// Order Status
	Status string `json:"X"`

	// Order Last Filled Quantity
	LastFilledQty string `json:"l"`

	// Order Filled Accumulated Quantity
	CumQty string `json:"z"`

	// Order Trade Time
	TradeTime int64 `json:"T"`
}

type FuturesWS_LiquidationOrder_Socket struct {
	Handler *Futures_Websocket
}

func (*FuturesWS_LiquidationOrder_Socket) CreateStreamName(symbol ...string) []string {
	streams := make([]string, len(symbol))
	for i := range streams {
		streams[i] = strings.ToLower(symbol[i])
	}

	return streams
}

func (socket *FuturesWS_LiquidationOrder_Socket) Subscribe(symbol ...string) (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	streamNames := socket.CreateStreamName(symbol...)
	return socket.Handler.Subscribe(streamNames...)
}

func (socket *FuturesWS_LiquidationOrder_Socket) Unsubscribe(symbol ...string) (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	streamNames := socket.CreateStreamName(symbol...)
	return socket.Handler.Unsubscribe(streamNames...)
}

func (futures_ws *Futures_Websockets) LiquidationOrders(publicOnMessage func(liquidationOrder *FuturesWS_LiquidationOrder), symbol ...string) (*FuturesWS_LiquidationOrder_Socket, *Error) {
	var newSocket FuturesWS_LiquidationOrder_Socket

	streamNames := newSocket.CreateStreamName(symbol...)

	socket, err := futures_ws.CreateSocket(streamNames, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var liquidationOrder *FuturesWS_LiquidationOrder
		err := json.Unmarshal(msg, &liquidationOrder)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(liquidationOrder)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type FuturesWS_AllLiquidationOrders_Socket struct {
	Handler *Futures_Websocket
}

func (*FuturesWS_AllLiquidationOrders_Socket) CreateStreamName() string {
	return "!forceOrder@arr"
}

func (socket *FuturesWS_AllLiquidationOrders_Socket) Subscribe() (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	return socket.Handler.Subscribe(socket.CreateStreamName())
}

func (socket *FuturesWS_AllLiquidationOrders_Socket) Unsubscribe() (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	return socket.Handler.Unsubscribe(socket.CreateStreamName())
}

func (futures_ws *Futures_Websockets) AllLiquidationOrders(publicOnMessage func(liquidationOrder *FuturesWS_LiquidationOrder)) (*FuturesWS_AllLiquidationOrders_Socket, *Error) {
	var newSocket FuturesWS_AllLiquidationOrders_Socket

	streamName := newSocket.CreateStreamName()

	socket, err := futures_ws.CreateSocket([]string{streamName}, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var liquidationOrder *FuturesWS_LiquidationOrder
		err := json.Unmarshal(msg, &liquidationOrder)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(liquidationOrder)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type FuturesWS_PartialBookDepth struct {

	// Event Type
	Event string `json:"e"`

	// Event Time
	EventTime int64 `json:"E"`

	TransactTime int64 `json:"T"`

	Symbol string `json:"s"`

	FirstUpdateId int64 `json:"U"`

	LastUpdateId string `json:"u"`

	Previous_LastUpdateId string `json:"pu"`

	// Bids to be updated
	//
	// [
	//     [
	//       "7405.96",      // Price level to be
	//       "3.340"         // Quantity
	//     ],
	// ]
	Bids [][2]string `json:"b"`

	// Asks to be updated
	//
	// [
	//     [
	//       "7405.96",      // Price level to be
	//       "3.340"         // Quantity
	//     ],
	// ]
	Asks [][2]string `json:"a"`
}

type FuturesWS_PartialBookDepth_Socket struct {
	Handler *Futures_Websocket
}

type FuturesWS_PartialBookDepth_Params struct {
	Symbol string

	// Possible values: 5, 10 or 20
	Levels int

	// Possible values: "500ms", "250ms" or "100ms"
	//
	// Default: "250ms"
	UpdateSpeed string
}

func (*FuturesWS_PartialBookDepth_Socket) CreateStreamName(params ...FuturesWS_PartialBookDepth_Params) []string {
	streamNames := make([]string, len(params))
	for i := range params {
		streamNames[i] = strings.ToLower(params[i].Symbol) + "@depth" + fmt.Sprint(params[i].Levels)
		if params[i].UpdateSpeed != "" && params[i].UpdateSpeed != "250ms" {
			streamNames[i] += "@" + params[i].UpdateSpeed
		}
	}

	return streamNames
}

func (socket *FuturesWS_PartialBookDepth_Socket) Subscribe(params ...FuturesWS_PartialBookDepth_Params) (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	streamNames := socket.CreateStreamName(params...)
	return socket.Handler.Subscribe(streamNames...)
}

func (socket *FuturesWS_PartialBookDepth_Socket) Unsubscribe(params ...FuturesWS_PartialBookDepth_Params) (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	streamNames := socket.CreateStreamName(params...)
	return socket.Handler.Unsubscribe(streamNames...)
}

func (futures_ws *Futures_Websockets) PartialBookDepth(publicOnMessage func(partialBookDepth *FuturesWS_PartialBookDepth), params ...FuturesWS_PartialBookDepth_Params) (*FuturesWS_PartialBookDepth_Socket, *Error) {
	var newSocket FuturesWS_PartialBookDepth_Socket

	streamNames := newSocket.CreateStreamName(params...)

	socket, err := futures_ws.CreateSocket(streamNames, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var partialBookDepth *FuturesWS_PartialBookDepth
		err := json.Unmarshal(msg, &partialBookDepth)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(partialBookDepth)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type FuturesWS_DiffBookDepth struct {

	// Event Type
	Event string `json:"e"`

	// Event Time
	EventTime int64 `json:"E"`

	TransactTime int64 `json:"T"`

	Symbol string `json:"s"`

	FirstUpdateId int64 `json:"U"`

	LastUpdateId int64 `json:"u"`

	Previous_LastUpdateId int64 `json:"pu"`

	// Bids to be updated
	//
	// [
	//     [
	//       "7405.96",      // Price level to be
	//       "3.340"         // Quantity
	//     ],
	// ]
	Bids [][2]string `json:"b"`

	// Asks to be updated
	//
	// [
	//     [
	//       "7405.96",      // Price level to be
	//       "3.340"         // Quantity
	//     ],
	// ]
	Asks [][2]string `json:"a"`
}

type FuturesWS_DiffBookDepth_Socket struct {
	Handler *Futures_Websocket
}

type FuturesWS_DiffBookDepth_Params struct {
	Symbol string

	// Possible values: "500ms", "250ms" or "100ms"
	//
	// Default: "250ms"
	UpdateSpeed string
}

func (*FuturesWS_DiffBookDepth_Socket) CreateStreamName(params ...FuturesWS_DiffBookDepth_Params) []string {
	streamNames := make([]string, len(params))
	for i := range params {
		streamNames[i] = strings.ToLower(params[i].Symbol) + "@depth"
		if params[i].UpdateSpeed != "" && params[i].UpdateSpeed != "250ms" {
			streamNames[i] += "@" + params[i].UpdateSpeed
		}
	}

	return streamNames
}

func (socket *FuturesWS_DiffBookDepth_Socket) Subscribe(params ...FuturesWS_DiffBookDepth_Params) (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	streamNames := socket.CreateStreamName(params...)
	return socket.Handler.Subscribe(streamNames...)
}

func (socket *FuturesWS_DiffBookDepth_Socket) Unsubscribe(params ...FuturesWS_DiffBookDepth_Params) (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	streamNames := socket.CreateStreamName(params...)
	return socket.Handler.Unsubscribe(streamNames...)
}

func (futures_ws *Futures_Websockets) DiffBookDepth(publicOnMessage func(diffBookDepth *FuturesWS_DiffBookDepth), params ...FuturesWS_DiffBookDepth_Params) (*FuturesWS_DiffBookDepth_Socket, *Error) {
	var newSocket FuturesWS_DiffBookDepth_Socket

	streamNames := newSocket.CreateStreamName(params...)

	socket, err := futures_ws.CreateSocket(streamNames, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var diffBookDepth *FuturesWS_DiffBookDepth
		err := json.Unmarshal(msg, &diffBookDepth)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(diffBookDepth)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type FuturesWS_CompositeIndexSymbolInfo struct {

	// Event Type
	Event string `json:"e"`

	// Event Time
	EventTime int64 `json:"E"`

	Symbol string `json:"s"`

	Price int64 `json:"p"`

	BaseAsset string `json:"C"`

	Composition []*FuturesWS_CompositeIndexSymbolInfo_Composition
}

type FuturesWS_CompositeIndexSymbolInfo_Composition struct {

	// Base asset
	BaseAsset string `json:"b"`

	// Quote asset
	QuoteAsset string `json:"q"`

	// Weight in quantity
	Weight string `json:"w"`

	// Weight in percentage
	WeightPercent string `json:"W"`

	// Index Price
	IndexPrice string `json:"i"`
}

type FuturesWS_CompositeIndexSymbolInfo_Socket struct {
	Handler *Futures_Websocket
}

func (*FuturesWS_CompositeIndexSymbolInfo_Socket) CreateStreamName(symbol ...string) []string {
	streamNames := make([]string, len(symbol))
	for i := range symbol {
		streamNames[i] = strings.ToLower(symbol[i]) + "@compositeIndex"
	}

	return streamNames
}

func (socket *FuturesWS_CompositeIndexSymbolInfo_Socket) Subscribe(symbol ...string) (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	streamNames := socket.CreateStreamName(symbol...)
	return socket.Handler.Subscribe(streamNames...)
}

func (socket *FuturesWS_CompositeIndexSymbolInfo_Socket) Unsubscribe(symbol ...string) (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	streamNames := socket.CreateStreamName(symbol...)
	return socket.Handler.Unsubscribe(streamNames...)
}

func (futures_ws *Futures_Websockets) CompositeIndexSymbolInfo(publicOnMessage func(compositeIndexSymbolInfo *FuturesWS_CompositeIndexSymbolInfo), symbol ...string) (*FuturesWS_CompositeIndexSymbolInfo_Socket, *Error) {
	var newSocket FuturesWS_CompositeIndexSymbolInfo_Socket

	streamNames := newSocket.CreateStreamName(symbol...)

	socket, err := futures_ws.CreateSocket(streamNames, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var compositeIndexSymbolInfo *FuturesWS_CompositeIndexSymbolInfo
		err := json.Unmarshal(msg, &compositeIndexSymbolInfo)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(compositeIndexSymbolInfo)
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

type FuturesWS_MultiAssetsModeAssetIndex struct {
	Event string `json:"e"`

	EventTime int64 `json:"E"`

	// asset index symbol
	Symbol string `json:"s"`

	// index price
	IndexPrice string `json:"i"`

	// bid buffer
	BidBuffer string `json:"b"`

	// ask buffer
	AskBuffer string `json:"a"`

	// bid rate
	BidRate string `json:"B"`

	// ask rate
	AskRate string `json:"A"`

	// auto exchange bid buffer
	AutoExchange_BidBuffer string `json:"q"`

	// auto exchange ask buffer
	AutoExchange_AskBuffer string `json:"g"`

	// auto exchange bid rate
	AutoExchange_BidRate string `json:"Q"`

	// auto exchange ask rate
	AutoExchange_AskRate string `json:"G"`
}

type FuturesWS_MultiAssetsModeAssetIndex_Socket struct {
	Handler *Futures_Websocket
}

func (*FuturesWS_MultiAssetsModeAssetIndex_Socket) CreateStreamName(assetSymbol ...string) []string {
	streamNames := make([]string, len(assetSymbol))
	for i := range assetSymbol {
		streamNames[i] = strings.ToLower(assetSymbol[i]) + "@assetIndex"
	}

	return streamNames
}

func (socket *FuturesWS_MultiAssetsModeAssetIndex_Socket) Subscribe(assetSymbol ...string) (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	streamNames := socket.CreateStreamName(assetSymbol...)
	return socket.Handler.Subscribe(streamNames...)
}

func (socket *FuturesWS_MultiAssetsModeAssetIndex_Socket) Unsubscribe(assetSymbol ...string) (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	streamNames := socket.CreateStreamName(assetSymbol...)
	return socket.Handler.Unsubscribe(streamNames...)
}

func (futures_ws *Futures_Websockets) MultiAssetsModeAssetIndex(publicOnMessage func(assetIndexes []*FuturesWS_MultiAssetsModeAssetIndex), assetSymbol ...string) (*FuturesWS_MultiAssetsModeAssetIndex_Socket, *Error) {
	var newSocket FuturesWS_MultiAssetsModeAssetIndex_Socket

	streamNames := newSocket.CreateStreamName(assetSymbol...)

	socket, err := futures_ws.CreateSocket(streamNames, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var assetIndexes []*FuturesWS_MultiAssetsModeAssetIndex
		err := json.Unmarshal(msg, &assetIndexes)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(assetIndexes)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type FuturesWS_AllMultiAssetsModeAssetIndexes_Socket struct {
	Handler *Futures_Websocket
}

func (*FuturesWS_AllMultiAssetsModeAssetIndexes_Socket) CreateStreamName() string {
	return "!assetIndex@arr"
}

func (socket *FuturesWS_AllMultiAssetsModeAssetIndexes_Socket) Subscribe() (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	return socket.Handler.Subscribe(socket.CreateStreamName())
}

func (socket *FuturesWS_AllMultiAssetsModeAssetIndexes_Socket) Unsubscribe() (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	return socket.Handler.Unsubscribe(socket.CreateStreamName())
}

func (futures_ws *Futures_Websockets) AllMultiAssetsModeAssetIndexes(publicOnMessage func(assetIndexes []*FuturesWS_MultiAssetsModeAssetIndex), assetSymbol ...string) (*FuturesWS_AllMultiAssetsModeAssetIndexes_Socket, *Error) {
	var newSocket FuturesWS_AllMultiAssetsModeAssetIndexes_Socket

	streamName := newSocket.CreateStreamName()

	socket, err := futures_ws.CreateSocket([]string{streamName}, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var assetIndexes []*FuturesWS_MultiAssetsModeAssetIndex
		err := json.Unmarshal(msg, &assetIndexes)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(assetIndexes)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Futures_ManagedOrderbook struct {
	Symbol       string
	LastUpdateId int64
	Bids         [][2]string
	Asks         [][2]string

	isReadyToUpdate bool
	isFetching      bool

	previousEvent *FuturesWS_DiffBookDepth
}

func (managedOrderBook *Futures_ManagedOrderbook) addEvent(event *FuturesWS_DiffBookDepth) {
	for _, newBid := range event.Bids {
		found := false
		for i, bid := range managedOrderBook.Bids {
			if bid[0] == newBid[0] {
				found = true

				float, err := strconv.ParseFloat(newBid[1], 64)
				if err == nil && float == 0 {
					managedOrderBook.Bids = slices.Delete(managedOrderBook.Bids, i, i+1)
				} else {
					managedOrderBook.Bids[i][1] = newBid[1]
				}
			}
		}

		if !found {
			float, err := strconv.ParseFloat(newBid[1], 64)
			if err != nil || float != 0 {
				inserted := false
				for i := range managedOrderBook.Bids {
					// Compare the price of the new bid with the existing bid
					if managedOrderBook.Bids[i][0] < newBid[0] {
						// Insert the new bid before this index
						managedOrderBook.Bids = append(managedOrderBook.Bids[:i], append([][2]string{newBid}, managedOrderBook.Bids[i:]...)...)
						inserted = true
						break
					}
				}

				// If the new bid is lower than all existing bids, append it to the end
				if !inserted {
					managedOrderBook.Bids = append(managedOrderBook.Bids, newBid)
				}
			}
		}
	}

	for _, newAsk := range event.Asks {
		found := false
		for i, ask := range managedOrderBook.Asks {
			if ask[0] == newAsk[0] {

				float, err := strconv.ParseFloat(newAsk[1], 64)
				if err == nil && float == 0 {
					managedOrderBook.Asks = slices.Delete(managedOrderBook.Asks, i, i+1)
				} else {
					managedOrderBook.Asks[i][1] = newAsk[1]
				}
				found = true
			}
		}

		if !found {
			float, err := strconv.ParseFloat(newAsk[1], 64)
			if err != nil || float != 0 {
				inserted := false
				for i := range managedOrderBook.Asks {
					// Compare the price of the new ask with the existing ask
					if managedOrderBook.Asks[i][0] > newAsk[0] {
						// Insert the new ask before this index
						managedOrderBook.Asks = append(managedOrderBook.Asks[:i], append([][2]string{newAsk}, managedOrderBook.Asks[i:]...)...)
						inserted = true
						break
					}
				}

				// If the new ask is higher than all existing asks, append it to the end
				if !inserted {
					managedOrderBook.Asks = append(managedOrderBook.Asks, newAsk)
				}
			}
		}

	}

	// start := time.Now().UnixMicro()
	sort.SliceStable(managedOrderBook.Asks, func(i, j int) bool {
		val1, _ := ParseFloat(managedOrderBook.Asks[i][0])
		val2, _ := ParseFloat(managedOrderBook.Asks[j][0])
		return val1 < val2
	})

	sort.SliceStable(managedOrderBook.Bids, func(i, j int) bool {
		val1, _ := ParseFloat(managedOrderBook.Bids[i][0])
		val2, _ := ParseFloat(managedOrderBook.Bids[j][0])
		return val1 > val2
	})
	// fmt.Println(time.Now().UnixMicro()-start, "micro")

	managedOrderBook.previousEvent = event
}

type FuturesWS_ManagedOrderBook_Handler struct {
	DiffBookDepth_Socket *FuturesWS_DiffBookDepth_Socket

	Orderbooks struct {
		mu      sync.Mutex
		Symbols map[string]struct {
			bufferedEvents []*FuturesWS_DiffBookDepth
			Orderbook      *Futures_ManagedOrderbook
		}
	}
}

func (handler *FuturesWS_ManagedOrderBook_Handler) handleWSMessage(futures_ws *Futures_Websockets, diffBookDepth *FuturesWS_DiffBookDepth) (shouldPushEvent bool, managedOrderBook *Futures_ManagedOrderbook) {
	handler.Orderbooks.mu.Lock()
	defer handler.Orderbooks.mu.Unlock()
	Orderbook_symbol, exists := handler.Orderbooks.Symbols[diffBookDepth.Symbol]
	if !exists {
		LOG_WS_MESSAGES(fmt.Sprintf("Orderbook data for %s weren't found in handleWSMessage", diffBookDepth.Symbol))
		return false, nil
	}

	Orderbook_symbol.bufferedEvents = append(Orderbook_symbol.bufferedEvents, diffBookDepth)

	if !Orderbook_symbol.Orderbook.isReadyToUpdate && Orderbook_symbol.Orderbook.isFetching {
		LOG_WS_VERBOSE(fmt.Sprintf("[MANAGEDORDERBOOK] '%s' isNotReadyToUpdate and isFetching", diffBookDepth.Symbol))
		return false, nil
	}

	if !Orderbook_symbol.Orderbook.isReadyToUpdate {
		LOG_WS_VERBOSE(fmt.Sprintf("[MANAGEDORDERBOOK] Fetching '%s' snapshot", diffBookDepth.Symbol))
		Orderbook_symbol.Orderbook.isFetching = true

		newOrderBook, _, err := futures_ws.binance.Futures.OrderBook(diffBookDepth.Symbol, 1000)
		Orderbook_symbol.Orderbook.isFetching = false
		if err != nil {
			LOG_WS_ERRORS(fmt.Sprintf("Failed to refetch orderbook data for %s: %s", diffBookDepth.Symbol, err.Error()))
			return false, nil
		}

		Orderbook_symbol.Orderbook.LastUpdateId = newOrderBook.LastUpdateId
		Orderbook_symbol.Orderbook.Bids = newOrderBook.Bids
		Orderbook_symbol.Orderbook.Asks = newOrderBook.Asks

		Orderbook_symbol.Orderbook.isReadyToUpdate = true
	}

	if len(Orderbook_symbol.bufferedEvents) == 0 {
		LOG_WS_VERBOSE(fmt.Sprintf("[MANAGEDORDERBOOK] '%s' bufferedEvents empty", diffBookDepth.Symbol))
		return false, nil
	}

	for len(Orderbook_symbol.bufferedEvents) > 0 {

		event := Orderbook_symbol.bufferedEvents[0]
		Orderbook_symbol.bufferedEvents = Orderbook_symbol.bufferedEvents[1:]

		if event.LastUpdateId < Orderbook_symbol.Orderbook.LastUpdateId {
			LOG_WS_VERBOSE(fmt.Sprintf("[MANAGEDORDERBOOK] '%s' ignored since 'u' < lastUpdateId: %d < %d", diffBookDepth.Symbol, event.LastUpdateId, Orderbook_symbol.Orderbook.LastUpdateId))
			continue
		}

		if Orderbook_symbol.Orderbook.previousEvent == nil {
			if event.FirstUpdateId <= Orderbook_symbol.Orderbook.LastUpdateId && Orderbook_symbol.Orderbook.LastUpdateId <= event.LastUpdateId {
				LOG_WS_VERBOSE(fmt.Sprintf("[MANAGEDORDERBOOK] '%s' first valid event found", diffBookDepth.Symbol))
				Orderbook_symbol.Orderbook.addEvent(event)
			} else {
				LOG_WS_VERBOSE(fmt.Sprintf("[MANAGEDORDERBOOK] '%s' INVALID FIRST EVENT FOUND, U: %d, lastUpdateId: %d, u: %d", diffBookDepth.Symbol, event.FirstUpdateId, Orderbook_symbol.Orderbook.LastUpdateId, event.LastUpdateId))
			}
		} else if event.Previous_LastUpdateId == Orderbook_symbol.Orderbook.previousEvent.LastUpdateId {
			LOG_WS_VERBOSE(fmt.Sprintf("[MANAGEDORDERBOOK] '%s' good event found, pu == [lastEvent].u, %d == %d", diffBookDepth.Symbol, event.Previous_LastUpdateId, Orderbook_symbol.Orderbook.previousEvent.LastUpdateId))
			Orderbook_symbol.Orderbook.addEvent(event)
		} else {
			Orderbook_symbol.Orderbook.isReadyToUpdate = false
			Orderbook_symbol.Orderbook.isFetching = false
			Orderbook_symbol.Orderbook.previousEvent = nil

			return false, nil
		}
	}

	return true, Orderbook_symbol.Orderbook
}

func (handler *FuturesWS_ManagedOrderBook_Handler) openNewSymbols(params ...FuturesWS_DiffBookDepth_Params) {
	handler.Orderbooks.mu.Lock()
	defer handler.Orderbooks.mu.Unlock()

	for _, param := range params {
		symbol := param.Symbol

		_, exists := handler.Orderbooks.Symbols[symbol]
		if exists {
			continue
		}

		handler.Orderbooks.Symbols[symbol] = struct {
			bufferedEvents []*FuturesWS_DiffBookDepth
			Orderbook      *Futures_ManagedOrderbook
		}{
			bufferedEvents: make([]*FuturesWS_DiffBookDepth, 0),
			Orderbook: &Futures_ManagedOrderbook{
				Symbol: symbol,

				isFetching:      false,
				isReadyToUpdate: false,
			},
		}
	}
}

func (handler *FuturesWS_ManagedOrderBook_Handler) removeSymbols(params ...FuturesWS_DiffBookDepth_Params) {
	handler.Orderbooks.mu.Lock()
	defer handler.Orderbooks.mu.Unlock()

	for _, param := range params {
		symbol := param.Symbol

		_, exists := handler.Orderbooks.Symbols[symbol]
		if !exists {
			continue
		}

		delete(handler.Orderbooks.Symbols, symbol)
	}
}

func (handler *FuturesWS_ManagedOrderBook_Handler) Unsubscribe(params ...FuturesWS_DiffBookDepth_Params) (resp *FuturesWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	resp, hasTimedOut, err = handler.DiffBookDepth_Socket.Unsubscribe(params...)
	if err != nil {
		return resp, hasTimedOut, err
	}

	handler.removeSymbols(params...)

	return resp, false, nil
}

func (handler *FuturesWS_ManagedOrderBook_Handler) Subscribe(params ...FuturesWS_DiffBookDepth_Params) (resp *FuturesWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	resp, hasTimedOut, err = handler.DiffBookDepth_Socket.Subscribe(params...)
	if err != nil {
		return resp, hasTimedOut, err
	}

	handler.openNewSymbols(params...)

	return resp, false, nil
}

func (futures_ws *Futures_Websockets) Managed_OrderBook(publicOnMessage func(orderBook *Futures_ManagedOrderbook), params ...FuturesWS_DiffBookDepth_Params) (*FuturesWS_ManagedOrderBook_Handler, *Error) {
	handler := &FuturesWS_ManagedOrderBook_Handler{}

	var newSocket FuturesWS_DiffBookDepth_Socket
	streamNames := newSocket.CreateStreamName(params...)
	socket, err := futures_ws.CreateSocket(streamNames, false)
	if err != nil {
		return nil, err
	}
	newSocket.Handler = socket

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var diffBookDepth *FuturesWS_DiffBookDepth
		err := json.Unmarshal(msg, &diffBookDepth)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}

		shouldPushEvent, managedOrderbook := handler.handleWSMessage(futures_ws, diffBookDepth)
		if shouldPushEvent {
			publicOnMessage(managedOrderbook)
		}
	}

	handler.DiffBookDepth_Socket = &newSocket
	handler.Orderbooks.Symbols = make(map[string]struct {
		bufferedEvents []*FuturesWS_DiffBookDepth
		Orderbook      *Futures_ManagedOrderbook
	})
	handler.openNewSymbols(params...)

	return handler, nil
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
			LOG_WS_ERRORS("[PRIVATEMESSAGEVALIDATOR ERR] WS Message is the following:", string(msg))
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
	data, timedOut, err := futures_ws.Websocket.SendRequest_sync(requestObj, timeout_sec...)
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
