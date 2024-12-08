package Binance

type BinanceConfig struct {
	timestamp_offset int64
}

func (config *BinanceConfig) init() {
	config.timestamp_offset = 0
}
