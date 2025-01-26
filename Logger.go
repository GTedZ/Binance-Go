package Binance

import "fmt"

func LOG_HTTP_QUERIES(a ...any) {
	if DevOptions.PRINT_HTTP_QUERIES {
		fmt.Println("[LIB][PRINT_HTTP_QUERIES]", fmt.Sprintln(a...))
	}
}

func LOG_HTTP_RESPONSES(a ...any) {
	if DevOptions.PRINT_HTTP_RESPONSES {
		fmt.Println("[LIB][PRINT_HTTP_RESPONSES]", fmt.Sprintln(a...))
	}
}

func LOG_ERRORS(a ...any) {
	if DevOptions.PRINT_ERRORS {
		fmt.Println("[LIB][PRINT_ERRORS]", fmt.Sprintln(a...))
	}
}

func LOG_ALL_ERRORS(a ...any) {
	if DevOptions.PRINT_ALL_ERRORS {
		fmt.Println("[LIB][PRINT_ALL_ERRORS]", fmt.Sprintln(a...))
	}
}

func LOG_WS_VERBOSE(a ...any) {
	if DevOptions.WS_VERBOSE {
		fmt.Println("[LIB][WS_VERBOSE]", fmt.Sprintln(a...))
	}
}

func LOG_WS_VERBOSE_FULL(a ...any) {
	if DevOptions.WS_VERBOSE_FULL {
		fmt.Println("[LIB][WS_VERBOSE_FULL]", fmt.Sprintln(a...))
	}
}

func LOG_WS_ERRORS(a ...any) {
	if DevOptions.PRINT_WS_ERRORS {
		fmt.Println("[LIB][PRINT_WS_ERRORS]", fmt.Sprintln(a...))
	}
}

func LOG_WS_MESSAGES(a ...any) {
	if DevOptions.PRINT_WS_MESSAGES {
		fmt.Println("[LIB][PRINT_WS_MESSAGES]", fmt.Sprintln(a...))
	}
}
