package main

import (
	"Binance"
	"fmt"
)

func main() {
	binance := Binance.CreateReadClient()
	socket, _ := binance.Spot.Websockets.AggTrade(func(aggTrade *Binance.SpotWS_AggTrade) { fmt.Println(aggTrade) }, "BTCUSDT", "ETHUSDT", "XRPUSDT")

	socket.Unsubscribe("BTCUSDT", "ETHUSDT")
	socket.Subscribe("XAIUSDT", "REIUSDT", "GOUSDT")
	socket.Unsubscribe("XAIUSDT", "XRPUSDT")
	socket.Subscribe("BTCUSDT", "ETHUSDT", "XRPUSDT", "XAIUSDT", "FLOKIUSDT")

	select {}
}
