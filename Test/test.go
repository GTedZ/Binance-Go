package main

import (
	"Binance"
	"fmt"
	"time"
)

func main() {
	socket, err := Binance.CreateSocket("wss://fstream.binance.com", []string{"btcusdt@kline_1m"}, false)

	if err != nil {
		fmt.Println("There was an error opening the socket:", err)
	}

	time.Sleep(time.Second * 2)

	socket.Reconnect()

	select {}
}
