package Binance

import (
	"fmt"
	"strconv"
)

type Error struct {
	// false => Error originating from binance's side
	// true =>
	IsLocalError bool

	StatusCode int

	Code int

	Message string
}

// Implement the `Error` method to satisfy the `error` interface
func (e *Error) Error() string {
	str := "[BINANCE ERROR] StatusCode " + strconv.Itoa(e.StatusCode)
	if e.IsLocalError {
		str = "[LIB ERROR]"
	}
	return fmt.Sprintf("%s - Code %d: \"%s\"", str, e.Code, e.Message)
}

const (
	HTTP_REQUEST_ERR         = -1
	RESPONSEBODY_READING_ERR = -2
	ERROR_PROCESSING_ERR     = -3
	RESPONSE_PROCESSING_ERR  = -4
)

func LocalError(code int, msg string) *Error {
	return &Error{
		IsLocalError: true,
		Code:         code,
		Message:      msg,
	}
}

type BinanceErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func BinanceError(resp *Response) (BinanceError *Error, UnmarshallError *Error) {
	var errResponse BinanceErrorResponse
	err := json.Unmarshal(resp.Body, &errResponse)
	if err != nil {
		return nil,
			LocalError(ERROR_PROCESSING_ERR, err.Error())
	}

	return &Error{
			IsLocalError: false,
			StatusCode:   resp.StatusCode,
			Code:         errResponse.Code,
			Message:      errResponse.Msg,
		},
		nil
}
