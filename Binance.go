package Binance

type Binance struct {
	configs BinanceConfig
	Opts    BinanceOptions

	API APIKEYS

	BinaryUtils BinaryUtils

	Spot    Spot
	Futures Futures
}

func CreateReadClient() Binance {
	binance := Binance{}

	binance.configs.init()
	binance.Opts.init()

	binance.Spot.init(&binance)
	binance.Futures.init(&binance)

	return binance
}

func CreateClient(APIKEY string, APISECRET string) Binance {
	binance := CreateReadClient()
	binance.API.Set(APIKEY, APISECRET)
	return binance
}

func CreateClientWithOptions(APIKEY string, APISECRET string, recvWindow int64) Binance {
	binance := CreateClient(APIKEY, APISECRET)

	binance.Opts.Set_recvWindow(recvWindow)

	return binance
}
