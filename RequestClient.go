package Binance

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

type RequestClient struct {
	options *BinanceOptions
	configs *BinanceConfig

	client http.Client
	api    APIKEYS
}

type Response struct {
	Status     string // e.g. "200 OK"
	StatusCode int    // e.g. 200
	Proto      string // e.g. "HTTP/1.0"
	ProtoMajor int    // e.g. 1
	ProtoMinor int    // e.g. 0

	// Header maps header keys to values. If the response had multiple
	// headers with the same key, they may be concatenated, with comma
	// delimiters.  (RFC 7230, section 3.2.2 requires that multiple headers
	// be semantically equivalent to a comma-delimited sequence.) When
	// Header values are duplicated by other fields in this struct (e.g.,
	// ContentLength, TransferEncoding, Trailer), the field values are
	// authoritative.
	//
	// Keys in the map are canonicalized (see CanonicalHeaderKey).
	Header http.Header

	Body []byte

	// ContentLength records the length of the associated content. The
	// value -1 indicates that the length is unknown. Unless Request.Method
	// is "HEAD", values >= 0 indicate that the given number of bytes may
	// be read from Body.
	ContentLength int64

	// Contains transfer encodings from outer-most to inner-most. Value is
	// nil, means that "identity" encoding is used.
	TransferEncoding []string

	// Close records whether the header directed that the connection be
	// closed after reading Body. The value is advice for clients: neither
	// ReadResponse nor Response.Write ever closes a connection.
	Close bool

	// Uncompressed reports whether the response was sent compressed but
	// was decompressed by the http package. When true, reading from
	// Body yields the uncompressed content instead of the compressed
	// content actually set from the server, ContentLength is set to -1,
	// and the "Content-Length" and "Content-Encoding" fields are deleted
	// from the responseHeader. To get the original response from
	// the server, set Transport.DisableCompression to true.
	Uncompressed bool

	// Trailer maps trailer keys to values in the same
	// format as Header.
	//
	// The Trailer initially contains only nil values, one for
	// each key specified in the server's "Trailer" header
	// value. Those values are not added to Header.
	//
	// Trailer must not be accessed concurrently with Read calls
	// on the Body.
	//
	// After Body.Read has returned io.EOF, Trailer will contain
	// any trailer values sent by the server.
	Trailer http.Header

	// Request is the request that was sent to obtain this Response.
	// Request's Body is nil (having already been consumed).
	// This is only populated for Client requests.
	Request *http.Request

	// TLS contains information about the TLS connection on which the
	// response was received. It is nil for unencrypted responses.
	// The pointer is shared between responses and should not be
	// modified.
	TLS *tls.ConnectionState
}

// # Fetches the current used weight returned the request.
//
// interval: "1m", "3m", "1d", "1W", "1M", or simply ""
//
// But most common and only one used as of writing this is "1m"
//
// Returns an error if the header is not found
func (resp *Response) GetUsedWeight(interval string) (int64, *Error) {
	key := "X-Mbx-Used-Weight"
	if interval != "" {
		key += "-" + interval
	}

	strValue := resp.Header.Get(key)

	if strValue == "" {
		errStr := "No Used Weight was found for this interval"
		if PRINT_ERRORS {
			fmt.Println(errStr)
		}
		return 0, LocalError(RESPONSE_HEADER_NOT_FOUND, "")
	}

	// Parses the value to int64
	value, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		fmt.Println("Error parsing header value:", err)
		return 0, LocalError(PARSING_ERROR, err.Error())
	}

	return value, nil
}

func (resp *Response) GetRequestTime() (time.Time, *Error) {
	key := "Date"

	strValue := resp.Header.Get(key)

	if strValue == "" {
		errStr := "No Date header was found for this request"
		if PRINT_ERRORS {
			fmt.Println(errStr)
		}
		return time.Now(), LocalError(RESPONSE_HEADER_NOT_FOUND, "")
	}

	parsedTime, err := time.Parse(time.RFC1123, strValue)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return time.Now(), LocalError(PARSING_ERROR, "There was an error parsing the date from request headers")
	}

	return parsedTime, nil
}

//

func (requestClient *RequestClient) init(options *BinanceOptions, configs *BinanceConfig) {
	requestClient.options = options
	requestClient.configs = configs
}

func (requestClient *RequestClient) Set_APIKEY(APIKEY string, APISECRET string) {
	requestClient.api.KEY = APIKEY
	requestClient.api.SECRET = APISECRET
}

//

func readResponseBody(rawResponse *http.Response) (*Response, error) {
	var resp Response

	data, err := io.ReadAll(rawResponse.Body)
	if err != nil {
		return nil, err
	}

	resp.Body = data
	resp.Close = rawResponse.Close
	resp.ContentLength = rawResponse.ContentLength
	resp.Header = rawResponse.Header
	resp.Proto = rawResponse.Proto
	resp.ProtoMajor = rawResponse.ProtoMajor
	resp.ProtoMinor = rawResponse.ProtoMinor
	resp.Request = rawResponse.Request
	resp.Status = rawResponse.Status
	resp.StatusCode = rawResponse.StatusCode
	resp.TLS = rawResponse.TLS
	resp.Trailer = rawResponse.Trailer
	resp.TransferEncoding = rawResponse.TransferEncoding
	resp.Uncompressed = rawResponse.Uncompressed

	return &resp, nil
}

// createQueryString transforms a map[string]interface{} into a query string
func createQueryString(params map[string]interface{}, sorted bool) string {
	if params == nil {
		return ""
	}

	// Extract keys to sort them if `sorted` is true
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}

	if sorted {
		sort.Strings(keys)
	}

	query := url.Values{}

	// Helper function to process values
	var addToQuery func(key string, value interface{})
	addToQuery = func(key string, value interface{}) {
		switch v := value.(type) {
		case string:
			query.Add(key, v)
		case []string:
			// Encode slices as JSON arrays
			jsonValue, err := json.Marshal(v)
			if err != nil {
				if VERBOSE {
					fmt.Printf("[VERBOSE] Error marshaling slice for key %s: %v\n", key, err)
				}
				return
			}
			query.Add(key, string(jsonValue)) // Add JSON-encoded array
		case []interface{}:
			for _, item := range v {
				addToQuery(key, item) // Recursively handle each item
			}
		case map[string]interface{}:
			// Handle nested maps with dot notation
			for subKey, subValue := range v {
				addToQuery(key+"."+subKey, subValue)
			}
		case int, int64, float64, bool: // Convert basic types to string
			query.Add(key, fmt.Sprintf("%v", v))
		default:
			if VERBOSE {
				fmt.Println("[VERBOSE] Error adding parameter, invalid type detected: received:", v)
			}
		}
	}

	// Process each key-value pair
	for _, key := range keys {
		addToQuery(key, params[key])
	}

	return query.Encode()
}

//

func (requestClient *RequestClient) Unsigned(method string, baseURL string, URL string, params map[string]interface{}) (*Response, *Error) {
	var err error
	var rawResponse *http.Response

	paramString := createQueryString(params, false)

	fullQuery := baseURL + URL + "?" + paramString

	switch method {
	case Constants.Methods.GET:
		rawResponse, err = http.Get(fullQuery)

	default:
		panic(fmt.Sprintf("Method passed to Unsigned Request function is invalid, received: '%s'\nSupported methods are ('%s', '%s', '%s', '%s', '%s')", method, Constants.Methods.GET, Constants.Methods.POST, Constants.Methods.PUT, Constants.Methods.PATCH, Constants.Methods.DELETE))
	}
	if err != nil {
		if PRINT_ERRORS {
			fmt.Println("[VERBOSE] Request error:", err)
		}
		Err := Error{
			IsLocalError: true,
			Code:         HTTP_REQUEST_ERR,
			Message:      err.Error(),
		}
		return nil, &Err
	}
	defer rawResponse.Body.Close()

	resp, err := readResponseBody(rawResponse)
	if err != nil {
		if PRINT_ERRORS {
			fmt.Println("[VERBOSE] Error reading response body:", err)
		}
		Err := Error{
			IsLocalError: true,
			Code:         RESPONSEBODY_READING_ERR,
			Message:      err.Error(),
		}
		return nil, &Err
	}

	if PRINT_HTTP_RESPONSES {
		fmt.Printf("%s %s: %s =>\nResponse: %s\n", resp.Request.Method, resp.Status, fullQuery, string(resp.Body))
	} else if PRINT_HTTP_QUERIES {
		fmt.Printf("%s %s: %s\n", resp.Request.Method, resp.Status, fullQuery)
	}

	if resp.StatusCode >= 400 {
		Err, UnmarshallErr := BinanceError(resp)
		if UnmarshallErr != nil {
			if PRINT_ERRORS {
				fmt.Println("[VERBOSE] Error processing error response body:", UnmarshallErr)
			}
			return nil, UnmarshallErr
		}

		return resp, Err
	}

	return resp, nil
}

// func (requestClient *RequestClient) UserStream() (resp *http.Response, err error) {
// 	requestClient.client.Do()

// 	return resp, nil
// }

// func (requestClient *RequestClient) Signed() (resp *http.Response, err error) {

// }
