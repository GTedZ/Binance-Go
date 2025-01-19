package Binance

import (
	"fmt"
	"os"
	"strconv"
	"time"
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

func newError(isLocal bool, statusCode int, code int, message string) *Error {
	err := &Error{
		IsLocalError: isLocal,
		StatusCode:   statusCode,
		Code:         code,
		Message:      message,
	}

	if PRINT_ERRORS {
		fmt.Println(err.Error())
	}

	if LOG_ERRORS && LOG_ERRORS_FILE != "" {
		logErrorToFile(err)
	}

	return err
}

func logErrorToFile(err *Error) {
	file, errOpen := os.OpenFile(LOG_ERRORS_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errOpen != nil {
		fmt.Println("Failed to open log file:", errOpen)
		return
	}
	defer file.Close()

	timestamp := time.Now()
	logEntry := fmt.Sprintf("%d - %s: %s\n",
		timestamp.UnixMilli(),
		timestamp.Format("02/01/2006 15:04:05.000"),
		err.Error(),
	)

	if _, errWrite := file.WriteString(logEntry); errWrite != nil {
		fmt.Println("Failed to write to log file:", errWrite)
	}
}

func LocalError(code int, msg string) *Error {
	return newError(true, 0, code, msg)
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

	err := newError(false, resp.StatusCode, errResponse.Code, errResponse.Msg)

	return err, nil
}
