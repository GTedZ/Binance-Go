package main

import (
	"fmt"
)

// 3287788908 97274.25000000 0.00014000 4175708079 4175708079 1732904929081 true true

func main() {
	binance := CreateClientWithOptions("Hello", "Hel", 5000)

	data, _, err := binance.Spot.Candlesticks("BTCUSDT", "1s", &Spot_Candlesticks_Params{Limit: 999})
	if err != nil {
		fmt.Println(err)
		return
	}

	// for _, aggTrade := range data {
	// 	fmt.Println(aggTrade)
	// }

	fmt.Println(len(data))
}
