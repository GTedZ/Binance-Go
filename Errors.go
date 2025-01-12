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
	HTTP_REQUEST_ERR = iota
	HTTP_SIGNATURE_ERR
	RESPONSEBODY_READING_ERR
	ERROR_PROCESSING_ERR
	RESPONSE_HEADER_NOT_FOUND
	PARSING_ERROR
	WS_OPEN_ERR
	WS_SEND_MESSAGE_ERR
	REQUEST_TIMEOUT_ERR
)

func LocalError(code int, msg string) *Error {
	err := &Error{
		IsLocalError: true,
		Code:         code,
		Message:      msg,
	}

	if PRINT_ERRORS {
		fmt.Println(err.Error())
	}

	return err
}

type BinanceErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Processes an erroneous 4XX HTTP Response
// Returns the library Error type
// In the case of an error parsing the error body, it returns a secondaly unmarshall error
func BinanceError(resp *Response) (BinanceError *Error, UnmarshallError *Error) {
	var errResponse BinanceErrorResponse

	unmarshallErr := json.Unmarshal(resp.Body, &errResponse)
	if unmarshallErr != nil {
		return nil,
			LocalError(ERROR_PROCESSING_ERR, unmarshallErr.Error())
	}

	err := &Error{
		IsLocalError: false,
		StatusCode:   resp.StatusCode,
		Code:         errResponse.Code,
		Message:      errResponse.Msg,
	}

	if PRINT_ERRORS {
		fmt.Println(err.Error())
	}

	return err, nil
}
