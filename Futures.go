package Binance

import (
	"fmt"
	"time"
)

type Futures struct {
	binance       *Binance
	requestClient RequestClient
	baseURL       string

	API APIKEYS

	// Websockets Futures_Websockets
}

func (futures *Futures) init(binance *Binance) {
	futures.binance = binance

	futures.requestClient.init(binance)
	futures.requestClient.Set_APIKEY(binance.API.KEY, binance.API.SECRET)

	futures.API.Set(binance.API.KEY, binance.API.SECRET)
}

/////////////////////////////////////////////////////////////////////////////////

// # Test connectivity to the Rest API.
//
// Weight: 1
//
// Data Source: Memory
func (futures *Futures) Ping() (latency int64, request *Response, err *Error) {
	startTime := time.Now().UnixMilli()
	httpResp, err := futures.makeFuturesRequest(&FuturesRequest{
		securityType: FUTURES_Constants.SecurityTypes.NONE,
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

/////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////

type FuturesRequest struct {
	method       string
	url          string
	params       map[string]interface{}
	securityType string
}

func (futures *Futures) makeFuturesRequest(request *FuturesRequest) (*Response, *Error) {

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
