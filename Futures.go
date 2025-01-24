package Binance

import (
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type Futures struct {
	binance       *Binance
	requestClient RequestClient
	baseURL       string

	API APIKEYS

	Websockets Futures_Websockets
}

func (futures *Futures) init(binance *Binance) {
	futures.binance = binance

	futures.requestClient.init(binance)
	futures.requestClient.Set_APIKEY(binance.API.KEY, binance.API.SECRET)

	futures.baseURL = FUTURES_Constants.URLs[0]

	futures.API.Set(binance.API.KEY, binance.API.SECRET)

	futures.Websockets.binance = binance
}

/////////////////////////////////////////////////////////////////////////////////

// # Test connectivity to the Rest API.
//
// Weight: 1
//
// Data Source: Memory
func (futures *Futures) Ping() (latency int64, request *Response, err *Error) {
	startTime := time.Now().UnixMilli()
	httpResp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/fapi/v1/ping",
	})
	diff := time.Now().UnixMilli() - startTime
	if err != nil {
		return diff, httpResp, err
	}

	return diff, httpResp, nil
}

func (futures *Futures) ServerTime() (*Futures_Time, *Response, *Error) {
	var futuresTime Futures_Time

	startTime := time.Now().UnixMilli()
	httpResp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/fapi/v1/time",
	})
	diff := time.Now().UnixMilli() - startTime
	futuresTime.Latency = diff
	if err != nil {
		return &futuresTime, httpResp, err
	}

	processingErr := json.Unmarshal(httpResp.Body, &futuresTime)
	if processingErr != nil {
		return &futuresTime, httpResp, LocalError(PARSING_ERROR, processingErr.Error())
	}

	return &futuresTime, httpResp, nil
}

/////////////////////////////////////////////////////////////////////////////////

//////// ExchangeInfo \\

func (futures *Futures) ExchangeInfo() (*Futures_ExchangeInfo, *Response, *Error) {
	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/fapi/v1/exchangeInfo",
	})
	if err != nil {
		return nil, resp, err
	}

	exchangeInfo, err := ParseFuturesExchangeInfo(resp)
	if err != nil {
		return nil, resp, err
	}

	return exchangeInfo, resp, nil
}

func ParseFuturesExchangeInfo(exchangeInfo_response *Response) (*Futures_ExchangeInfo, *Error) {
	var exchangeInfo Futures_ExchangeInfo

	err := json.Unmarshal(exchangeInfo_response.Body, &exchangeInfo)
	if err != nil {
		return nil, LocalError(PARSING_ERROR, err.Error())
	}

	exchangeInfo.Symbols = make(map[string]*Futures_Symbol)
	for _, symbol_obj := range exchangeInfo.Symbols_arr {
		exchangeInfo.Symbols[symbol_obj.Symbol] = symbol_obj
	}

	exchangeInfo.Assets = make(map[string]*Futures_Asset)
	for _, asset_obj := range exchangeInfo.Assets_arr {
		exchangeInfo.Assets[asset_obj.Asset] = asset_obj
	}

	return &exchangeInfo, nil
}

func (symbol *Futures_Symbol) UnmarshalJSON(data []byte) error {
	// Define an alias type to avoid recursive calls to UnmarshalJSON
	type Alias Futures_Symbol

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
			if VERBOSE {
				fmt.Println("Error unmarshalling into temp map:", err)
			}
			continue
		}

		switch tempObj["filterType"] {
		case FUTURES_Constants.SymbolFilterTypes.PRICE_FILTER:
			symbol.Filters.PRICE_FILTER = &Futures_SymbolFilter_PRICE_FILTER{}
			err = json.Unmarshal(filter, &symbol.Filters.PRICE_FILTER)

		case FUTURES_Constants.SymbolFilterTypes.LOT_SIZE:
			symbol.Filters.LOT_SIZE = &Futures_SymbolFilter_LOT_SIZE{}
			err = json.Unmarshal(filter, &symbol.Filters.LOT_SIZE)

		case FUTURES_Constants.SymbolFilterTypes.MARKET_LOT_SIZE:
			symbol.Filters.MARKET_LOT_SIZE = &Futures_SymbolFilter_MARKET_LOT_SIZE{}
			err = json.Unmarshal(filter, &symbol.Filters.MARKET_LOT_SIZE)

		case FUTURES_Constants.SymbolFilterTypes.MAX_NUM_ORDERS:
			symbol.Filters.MAX_NUM_ORDERS = &Futures_SymbolFilter_MAX_NUM_ORDERS{}
			err = json.Unmarshal(filter, &symbol.Filters.MAX_NUM_ORDERS)

		case FUTURES_Constants.SymbolFilterTypes.MAX_NUM_ALGO_ORDERS:
			symbol.Filters.MAX_NUM_ALGO_ORDERS = &Futures_SymbolFilter_MAX_NUM_ALGO_ORDERS{}
			err = json.Unmarshal(filter, &symbol.Filters.MAX_NUM_ALGO_ORDERS)

		case FUTURES_Constants.SymbolFilterTypes.PERCENT_PRICE:
			symbol.Filters.PERCENT_PRICE = &Futures_SymbolFilter_PERCENT_PRICE{}
			err = json.Unmarshal(filter, &symbol.Filters.PERCENT_PRICE)

		case FUTURES_Constants.SymbolFilterTypes.MIN_NOTIONAL:
			symbol.Filters.MIN_NOTIONAL = &Futures_SymbolFilter_MIN_NOTIONAL{}
			err = json.Unmarshal(filter, &symbol.Filters.MIN_NOTIONAL)
		default:
			if VERBOSE {
				fmt.Println("A missing field was intercepted of value", tempObj["filterType"], "in the", symbol.Symbol, "symbol's info.")
			}
		}
		if err != nil {
			if VERBOSE {
				fmt.Println("There was an error parsing", tempObj["filterType"], "in the", symbol.Symbol, "symbol's info =>", err)
			}
		}

	}

	return nil
}

//////// ExchangeInfo //
/////////////////////////////////////////////////////////////////////////////////

func (futures *Futures) OrderBook(symbol string, limit ...int64) (*Futures_OrderBook, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol

	if len(limit) != 0 {
		opts["limit"] = limit[0]
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/fapi/v1/depth",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var orderBook Futures_OrderBook

	unmarshallErr := json.Unmarshal(resp.Body, &orderBook)
	if unmarshallErr != nil {
		return nil, resp, LocalError(PARSING_ERROR, unmarshallErr.Error())
	}

	return &orderBook, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

func (futures *Futures) Trades(symbol string, limit ...int64) ([]*Futures_Trade, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol

	var lim int64 = 500

	if len(limit) != 0 {
		opts["limit"] = limit[0]
		lim = limit[0]
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/fapi/v1/trades",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	trades := make([]*Futures_Trade, lim)

	unmarshallErr := json.Unmarshal(resp.Body, &trades)
	if unmarshallErr != nil {
		return nil, resp, LocalError(PARSING_ERROR, unmarshallErr.Error())
	}

	return trades, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

type Futures_HistoricalTrades_Params struct {
	// Default 100, max 500.
	Limit int64
	// TradeId to fetch from. Default gets the most recent trades.
	FromId int64
}

func (futures *Futures) HistoricalTrades(symbol string, opt_params ...Futures_HistoricalTrades_Params) ([]*Futures_Trade, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol

	var limit int64 = 100

	if len(opt_params) != 0 {
		params := opt_params[0]
		if params.Limit != 0 {
			opts["limit"] = params.Limit
			limit = params.Limit
		}
		if params.FromId != 0 {
			opts["fromId"] = params.FromId
		}
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/fapi/v1/historicalTrades",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	historicalTrades := make([]*Futures_Trade, limit)

	unmarshallErr := json.Unmarshal(resp.Body, &historicalTrades)
	if unmarshallErr != nil {
		return nil, resp, LocalError(PARSING_ERROR, unmarshallErr.Error())
	}

	return historicalTrades, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

type Futures_AggTrade_Params struct {
	// ID to get aggregate trades from INCLUSIVE.
	FromId int64
	// Timestamp in ms to get aggregate trades from INCLUSIVE.
	StartTime int64
	// Timestamp in ms to get aggregate trades until INCLUSIVE.
	EndTime int64
	// Default 500; max 1000.
	Limit int64
}

func (futures *Futures) AggTrades(symbol string, opt_params ...Futures_AggTrade_Params) ([]*Futures_AggTrade_Params, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol

	var limit int64 = 100

	if len(opt_params) != 0 {
		params := opt_params[0]
		if params.FromId != 0 {
			opts["fromId"] = params.FromId
		}
		if params.StartTime != 0 {
			opts["startTime"] = params.StartTime
		}
		if params.EndTime != 0 {
			opts["endTime"] = params.EndTime
		}
		if params.Limit != 0 {
			opts["limit"] = params.Limit
			limit = params.Limit
		}
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/fapi/v1/aggTrades",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	aggTrades := make([]*Futures_AggTrade_Params, limit)

	unmarshallErr := json.Unmarshal(resp.Body, &aggTrades)
	if unmarshallErr != nil {
		return nil, resp, LocalError(PARSING_ERROR, unmarshallErr.Error())
	}

	return aggTrades, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

type Futures_Candlesticks_Params struct {
	StartTime int64
	EndTime   int64
	// Default 500; max 1500.
	Limit int64
}

func (futures *Futures) Candlesticks(symbol string, interval string, opt_params ...*Futures_Candlesticks_Params) ([]*Futures_Candlestick, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol
	opts["interval"] = interval
	if len(opt_params) != 0 {
		params := opt_params[0]
		if params.Limit != 0 {
			opts["limit"] = params.Limit
		}
		if params.StartTime != 0 {
			opts["startTime"] = params.StartTime
		}
		if params.EndTime != 0 {
			opts["endTime"] = params.EndTime
		}
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/fapi/v1/klines",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	// Unmarshal the data
	var rawCandlesticks [][]interface{}
	processingErr := json.Unmarshal(resp.Body, &rawCandlesticks)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERROR, processingErr.Error())
	}

	// Convert the raw data to Spot_Candlestick slice
	candlesticks := make([]*Futures_Candlestick, len(rawCandlesticks))
	for i, raw := range rawCandlesticks {
		candlesticks[i] = &Futures_Candlestick{
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

// # Kline/candlestick bars for a specific contract type.
//
// Klines are uniquely identified by their open time.
//
// Contract Types:
// "PERPETUAL" | "CURRENT_QUARTER" | "NEXT_QUARTER"
func (futures *Futures) ContinuousContractCandlesticks(symbol string, contractType string, interval string, opt_params ...*Futures_Candlesticks_Params) ([]*Futures_Candlestick, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol
	opts["contractType"] = contractType
	opts["interval"] = interval
	if len(opt_params) != 0 {
		params := opt_params[0]
		if params.Limit != 0 {
			opts["limit"] = params.Limit
		}
		if params.StartTime != 0 {
			opts["startTime"] = params.StartTime
		}
		if params.EndTime != 0 {
			opts["endTime"] = params.EndTime
		}
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/fapi/v1/continuousKlines",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	// Unmarshal the data
	var rawCandlesticks [][]interface{}
	processingErr := json.Unmarshal(resp.Body, &rawCandlesticks)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERROR, processingErr.Error())
	}

	// Convert the raw data to Spot_Candlestick slice
	candlesticks := make([]*Futures_Candlestick, len(rawCandlesticks))
	for i, raw := range rawCandlesticks {
		candlesticks[i] = &Futures_Candlestick{
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

type Futures_PriceCandlesticks_Params struct {
	StartTime int64
	EndTime   int64
	// Default 500; max 1500.
	Limit int64
}

func (futures *Futures) IndexPriceCandlesticks(symbol string, interval string, opt_params ...*Futures_PriceCandlesticks_Params) ([]*Futures_PriceCandlestick, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol
	opts["interval"] = interval
	if len(opt_params) != 0 {
		params := opt_params[0]
		if params.Limit != 0 {
			opts["limit"] = params.Limit
		}
		if params.StartTime != 0 {
			opts["startTime"] = params.StartTime
		}
		if params.EndTime != 0 {
			opts["endTime"] = params.EndTime
		}
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/fapi/v1/indexPriceKlines",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	// Unmarshal the data
	var rawCandlesticks [][]interface{}
	processingErr := json.Unmarshal(resp.Body, &rawCandlesticks)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERROR, processingErr.Error())
	}

	// Convert the raw data to Spot_Candlestick slice
	candlesticks := make([]*Futures_PriceCandlestick, len(rawCandlesticks))
	for i, raw := range rawCandlesticks {
		candlesticks[i] = &Futures_PriceCandlestick{
			OpenTime:  int64(raw[0].(float64)),
			Open:      raw[1].(string),
			High:      raw[2].(string),
			Low:       raw[3].(string),
			Close:     raw[4].(string),
			Ignore1:   raw[5].(string),
			CloseTime: int64(raw[6].(float64)),
			Ignore2:   raw[7].(string),
			Ignore3:   int64(raw[8].(float64)),
			Ignore4:   raw[9].(string),
			Ignore5:   raw[10].(string),
			Unused:    raw[11].(string),
		}
	}

	return candlesticks, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

func (futures *Futures) MarkPriceCandlesticks(symbol string, contractType string, interval string, opt_params ...*Futures_PriceCandlesticks_Params) ([]*Futures_PriceCandlestick, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol
	opts["interval"] = interval
	if len(opt_params) != 0 {
		params := opt_params[0]
		if params.Limit != 0 {
			opts["limit"] = params.Limit
		}
		if params.StartTime != 0 {
			opts["startTime"] = params.StartTime
		}
		if params.EndTime != 0 {
			opts["endTime"] = params.EndTime
		}
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/fapi/v1/markPriceKlines",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	// Unmarshal the data
	var rawCandlesticks [][]interface{}
	processingErr := json.Unmarshal(resp.Body, &rawCandlesticks)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERROR, processingErr.Error())
	}

	// Convert the raw data to Spot_Candlestick slice
	candlesticks := make([]*Futures_PriceCandlestick, len(rawCandlesticks))
	for i, raw := range rawCandlesticks {
		candlesticks[i] = &Futures_PriceCandlestick{
			OpenTime:  int64(raw[0].(float64)),
			Open:      raw[1].(string),
			High:      raw[2].(string),
			Low:       raw[3].(string),
			Close:     raw[4].(string),
			Ignore1:   raw[5].(string),
			CloseTime: int64(raw[6].(float64)),
			Ignore2:   raw[7].(string),
			Ignore3:   int64(raw[8].(float64)),
			Ignore4:   raw[9].(string),
			Ignore5:   raw[10].(string),
			Unused:    raw[11].(string),
		}
	}

	return candlesticks, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

func (futures *Futures) PremiumIndexCandlesticks(symbol string, contractType string, interval string, opt_params ...*Futures_PriceCandlesticks_Params) ([]*Futures_PriceCandlestick, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol
	opts["interval"] = interval
	if len(opt_params) != 0 {
		params := opt_params[0]
		if params.Limit != 0 {
			opts["limit"] = params.Limit
		}
		if params.StartTime != 0 {
			opts["startTime"] = params.StartTime
		}
		if params.EndTime != 0 {
			opts["endTime"] = params.EndTime
		}
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/fapi/v1/premiumIndexKlines",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	// Unmarshal the data
	var rawCandlesticks [][]interface{}
	processingErr := json.Unmarshal(resp.Body, &rawCandlesticks)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERROR, processingErr.Error())
	}

	// Convert the raw data to Spot_Candlestick slice
	candlesticks := make([]*Futures_PriceCandlestick, len(rawCandlesticks))
	for i, raw := range rawCandlesticks {
		candlesticks[i] = &Futures_PriceCandlestick{
			OpenTime:  int64(raw[0].(float64)),
			Open:      raw[1].(string),
			High:      raw[2].(string),
			Low:       raw[3].(string),
			Close:     raw[4].(string),
			Ignore1:   raw[5].(string),
			CloseTime: int64(raw[6].(float64)),
			Ignore2:   raw[7].(string),
			Ignore3:   int64(raw[8].(float64)),
			Ignore4:   raw[9].(string),
			Ignore5:   raw[10].(string),
			Unused:    raw[11].(string),
		}
	}

	return candlesticks, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

func (futures *Futures) MarkPrice(symbol ...string) ([]*Futures_MarkPrice, *Response, *Error) {
	opts := make(map[string]interface{})

	if len(symbol) != 0 {
		opts["symbol"] = symbol[0]
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/fapi/v1/premiumIndex",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	if len(symbol) != 0 {
		var markPrice Futures_MarkPrice

		unmarshallErr := json.Unmarshal(resp.Body, &markPrice)
		if unmarshallErr != nil {
			return nil, resp, LocalError(PARSING_ERROR, unmarshallErr.Error())
		}

		return []*Futures_MarkPrice{&markPrice}, resp, nil
	} else {
		var markPrices []*Futures_MarkPrice

		unmarshallErr := json.Unmarshal(resp.Body, &markPrices)
		if unmarshallErr != nil {
			return nil, resp, LocalError(PARSING_ERROR, unmarshallErr.Error())
		}

		return markPrices, resp, nil
	}
}

/////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////

type FuturesRequest struct {
	method       string
	url          string
	params       map[string]interface{}
	securityType string
}

func (futures *Futures) makeRequest(request *FuturesRequest) (*Response, *Error) {

	switch request.securityType {
	case FUTURES_Constants.SecurityTypes.NONE:
		return futures.requestClient.Unsigned(request.method, futures.baseURL, request.url, request.params)
	case FUTURES_Constants.SecurityTypes.MARKET_DATA:
		return futures.requestClient.APIKEY_only(request.method, futures.baseURL, request.url, request.params)
	case FUTURES_Constants.SecurityTypes.USER_STREAM:
		return futures.requestClient.APIKEY_only(request.method, futures.baseURL, request.url, request.params)

	case FUTURES_Constants.SecurityTypes.TRADE:
		return futures.requestClient.Signed(request.method, futures.baseURL, request.url, request.params)
	case FUTURES_Constants.SecurityTypes.USER_DATA:
		return futures.requestClient.Signed(request.method, futures.baseURL, request.url, request.params)

	default:
		panic(fmt.Sprintf("Security Type passed to Request function is invalid, received: '%s'\nSupported methods are ('%s', '%s', '%s', '%s')", request.securityType, FUTURES_Constants.SecurityTypes.NONE, FUTURES_Constants.SecurityTypes.USER_STREAM, FUTURES_Constants.SecurityTypes.TRADE, FUTURES_Constants.SecurityTypes.USER_DATA))
	}

}
