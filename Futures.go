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
			LOG_ERRORS("Error unmarshalling into temp map:", err)
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

	// Convert the raw data to Futures_Candlestick slice
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

	// Convert the raw data to Futures_Candlestick slice
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

	// Convert the raw data to Futures_Candlestick slice
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

	// Convert the raw data to Futures_Candlestick slice
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

	// Convert the raw data to Futures_Candlestick slice
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

type Futures_FundingRate_Params struct {
	Symbol    string
	StartTime int64
	EndTime   int64
	Limit     int64
}

func (futures *Futures) FundingRateHistory(opt_params ...Futures_FundingRate_Params) ([]*Futures_FundingRate, *Response, *Error) {
	opts := make(map[string]interface{})

	if len(opt_params) != 0 {
		params := opt_params[0]
		if IsDifferentFromDefault(params.Symbol) {
			opts["symbol"] = params.Symbol
		}
		if IsDifferentFromDefault(params.StartTime) {
			opts["startTime"] = params.StartTime
		}
		if IsDifferentFromDefault(params.EndTime) {
			opts["endTime"] = params.EndTime
		}
		if IsDifferentFromDefault(params.Limit) {
			opts["limit"] = params.Limit
		}
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/fapi/v1/fundingRate",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var fundingRates []*Futures_FundingRate

	unmarshallErr := json.Unmarshal(resp.Body, &fundingRates)
	if unmarshallErr != nil {
		return nil, resp, LocalError(PARSING_ERROR, unmarshallErr.Error())
	}

	return fundingRates, resp, nil
}

func (futures *Futures) FundingRate() ([]*Futures_FundingRate, *Response, *Error) {
	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/fapi/v1/fundingInfo",
	})
	if err != nil {
		return nil, resp, err
	}

	var fundingRates []*Futures_FundingRate

	unmarshallErr := json.Unmarshal(resp.Body, &fundingRates)
	if unmarshallErr != nil {
		return nil, resp, LocalError(PARSING_ERROR, unmarshallErr.Error())
	}

	return fundingRates, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

// If the symbol is not sent, bookTickers for all symbols will be returned in an array.
func (futures *Futures) Ticker24h(symbol ...string) ([]*Futures_24hTicker, *Response, *Error) {
	opts := make(map[string]interface{})

	if len(symbol) != 0 {
		opts["symbol"] = symbol[0]
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/fapi/v1/ticker/24hr",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	if len(symbol) == 0 {
		var ticker Futures_24hTicker
		unmarshallErr := json.Unmarshal(resp.Body, &ticker)
		if unmarshallErr != nil {
			return nil, resp, LocalError(PARSING_ERROR, unmarshallErr.Error())
		}

		return []*Futures_24hTicker{&ticker}, resp, nil
	} else {
		var tickers []*Futures_24hTicker
		unmarshallErr := json.Unmarshal(resp.Body, &tickers)
		if unmarshallErr != nil {
			return nil, resp, LocalError(PARSING_ERROR, unmarshallErr.Error())
		}

		return tickers, resp, nil
	}
}

// If the symbol is not sent, bookTickers for all symbols will be returned in an array.
func (futures *Futures) PriceTicker_v1(symbol ...string) ([]*Futures_PriceTicker, *Response, *Error) {
	opts := make(map[string]interface{})

	if len(symbol) != 0 {
		opts["symbol"] = symbol[0]
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/fapi/v1/ticker/price",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	if len(symbol) == 0 {
		var priceTicker Futures_PriceTicker
		unmarshallErr := json.Unmarshal(resp.Body, &priceTicker)
		if unmarshallErr != nil {
			return nil, resp, LocalError(PARSING_ERROR, unmarshallErr.Error())
		}

		return []*Futures_PriceTicker{&priceTicker}, resp, nil
	} else {
		var priceTickers []*Futures_PriceTicker
		unmarshallErr := json.Unmarshal(resp.Body, &priceTickers)
		if unmarshallErr != nil {
			return nil, resp, LocalError(PARSING_ERROR, unmarshallErr.Error())
		}

		return priceTickers, resp, nil
	}
}

// If the symbol is not sent, bookTickers for all symbols will be returned in an array.
// The field X-MBX-USED-WEIGHT-1M in response header is not accurate from this endpoint, please ignore.
func (futures *Futures) PriceTicker(symbol ...string) ([]*Futures_PriceTicker, *Response, *Error) {
	opts := make(map[string]interface{})

	if len(symbol) != 0 {
		opts["symbol"] = symbol[0]
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/fapi/v2/ticker/price",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	if len(symbol) == 0 {
		var priceTicker Futures_PriceTicker
		unmarshallErr := json.Unmarshal(resp.Body, &priceTicker)
		if unmarshallErr != nil {
			return nil, resp, LocalError(PARSING_ERROR, unmarshallErr.Error())
		}

		return []*Futures_PriceTicker{&priceTicker}, resp, nil
	} else {
		var priceTickers []*Futures_PriceTicker
		unmarshallErr := json.Unmarshal(resp.Body, &priceTickers)
		if unmarshallErr != nil {
			return nil, resp, LocalError(PARSING_ERROR, unmarshallErr.Error())
		}

		return priceTickers, resp, nil
	}
}

/////////////////////////////////////////////////////////////////////////////////

// If the symbol is not sent, bookTickers for all symbols will be returned in an array.
// The field X-MBX-USED-WEIGHT-1M in response header is not accurate from this endpoint, please ignore.
func (futures *Futures) BookTicker(symbol ...string) ([]*Futures_BookTicker, *Response, *Error) {
	opts := make(map[string]interface{})

	if len(symbol) != 0 {
		opts["symbol"] = symbol[0]
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/fapi/v1/ticker/bookTicker",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	if len(symbol) == 0 {
		var bookTicker Futures_BookTicker
		unmarshallErr := json.Unmarshal(resp.Body, &bookTicker)
		if unmarshallErr != nil {
			return nil, resp, LocalError(PARSING_ERROR, unmarshallErr.Error())
		}

		return []*Futures_BookTicker{&bookTicker}, resp, nil
	} else {
		var bookTickers []*Futures_BookTicker
		unmarshallErr := json.Unmarshal(resp.Body, &bookTickers)
		if unmarshallErr != nil {
			return nil, resp, LocalError(PARSING_ERROR, unmarshallErr.Error())
		}

		return bookTickers, resp, nil
	}
}

/////////////////////////////////////////////////////////////////////////////////

func (futures *Futures) DeliveryPrice(pair string) ([]*Futures_DeliveryPrice, *Response, *Error) {
	opts := make(map[string]interface{})
	opts["pair"] = pair

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/futures/data/delivery-price",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var deliveryPrices []*Futures_DeliveryPrice
	unmarshallErr := json.Unmarshal(resp.Body, &deliveryPrices)
	if unmarshallErr != nil {
		return nil, resp, LocalError(PARSING_ERROR, unmarshallErr.Error())
	}

	return deliveryPrices, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

func (futures *Futures) OpenInterest(symbol string) (*Futures_OpenInterest, *Response, *Error) {
	opts := make(map[string]interface{})
	opts["symbol"] = symbol

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/fapi/v1/openInterest",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var openInterest *Futures_OpenInterest
	unmarshallErr := json.Unmarshal(resp.Body, &openInterest)
	if unmarshallErr != nil {
		return nil, resp, LocalError(PARSING_ERROR, unmarshallErr.Error())
	}

	return openInterest, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

type Futures_OpenInterestStatistics_Params struct {
	Limit     int64
	StartTime int64
	EndTime   int64
}

func (futures *Futures) OpenInterestStatistics(symbol string, period string, opt_params ...Futures_OpenInterestStatistics_Params) ([]*Futures_OpenInterestStatistics, *Response, *Error) {
	opts := make(map[string]interface{})
	opts["symbol"] = symbol

	if len(opt_params) != 0 {
		params := opt_params[0]
		if IsDifferentFromDefault(params.Limit) {
			opts["limit"] = params.Limit
		}
		if IsDifferentFromDefault(params.StartTime) {
			opts["startTime"] = params.StartTime
		}
		if IsDifferentFromDefault(params.EndTime) {
			opts["endTime"] = params.EndTime
		}
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
		method:       Constants.Methods.GET,
		url:          "/futures/data/openInterestHist",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var openInterestStatistics []*Futures_OpenInterestStatistics
	unmarshallErr := json.Unmarshal(resp.Body, &openInterestStatistics)
	if unmarshallErr != nil {
		return nil, resp, LocalError(PARSING_ERROR, unmarshallErr.Error())
	}

	return openInterestStatistics, resp, nil
}

// //////////////////////////// Orders \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\

// # DISCLAIMER
//
// Maybe someone will have the patience to create a unified structure for the functions
// that is similar to the one I used for my javascript library "binance-lib" that is still performant
// specifically in https://github.com/GTedZ/Binance-lib/blob/main/Futures/RESTful.js#L589
func (futures *Futures) newOrder(opts map[string]interface{}) (*Futures_Order, *Response, *Error) {
	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.TRADE,
		method:       Constants.Methods.POST,
		url:          "/fapi/v1/order",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var order *Futures_Order
	processingErr := json.Unmarshal(resp.Body, &order)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERROR, processingErr.Error())
	}
	return order, resp, nil
}

type Futures_Order_Params struct {
	PositionSide            string `json:"positionSide"`
	TimeInForce             string `json:"timeInForce"`
	Quantity                string `json:"quantity"`
	ReduceOnly              bool   `json:"reduceOnly"`
	Price                   string `json:"price"`
	NewClientOrderId        string `json:"newClientOrderId"`
	StopPrice               string `json:"stopPrice"`
	ClosePosition           string `json:"closePosition"`
	ActivationPrice         string `json:"activationPrice"`
	CallbackRate            string `json:"callbackRate"`
	WorkingType             string `json:"workingType"`
	PriceProtect            string `json:"priceProtect"`
	NewOrderRespType        string `json:"newOrderRespType"`
	PriceMatch              string `json:"priceMatch"`
	SelfTradePreventionMode string `json:"selfTradePreventionMode"`
	GoodTillDate            int64  `json:"goodTillDate"`
	RecvWindow              int64  `json:"recvWindow"`
}

func (futures *Futures) NewOrder(symbol string, side string, Type string, opt_params ...Futures_Order_Params) (*Futures_Order, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol
	opts["side"] = side
	opts["type"] = Type

	if len(opt_params) != 0 {
		params := opt_params[0]
		if IsDifferentFromDefault(params.PositionSide) {
			opts["positionSide"] = params.PositionSide
		}
		if IsDifferentFromDefault(params.TimeInForce) {
			opts["timeInForce"] = params.TimeInForce
		}
		if IsDifferentFromDefault(params.Quantity) {
			opts["quantity"] = params.Quantity
		}
		if IsDifferentFromDefault(params.ReduceOnly) {
			opts["reduceOnly"] = params.ReduceOnly
		}
		if IsDifferentFromDefault(params.Price) {
			opts["price"] = params.Price
		}
		if IsDifferentFromDefault(params.NewClientOrderId) {
			opts["newClientOrderId"] = params.NewClientOrderId
		}
		if IsDifferentFromDefault(params.StopPrice) {
			opts["stopPrice"] = params.StopPrice
		}
		if IsDifferentFromDefault(params.ClosePosition) {
			opts["closePosition"] = params.ClosePosition
		}
		if IsDifferentFromDefault(params.ActivationPrice) {
			opts["activationPrice"] = params.ActivationPrice
		}
		if IsDifferentFromDefault(params.CallbackRate) {
			opts["callbackRate"] = params.CallbackRate
		}
		if IsDifferentFromDefault(params.WorkingType) {
			opts["workingType"] = params.WorkingType
		}
		if IsDifferentFromDefault(params.PriceProtect) {
			opts["priceProtect"] = params.PriceProtect
		}
		if IsDifferentFromDefault(params.NewOrderRespType) {
			opts["newOrderRespType"] = params.NewOrderRespType
		}
		if IsDifferentFromDefault(params.PriceMatch) {
			opts["priceMatch"] = params.PriceMatch
		}
		if IsDifferentFromDefault(params.SelfTradePreventionMode) {
			opts["selfTradePreventionMode"] = params.SelfTradePreventionMode
		}
		if IsDifferentFromDefault(params.GoodTillDate) {
			opts["goodTillDate"] = params.GoodTillDate
		}
		if IsDifferentFromDefault(params.RecvWindow) {
			opts["recvWindow"] = params.RecvWindow
		}
	}

	return futures.newOrder(opts)
}

///////////////////////// LIMIT \\\\\\\\\\\\\\\\\\\\\\\\\\\\

type Futures_LimitOrder_Params struct {
	PositionSide            string
	ReduceOnly              bool
	NewClientOrderId        string
	WorkingType             string
	NewOrderRespType        string
	PriceMatch              string
	SelfTradePreventionMode string
	GoodTillDate            int64
	RecvWindow              int64
}

func (futures *Futures) LimitOrder(symbol string, side string, price string, quantity string, timeInForce string, opt_params ...Futures_LimitOrder_Params) (*Futures_Order, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol
	opts["side"] = side
	opts["type"] = "LIMIT"
	opts["price"] = price
	opts["quantity"] = quantity
	opts["timeInForce"] = timeInForce

	if len(opt_params) != 0 {
		params := opt_params[0]
		if IsDifferentFromDefault(params.PositionSide) {
			opts["positionSide"] = params.PositionSide
		}
		if IsDifferentFromDefault(params.ReduceOnly) {
			opts["reduceOnly"] = params.ReduceOnly
		}
		if IsDifferentFromDefault(params.NewClientOrderId) {
			opts["newClientOrderId"] = params.NewClientOrderId
		}
		if IsDifferentFromDefault(params.WorkingType) {
			opts["workingType"] = params.WorkingType
		}
		if IsDifferentFromDefault(params.NewOrderRespType) {
			opts["newOrderRespType"] = params.NewOrderRespType
		}
		if IsDifferentFromDefault(params.PriceMatch) {
			opts["priceMatch"] = params.PriceMatch
			delete(opts, "price")
		}
		if IsDifferentFromDefault(params.SelfTradePreventionMode) {
			opts["selfTradePreventionMode"] = params.SelfTradePreventionMode
		}
		if IsDifferentFromDefault(params.GoodTillDate) {
			opts["goodTillDate"] = params.GoodTillDate
		}
		if IsDifferentFromDefault(params.RecvWindow) {
			opts["recvWindow"] = params.RecvWindow
		}
	}

	return futures.newOrder(opts)
}

func (futures *Futures) LimitBuy(symbol string, price string, quantity string, timeInForce string, opt_params ...Futures_LimitOrder_Params) (*Futures_Order, *Response, *Error) {
	return futures.LimitOrder(symbol, "BUY", price, quantity, timeInForce, opt_params...)
}

func (futures *Futures) LimitSell(symbol string, price string, quantity string, timeInForce string, opt_params ...Futures_LimitOrder_Params) (*Futures_Order, *Response, *Error) {
	return futures.LimitOrder(symbol, "SELL", price, quantity, timeInForce, opt_params...)
}

///////////////////////// LIMIT ////////////////////////////

///////////////////////// MARKET \\\\\\\\\\\\\\\\\\\\\\\\\\\

type Futures_MarketOrder_Params struct {
	PositionSide            string
	ReduceOnly              bool
	NewClientOrderId        string
	WorkingType             string
	NewOrderRespType        string
	SelfTradePreventionMode string
	RecvWindow              int64
}

func (futures *Futures) MarketOrder(symbol string, side string, quantity string, opt_params ...Futures_MarketOrder_Params) (*Futures_Order, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol
	opts["side"] = side
	opts["type"] = "MARKET"
	opts["quantity"] = quantity

	if len(opt_params) != 0 {
		params := opt_params[0]
		if IsDifferentFromDefault(params.PositionSide) {
			opts["positionSide"] = params.PositionSide
		}
		if IsDifferentFromDefault(params.ReduceOnly) {
			opts["reduceOnly"] = params.ReduceOnly
		}
		if IsDifferentFromDefault(params.NewClientOrderId) {
			opts["newClientOrderId"] = params.NewClientOrderId
		}
		if IsDifferentFromDefault(params.WorkingType) {
			opts["workingType"] = params.WorkingType
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

	return futures.newOrder(opts)
}

func (futures *Futures) MarketBuy(symbol string, quantity string, opt_params ...Futures_MarketOrder_Params) (*Futures_Order, *Response, *Error) {
	return futures.MarketOrder(symbol, "BUY", quantity, opt_params...)
}

func (futures *Futures) MarketSell(symbol string, quantity string, opt_params ...Futures_MarketOrder_Params) (*Futures_Order, *Response, *Error) {
	return futures.MarketOrder(symbol, "SELL", quantity, opt_params...)
}

///////////////////////// MARKET ///////////////////////////

// \\\\\\\\\\\\\\\\\\\\\\\\\\\ Orders ////////////////////////////////////////

// func (futures *Futures) BatchOrders(batchOrders []Futures_Order_Params) {

// }

func (futures *Futures) ChangeMarginType(symbol string, marginType string, recvWindow ...int64) (*Futures_ChangeMarginType_Response, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol
	opts["marginType"] = marginType

	if len(recvWindow) != 0 {
		opts["recvWindow"] = recvWindow[0]
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.TRADE,
		method:       Constants.Methods.POST,
		url:          "/fapi/v1/marginType",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var response *Futures_ChangeMarginType_Response
	processingErr := json.Unmarshal(resp.Body, &response)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERROR, processingErr.Error())
	}
	return response, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////////

func (futures *Futures) ChangePositionMode(toHedgeMode bool, recvWindow ...int64) (*Futures_ChangePositionMode_Response, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["dualSidePosition"] = toHedgeMode

	if len(recvWindow) != 0 {
		opts["recvWindow"] = recvWindow[0]
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.TRADE,
		method:       Constants.Methods.POST,
		url:          "/fapi/v1/positionSide/dual",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var response *Futures_ChangePositionMode_Response
	processingErr := json.Unmarshal(resp.Body, &response)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERROR, processingErr.Error())
	}
	return response, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

func (futures *Futures) ChangeInitialLeverage(symbol string, leverage int, recvWindow ...int64) (*Futures_ChangeInitialLeverage_Response, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["symbol"] = symbol
	opts["leverage"] = leverage

	if len(recvWindow) != 0 {
		opts["recvWindow"] = recvWindow[0]
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.TRADE,
		method:       Constants.Methods.POST,
		url:          "/fapi/v1/leverage",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var response *Futures_ChangeInitialLeverage_Response
	processingErr := json.Unmarshal(resp.Body, &response)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERROR, processingErr.Error())
	}
	return response, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

func (futures *Futures) ChangeMultiAssetsMode(multiAssetsMargin bool, recvWindow ...int64) (*Futures_ChangeMultiAssetsMode_Response, *Response, *Error) {
	opts := make(map[string]interface{})

	opts["multiAssetsMargin"] = multiAssetsMargin

	if len(recvWindow) != 0 {
		opts["recvWindow"] = recvWindow[0]
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.TRADE,
		method:       Constants.Methods.POST,
		url:          "/fapi/v1/multiAssetsMargin",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var response *Futures_ChangeMultiAssetsMode_Response
	processingErr := json.Unmarshal(resp.Body, &response)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERROR, processingErr.Error())
	}
	return response, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////

func (futures *Futures) AccountInfo(recvWindow ...int64) (*Futures_AccountInfo, *Response, *Error) {
	opts := make(map[string]interface{})

	if len(recvWindow) != 0 {
		opts["recvWindow"] = recvWindow[0]
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.USER_DATA,
		method:       Constants.Methods.GET,
		url:          "/fapi/v3/account",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var accountInfo *Futures_AccountInfo
	processingErr := json.Unmarshal(resp.Body, &accountInfo)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERROR, processingErr.Error())
	}
	return accountInfo, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

func (futures *Futures) AccountConfiguration(recvWindow ...int64) (*Futures_AccountConfiguration, *Response, *Error) {
	opts := make(map[string]interface{})

	if len(recvWindow) != 0 {
		opts["recvWindow"] = recvWindow[0]
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.USER_DATA,
		method:       Constants.Methods.GET,
		url:          "/fapi/v1/accountConfig",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	var accountConfig *Futures_AccountConfiguration
	processingErr := json.Unmarshal(resp.Body, &accountConfig)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERROR, processingErr.Error())
	}
	return accountConfig, resp, nil
}

/////////////////////////////////////////////////////////////////////////////////

func (futures *Futures) LeverageBrackets(symbol ...string) ([]*Futures_LeverageBrackets, *Response, *Error) {
	opts := make(map[string]interface{})

	if len(symbol) != 0 {
		opts["symbol"] = symbol[0]
	}

	resp, err := futures.makeRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.USER_DATA,
		method:       Constants.Methods.GET,
		url:          "/fapi/v1/leverageBracket",
		params:       opts,
	})
	if err != nil {
		return nil, resp, err
	}

	if len(symbol) != 0 {
		var leverageBrackets *Futures_LeverageBrackets
		processingErr := json.Unmarshal(resp.Body, &leverageBrackets)
		if processingErr != nil {
			return nil, resp, LocalError(PARSING_ERROR, processingErr.Error())
		}
		return []*Futures_LeverageBrackets{leverageBrackets}, resp, nil
	}

	var leverageBrackets []*Futures_LeverageBrackets
	processingErr := json.Unmarshal(resp.Body, &leverageBrackets)
	if processingErr != nil {
		return nil, resp, LocalError(PARSING_ERROR, processingErr.Error())
	}
	return leverageBrackets, resp, nil
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
