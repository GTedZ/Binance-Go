package Binance

type BinanceOptions struct {
	updateTimestampOffset bool
	recvWindow            int64
}

func (options *BinanceOptions) init() {
	options.updateTimestampOffset = false
	options.recvWindow = 5000
}

func (options *BinanceOptions) Set_UpdateTimestampOffset(value bool) {
	options.updateTimestampOffset = value
}

func (options *BinanceOptions) Set_recvWindow(recvWindow int64) {
	options.recvWindow = recvWindow
}
