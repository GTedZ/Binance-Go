package main

type Binance struct {
	configs BinanceConfig
	Options BinanceOptions

	API APIKEYS

	Spot    Spot
	Futures Futures
}

func CreateClient(APIKEY string, APISECRET string) Binance {
	binance := Binance{}

	binance.configs.init()

	binance.Options.init()
	binance.API.Set(APIKEY, APISECRET)

	binance.Spot.init(&binance)
	binance.Futures.init(&binance)

	return binance
}

func CreateClientWithOptions(APIKEY string, APISECRET string, recvWindow int64) Binance {
	binance := CreateClient(APIKEY, APISECRET)

	binance.Options.Set_recvWindow(recvWindow)

	return binance
}
