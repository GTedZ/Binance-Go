package Binance

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
)

type Spot_Websockets struct {
	binance *Binance
}

type Spot_Websocket struct {
	Websocket *Websocket
	Conn      *ws.Conn
	// Host server's URL
	BaseURL string

	pendingRequests map[string]chan []byte
}

func (spot_ws *Spot_Websocket) Close() error {
	return spot_ws.Websocket.Close()
}

// Forcefully reconnects the socket
// Also makes it a reconnecting socket if it weren't before
// Useless, but there nonetheless...
func (spot_ws *Spot_Websocket) Reconnect() {
	spot_ws.Websocket.Reconnect()
}

func (spot_ws *Spot_Websocket) SetMessageListener(f func(messageType int, msg []byte)) {
	spot_ws.Websocket.OnMessage = f
}

func (spot_ws *Spot_Websocket) SetPingListener(f func(appData string)) {
	spot_ws.Websocket.OnPing = f
}

func (spot_ws *Spot_Websocket) SetPongListener(f func(appData string)) {
	spot_ws.Websocket.OnPong = f
}

// This is called when socket has been disconnected
// Called when the detected a disconnection and wants to reconnect afterwards
// Usually called right before the 'ReconnectingListener'
func (spot_ws *Spot_Websocket) SetDisconnectListener(f func(code int, text string)) {
	spot_ws.Websocket.OnDisconnect = f
}

// This is called when socket began reconnecting
func (spot_ws *Spot_Websocket) SetReconnectingListener(f func()) {
	spot_ws.Websocket.OnReconnecting = f
}

// This is called when the socket has successfully reconnected after a disconnection
func (spot_ws *Spot_Websocket) SetReconnectListener(f func()) {
	spot_ws.Websocket.OnReconnect = f
}

// This is called when the websocket closes indefinitely
// Meaning when you invoke the 'Close()' method
// Or any other way a websocket is set to never reconnect on a disconnection
func (spot_ws *Spot_Websocket) SetCloseListener(f func(code int, text string)) {
	spot_ws.Websocket.OnClose = f
}

type SpotWS_PrivateMessage struct {
	Id string `json:"id"`
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type SpotWS_AggTrade struct {
	Event     string `json:"e"`
	EventTime int64  `json:"E"`
	Symbol    string `json:"s"`
	// Trade ID
	AggTradeId int64  `json:"a"`
	Price      string `json:"p"`
	Quantity   string `json:"q"`
	// First Trade ID
	FirstTradeId int64 `json:"f"`
	// Last Trade ID
	LastTradeId int64 `json:"l"`
	// Trade time
	Timestamp int64 `json:"T"`
	// Is the buyer the market maker?
	IsMaker bool `json:"m"`
	// Ignore
	Ignore bool `json:"M"`
}

type SpotWS_AggTrade_Socket struct {
	Handler *Spot_Websocket
}

func (*SpotWS_AggTrade_Socket) CreateStreamName(symbol string) string {
	return strings.ToLower(symbol) + "@aggTrade"
}

func (socket *SpotWS_AggTrade_Socket) Subscribe(symbol ...string) (resp *SpotWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	for i := range symbol {
		symbol[i] = socket.CreateStreamName(symbol[i])
	}

	return socket.Handler.Subscribe(symbol...)
}
func (socket *SpotWS_AggTrade_Socket) Unsubscribe(symbol ...string) (resp *SpotWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	for i := range symbol {
		symbol[i] = socket.CreateStreamName(symbol[i])
	}
	return socket.Handler.Unsubscribe(symbol...)
}

func (spot_ws *Spot_Websockets) AggTrade(publicOnMessage func(aggTrade *SpotWS_AggTrade), symbol ...string) (*SpotWS_AggTrade_Socket, *Error) {
	var newSocket SpotWS_AggTrade_Socket
	for i := range symbol {
		symbol[i] = newSocket.CreateStreamName(symbol[i])
	}
	socket, err := spot_ws.CreateSocket(symbol, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var aggTrade SpotWS_AggTrade
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

type SpotWS_Trade struct {
	Event     string `json:"e"`
	EventTime int64  `json:"E"`
	Symbol    string `json:"s"`
	// Trade ID
	TradeID  int64  `json:"t"`
	Price    string `json:"p"`
	Quantity string `json:"q"`
	// Trade time
	Timestamp int64 `json:"T"`
	// Is the buyer the market maker?
	IsMaker bool `json:"m"`
	// Ignore
	Ignore bool `json:"M"`
}

type SpotWS_Trade_Socket struct {
	Handler *Spot_Websocket
}

func (*SpotWS_Trade_Socket) CreateStreamName(symbol string) string {
	return strings.ToLower(symbol) + "@trade"
}

func (socket *SpotWS_Trade_Socket) Subscribe(symbol ...string) (resp *SpotWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	for i := range symbol {
		symbol[i] = socket.CreateStreamName(symbol[i])
	}

	return socket.Handler.Subscribe(symbol...)
}
func (socket *SpotWS_Trade_Socket) Unsubscribe(symbol ...string) (resp *SpotWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	for i := range symbol {
		symbol[i] = socket.CreateStreamName(symbol[i])
	}
	return socket.Handler.Unsubscribe(symbol...)
}

func (spot_ws *Spot_Websockets) Trade(publicOnMessage func(aggTrade *SpotWS_Trade), symbol ...string) (*SpotWS_Trade_Socket, *Error) {
	var newSocket SpotWS_Trade_Socket
	for i := range symbol {
		symbol[i] = newSocket.CreateStreamName(symbol[i])
	}
	socket, err := spot_ws.CreateSocket(symbol, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var trade SpotWS_Trade
		err := json.Unmarshal(msg, &trade)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(&trade)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type SpotWS_Candlestick_MSG struct {

	// Event type
	Event string `json:"e"`

	// Event time
	EventTime int64 `json:"E"`

	// Symbol
	Symbol string `json:"s"`

	Candle *SpotWS_Candlestick `json:"k"`
}
type SpotWS_Candlestick struct {

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

type SpotWS_Candlestick_StreamIdentifier struct {
	Symbol   string
	Interval string
}

type SpotWS_Candlestick_Socket struct {
	Handler *Spot_Websocket
}

func (*SpotWS_Candlestick_Socket) CreateStreamName(symbol string, interval string) string {
	return strings.ToLower(symbol) + "@kline_" + interval
}

func (socket *SpotWS_Candlestick_Socket) Subscribe(identifiers ...SpotWS_Candlestick_StreamIdentifier) (resp *SpotWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	streams := make([]string, len(identifiers))
	for i, id := range identifiers {
		streams[i] = socket.CreateStreamName(id.Symbol, id.Interval)
	}

	return socket.Handler.Subscribe(streams...)
}
func (socket *SpotWS_Candlestick_Socket) Unsubscribe(identifiers ...SpotWS_Candlestick_StreamIdentifier) (resp *SpotWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	streams := make([]string, len(identifiers))
	for i, id := range identifiers {
		streams[i] = socket.CreateStreamName(id.Symbol, id.Interval)
	}

	return socket.Handler.Unsubscribe(streams...)
}

func (spot_ws *Spot_Websockets) Candlesticks(publicOnMessage func(candlestick_msg *SpotWS_Candlestick_MSG), identifiers ...SpotWS_Candlestick_StreamIdentifier) (*SpotWS_Candlestick_Socket, *Error) {
	var newSocket SpotWS_Candlestick_Socket

	streams := make([]string, len(identifiers))
	for i, id := range identifiers {
		streams[i] = newSocket.CreateStreamName(id.Symbol, id.Interval)
	}

	socket, err := spot_ws.CreateSocket(streams, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var candlestick_msg SpotWS_Candlestick_MSG
		err := json.Unmarshal(msg, &candlestick_msg)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(&candlestick_msg)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type SpotWS_Candlestick_TimezoneOffset_Socket struct {
	Handler *Spot_Websocket
}

func (*SpotWS_Candlestick_TimezoneOffset_Socket) CreateStreamName(symbol string, interval string) string {
	return strings.ToLower(symbol) + "@kline_" + interval + "@+08:00"
}

func (socket *SpotWS_Candlestick_TimezoneOffset_Socket) Subscribe(identifiers ...SpotWS_Candlestick_StreamIdentifier) (resp *SpotWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	streams := make([]string, len(identifiers))
	for i, id := range identifiers {
		streams[i] = socket.CreateStreamName(id.Symbol, id.Interval)
	}

	return socket.Handler.Subscribe(streams...)
}
func (socket *SpotWS_Candlestick_TimezoneOffset_Socket) Unsubscribe(identifiers ...SpotWS_Candlestick_StreamIdentifier) (resp *SpotWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	streams := make([]string, len(identifiers))
	for i, id := range identifiers {
		streams[i] = socket.CreateStreamName(id.Symbol, id.Interval)
	}

	return socket.Handler.Unsubscribe(streams...)
}

func (spot_ws *Spot_Websockets) Candlestick_WithOffset(publicOnMessage func(candlestick_msg *SpotWS_Candlestick_MSG), identifiers ...SpotWS_Candlestick_StreamIdentifier) (*SpotWS_Candlestick_TimezoneOffset_Socket, *Error) {
	var newSocket SpotWS_Candlestick_TimezoneOffset_Socket

	streams := make([]string, len(identifiers))
	for i, id := range identifiers {
		streams[i] = newSocket.CreateStreamName(id.Symbol, id.Interval)
	}
	socket, err := spot_ws.CreateSocket(streams, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var candlestick_msg SpotWS_Candlestick_MSG
		err := json.Unmarshal(msg, &candlestick_msg)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(&candlestick_msg)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type SpotWS_MiniTicker struct {

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

type SpotWS_MiniTicker_Socket struct {
	Handler *Spot_Websocket
}

func (*SpotWS_MiniTicker_Socket) CreateStreamName(symbol string) string {
	return strings.ToLower(symbol) + "@miniTicker"
}

func (socket *SpotWS_MiniTicker_Socket) Subscribe(symbol ...string) (resp *SpotWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	for i := range symbol {
		symbol[i] = socket.CreateStreamName(symbol[i])
	}

	return socket.Handler.Subscribe(symbol...)
}
func (socket *SpotWS_MiniTicker_Socket) Unsubscribe(symbol ...string) (resp *SpotWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	for i := range symbol {
		symbol[i] = socket.CreateStreamName(symbol[i])
	}

	return socket.Handler.Unsubscribe(symbol...)
}

func (spot_ws *Spot_Websockets) MiniTicker(publicOnMessage func(miniTicker *SpotWS_MiniTicker), symbol ...string) (*SpotWS_MiniTicker_Socket, *Error) {
	var newSocket SpotWS_MiniTicker_Socket
	for i := range symbol {
		symbol[i] = newSocket.CreateStreamName(symbol[i])
	}
	socket, err := spot_ws.CreateSocket(symbol, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var miniTicker SpotWS_MiniTicker
		err := json.Unmarshal(msg, &miniTicker)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(&miniTicker)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type SpotWS_AllMiniTickers_Socket struct {
	Handler *Spot_Websocket
}

func (*SpotWS_AllMiniTickers_Socket) CreateStreamName() string {
	return "!miniTicker@arr"
}

func (socket *SpotWS_AllMiniTickers_Socket) Subscribe() (resp *SpotWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	return socket.Handler.Subscribe(socket.CreateStreamName())
}
func (socket *SpotWS_AllMiniTickers_Socket) Unsubscribe() (resp *SpotWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	return socket.Handler.Unsubscribe(socket.CreateStreamName())
}

func (spot_ws *Spot_Websockets) AllMiniTickers(publicOnMessage func(miniTickers []*SpotWS_MiniTicker)) (*SpotWS_AllMiniTickers_Socket, *Error) {
	var newSocket SpotWS_AllMiniTickers_Socket
	socket, err := spot_ws.CreateSocket([]string{newSocket.CreateStreamName()}, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var miniTickers []*SpotWS_MiniTicker
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

type SpotWS_Ticker struct {

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

	// First trade(F)-1 price (first trade before the 24hr rolling window)
	PreviousFirstTradePrice string `json:"x"`

	// Last price
	LastPrice string `json:"c"`

	// Last quantity
	LastQty string `json:"Q"`

	// Best bid price
	Bid string `json:"b"`

	// Best bid quantity
	BidQty string `json:"B"`

	// Best ask price
	Ask string `json:"a"`

	// Best ask quantity
	AskQty string `json:"A"`

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

type SpotWS_Ticker_Socket struct {
	Handler *Spot_Websocket
}

func (*SpotWS_Ticker_Socket) CreateStreamName(symbol string) string {
	return strings.ToLower(symbol) + "@ticker"
}

func (socket *SpotWS_Ticker_Socket) Subscribe(symbol ...string) (resp *SpotWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	for i := range symbol {
		symbol[i] = socket.CreateStreamName(symbol[i])
	}

	return socket.Handler.Subscribe(symbol...)
}
func (socket *SpotWS_Ticker_Socket) Unsubscribe(symbol ...string) (resp *SpotWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	for i := range symbol {
		symbol[i] = socket.CreateStreamName(symbol[i])
	}

	return socket.Handler.Unsubscribe(symbol...)
}

func (spot_ws *Spot_Websockets) Ticker(publicOnMessage func(ticker *SpotWS_Ticker), symbol ...string) (*SpotWS_Ticker_Socket, *Error) {
	var newSocket SpotWS_Ticker_Socket
	for i := range symbol {
		symbol[i] = newSocket.CreateStreamName(symbol[i])
	}
	socket, err := spot_ws.CreateSocket(symbol, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var ticker SpotWS_Ticker
		err := json.Unmarshal(msg, &ticker)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(&ticker)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type SpotWS_AllTickers_Socket struct {
	Handler *Spot_Websocket
}

func (*SpotWS_AllTickers_Socket) CreateStreamName() string {
	return "!ticker@arr"
}

func (socket *SpotWS_AllTickers_Socket) Subscribe() (resp *SpotWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	return socket.Handler.Subscribe(socket.CreateStreamName())
}
func (socket *SpotWS_AllTickers_Socket) Unsubscribe() (resp *SpotWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	return socket.Handler.Unsubscribe(socket.CreateStreamName())
}

func (spot_ws *Spot_Websockets) AllTickers(publicOnMessage func(tickers []*SpotWS_Ticker)) (*SpotWS_AllTickers_Socket, *Error) {
	var newSocket SpotWS_AllTickers_Socket
	socket, err := spot_ws.CreateSocket([]string{newSocket.CreateStreamName()}, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var tickers []*SpotWS_Ticker
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

type SpotWS_RollingWindowStatistic struct {

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

	// Open price
	Open string `json:"o"`

	// High price
	High string `json:"h"`

	// Low price
	Low string `json:"l"`

	// Last price
	LastPrice string `json:"c"`

	// Weighted average price
	WeightedAveragePrice string `json:"w"`

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

type SpotWS_RollingWindowStatistics_StreamIdentifier struct {
	Symbol     string
	WindowSize string
}

type SpotWS_RollingWindowStatistics_Socket struct {
	Handler *Spot_Websocket
}

func (*SpotWS_RollingWindowStatistics_Socket) CreateStreamName(symbol string, windowSize string) string {
	return strings.ToLower(symbol) + "@ticker_" + windowSize
}

func (socket *SpotWS_RollingWindowStatistics_Socket) Subscribe(identifiers ...SpotWS_RollingWindowStatistics_StreamIdentifier) (resp *SpotWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	streams := make([]string, len(identifiers))
	for i, id := range identifiers {
		streams[i] = socket.CreateStreamName(id.Symbol, id.WindowSize)
	}

	return socket.Handler.Subscribe(streams...)
}
func (socket *SpotWS_RollingWindowStatistics_Socket) Unsubscribe(identifiers ...SpotWS_RollingWindowStatistics_StreamIdentifier) (resp *SpotWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	streams := make([]string, len(identifiers))
	for i, id := range identifiers {
		streams[i] = socket.CreateStreamName(id.Symbol, id.WindowSize)
	}

	return socket.Handler.Unsubscribe(streams...)
}

func (spot_ws *Spot_Websockets) RollingWindowStatistics(publicOnMessage func(rwStat *SpotWS_RollingWindowStatistic), identifiers ...SpotWS_RollingWindowStatistics_StreamIdentifier) (*SpotWS_RollingWindowStatistics_Socket, *Error) {
	var newSocket SpotWS_RollingWindowStatistics_Socket

	streams := make([]string, len(identifiers))
	for i, id := range identifiers {
		streams[i] = newSocket.CreateStreamName(id.Symbol, id.WindowSize)
	}

	socket, err := spot_ws.CreateSocket(streams, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var rwStat SpotWS_RollingWindowStatistic
		err := json.Unmarshal(msg, &rwStat)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(&rwStat)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type SpotWS_AllRollingWindowStatistics_Socket struct {
	Handler *Spot_Websocket
}

func (*SpotWS_AllRollingWindowStatistics_Socket) CreateStreamName(WindowSize string) string {
	return "!ticker_" + WindowSize + "@arr"
}

func (socket *SpotWS_AllRollingWindowStatistics_Socket) Subscribe(WindowSize ...string) (resp *SpotWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	for i := range WindowSize {
		WindowSize[i] = socket.CreateStreamName(WindowSize[i])
	}

	return socket.Handler.Subscribe(WindowSize...)
}
func (socket *SpotWS_AllRollingWindowStatistics_Socket) Unsubscribe(WindowSize ...string) (resp *SpotWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	for i := range WindowSize {
		WindowSize[i] = socket.CreateStreamName(WindowSize[i])
	}

	return socket.Handler.Unsubscribe(WindowSize...)
}

func (spot_ws *Spot_Websockets) AllRollingWindowStatistics(publicOnMessage func(rwStats []*SpotWS_RollingWindowStatistic), WindowSize ...string) (*SpotWS_AllRollingWindowStatistics_Socket, *Error) {
	var newSocket SpotWS_AllRollingWindowStatistics_Socket

	for i := range WindowSize {
		WindowSize[i] = newSocket.CreateStreamName(WindowSize[i])
	}

	socket, err := spot_ws.CreateSocket(WindowSize, false)
	if err != nil {
		return nil, err
	}

	socket.Websocket.OnMessage = func(messageType int, msg []byte) {
		var rwStats []*SpotWS_RollingWindowStatistic
		err := json.Unmarshal(msg, &rwStats)
		if err != nil {
			LocalError(PARSING_ERROR, err.Error())
			return
		}
		publicOnMessage(rwStats)
	}

	newSocket.Handler = socket
	return &newSocket, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (*Spot_Websockets) CreateSocket(streams []string, isCombined bool) (*Spot_Websocket, *Error) {
	baseURL := SPOT_Constants.Websocket.URLs[0]

	socket, err := CreateSocket(baseURL, streams, isCombined)
	if err != nil {
		return nil, err
	}

	socket.privateMessageValidator = func(msg []byte) (isPrivate bool, Id string) {

		if len(msg) > 0 && msg[0] == '[' {
			return false, ""
		}

		var privateMessage SpotWS_PrivateMessage
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

	ws := &Spot_Websocket{
		Websocket:       socket,
		Conn:            socket.Conn,
		BaseURL:         baseURL,
		pendingRequests: make(map[string]chan []byte),
	}

	return ws, nil
}

func createRequestObject() map[string]interface{} {
	requestObj := make(map[string]interface{})

	requestObj["id"] = uuid.New().String()
	return requestObj
}

type SpotWS_ListSubscriptions_Response struct {
	Id     string   `json:"id"`
	Result []string `json:"result"`
}

func (spot_ws *Spot_Websocket) ListSubscriptions(timeout_sec ...int) (resp *SpotWS_ListSubscriptions_Response, hasTimedOut bool, err *Error) {
	requestObj := createRequestObject()
	requestObj["method"] = "LIST_SUBSCRIPTIONS"
	data, timedOut, err := spot_ws.Websocket.SendRequest_sync(requestObj, timeout_sec[0])
	if err != nil || timedOut {
		return nil, timedOut, err
	}

	var response SpotWS_ListSubscriptions_Response
	unmarshallErr := json.Unmarshal(data, &response)
	if unmarshallErr != nil {
		return nil, false, LocalError(ERROR_PROCESSING_ERR, unmarshallErr.Error())
	}

	return &response, false, nil
}

type SpotWS_Subscribe_Response struct {
	Id string `json:"id"`
}

func (spot_ws *Spot_Websocket) Subscribe(stream ...string) (resp *SpotWS_Subscribe_Response, hasTimedOut bool, err *Error) {
	requestObj := createRequestObject()
	requestObj["method"] = "SUBSCRIBE"
	requestObj["params"] = stream
	data, timedOut, err := spot_ws.Websocket.SendRequest_sync(requestObj)
	if err != nil || timedOut {
		return nil, timedOut, err
	}

	var response SpotWS_Subscribe_Response
	unmarshallErr := json.Unmarshal(data, &response)
	if unmarshallErr != nil {
		return nil, false, LocalError(ERROR_PROCESSING_ERR, unmarshallErr.Error())
	}

	spot_ws.Websocket.Streams = append(spot_ws.Websocket.Streams, stream...)

	LOG_WS_VERBOSE("Successfully Subscribed to", stream)

	return &response, false, nil
}

type SpotWS_Unsubscribe_Response struct {
	Id string `json:"id"`
}

func (spot_ws *Spot_Websocket) Unsubscribe(stream ...string) (resp *SpotWS_Unsubscribe_Response, hasTimedOut bool, err *Error) {
	requestObj := createRequestObject()
	requestObj["method"] = "UNSUBSCRIBE"
	requestObj["params"] = stream
	data, timedOut, err := spot_ws.Websocket.SendRequest_sync(requestObj)
	if err != nil || timedOut {
		return nil, timedOut, err
	}

	var response SpotWS_Unsubscribe_Response
	unmarshallErr := json.Unmarshal(data, &response)
	if unmarshallErr != nil {
		return nil, false, LocalError(ERROR_PROCESSING_ERR, unmarshallErr.Error())
	}

	// Filter out the streams to remove from spot_ws.Websocket.Streams
	streamMap := make(map[string]bool)
	for _, s := range stream {
		streamMap[s] = true
	}

	var updatedStreams []string
	for _, existingStream := range spot_ws.Websocket.Streams {
		if !streamMap[existingStream] {
			updatedStreams = append(updatedStreams, existingStream)
		}
	}
	spot_ws.Websocket.Streams = updatedStreams

	LOG_WS_VERBOSE("Successfully Unsubscribed from", stream)

	return &response, false, nil
}
