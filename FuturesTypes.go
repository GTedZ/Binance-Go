package Binance

import "fmt"

var FUTURES_Constants = struct {
	URLs [1]string

	SecurityTypes Futures_SecurityTypes_ENUM

	SymbolTypes      Futures_SymbolTypes_ENUM
	ContractTypes    Futures_ContractTypes_ENUM
	ContractStatuses Futures_ContractStatuses_ENUM

	OrderStatuses Futures_OrderStatuses_ENUM
	OrderTypes    Futures_OrderTypes_ENUM
	OrderSides    Futures_OrderSides_ENUM

	PositionSides Futures_PositionSides_ENUM

	TimeInForce  Futures_TimeInForce_ENUM
	WorkingTypes Futures_WorkingTypes_ENUM

	NewOrderRespTypes Futures_NewOrderRespTypes_ENUM

	ChartIntervals Futures_ChartIntervals_ENUM

	STPModes   Futures_STPModes_ENUM
	PriceMatch Futures_PriceMatch_ENUM

	SymbolFilterTypes  FUTURES_Symbol_FilterTypes_ENUM
	RateLimitTypes     Futures_RateLimitTypes_ENUM
	RateLimitIntervals Futures_RateLimitIntervals_ENUM

	Websocket Futures_Websocket_Constants
}{
	URLs: [1]string{"https://fapi.binance.com"},
	SecurityTypes: Futures_SecurityTypes_ENUM{
		NONE:        "NONE",
		MARKET_DATA: "MARKET_DATA",
		USER_STREAM: "USER_STREAM",
		TRADE:       "TRADE",
		USER_DATA:   "USER_DATA",
	},
	SymbolTypes: Futures_SymbolTypes_ENUM{
		FUTURE: "FUTURE",
	},
	ContractTypes: Futures_ContractTypes_ENUM{
		PERPETUAL:            "PERPETUAL",
		CURRENT_MONTH:        "CURRENT_MONTH",
		NEXT_MONTH:           "NEXT_MONTH",
		CURRENT_QUARTER:      "CURRENT_QUARTER",
		NEXT_QUARTER:         "NEXT_QUARTER",
		PERPETUAL_DELIVERING: "PERPETUAL_DELIVERING",
	},
	ContractStatuses: Futures_ContractStatuses_ENUM{
		PENDING_TRADING: "PENDING_TRADING",
		TRADING:         "TRADING",
		PRE_DELIVERING:  "PRE_DELIVERING",
		DELIVERING:      "DELIVERING",
		DELIVERED:       "DELIVERED",
		PRE_SETTLE:      "PRE_SETTLE",
		SETTLING:        "SETTLING",
		CLOSE:           "CLOSE",
	},
	OrderStatuses: Futures_OrderStatuses_ENUM{
		NEW:              "NEW",
		PARTIALLY_FILLED: "PARTIALLY_FILLED",
		FILLED:           "FILLED",
		CANCELED:         "CANCELED",
		REJECTED:         "REJECTED",
		EXPIRED:          "EXPIRED",
		EXPIRED_IN_MATCH: "EXPIRED_IN_MATCH",
	},
	OrderTypes: Futures_OrderTypes_ENUM{
		LIMIT:                "LIMIT",
		MARKET:               "MARKET",
		STOP:                 "STOP",
		STOP_MARKET:          "STOP_MARKET",
		TAKE_PROFIT:          "TAKE_PROFIT",
		TAKE_PROFIT_MARKET:   "TAKE_PROFIT_MARKET",
		TRAILING_STOP_MARKET: "TRAILING_STOP_MARKET",
	},
	OrderSides: Futures_OrderSides_ENUM{
		BUY:  "BUY",
		SELL: "SELL",
	},
	PositionSides: Futures_PositionSides_ENUM{
		BOTH:  "BOTH",
		LONG:  "LONG",
		SHORT: "SHORT",
	},
	TimeInForce: Futures_TimeInForce_ENUM{
		GTC: "GTC",
		IOC: "IOC",
		FOK: "FOK",
	},
	WorkingTypes: Futures_WorkingTypes_ENUM{
		MARK_PRICE:     "MARK_PRICE",
		CONTRACT_PRICE: "CONTRACT_PRICE",
	},
	NewOrderRespTypes: Futures_NewOrderRespTypes_ENUM{
		ACK:    "ACK",
		RESULT: "RESULT",
	},
	ChartIntervals: Futures_ChartIntervals_ENUM{
		MIN:      "1m",
		MINS_3:   "3m",
		MINS_5:   "5m",
		MINS_15:  "15m",
		MINS_30:  "30m",
		HOUR:     "1h",
		HOURS_2:  "2h",
		HOURS_4:  "4h",
		HOURS_6:  "6h",
		HOURS_8:  "8h",
		HOURS_12: "12h",
		DAY:      "1d",
		DAYS_3:   "3d",
		WEEK:     "1w",
		MONTH:    "1M",
	},
	STPModes: Futures_STPModes_ENUM{
		NONE:         "NONE",
		EXPIRE_TAKER: "EXPIRE_TAKER",
		EXPIRE_BOTH:  "EXPIRE_BOTH",
		EXPIRE_MAKER: "EXPIRE_MAKER",
	},
	PriceMatch: Futures_PriceMatch_ENUM{
		NONE:        "NONE",
		OPPONENT:    "OPPONENT",
		OPPONENT_5:  "OPPONENT_5",
		OPPONENT_10: "OPPONENT_10",
		OPPONENT_20: "OPPONENT_20",
		QUEUE:       "QUEUE",
		QUEUE_5:     "QUEUE_5",
		QUEUE_10:    "QUEUE_10",
		QUEUE_20:    "QUEUE_20",
	},
	SymbolFilterTypes: FUTURES_Symbol_FilterTypes_ENUM{
		PRICE_FILTER:        "PRICE_FILTER",
		LOT_SIZE:            "LOT_SIZE",
		MARKET_LOT_SIZE:     "MARKET_LOT_SIZE",
		MAX_NUM_ORDERS:      "MAX_NUM_ORDERS",
		MAX_NUM_ALGO_ORDERS: "MAX_NUM_ALGO_ORDERS",
		PERCENT_PRICE:       "PERCENT_PRICE",
		MIN_NOTIONAL:        "MIN_NOTIONAL",
	},
	RateLimitTypes: Futures_RateLimitTypes_ENUM{
		REQUEST_WEIGHT: "REQUEST_WEIGHT",
		ORDERS:         "ORDERS",
	},
	RateLimitIntervals: Futures_RateLimitIntervals_ENUM{
		SECOND: "SECOND",
		MINUTE: "MINUTE",
		DAY:    "DAY",
	},
	Websocket: Futures_Websocket_Constants{
		URLs: []string{"wss://fstream.binance.com"},
	},
}

////////////////////////////////////////////////////////////////////////////////////////////////////////// Declarations
//////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////// Definitions

type Futures_SecurityTypes_ENUM struct {
	NONE        string
	MARKET_DATA string
	USER_STREAM string
	TRADE       string
	USER_DATA   string
}

type Futures_SymbolTypes_ENUM struct {
	FUTURE string
}

type Futures_ContractTypes_ENUM struct {
	PERPETUAL            string
	CURRENT_MONTH        string
	NEXT_MONTH           string
	CURRENT_QUARTER      string
	NEXT_QUARTER         string
	PERPETUAL_DELIVERING string
}

type Futures_ContractStatuses_ENUM struct {
	PENDING_TRADING string
	TRADING         string
	PRE_DELIVERING  string
	DELIVERING      string
	DELIVERED       string
	PRE_SETTLE      string
	SETTLING        string
	CLOSE           string
}

type Futures_OrderStatuses_ENUM struct {
	NEW              string
	PARTIALLY_FILLED string
	FILLED           string
	CANCELED         string
	REJECTED         string
	EXPIRED          string
	EXPIRED_IN_MATCH string
}

type Futures_OrderTypes_ENUM struct {
	LIMIT                string
	MARKET               string
	STOP                 string
	STOP_MARKET          string
	TAKE_PROFIT          string
	TAKE_PROFIT_MARKET   string
	TRAILING_STOP_MARKET string
}

type Futures_OrderSides_ENUM struct {
	BUY  string
	SELL string
}

type Futures_PositionSides_ENUM struct {
	BOTH  string
	LONG  string
	SHORT string
}

type Futures_TimeInForce_ENUM struct {
	GTC string
	IOC string
	FOK string
}

type Futures_WorkingTypes_ENUM struct {
	MARK_PRICE     string
	CONTRACT_PRICE string
}

type Futures_NewOrderRespTypes_ENUM struct {
	ACK    string
	RESULT string
}

type Futures_ChartIntervals_ENUM struct {
	MIN      string
	MINS_3   string
	MINS_5   string
	MINS_15  string
	MINS_30  string
	HOUR     string
	HOURS_2  string
	HOURS_4  string
	HOURS_6  string
	HOURS_8  string
	HOURS_12 string
	DAY      string
	DAYS_3   string
	WEEK     string
	MONTH    string
}

type Futures_STPModes_ENUM struct {
	NONE         string
	EXPIRE_TAKER string
	EXPIRE_BOTH  string
	EXPIRE_MAKER string
}

type Futures_PriceMatch_ENUM struct {
	NONE        string
	OPPONENT    string
	OPPONENT_5  string
	OPPONENT_10 string
	OPPONENT_20 string
	QUEUE       string
	QUEUE_5     string
	QUEUE_10    string
	QUEUE_20    string
}

type FUTURES_Symbol_FilterTypes_ENUM struct {
	PRICE_FILTER        string
	LOT_SIZE            string
	MARKET_LOT_SIZE     string
	MAX_NUM_ORDERS      string
	MAX_NUM_ALGO_ORDERS string
	PERCENT_PRICE       string
	MIN_NOTIONAL        string
}

type Futures_RateLimitTypes_ENUM struct {
	REQUEST_WEIGHT string
	ORDERS         string
}

type Futures_RateLimitIntervals_ENUM struct {
	SECOND string
	MINUTE string
	DAY    string
}

type Futures_Websocket_Constants struct {
	URLs []string
}

type Futures_RateLimitType struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
}

////////////////////////////////////////////////////////////////////////////////////////////////////////// Definitions
//////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////// Response Types

type Futures_Asset struct {
	Asset string `json:"asset"`
	// whether the asset can be used as margin in Multi-Assets mode
	MarginAvailable bool `json:"marginAvailable"`
	// auto-exchange threshold in Multi-Assets margin mode
	AutoAssetExchange string `json:"autoAssetExchange"`
}

type Futures_Symbol struct {
	Symbol       string `json:"symbol"`
	Pair         string `json:"pair"`
	ContractType string `json:"contractType"`
	DeliveryDate int64  `json:"deliveryDate"`
	OnboardDate  int64  `json:"onboardDate"`
	Status       string `json:"status"`
	// ignore
	MaintMarginPercent string `json:"maintMarginPercent"`
	// ignore
	RequiredMarginPercent string   `json:"requiredMarginPercent"`
	BaseAsset             string   `json:"baseAsset"`
	QuoteAsset            string   `json:"quoteAsset"`
	PricePrecision        int64    `json:"pricePrecision"`
	QuantityPrecision     int64    `json:"quantityPrecision"`
	BaseAssetPrecision    int64    `json:"baseAssetPrecision"`
	QuoteAssetPrecision   int64    `json:"quoteAssetPrecision"`
	UnderlyingType        string   `json:"underlyingType"`
	UnderlyingSubType     []string `json:"underlyingSubType"`
	SettlePlan            int64    `json:"settlePlan"`
	TriggerProtect        string   `json:"triggerProtect"`
	Filters               Futures_SymbolFilters
	OrderType             []string `json:"orderType"`
	TimeInForce           []string `json:"timeInForce"`
	LiquidationFee        string   `json:"liquidationFee"`
	MarketTakeBound       string   `json:"marketTakeBound"`
}

// # Truncates a price string to the last significant digit
//
// Symbol Filters rule "PRICE_FILTER" defines the highest precision the symbol accepts
// i.e: BTCUSDT has a precision of 2, meaning if you want to buy BTCUSDT at "123_456.7891",
// it would be truncated down to "123_456.78"
func (spotSymbol *Futures_Symbol) TruncPrice_float64(price float64) string {
	return spotSymbol.TruncPrice(fmt.Sprint(price))
}

// # Truncates a price string to the last significant digit
//
// Symbol Filters rule "PRICE_FILTER" defines the highest precision the symbol accepts
// i.e: BTCUSDT has a precision of 2, meaning if you want to buy BTCUSDT at "123_456.7891",
// it would be truncated down to "123_456.78"
func (spotSymbol *Futures_Symbol) TruncPrice(priceStr string) string {
	if spotSymbol.Filters.PRICE_FILTER == nil || spotSymbol.Filters.PRICE_FILTER.TickSize == "" {
		return priceStr
	}

	return ""
}

type Futures_SymbolFilters struct {
	PRICE_FILTER        *Futures_SymbolFilter_PRICE_FILTER
	LOT_SIZE            *Futures_SymbolFilter_LOT_SIZE
	MARKET_LOT_SIZE     *Futures_SymbolFilter_MARKET_LOT_SIZE
	MAX_NUM_ORDERS      *Futures_SymbolFilter_MAX_NUM_ORDERS
	MAX_NUM_ALGO_ORDERS *Futures_SymbolFilter_MAX_NUM_ALGO_ORDERS
	PERCENT_PRICE       *Futures_SymbolFilter_PERCENT_PRICE
	MIN_NOTIONAL        *Futures_SymbolFilter_MIN_NOTIONAL
}

type Futures_ExchangeInfo_SORS struct {
	BaseAsset string   `json:"baseAsset"`
	Symbols   []string `json:"symbols"`
}

type Futures_SymbolFilter_PRICE_FILTER struct {
	FilterType string `json:"filterType"`
	MinPrice   string `json:"minPrice"`
	MaxPrice   string `json:"maxPrice"`
	TickSize   string `json:"tickSize"`
}

type Futures_SymbolFilter_LOT_SIZE struct {
	FilterType string `json:"filterType"`
	MinQty     string `json:"minQty"`
	MaxQty     string `json:"maxQty"`
	StepSize   string `json:"stepSize"`
}

type Futures_SymbolFilter_MARKET_LOT_SIZE struct {
	FilterType string `json:"filterType"`
	MinQty     string `json:"minQty"`
	MaxQty     string `json:"maxQty"`
	StepSize   string `json:"stepSize"`
}

type Futures_SymbolFilter_MAX_NUM_ORDERS struct {
	FilterType string `json:"filterType"`
	Limit      int64  `json:"limit"`
}

type Futures_SymbolFilter_MAX_NUM_ALGO_ORDERS struct {
	FilterType string `json:"filterType"`
	Limit      int64  `json:"limit"`
}

type Futures_SymbolFilter_PERCENT_PRICE struct {
	FilterType        string `json:"filterType"`
	MultiplierUp      string `json:"multiplierUp"`
	MultiplierDown    string `json:"multiplierDown"`
	MultiplierDecimal string `json:"multiplierDecimal"`
}

type Futures_SymbolFilter_MIN_NOTIONAL struct {
	FilterType string `json:"filterType"`
	Notional   string `json:"notional"`
}

//

type Futures_Time struct {
	ServerTime int64 `json:"serverTime"`

	Latency int64
}

type Futures_ExchangeInfo struct {
	// Not used by binance
	ExchangeFilters any                      `json:"exchangeFilters"`
	RateLimits      []*Futures_RateLimitType `json:"rateLimits"`
	ServerTime      int64                    `json:"serverTime"`
	Assets_arr      []*Futures_Asset         `json:"assets"`
	Symbols_arr     []*Futures_Symbol        `json:"symbols"`
	Timezone        string                   `json:"timezone"`

	Assets  map[string]*Futures_Asset
	Symbols map[string]*Futures_Symbol
}

type Futures_OrderBook struct {
	LastUpdateId int64       `json:"lastUpdateId"`
	Time         int64       `json:"E"`
	TransactTime int64       `json:"T"`
	Bids         [][2]string `json:"bids"`
	Asks         [][2]string `json:"asks"`
}

type Futures_Trade struct {
	Id           int64  `json:"id"`
	Price        string `json:"price"`
	Qty          string `json:"qty"`
	QuoteQty     string `json:"quoteQty"`
	Timestamp    int64  `json:"time"`
	IsBuyerMaker bool   `json:"isBuyerMaker"`
}

type Futures_AggTrade struct {
	AggTradeId   int64  `json:"a"`
	Price        string `json:"p"`
	Qty          string `json:"q"`
	FirstTradeId int64  `json:"f"`
	LastTradeId  int64  `json:"l"`
	Timestamp    int64  `json:"T"`
	IsBuyerMaker bool   `json:"m"`
}

type Futures_Candlestick struct {
	// Kline open time
	OpenTime int64
	// Open price
	Open string
	// High price
	High string
	// Low price
	Low string
	// Close price
	Close string
	// Volume
	Volume string
	// Kline Close time
	CloseTime int64
	// Quote asset volume
	QuoteAssetVolume string
	// Number of trades
	TradeCount int64
	// Taker buy base asset volume
	TakerBuyBaseAssetVolume string
	// Taker buy quote asset volume
	TakerBuyQuoteAssetVolume string
	// Unused field, ignore.
	Unused string
}

type Futures_PriceCandlestick struct {
	// Kline open time
	OpenTime int64
	// Open price
	Open string
	// High price
	High string
	// Low price
	Low string
	// Close price
	Close string
	// Volume
	Ignore1 string
	// Kline Close time
	CloseTime int64
	// Quote asset volume
	Ignore2 string
	// Number of trades
	Ignore3 int64
	// Taker buy base asset volume
	Ignore4 string
	// Taker buy quote asset volume
	Ignore5 string
	// Unused field, ignore.
	Unused string
}

type Futures_MarkPrice struct {
	Symbol               string `json:"symbol"`
	MarkPrice            string `json:"markPrice"`
	IndexPrice           string `json:"indexPrice"`
	EstimatedSettlePrice string `json:"estimatedSettlePrice"`
	LastFundingRate      string `json:"lastFundingRate"`
	NextFundingTime      int64  `json:"nextFundingTime"`
	InterestRate         string `json:"interestRate"`
	Time                 int64  `json:"time"`
}
