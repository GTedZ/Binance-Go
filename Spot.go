package Binance

import (
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type Spot struct {
	binance       *Binance
	requestClient RequestClient
	baseURL       string

	API APIKEYS

	Websockets Spot_Websockets
}

func (spot *Spot) init(binance *Binance) {
	spot.binance = binance

	spot.requestClient.init(binance)
	spot.requestClient.Set_APIKEY(binance.API.KEY, binance.API.SECRET)
	spot.baseURL = SPOT_Constants.URLs[0]

	spot.API.Set(binance.API.KEY, binance.API.SECRET)

	spot.Websockets.binance = binance
}

/////////////////////////////////////////////////////////////////////////////////

// # Test connectivity to the Rest API.
//
// Weight: 1
//
// Data Source: Memory
func (spot *Spot) Ping() (latency int64, request *Response, err *Error) {
	startTime := time.Now().UnixMilli()
	httpResp, err := spot.makeRequest(&SpotRequest{
		securityType: SPOT_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/api/v3/ping",
	})
	diff := time.Now().UnixMilli() - startTime
	if err != nil {
		return diff, httpResp, err
	}

	return diff, httpResp, nil
}

/////////////////////////////////////////////////////////////////////////////////

// # Check server time
//
// Test connectivity to the Rest API and get the current server time.
//
// Weight: 1
//
// Data Source: Memory
func (spot *Spot) ServerTime() (*Spot_Time, *Response, *Error) {
	var spotTime Spot_Time

	startTime := time.Now().UnixMilli()
	httpResp, err := spot.makeRequest(&SpotRequest{
		securityType: SPOT_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/api/v3/time",
	})
	diff := time.Now().UnixMilli() - startTime
	spotTime.Latency = diff
	if err != nil {
		return &spotTime, httpResp, err
	}

	processingErr := json.Unmarshal(httpResp.Body, &spotTime)
	if processingErr != nil {
		return &spotTime, httpResp, LocalError(PARSING_ERR, processingErr.Error())
	}

	return &spotTime, httpResp, nil
}

/////////////////////////////////////////////////////////////////////////////////

//////// ExchangeInfo \\

type Spot_ExchangeInfo_Params struct {
	Symbol       string
	Symbols      []string
	Permissions  []string
	SymbolStatus string
	// The logic is flipped with "Dont Show" here
	// Because bools are always initialized as "false" while the exchange default is "true"
	DontShowPermissionSets bool
}

// # Exchange information
//
// Current exchange trading rules
// and symbol information
// with optional parameters
//
// Weight: 20
//
// usage:
//
// data, _, err := binance.Spot.ExchangeInfo_Params(&Spot_ExchangeInfo_Params{SymbolStatus: "TRADING", Permissions: []string{"SPOT"}})
//
// Parameters:
//
//	type Spot_ExchangeInfo_Params struct {
//		Symbol       string
//		Symbols      []string
//		Permissions  []string
//		SymbolStatus string
//		// The logic is flipped with "Dont Show" here
//		// Because bools are always initialize as "false" while the exchange default is "true"
//		DontShowPermissionSets bool
//	}
func (spot *Spot) ExchangeInfo_Params(params *Spot_ExchangeInfo_Params) (*Spot_ExchangeInfo, *Response, *Error) {
	opts := make(map[string]interface{})

	if len(params.Symbols) != 0 {
		opts["symbols"] = params.Symbols
	} else if params.Symbol != "" {
		opts["symbol"] = params.Symbol
	} else {
		if len(params.Permissions) != 0 {
			opts["permissions"] = params.Permissions
		}
		if params.SymbolStatus != "" {
			opts["symbolStatus"] = params.SymbolStatus
		}
		// if len(params.SymbolStatus) != 0 {
		// 	opts["symbolStatus"] = params.SymbolStatus
		// }
	}
	if params.DontShowPermissionSets {
		opts["showPermissionSets"] = !params.DontShowPermissionSets
	}

	resp, err := spot.makeRequest(&SpotRequest{
		securityType: SPOT_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/api/v3/exchangeInfo",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	exchangeInfo, err := ParseSpotExchangeInfo(resp)
	if err != nil {
		return nil, resp, err
	}

	return exchangeInfo, resp, nil
}

// # Exchange information
//
// Current exchange trading rules
// and symbol information
//
// Weight: 20
func (spot *Spot) ExchangeInfo() (*Spot_ExchangeInfo, *Response, *Error) {
	resp, err := spot.makeRequest(&SpotRequest{
		securityType: SPOT_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/api/v3/exchangeInfo",
	})
	if err != nil {
		return nil, resp, err
	}

	exchangeInfo, err := ParseSpotExchangeInfo(resp)
	if err != nil {
		return nil, resp, err
	}

	return exchangeInfo, resp, nil
}

func ParseSpotExchangeInfo(exchangeInfo_response *Response) (*Spot_ExchangeInfo, *Error) {
	var exchangeInfo Spot_ExchangeInfo

	err := json.Unmarshal(exchangeInfo_response.Body, &exchangeInfo)
	if err != nil {
		return nil, LocalError(PARSING_ERR, err.Error())
	}

	exchangeInfo.Symbols.Map = make(map[string]*Spot_Symbol)

	for _, symbol_obj := range exchangeInfo.Symbols_arr {
		exchangeInfo.Symbols.Map[symbol_obj.Symbol] = symbol_obj
	}

	return &exchangeInfo, nil
}

func (exchangeInfo *Spot_ExchangeInfo) UnmarshalJSON(data []byte) error {
	type Alias Spot_ExchangeInfo

	// Create an anonymous struct embedding the Alias type
	aux := &struct {
		Filters []jsoniter.RawMessage `json:"exchangeFilters"` // Capture filters as raw JSON
		*Alias
	}{
		Alias: (*Alias)(exchangeInfo), // Casting exchangeInfo to the alias so that unmarshall doesnt recursively call it again when we unmarshal it via this struct's unmarshall
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return LocalError(PARSING_ERR, err.Error())
	}

	for _, filter := range aux.Filters {
		var tempObj map[string]interface{}
		if err := json.Unmarshal(filter, &tempObj); err != nil {
			LOG_ERRORS("Error unmarshalling into temp map:", err)
			continue
		}

		switch tempObj["filterType"] {
		case SPOT_Constants.ExchangeFilterTypes.EXCHANGE_MAX_NUM_ORDERS:
			exchangeInfo.ExchangeFilters.EXCHANGE_MAX_NUM_ORDERS = &Spot_ExchangeFilter_EXCHANGE_MAX_NUM_ORDERS{}
			err = json.Unmarshal(filter, &exchangeInfo.ExchangeFilters.EXCHANGE_MAX_NUM_ORDERS)

		case SPOT_Constants.ExchangeFilterTypes.EXCHANGE_MAX_NUM_ALGO_ORDERS:
			exchangeInfo.ExchangeFilters.EXCHANGE_MAX_NUM_ALGO_ORDERS = &Spot_ExchangeFilter_EXCHANGE_MAX_NUM_ALGO_ORDERS{}
			err = json.Unmarshal(filter, &exchangeInfo.ExchangeFilters.EXCHANGE_MAX_NUM_ALGO_ORDERS)

		case SPOT_Constants.ExchangeFilterTypes.EXCHANGE_MAX_NUM_ICEBERG_ORDERS:
			exchangeInfo.ExchangeFilters.EXCHANGE_MAX_NUM_ICEBERG_ORDERS = &Spot_ExchangeFilter_EXCHANGE_MAX_NUM_ICEBERG_ORDERS{}
			err = json.Unmarshal(filter, &exchangeInfo.ExchangeFilters.EXCHANGE_MAX_NUM_ICEBERG_ORDERS)
		default:
			LOG_ERRORS("A missing field was intercepted of value", tempObj["filterType"], "in exchangeInfo.")
		}
		if err != nil {
			LOG_ERRORS("There was an error parsing", tempObj["filterType"], "in exchangeInfo =>", err)
		}

	}

	return nil
}

func (symbol *Spot_Symbol) UnmarshalJSON(data []byte) error {
	// Define an alias type to avoid recursive calls to UnmarshalJSON
	type Alias Spot_Symbol

	// Create an anonymous struct embedding the Alias type
	aux := &struct {
		Filters []jsoniter.RawMessage `json:"filters"` // Capture filters as raw JSON
		*Alias
	}{
		Alias: (*Alias)(symbol), // Casting symbol to the alias so that unmarshall doesnt recursively call it again when we unmarshal it via this struct's unmarshall
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	for _, filter := range aux.Filters {
		var tempObj map[string]interface{}
		if err := json.Unmarshal(filter, &tempObj); err != nil {
			LOG_ERRORS("Error unmarshalling into temp map:", err)
			continue
		}

		switch tempObj["filterType"] {
		case SPOT_Constants.SymbolFilterTypes.PRICE_FILTER:
			symbol.Filters.PRICE_FILTER = &Spot_SymbolFilter_PRICE_FILTER{}
			err = json.Unmarshal(filter, &symbol.Filters.PRICE_FILTER)

		case SPOT_Constants.SymbolFilterTypes.PERCENT_PRICE:
			symbol.Filters.PERCENT_PRICE = &Spot_SymbolFilter_PERCENT_PRICE{}
			err = json.Unmarshal(filter, &symbol.Filters.PERCENT_PRICE)

		case SPOT_Constants.SymbolFilterTypes.PERCENT_PRICE_BY_SIDE:
			symbol.Filters.PERCENT_PRICE_BY_SIDE = &Spot_SymbolFilter_PERCENT_PRICE_BY_SIDE{}
			err = json.Unmarshal(filter, &symbol.Filters.PERCENT_PRICE_BY_SIDE)

		case SPOT_Constants.SymbolFilterTypes.LOT_SIZE:
			symbol.Filters.LOT_SIZE = &Spot_SymbolFilter_LOT_SIZE{}
			err = json.Unmarshal(filter, &symbol.Filters.LOT_SIZE)

		case SPOT_Constants.SymbolFilterTypes.MIN_NOTIONAL:
			symbol.Filters.MIN_NOTIONAL = &Spot_SymbolFilter_MIN_NOTIONAL{}
			err = json.Unmarshal(filter, &symbol.Filters.MIN_NOTIONAL)

		case SPOT_Constants.SymbolFilterTypes.NOTIONAL:
			symbol.Filters.NOTIONAL = &Spot_SymbolFilter_NOTIONAL{}
			err = json.Unmarshal(filter, &symbol.Filters.NOTIONAL)

		case SPOT_Constants.SymbolFilterTypes.ICEBERG_PARTS:
			symbol.Filters.ICEBERG_PARTS = &Spot_SymbolFilter_ICEBERG_PARTS{}
			err = json.Unmarshal(filter, &symbol.Filters.ICEBERG_PARTS)

		case SPOT_Constants.SymbolFilterTypes.MARKET_LOT_SIZE:
			symbol.Filters.MARKET_LOT_SIZE = &Spot_SymbolFilter_MARKET_LOT_SIZE{}
			err = json.Unmarshal(filter, &symbol.Filters.MARKET_LOT_SIZE)

		case SPOT_Constants.SymbolFilterTypes.MAX_NUM_ORDERS:
			symbol.Filters.MAX_NUM_ORDERS = &Spot_SymbolFilter_MAX_NUM_ORDERS{}
			err = json.Unmarshal(filter, &symbol.Filters.MAX_NUM_ORDERS)

		case SPOT_Constants.SymbolFilterTypes.MAX_NUM_ALGO_ORDERS:
			symbol.Filters.MAX_NUM_ALGO_ORDERS = &Spot_SymbolFilter_MAX_NUM_ALGO_ORDERS{}
			err = json.Unmarshal(filter, &symbol.Filters.MAX_NUM_ALGO_ORDERS)

		case SPOT_Constants.SymbolFilterTypes.MAX_NUM_ICEBERG_ORDERS:
			symbol.Filters.MAX_NUM_ICEBERG_ORDERS = &Spot_SymbolFilter_MAX_NUM_ICEBERG_ORDERS{}
			err = json.Unmarshal(filter, &symbol.Filters.MAX_NUM_ICEBERG_ORDERS)

		case SPOT_Constants.SymbolFilterTypes.MAX_POSITION:
			symbol.Filters.MAX_POSITION = &Spot_SymbolFilter_MAX_POSITION{}
			err = json.Unmarshal(filter, &symbol.Filters.MAX_POSITION)

		case SPOT_Constants.SymbolFilterTypes.TRAILING_DELTA:
			symbol.Filters.TRAILING_DELTA = &Spot_SymbolFilter_TRAILING_DELTA{}
			err = json.Unmarshal(filter, &symbol.Filters.TRAILING_DELTA)
		default:
			LOG_ERRORS("A missing field was intercepted of value", tempObj["filterType"], "in the", symbol.Symbol, "symbol's info.")
		}
		if err != nil {
			LOG_ERRORS("There was an error parsing", tempObj["filterType"], "in the", symbol.Symbol, "symbol's info =>", err)
		}

	}

	return nil
}

//////// ExchangeInfo //
/////////////////////////////////////////////////////////////////////////////////

// # Order Book
//
// Weight adjusted based on the limit:
//
// | ------------------------------ |
//
// # | Limit         Request Weight |
//
// | ------------------------------ |
//
// | 1-100         => 5
//
// | 101-500       => 25
//
// | 501-1000      => 50
//
// | 1001-5000     => 250
func (spot *Spot) OrderBook(symbol string, limit ...int64) (*Spot_OrderBook, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol
	if len(limit) != 0 {
		opts["limit"] = limit[0]
	}

	resp, err := spot.makeRequest(&SpotRequest{
		securityType: SPOT_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/api/v3/depth",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var orderBook *Spot_OrderBook

	processingErr := json.Unmarshal(resp.Body, &orderBook)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERR, processingErr.Error())
	}

	return orderBook, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

// # Recent trades list
//
// Get recent trades.
//
// Weight: 25
//
// limit's default is 500, nax is 1000
func (spot *Spot) RecentTrades(symbol string, limit ...int64) ([]*Spot_Trade, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol
	if len(limit) != 0 {
		opts["limit"] = limit[0]
	}

	resp, err := spot.makeRequest(&SpotRequest{
		securityType: SPOT_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/api/v3/trades",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var trades []*Spot_Trade

	processingErr := json.Unmarshal(resp.Body, &trades)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERR, processingErr.Error())
	}

	return trades, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

type Spot_OldTrades_Params struct {
	// Default 500; max 1000.
	Limit int64
	// TradeId to fetch from. Default gets most recent trades.
	FromId int64
}

// # Old trade lookup
//
// Get older trades.
//
// Weight: 25
//
// Parameters:
//
//	type Spot_OldTrades_Params struct {
//		// Default 500; max 1000.
//		Limit int64
//		// TradeId to fetch from. Default gets most recent trades.
//		FromId int64
//	}
func (spot *Spot) OldTrades(symbol string, opt_params ...*Spot_OldTrades_Params) ([]*Spot_Trade, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol
	if len(opt_params) != 0 {
		params := opt_params[0]
		if params.Limit != 0 {
			opts["limit"] = params.Limit
		}
		if params.FromId != 0 {
			opts["fromId"] = params.FromId
		}
	}

	resp, err := spot.makeRequest(&SpotRequest{
		securityType: SPOT_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/api/v3/historicalTrades",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var trades []*Spot_Trade

	processingErr := json.Unmarshal(resp.Body, &trades)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERR, processingErr.Error())
	}

	return trades, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

type Spot_AggTrades_Params struct {
	// Default 500; max 1000.
	Limit int64
	// ID to get aggregate trades from INCLUSIVE.
	FromId int64
	// Timestamp in ms to get aggregate trades from INCLUSIVE.
	StartTime int64
	// Timestamp in ms to get aggregate trades until INCLUSIVE.
	EndTime int64
}

// #Compressed/Aggregate trades list
//
// Get compressed, aggregate trades. Trades that fill at the time, from the same taker order, with the same price will have the quantity aggregated.
//
// Weight: 2
//
// Parameters:
//
//	type Spot_AggTrades_Params struct {
//		// Default 500; max 1000.
//		Limit int64
//		// ID to get aggregate trades from INCLUSIVE.
//		FromId int64
//		// Timestamp in ms to get aggregate trades from INCLUSIVE.
//		StartTime int64
//		// Timestamp in ms to get aggregate trades until INCLUSIVE.
//		EndTime int64
//	}
func (spot *Spot) AggTrades(symbol string, opt_params ...*Spot_AggTrades_Params) ([]*Spot_AggTrade, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol
	if len(opt_params) != 0 {
		params := opt_params[0]
		if params.Limit != 0 {
			opts["limit"] = params.Limit
		}
		if params.FromId != 0 {
			opts["fromId"] = params.FromId
		}
		if params.StartTime != 0 {
			opts["startTime"] = params.StartTime
		}
		if params.EndTime != 0 {
			opts["endTime"] = params.EndTime
		}
	}

	resp, err := spot.makeRequest(&SpotRequest{
		securityType: SPOT_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/api/v3/aggTrades",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var aggTrades []*Spot_AggTrade

	processingErr := json.Unmarshal(resp.Body, &aggTrades)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERR, processingErr.Error())
	}

	return aggTrades, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

type Spot_Candlesticks_Params struct {
	// Default: 0 (UTC)
	TimeZone  string
	StartTime int64
	EndTime   int64
	// Default 500; max 1000.
	Limit int64
}

// # Kline/Candlestick data
//
// Kline/candlestick bars for a symbol. Klines are uniquely identified by their open time.
//
// Weight: 2
//
// Parameters:
//
//	type Spot_Candlesticks_Params struct {
//		// Default: 0 (UTC)
//		TimeZone  string
//		StartTime int64
//		EndTime   int64
//		// Default 500; max 1000.
//		// # Interval	interval value
//		//
//		// seconds:	"1s"
//		//
//		// minutes:	"1m", "3m", "5m", "15m", "30m"
//		//
//		// hours:	"1h", "2h", "4h", "6h", "8h", "12h"
//		//
//		// days:	"1d", "3d"
//		//
//		// weeks:	"1w"
//		//
//		// months:	"1M"
//		Limit int64
//	}
//
// # Supported kline intervals (case-sensitive):
//
// # Interval	interval value
//
// seconds:	"1s"
//
// minutes:	"1m", "3m", "5m", "15m", "30m"
//
// hours:	"1h", "2h", "4h", "6h", "8h", "12h"
//
// days:	"1d", "3d"
//
// weeks:	"1w"
//
// months:	"1M"
func (spot *Spot) Candlesticks(symbol string, interval string, opt_params ...*Spot_Candlesticks_Params) ([]*Spot_Candlestick, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol
	opts["interval"] = interval
	if len(opt_params) != 0 {
		params := opt_params[0]
		if params.Limit != 0 {
			opts["limit"] = params.Limit
		}
		if params.TimeZone != "" {
			opts["timeZone"] = params.TimeZone
		}
		if params.StartTime != 0 {
			opts["startTime"] = params.StartTime
		}
		if params.EndTime != 0 {
			opts["endTime"] = params.EndTime
		}
	}

	resp, err := spot.makeRequest(&SpotRequest{
		securityType: SPOT_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/api/v3/klines",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	// Unmarshal the data
	var rawCandlesticks [][]interface{}
	processingErr := json.Unmarshal(resp.Body, &rawCandlesticks)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERR, processingErr.Error())
	}

	// Convert the raw data to Spot_Candlestick slice
	candlesticks := make([]*Spot_Candlestick, len(rawCandlesticks))
	for i, raw := range rawCandlesticks {
		candlesticks[i] = &Spot_Candlestick{
			OpenTime:                 int64(raw[0].(float64)),
			Open:                     raw[1].(string),
			High:                     raw[2].(string),
			Low:                      raw[3].(string),
			Close:                    raw[4].(string),
			Volume:                   raw[5].(string),
			CloseTime:                int64(raw[6].(float64)),
			QuoteAssetVolume:         raw[7].(string),
			TradeCount:               int64(raw[8].(float64)),
			TakerBuyBaseAssetVolume:  raw[9].(string),
			TakerBuyQuoteAssetVolume: raw[10].(string),
			Unused:                   raw[11].(string),
		}
	}

	return candlesticks, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

func (spot *Spot) UIKlines(symbol string, interval string, opt_params ...*Spot_Candlesticks_Params) ([]*Spot_Candlestick, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol
	opts["interval"] = interval
	if len(opt_params) != 0 {
		params := opt_params[0]
		if params.Limit != 0 {
			opts["limit"] = params.Limit
		}
		if params.TimeZone != "" {
			opts["timeZone"] = params.TimeZone
		}
		if params.StartTime != 0 {
			opts["startTime"] = params.StartTime
		}
		if params.EndTime != 0 {
			opts["endTime"] = params.EndTime
		}
	}

	resp, err := spot.makeRequest(&SpotRequest{
		securityType: SPOT_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/api/v3/uiKlines",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	// Unmarshal the data
	var rawCandlesticks [][]interface{}
	processingErr := json.Unmarshal(resp.Body, &rawCandlesticks)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERR, processingErr.Error())
	}

	// Convert the raw data to Spot_Candlestick slice
	candlesticks := make([]*Spot_Candlestick, len(rawCandlesticks))
	for i, raw := range rawCandlesticks {
		candlesticks[i] = &Spot_Candlestick{
			OpenTime:                 int64(raw[0].(float64)),
			Open:                     raw[1].(string),
			High:                     raw[2].(string),
			Low:                      raw[3].(string),
			Close:                    raw[4].(string),
			Volume:                   raw[5].(string),
			CloseTime:                int64(raw[6].(float64)),
			QuoteAssetVolume:         raw[7].(string),
			TradeCount:               int64(raw[8].(float64)),
			TakerBuyBaseAssetVolume:  raw[9].(string),
			TakerBuyQuoteAssetVolume: raw[10].(string),
			Unused:                   raw[11].(string),
		}
	}

	return candlesticks, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

func (spot *Spot) AveragePrice(symbol string) (*Spot_AveragePrice, *Response, *Error) {
	opts := make(map[string]interface{})
	opts["symbol"] = symbol

	resp, err := spot.makeRequest(&SpotRequest{
		securityType: SPOT_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/api/v3/avgPrice",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	// Unmarshal the data
	var avgPrice Spot_AveragePrice
	processingErr := json.Unmarshal(resp.Body, &avgPrice)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERR, processingErr.Error())
	}

	return &avgPrice, resp, nil

}

/////////////////////////////////////////////////////////////////////////////////

func (spot *Spot) Ticker_RollingWindow24h(symbol ...string) ([]*Spot_Ticker_RollingWindow24h, *Response, *Error) {
	opts := make(map[string]interface{})

	if len(symbol) != 0 {
		opts["symbols"] = symbol
	}

	resp, err := spot.makeRequest(&SpotRequest{
		securityType: SPOT_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/api/v3/ticker/24hr",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var tickers []*Spot_Ticker_RollingWindow24h
	processingErr := json.Unmarshal(resp.Body, &tickers)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERR, processingErr.Error())
	}
	return tickers, resp, nil

}

/////////////////////////////////////////////////////////////////////////////////

func (spot *Spot) MiniTicker_RollingWindow24h(symbol ...string) ([]*Spot_MiniTicker_RollingWindow24h, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["type"] = "MINI"

	if len(symbol) != 0 {
		opts["symbols"] = symbol
	}

	resp, err := spot.makeRequest(&SpotRequest{
		securityType: SPOT_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/api/v3/ticker/24hr",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var miniTickers []*Spot_MiniTicker_RollingWindow24h
	processingErr := json.Unmarshal(resp.Body, &miniTickers)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERR, processingErr.Error())
	}
	return miniTickers, resp, nil

}

/////////////////////////////////////////////////////////////////////////////////

type Spot_Ticker_RollingWindow_Params struct {
	Symbols    []string
	WindowSize string
}

func (spot *Spot) Ticker_RollingWindow(opt_params *Spot_Ticker_RollingWindow_Params) ([]*Spot_Ticker_RollingWindow, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["windowSize"] = opt_params.WindowSize
	if len(opt_params.Symbols) != 0 {
		opts["symbols"] = opt_params.Symbols
	}

	resp, err := spot.makeRequest(&SpotRequest{
		securityType: SPOT_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/api/v3/ticker",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var tickers []*Spot_Ticker_RollingWindow
	processingErr := json.Unmarshal(resp.Body, &tickers)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERR, processingErr.Error())
	}
	return tickers, resp, nil

}

/////////////////////////////////////////////////////////////////////////////////

func (spot *Spot) MiniTicker_RollingWindow(opt_params *Spot_Ticker_RollingWindow_Params) ([]*Spot_MiniTicker_RollingWindow, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["windowSize"] = opt_params.WindowSize
	opts["type"] = "MINI"
	if len(opt_params.Symbols) != 0 {
		opts["symbols"] = opt_params.Symbols
	}

	resp, err := spot.makeRequest(&SpotRequest{
		securityType: SPOT_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/api/v3/ticker",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var miniTickers []*Spot_MiniTicker_RollingWindow
	processingErr := json.Unmarshal(resp.Body, &miniTickers)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERR, processingErr.Error())
	}
	return miniTickers, resp, nil

}

/////////////////////////////////////////////////////////////////////////////////

type Spot_Ticker_Params struct {
	Symbols  []string
	Timezone string
}

func (spot *Spot) Ticker(opt_params *Spot_Ticker_Params) ([]*Spot_Ticker, *Response, *Error) {
	opts := make(map[string]interface{})

	if len(opt_params.Symbols) != 0 {
		opts["symbols"] = opt_params.Symbols
	}
	if opt_params.Timezone != "" {
		opts["timezone"] = opt_params.Timezone
	}

	resp, err := spot.makeRequest(&SpotRequest{
		securityType: SPOT_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/api/v3/ticker/tradingDay",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var tickers []*Spot_Ticker
	processingErr := json.Unmarshal(resp.Body, &tickers)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERR, processingErr.Error())
	}
	return tickers, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

func (spot *Spot) MiniTicker(opt_params *Spot_Ticker_Params) ([]*Spot_MiniTicker, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["type"] = "MINI"
	if len(opt_params.Symbols) != 0 {
		opts["symbols"] = opt_params.Symbols
	}
	if opt_params.Timezone != "" {
		opts["timezone"] = opt_params.Timezone
	}

	resp, err := spot.makeRequest(&SpotRequest{
		securityType: SPOT_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/api/v3/ticker/tradingDay",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var miniTickers []*Spot_MiniTicker
	processingErr := json.Unmarshal(resp.Body, &miniTickers)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERR, processingErr.Error())
	}
	return miniTickers, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

func (spot *Spot) PriceTicker(symbol ...string) ([]*Spot_PriceTicker, *Response, *Error) {
	opts := make(map[string]interface{})

	if len(symbol) != 0 {
		opts["symbols"] = symbol
	}

	resp, err := spot.makeRequest(&SpotRequest{
		securityType: SPOT_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/api/v3/ticker/price",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var priceTickers []*Spot_PriceTicker
	processingErr := json.Unmarshal(resp.Body, &priceTickers)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERR, processingErr.Error())
	}
	return priceTickers, resp, nil

}

/////////////////////////////////////////////////////////////////////////////////

func (spot *Spot) BookTicker(symbol ...string) ([]*Spot_BookTicker, *Response, *Error) {
	opts := make(map[string]interface{})

	if len(symbol) != 0 {
		opts["symbols"] = symbol
	}

	resp, err := spot.makeRequest(&SpotRequest{
		securityType: SPOT_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/api/v3/ticker/bookTicker",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var bookTickers []*Spot_BookTicker
	processingErr := json.Unmarshal(resp.Body, &bookTickers)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERR, processingErr.Error())
	}
	return bookTickers, resp, nil

}

/////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////

// //////////////////////////// Orders \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\

// # DISCLAIMER
//
// Maybe someone will have the patience to create a unified structure for the functions
// that is similar to the one I used for my javascript library "binance-lib" that is still performant
// specifically in https://github.com/GTedZ/Binance-lib/blob/main/Spot/RESTful.js#L589
func (spot *Spot) newOrder(opts map[string]interface{}) (*Spot_Order, *Response, *Error) {
	resp, err := spot.makeRequest(&SpotRequest{
		securityType: SPOT_Constants.SecurityTypes.TRADE,
		method:       Constants.Methods.POST,
		url:          "/api/v3/order",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var order *Spot_Order
	processingErr := json.Unmarshal(resp.Body, &order)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERR, processingErr.Error())
	}
	return order, resp, nil
}

type Spot_Order_Params struct {
	TimeInForce             string
	Quantity                string
	QuoteOrderQty           string
	Price                   string
	NewClientOrderId        string
	StrategyId              int64
	StrategyType            int64
	StopPrice               string
	TrailingDelta           int64
	IcebergQty              string
	NewOrderRespType        string
	SelfTradePreventionMode string
	RecvWindow              int64
}

func (spot *Spot) NewOrder(symbol string, side string, Type string, opt_params ...Spot_Order_Params) (*Spot_Order, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol
	opts["side"] = side
	opts["type"] = Type

	if len(opt_params) != 0 {
		params := opt_params[0]
		if IsDifferentFromDefault(params.TimeInForce) {
			opts["timeInForce"] = params.TimeInForce
		}
		if IsDifferentFromDefault(params.Quantity) {
			opts["quantity"] = params.Quantity
		}
		if IsDifferentFromDefault(params.QuoteOrderQty) {
			opts["quoteOrderQty"] = params.QuoteOrderQty
		}
		if IsDifferentFromDefault(params.Price) {
			opts["price"] = params.Price
		}
		if IsDifferentFromDefault(params.NewClientOrderId) {
			opts["newClientOrderId"] = params.NewClientOrderId
		}
		if IsDifferentFromDefault(params.StrategyId) {
			opts["strategyId"] = params.StrategyId
		}
		if IsDifferentFromDefault(params.StrategyType) {
			opts["strategyType"] = params.StrategyType
		}
		if IsDifferentFromDefault(params.StopPrice) {
			opts["stopPrice"] = params.StopPrice
		}
		if IsDifferentFromDefault(params.TrailingDelta) {
			opts["trailingDelta"] = params.TrailingDelta
		}
		if IsDifferentFromDefault(params.IcebergQty) {
			opts["icebergQty"] = params.IcebergQty
		}
		if IsDifferentFromDefault(params.NewOrderRespType) {
			opts["newOrderRespType"] = params.NewOrderRespType
		}
		if IsDifferentFromDefault(params.SelfTradePreventionMode) {
			opts["selfTradePreventionMode"] = params.SelfTradePreventionMode
		}
		if IsDifferentFromDefault(params.RecvWindow) {
			opts["recvWindow"] = params.RecvWindow
		}
	}

	return spot.newOrder(opts)
}

///////////////////////// LIMIT \\\\\\\\\\\\\\\\\\\\\\\\\\\\

type Spot_LimitOrder_Params struct {
	TimeInForce             string
	NewClientOrderId        string
	StrategyId              int64
	StrategyType            int64
	IcebergQty              string
	NewOrderRespType        string
	SelfTradePreventionMode string
	RecvWindow              int64
}

func (spot *Spot) LimitOrder(symbol string, side string, price string, quantity string, opt_params ...Spot_LimitOrder_Params) (*Spot_Order, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol
	opts["side"] = side
	opts["type"] = "LIMIT"
	opts["price"] = price
	opts["quantity"] = quantity

	if len(opt_params) != 0 {
		params := opt_params[0]
		if IsDifferentFromDefault(params.TimeInForce) {
			opts["timeInForce"] = params.TimeInForce
		}
		if IsDifferentFromDefault(params.NewClientOrderId) {
			opts["newClientOrderId"] = params.NewClientOrderId
		}
		if IsDifferentFromDefault(params.StrategyId) {
			opts["strategyId"] = params.StrategyId
		}
		if IsDifferentFromDefault(params.StrategyType) {
			opts["strategyType"] = params.StrategyType
		}
		if IsDifferentFromDefault(params.IcebergQty) {
			opts["icebergQty"] = params.IcebergQty
		}
		if IsDifferentFromDefault(params.NewOrderRespType) {
			opts["newOrderRespType"] = params.NewOrderRespType
		}
		if IsDifferentFromDefault(params.SelfTradePreventionMode) {
			opts["selfTradePreventionMode"] = params.SelfTradePreventionMode
		}
		if IsDifferentFromDefault(params.RecvWindow) {
			opts["recvWindow"] = params.RecvWindow
		}
	}

	return spot.newOrder(opts)
}

func (spot *Spot) LimitBuy(symbol string, price string, quantity string, opt_params ...Spot_LimitOrder_Params) (*Spot_Order, *Response, *Error) {
	return spot.LimitOrder(symbol, "BUY", price, quantity, opt_params...)
}

func (spot *Spot) LimitSell(symbol string, price string, quantity string, opt_params ...Spot_LimitOrder_Params) (*Spot_Order, *Response, *Error) {
	return spot.LimitOrder(symbol, "SELL", price, quantity, opt_params...)
}

///////////////////////// LIMIT ////////////////////////////

///////////////////////// LIMIT_MAKER \\\\\\\\\\\\\\\\\\\\\\\\\\\

type Spot_LimitMakerOrder_Params struct {
	TimeInForce             string
	NewClientOrderId        string
	StrategyId              int64
	StrategyType            int64
	IcebergQty              string
	NewOrderRespType        string
	SelfTradePreventionMode string
	RecvWindow              int64
}

func (spot *Spot) LimitMakerOrder(symbol string, side string, quantity string, price string, opt_params ...Spot_LimitMakerOrder_Params) (*Spot_Order, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol
	opts["side"] = side
	opts["type"] = "LIMIT_MAKER"

	if len(opt_params) != 0 {
		params := opt_params[0]
		if IsDifferentFromDefault(params.TimeInForce) {
			opts["timeInForce"] = params.TimeInForce
		}
		if IsDifferentFromDefault(params.NewClientOrderId) {
			opts["newClientOrderId"] = params.NewClientOrderId
		}
		if IsDifferentFromDefault(params.StrategyId) {
			opts["strategyId"] = params.StrategyId
		}
		if IsDifferentFromDefault(params.StrategyType) {
			opts["strategyType"] = params.StrategyType
		}
		if IsDifferentFromDefault(params.IcebergQty) {
			opts["icebergQty"] = params.IcebergQty
		}
		if IsDifferentFromDefault(params.NewOrderRespType) {
			opts["newOrderRespType"] = params.NewOrderRespType
		}
		if IsDifferentFromDefault(params.SelfTradePreventionMode) {
			opts["selfTradePreventionMode"] = params.SelfTradePreventionMode
		}
		if IsDifferentFromDefault(params.RecvWindow) {
			opts["recvWindow"] = params.RecvWindow
		}
	}

	return spot.newOrder(opts)
}

func (spot *Spot) LimitMakerBuy(symbol string, quantity string, price string, opt_params ...Spot_LimitMakerOrder_Params) (*Spot_Order, *Response, *Error) {
	return spot.LimitMakerOrder(symbol, "BUY", quantity, price, opt_params...)
}

func (spot *Spot) LimitMakerSell(symbol string, side string, quantity string, price string, opt_params ...Spot_LimitMakerOrder_Params) (*Spot_Order, *Response, *Error) {
	return spot.LimitMakerOrder(symbol, "SELL", quantity, price, opt_params...)
}

///////////////////////// LIMIT_MAKER ///////////////////////////

///////////////////////// MARKET \\\\\\\\\\\\\\\\\\\\\\\\\\\

type Spot_MarketOrder_Params struct {
	NewClientOrderId        string
	StrategyId              int64
	StrategyType            int64
	NewOrderRespType        string
	SelfTradePreventionMode string
	RecvWindow              int64
}

func (spot *Spot) MarketOrder(symbol string, side string, orderValue string, is_OrderValue_in_BaseAsset bool, opt_params ...Spot_MarketOrder_Params) (*Spot_Order, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol
	opts["side"] = side
	opts["type"] = "MARKET"
	if is_OrderValue_in_BaseAsset {
		opts["quantity"] = orderValue
	} else {
		opts["quoteOrderQty"] = orderValue
	}

	if len(opt_params) != 0 {
		params := opt_params[0]
		if IsDifferentFromDefault(params.NewClientOrderId) {
			opts["newClientOrderId"] = params.NewClientOrderId
		}
		if IsDifferentFromDefault(params.StrategyId) {
			opts["strategyId"] = params.StrategyId
		}
		if IsDifferentFromDefault(params.StrategyType) {
			opts["strategyType"] = params.StrategyType
		}
		if IsDifferentFromDefault(params.NewOrderRespType) {
			opts["newOrderRespType"] = params.NewOrderRespType
		}
		if IsDifferentFromDefault(params.SelfTradePreventionMode) {
			opts["selfTradePreventionMode"] = params.SelfTradePreventionMode
		}
		if IsDifferentFromDefault(params.RecvWindow) {
			opts["recvWindow"] = params.RecvWindow
		}
	}

	return spot.newOrder(opts)
}

func (spot *Spot) MarketBuy(symbol string, side string, orderValue string, is_OrderValue_in_BaseAsset bool, opt_params ...Spot_MarketOrder_Params) (*Spot_Order, *Response, *Error) {
	return spot.MarketOrder(symbol, "BUY", orderValue, is_OrderValue_in_BaseAsset, opt_params...)
}

func (spot *Spot) MarketSell(symbol string, side string, orderValue string, is_OrderValue_in_BaseAsset bool, opt_params ...Spot_MarketOrder_Params) (*Spot_Order, *Response, *Error) {
	return spot.MarketOrder(symbol, "SELL", orderValue, is_OrderValue_in_BaseAsset, opt_params...)
}

///////////////////////// MARKET ///////////////////////////

// \\\\\\\\\\\\\\\\\\\\\\\\\\\ Orders ////////////////////////////////////////

type Spot_QueryOrder_Params struct {
	OrigClientOrderId string
	RecvWindow        int64
}

func (spot *Spot) QueryOrder(symbol string, orderId int64, opt_params ...Spot_QueryOrder_Params) (*Spot_Order, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol
	opts["orderId"] = orderId

	if len(opt_params) != 0 {
		params := opt_params[0]
		if IsDifferentFromDefault(params.OrigClientOrderId) {
			opts["origClientOrderId"] = params.OrigClientOrderId
			delete(opts, "orderId")
		}
		if IsDifferentFromDefault(params.RecvWindow) {
			opts["recvWindow"] = params.RecvWindow
		}
	}

	resp, err := spot.makeRequest(&SpotRequest{
		securityType: SPOT_Constants.SecurityTypes.USER_DATA,
		method:       Constants.Methods.GET,
		url:          "/api/v3/order",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var order *Spot_Order
	processingErr := json.Unmarshal(resp.Body, &order)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERR, processingErr.Error())
	}
	return order, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////

type Spot_AccountInfo_Params struct {
	OmitZeroBalances bool
	RecvWindow       int64
}

func (spot *Spot) AccountInfo(opt_params ...Spot_AccountInfo_Params) (*Spot_AccountInfo, *Response, *Error) {
	opts := make(map[string]interface{})

	if len(opt_params) != 0 {
		params := opt_params[0]
		opts["omitZeroBalances"] = params.OmitZeroBalances
		if params.RecvWindow != 0 {
			opts["recvWindow"] = params.RecvWindow
		}
	}

	resp, err := spot.makeRequest(&SpotRequest{
		securityType: SPOT_Constants.SecurityTypes.USER_DATA,
		method:       Constants.Methods.GET,
		url:          "/api/v3/account",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var accountInfo *Spot_AccountInfo
	processingErr := json.Unmarshal(resp.Body, &accountInfo)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERR, processingErr.Error())
	}
	return accountInfo, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////

type SpotRequest struct {
	method       string
	url          string
	params       map[string]interface{}
	securityType string
}

func (spot *Spot) makeRequest(request *SpotRequest) (*Response, *Error) {

	switch request.securityType {
	case SPOT_Constants.SecurityTypes.NONE:
		return spot.requestClient.Unsigned(request.method, SPOT_Constants.URL_Data_Only, request.url, request.params)
	case SPOT_Constants.SecurityTypes.USER_STREAM:
		return spot.requestClient.APIKEY_only(request.method, spot.baseURL, request.url, request.params)

	case SPOT_Constants.SecurityTypes.TRADE:
		return spot.requestClient.Signed(request.method, spot.baseURL, request.url, request.params)
	case SPOT_Constants.SecurityTypes.USER_DATA:
		return spot.requestClient.Signed(request.method, spot.baseURL, request.url, request.params)

	default:
		panic(fmt.Sprintf("Security Type passed to Request function is invalid, received: '%s'\nSupported methods are ('%s', '%s', '%s', '%s')", request.securityType, SPOT_Constants.SecurityTypes.NONE, SPOT_Constants.SecurityTypes.USER_STREAM, SPOT_Constants.SecurityTypes.TRADE, SPOT_Constants.SecurityTypes.USER_DATA))
	}

}
