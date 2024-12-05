package main

type BinanceOptions struct {
	printErrors           bool
	verbose               bool
	updateTimestampOffset bool
	recvWindow            int64
}

func (options *BinanceOptions) init() {
	options.printErrors = false
	options.verbose = false
	options.updateTimestampOffset = false
	options.recvWindow = 5000
}

func (options *BinanceOptions) Set_PrintErrors(value bool) {
	options.printErrors = value
}

func (options *BinanceOptions) Set_Verbose(value bool) {
	options.verbose = value
}

func (options *BinanceOptions) Set_UpdateTimestampOffset(value bool) {
	options.updateTimestampOffset = value
}

func (options *BinanceOptions) Set_recvWindow(recvWindow int64) {
	options.recvWindow = recvWindow
}
