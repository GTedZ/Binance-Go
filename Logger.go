package Binance

import "fmt"

func LOG_ERR(a ...any) {
	if PRINT_ERRORS {
		fmt.Println(a...)
	}
}

func LOG_VERBOSE(a ...any) {
	if VERBOSE {
		fmt.Println(a...)
	}
}

func LOG_WS_VERBOSE(a ...any) {
	if WS_VERBOSE {
		fmt.Println(a...)
	}
}

func LOG_WS_VERBOSE_FULL(a ...any) {
	if WS_VERBOSE_FULL {
		fmt.Println(a...)
	}
}
