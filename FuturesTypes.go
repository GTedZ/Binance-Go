package Binance

import (
	"fmt"
	"math"
	"strconv"
	"sync"
)

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

// # Checks if the price passes the "PRICE_FILTER"
//
// "reason" is returned on any failure, possible values are:
//
// - "minPrice" if the price < minPrice. 		"suggestion" will be returned with the value "minPrice".
//
// - "maxPrice" if the price > maxPrice. 		"suggestion" will be returned with the value "maxPrice".
//
// - "tickSize" if the price % tickSize != 0. 	"suggestion" will be returned with the corrected value.
//
// "suggestion" must be ignored if it is returned as 0.
// "suggestion" is always returned as "price" if it passes the filter.
func (futuresSymbol *Futures_Symbol) PRICE_FILTER(price float64) (isValid bool, reason string, suggestion float64, err *Error) {

	if futuresSymbol.Filters.PRICE_FILTER == nil {
		return true, "", price, nil
	}

	minPrice, parseErr := strconv.ParseFloat(futuresSymbol.Filters.PRICE_FILTER.MinPrice, 64)
	if parseErr != nil {
		return false, "", 0, LocalError(PARSING_ERR, parseErr.Error())
	}
	maxPrice, parseErr := strconv.ParseFloat(futuresSymbol.Filters.PRICE_FILTER.MaxPrice, 64)
	if parseErr != nil {
		return false, "", 0, LocalError(PARSING_ERR, parseErr.Error())
	}
	tickSize, parseErr := strconv.ParseFloat(futuresSymbol.Filters.PRICE_FILTER.TickSize, 64)
	if parseErr != nil {
		return false, "", 0, LocalError(PARSING_ERR, parseErr.Error())
	}

	if minPrice != 0 && price < minPrice {
		return false, "minPrice", minPrice, nil
	}

	if maxPrice != 0 && price > maxPrice {
		return false, "maxPrice", maxPrice, nil
	}

	if tickSize != 0 && math.Remainder(price, tickSize) != 0 {
		suggestion, parseErr := strconv.ParseFloat(Format_TickSize_str(fmt.Sprint(price), futuresSymbol.Filters.PRICE_FILTER.TickSize), 64)
		if parseErr != nil {
			return false, "", 0, LocalError(PARSING_ERR, parseErr.Error())
		}

		return false, "tickSize", suggestion, nil
	}

	return true, "", price, nil
}

// # Checks if the price passes the "PRICE_FILTER"
func (futuresSymbol *Futures_Symbol) PRICE_FILTER_COMPACT(price float64) (isValid bool, err *Error) {
	isValid, _, _, err = futuresSymbol.PRICE_FILTER(price)
	return isValid, err
}

// # Checks if the quantity passes the "LOT_SIZE"
//
// "reason" is returned on any failure, possible values are:
//
// - "minQty" if the quantity < minQty. "suggestion" will be returned with the value "minQty".
//
// - "maxQty" if the quantity > maxQty. "suggestion" will be returned with the value "maxQty".
//
// - "stepSize" if the quantity % stepSize != 0. "suggestion" will be returned with the corrected value.
//
// "suggestion" must be ignored if it is returned as 0.
// "suggestion" is always returned as "quantity" if it passes the filter.
func (futuresSymbol *Futures_Symbol) LOT_SIZE(quantity float64) (isValid bool, reason string, suggestion float64, err *Error) {

	if futuresSymbol.Filters.LOT_SIZE == nil {
		return true, "", quantity, nil
	}

	minQty, parseErr := strconv.ParseFloat(futuresSymbol.Filters.LOT_SIZE.MinQty, 64)
	if parseErr != nil {
		return false, "", 0, LocalError(PARSING_ERR, parseErr.Error())
	}
	maxQty, parseErr := strconv.ParseFloat(futuresSymbol.Filters.LOT_SIZE.MaxQty, 64)
	if parseErr != nil {
		return false, "", 0, LocalError(PARSING_ERR, parseErr.Error())
	}
	stepSize, parseErr := strconv.ParseFloat(futuresSymbol.Filters.LOT_SIZE.StepSize, 64)
	if parseErr != nil {
		return false, "", 0, LocalError(PARSING_ERR, parseErr.Error())
	}

	if minQty != 0 && quantity < minQty {
		return false, "minQty", minQty, nil
	}

	if maxQty != 0 && quantity > maxQty {
		return false, "maxQty", maxQty, nil
	}

	if stepSize != 0 && math.Remainder(quantity, stepSize) != 0 {
		suggestion, parseErr := strconv.ParseFloat(Format_TickSize_str(fmt.Sprint(quantity), futuresSymbol.Filters.LOT_SIZE.StepSize), 64)
		if parseErr != nil {
			return false, "", 0, LocalError(PARSING_ERR, parseErr.Error())
		}

		return false, "stepSize", suggestion, nil
	}

	return true, "", quantity, nil
}

// # Checks if the price passes the "LOT_SIZE"
func (futuresSymbol *Futures_Symbol) LOT_SIZE_COMPACT(price float64) (isValid bool, err *Error) {
	isValid, _, _, err = futuresSymbol.LOT_SIZE(price)
	return isValid, err
}

// # Checks if the quantity passes the "LOT_SIZE"
//
// "reason" is returned on any failure, possible values are:
//
// - "minQty" if the quantity < minQty. "suggestion" will be returned with the value "minQty".
//
// - "maxQty" if the quantity > maxQty. "suggestion" will be returned with the value "maxQty".
//
// - "stepSize" if the quantity % stepSize != 0. "suggestion" will be returned with the corrected value.
//
// "suggestion" must be ignored if it is returned as 0.
// "suggestion" is always returned as "quantity" if it passes the filter.
func (futuresSymbol *Futures_Symbol) MARKET_LOT_SIZE(quantity float64) (isValid bool, reason string, suggestion float64, err *Error) {

	if futuresSymbol.Filters.LOT_SIZE == nil {
		return true, "", quantity, nil
	}

	minQty, parseErr := strconv.ParseFloat(futuresSymbol.Filters.MARKET_LOT_SIZE.MinQty, 64)
	if parseErr != nil {
		return false, "", 0, LocalError(PARSING_ERR, parseErr.Error())
	}
	maxQty, parseErr := strconv.ParseFloat(futuresSymbol.Filters.MARKET_LOT_SIZE.MaxQty, 64)
	if parseErr != nil {
		return false, "", 0, LocalError(PARSING_ERR, parseErr.Error())
	}
	stepSize, parseErr := strconv.ParseFloat(futuresSymbol.Filters.MARKET_LOT_SIZE.StepSize, 64)
	if parseErr != nil {
		return false, "", 0, LocalError(PARSING_ERR, parseErr.Error())
	}

	if minQty != 0 && quantity < minQty {
		return false, "minQty", minQty, nil
	}

	if maxQty != 0 && quantity > maxQty {
		return false, "maxQty", maxQty, nil
	}

	if stepSize != 0 && math.Remainder(quantity, stepSize) != 0 {
		suggestion, parseErr := strconv.ParseFloat(Format_TickSize_str(fmt.Sprint(quantity), futuresSymbol.Filters.MARKET_LOT_SIZE.StepSize), 64)
		if parseErr != nil {
			return false, "", 0, LocalError(PARSING_ERR, parseErr.Error())
		}

		return false, "stepSize", suggestion, nil
	}

	return true, "", quantity, nil
}

// # Checks if the price passes the "MARKET_LOT_SIZE"
func (futuresSymbol *Futures_Symbol) MARKET_LOT_SIZE_COMPACT(price float64) (isValid bool, err *Error) {
	isValid, _, _, err = futuresSymbol.LOT_SIZE(price)
	return isValid, err
}

// # Truncates a price string to the last significant digit
//
// Symbol Filters rule "LOT_SIZE" defines the highest precision the symbol's Quantity (via base asset) accepts
// And for MARKET orders the "MARKET_LOT_SIZE" also applies
// i.e: BTCUSDT has a precision of 5, meaning if you want to buy "0.12345678" BTC,
// it would be truncated down to "0.12345" BTC
func (futuresSymbol *Futures_Symbol) TruncQuantity_float64(quantity float64, IsForMarketOrder bool) string {
	return futuresSymbol.TruncQuantity(fmt.Sprint(quantity), IsForMarketOrder)
}

func (futuresSymbol *Futures_Symbol) TruncQuantity(quantity string, IsForMarketOrder bool) string {
	truncQuantity := quantity
	if futuresSymbol.Filters.LOT_SIZE != nil && futuresSymbol.Filters.LOT_SIZE.StepSize != "" {
		truncQuantity = Format_TickSize_str(truncQuantity, futuresSymbol.Filters.LOT_SIZE.StepSize)
	}

	if IsForMarketOrder && futuresSymbol.Filters.MARKET_LOT_SIZE != nil && futuresSymbol.Filters.MARKET_LOT_SIZE.StepSize != "" {
		truncQuantity = Format_TickSize_str(truncQuantity, futuresSymbol.Filters.MARKET_LOT_SIZE.StepSize)
	}

	return truncQuantity
}

// # Truncates a price string to the last significant digit
//
// Symbol Filters rule "PRICE_FILTER" defines the highest precision the symbol accepts
// i.e: BTCUSDT has a precision of 2, meaning if you want to buy BTCUSDT at "123_456.7891",
// it would be truncated down to "123_456.78"
func (futuresSymbol *Futures_Symbol) TruncPrice_float64(price float64) string {
	return futuresSymbol.TruncPrice(fmt.Sprint(price))
}

// # Truncates a price string to the last significant digit
//
// Symbol Filters rule "PRICE_FILTER" defines the highest precision the symbol accepts
// i.e: BTCUSDT has a precision of 2, meaning if you want to buy BTCUSDT at "123_456.7891",
// it would be truncated down to "123_456.78"
func (futuresSymbol *Futures_Symbol) TruncPrice(priceStr string) string {
	if futuresSymbol.Filters.PRICE_FILTER == nil || futuresSymbol.Filters.PRICE_FILTER.TickSize == "" {
		return priceStr
	}

	return Format_TickSize_str(priceStr, futuresSymbol.Filters.PRICE_FILTER.TickSize)
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

	Assets struct {
		Mu  sync.Mutex
		Map map[string]*Futures_Asset
	}
	Symbols struct {
		Mu  sync.Mutex
		Map map[string]*Futures_Symbol
	}
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

type Futures_FundingRate struct {
	Symbol      string `json:"symbol"`
	FundingRate string `json:"fundingRate"`
	FundingTime int64  `json:"fundingTime"`
	MarkPrice   string `json:"markPrice"`
}

type Futures_24hTicker struct {
	Symbol             string `json:"symbol"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	WeightedAvgPrice   string `json:"weightedAvgPrice"`
	LastPrice          string `json:"lastPrice"`
	LastQty            string `json:"lastQty"`
	Open               string `json:"openPrice"`
	High               string `json:"highPrice"`
	Low                string `json:"lowPrice"`
	Volume             string `json:"volume"`
	QuoteVolume        string `json:"quoteVolume"`
	OpenTime           int64  `json:"openTime"`
	CloseTime          int64  `json:"closeTime"`
	FirstId            int64  `json:"firstId"`
	LastId             int64  `json:"lastId"`
	Count              int64  `json:"count"`
}

type Futures_PriceTicker struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
	Time   int64  `json:"time"`
}

type Futures_BookTicker struct {
	Symbol   string `json:"symbol"`
	BidPrice string `json:"bidPrice"`
	BidQty   string `json:"bidQty"`
	AskPrice string `json:"askPrice"`
	AskQty   string `json:"askQty"`
	Time     int64  `json:"time"`
}

type Futures_DeliveryPrice struct {
	DeliveryTime  int64 `json:"deliveryTime"`
	DeliveryPrice int64 `json:"deliveryPrice"`
}

type Futures_OpenInterest struct {
	Symbol       string `json:"symbol"`
	OpenInterest string `json:"openInterest"`
	Time         int64  `json:"time"`
}

type Futures_OpenInterestStatistics struct {
	Symbol               string `json:"symbol"`
	SumOpenInterest      string `json:"sumOpenInterest"`
	SumOpenInterestValue string `json:"sumOpenInterestValue"`
	Timestamp            string `json:"timestamp"`
}

//////////////////////////////////////
//////////////////////////////////////
//////////////////////////////////////

type Futures_Order struct {
	ClientOrderId string `json:"clientOrderId"`

	CumQty string `json:"cumQty"`

	CumQuote string `json:"cumQuote"`

	ExecutedQty string `json:"executedQty"`

	OrderId int64 `json:"orderId"`

	AvgPrice string `json:"avgPrice"`

	OrigQty string `json:"origQty"`

	Price string `json:"price"`

	ReduceOnly bool `json:"reduceOnly"`

	Side string `json:"side"`

	PositionSide string `json:"positionSide"`

	Status string `json:"status"`

	// please ignore when order type is "TRAILING_STOP_MARKET"
	StopPrice string `json:"stopPrice"`

	// if Close-All
	ClosePosition bool `json:"closePosition"`

	Symbol string `json:"symbol"`

	TimeInForce string `json:"timeInForce"`

	Type string `json:"type"`

	OrigType string `json:"origType"`

	// activation price, only return with "TRAILING_STOP_MARKET" order
	ActivatePrice string `json:"activatePrice"`

	// callback rate, only return with "TRAILING_STOP_MARKET" order
	PriceRate string `json:"priceRate"`

	UpdateTime int64 `json:"updateTime"`

	WorkingType string `json:"workingType"`

	// if conditional order trigger is protected
	PriceProtect bool `json:"priceProtect"`

	// price match mode
	PriceMatch string `json:"priceMatch"`

	// self trading preventation mode
	SelfTradePreventionMode string `json:"selfTradePreventionMode"`

	// order pre-set auto cancel time for "TIF" "GTD" order
	GoodTillDate int64 `json:"goodTillDate"`
}

type Futures_ChangeMarginType_Response struct {
	// 200 for success
	Code int `json:"code"`
	// "success"
	Msg string `json:"msg"`
}

func (response *Futures_ChangeMarginType_Response) IsAlreadyChanged(err *Error) bool {
	return !err.IsLocalError &&
		(err.Code == -4046 || err.Message == "No need to change margin type.")

}

type Futures_ChangePositionMode_Response struct {
	// 200 for success
	Code int `json:"code"`
	// "success"
	Msg string `json:"msg"`
}

func (*Futures_ChangePositionMode_Response) IsAlreadyChanged(err *Error) bool {
	return !err.IsLocalError &&
		(err.Code == -4059 || err.Message == "No need to change position side.")
}

type Futures_ChangeInitialLeverage_Response struct {
	Symbol           string `json:"symbol"`
	Leverage         int64  `json:"leverage"`
	MaxNotionalValue string `json:"maxNotionalValue"`
}

type Futures_ChangeMultiAssetsMode_Response struct {
	// 200 for success
	Code int `json:"code"`
	// "success"
	Msg string `json:"msg"`
}

func (*Futures_ChangeMultiAssetsMode_Response) IsAlreadyChanged(err *Error) bool {
	return !err.IsLocalError &&
		(err.Code == -4167 || err.Message == "Unable to adjust to Multi-Assets mode with symbols of USDⓈ-M Futures under isolated-margin mode.")
}

//////////////////////////////////////
//////////////////////////////////////
//////////////////////////////////////

type Futures_AccountInfo struct {
	TotalInitialMargin          string
	TotalMaintMargin            string
	TotalWalletBalance          string
	TotalUnrealizedProfit       string
	TotalMarginBalance          string
	TotalPositionInitialMargin  string
	TotalOpenOrderInitialMargin string
	TotalCrossWalletBalance     string
	TotalCrossUnPnl             string
	AvailableBalance            string
	MaxWithdrawAmount           string
	Assets                      []*Futures_AccountInfo_Asset
	Positions                   []*Futures_AccountInfo_Position
}

type Futures_AccountInfo_Asset struct {
	Asset                  string `json:"asset"`
	WalletBalance          string `json:"walletBalance"`
	UnrealizedProfit       string `json:"unrealizedProfit"`
	MarginBalance          string `json:"marginBalance"`
	MaintMargin            string `json:"maintMargin"`
	InitialMargin          string `json:"initialMargin"`
	PositionInitialMargin  string `json:"positionInitialMargin"`
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
	CrossWalletBalance     string `json:"crossWalletBalance"`
	CrossUnPnl             string `json:"crossUnPnl"`
	AvailableBalance       string `json:"availableBalance"`
	MaxWithdrawAmount      string `json:"maxWithdrawAmount"`
	UpdateTime             int64  `json:"updateTime"`
}

type Futures_AccountInfo_Position struct {
	Symbol           string `json:"symbol"`
	PositionSide     string `json:"positionSide"`
	PositionAmt      string `json:"positionAmt"`
	UnrealizedProfit string `json:"unrealizedProfit"`
	IsolatedMargin   string `json:"isolatedMargin"`
	Notional         string `json:"notional"`
	IsolatedWallet   string `json:"isolatedWallet"`
	InitialMargin    string `json:"initialMargin"`
	MaintMargin      string `json:"maintMargin"`
	UpdateTime       int64  `json:"updateTime"`
}

type Futures_UserCommissionRate struct {
	Symbol              string `json:"symbol"`
	MakerCommissionRate string `json:"makerCommissionRate"`
	TakerCommissionRate string `json:"takerCommissionRate"`
}

type Futures_AccountConfiguration struct {
	FeeTier           int64 `json:"feeTier"`
	CanTrade          bool  `json:"canTrade"`
	CanDeposit        bool  `json:"canDeposit"`
	CanWithdraw       bool  `json:"canWithdraw"`
	DualSidePosition  bool  `json:"dualSidePosition"`
	UpdateTime        int64 `json:"updateTime"`
	MultiAssetsMargin bool  `json:"multiAssetsMargin"`
	TradeGroupId      int64 `json:"tradeGroupId"`
}

type Futures_LeverageBrackets struct {
	Symbol string `json:"symbol"`

	// user symbol bracket multiplier, only appears when user's symbol bracket is adjusted
	NotionalCoef float64 `json:"notionalCoef"`

	Brackets []*Futures_LeverageBrackets_Bracket `json:"brackets"`
}
type Futures_LeverageBrackets_Bracket struct {

	// Notional bracket
	Bracket int64 `json:"bracket"`

	// Max initial leverage for this bracket
	InitialLeverage int64 `json:"initialLeverage"`

	// Cap notional of this bracket
	NotionalCap int64 `json:"notionalCap"`

	// Notional threshold of this bracket
	NotionalFloor int64 `json:"notionalFloor"`

	// Maintenance ratio for this bracket
	MaintMarginRatio float64 `json:"maintMarginRatio"`

	// Auxiliary number for quick calculation
	Cum float64 `json:"cum"`
}
