package Binance

type Futures struct {
	binance       *Binance
	requestClient RequestClient

	API APIKEYS
}

func (futures *Futures) init(binance *Binance) {
	futures.binance = binance

	futures.requestClient.init(&binance.Options, &binance.configs)
	futures.requestClient.Set_APIKEY(binance.API.KEY, binance.API.SECRET)

	futures.API.Set(binance.API.KEY, binance.API.SECRET)
}
