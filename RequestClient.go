package Binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

type RequestClient struct {
	binance *Binance

	client *http.Client
	api    APIKEYS
}

type Response struct {
	Status     string // e.g. "200 OK"
	StatusCode int    // e.g. 200
	Proto      string // e.g. "HTTP/1.0"
	ProtoMajor int    // e.g. 1
	ProtoMinor int    // e.g. 0

	// This is added by the 'Binance-Go' library
	// It is simply the elapsed time between sending the request and receiving the response
	Latency int64

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
		LOG_ERRORS(errStr)
		return 0, LocalError(RESPONSE_HEADER_NOT_FOUND, errStr)
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
		LOG_ERRORS(errStr)
		return time.Now(), LocalError(RESPONSE_HEADER_NOT_FOUND, errStr)
	}

	parsedTime, err := time.Parse(time.RFC1123, strValue)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return time.Now(), LocalError(PARSING_ERROR, "There was an error parsing the date from request headers")
	}

	return parsedTime, nil
}

func (resp *Response) GetLatency() (latency int64, err *Error) {
	if resp == nil {
		return 0, LocalError(PARSING_ERROR, "Cannot read latency from nil response")
	}
	return resp.Latency, nil
}

//

func (requestClient *RequestClient) init(binance *Binance) {
	requestClient.binance = binance
	requestClient.client = &http.Client{}
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
				LOG_ERRORS(fmt.Sprintf("[VERBOSE] Error marshaling slice for key %s: %v\n", key, err))
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
			LOG_ERRORS("[VERBOSE] Error adding parameter, invalid type detected: received:", v)
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

	startTime := time.Now().UnixMilli()
	switch method {
	case Constants.Methods.GET:
		rawResponse, err = http.Get(fullQuery)

	default:
		panic(fmt.Sprintf("Method passed to Unsigned Request function is invalid, received: '%s'\nSupported methods are ('%s', '%s', '%s', '%s', '%s')", method, Constants.Methods.GET, Constants.Methods.POST, Constants.Methods.PUT, Constants.Methods.PATCH, Constants.Methods.DELETE))
	}
	if err != nil {
		LOG_ERRORS("[VERBOSE] Request error:", err)
		return nil, LocalError(HTTP_REQUEST_ERR, err.Error())
	}
	defer rawResponse.Body.Close()
	latency := time.Now().UnixMilli() - startTime

	resp, err := readResponseBody(rawResponse)
	if err != nil {
		LOG_ERRORS("[VERBOSE] Error reading response body:", err)
		return nil, LocalError(RESPONSEBODY_READING_ERR, err.Error())
	}
	resp.Latency = latency

	LOG_HTTP_QUERIES(fmt.Sprintf("%s %s: %s\n", resp.Request.Method, resp.Status, fullQuery))
	LOG_HTTP_RESPONSES(fmt.Sprintf("%s %s: %s =>\nResponse: %s\n", resp.Request.Method, resp.Status, fullQuery, string(resp.Body)))

	if resp.StatusCode >= 400 {
		Err, UnmarshallErr := BinanceError(resp)
		if UnmarshallErr != nil {
			LOG_ERRORS("[VERBOSE] Error processing error response body:", UnmarshallErr)
			return nil, UnmarshallErr
		}

		return resp, Err
	}

	return resp, nil
}

func (requestClient *RequestClient) APIKEY_only(method string, baseURL string, URL string, params map[string]interface{}) (*Response, *Error) {

	paramString := createQueryString(params, false)

	fullQuery := baseURL + URL + "?" + paramString

	req, err := http.NewRequest(method, fullQuery, nil)
	if err != nil {
		return nil, LocalError(HTTP_REQUEST_ERR, err.Error())
	}

	req.Header.Set("X-MBX-APIKEY", requestClient.api.KEY)

	rawResponse, err := requestClient.client.Do(req)
	if err != nil {
		LOG_ERRORS("[VERBOSE] Request error:", err)
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
		LOG_ERRORS("[VERBOSE] Error reading response body:", err)
		Err := Error{
			IsLocalError: true,
			Code:         RESPONSEBODY_READING_ERR,
			Message:      err.Error(),
		}
		return nil, &Err
	}

	LOG_HTTP_QUERIES(fmt.Sprintf("%s %s: %s\n", resp.Request.Method, resp.Status, fullQuery))
	LOG_HTTP_RESPONSES(fmt.Sprintf("%s %s: %s =>\nResponse: %s\n", resp.Request.Method, resp.Status, fullQuery, string(resp.Body)))

	if resp.StatusCode >= 400 {
		Err, UnmarshallErr := BinanceError(resp)
		if UnmarshallErr != nil {
			LOG_ERRORS("[VERBOSE] Error processing error response body:", UnmarshallErr)
			return nil, UnmarshallErr
		}

		return resp, Err
	}

	return resp, nil
}

func (requestClient *RequestClient) Signed(method string, baseURL string, URL string, params map[string]interface{}) (*Response, *Error) {

	params["timestamp"] = time.Now().UnixMilli() + requestClient.binance.configs.timestamp_offset

	if requestClient.binance.Opts.recvWindow != 5000 && params["recvWindow"] == nil {
		params["recvWindow"] = requestClient.binance.Opts.recvWindow
	}

	paramString := createQueryString(params, false)

	h := hmac.New(sha256.New, []byte(requestClient.api.SECRET))
	_, err := h.Write([]byte(paramString))
	if err != nil {
		return nil, LocalError(HTTP_SIGNATURE_ERR, err.Error())
	}

	signature := hex.EncodeToString(h.Sum(nil))

	fullQuery := baseURL + URL + "?" + paramString + "&signature=" + signature

	req, err := http.NewRequest(method, fullQuery, nil)
	if err != nil {
		return nil, LocalError(HTTP_REQUEST_ERR, err.Error())
	}

	req.Header.Set("X-MBX-APIKEY", requestClient.api.KEY)

	rawResponse, err := requestClient.client.Do(req)
	if err != nil {
		Err := Error{
			IsLocalError: true,
			Code:         HTTP_REQUEST_ERR,
			Message:      err.Error(),
		}
		LOG_ERRORS("[VERBOSE] Request error:", Err.Error())
		return nil, &Err
	}
	defer rawResponse.Body.Close()

	resp, err := readResponseBody(rawResponse)
	if err != nil {
		Err := Error{
			IsLocalError: true,
			Code:         RESPONSEBODY_READING_ERR,
			Message:      err.Error(),
		}
		LOG_ERRORS("[VERBOSE] Error reading response body:", Err.Error())
		return nil, &Err
	}

	LOG_HTTP_RESPONSES(fmt.Sprintf("%s %s: %s =>\nResponse: %s\n", resp.Request.Method, resp.Status, fullQuery, string(resp.Body)))
	LOG_HTTP_QUERIES(fmt.Sprintf("%s %s: %s\n", resp.Request.Method, resp.Status, fullQuery))

	if resp.StatusCode >= 400 {
		Err, UnmarshallErr := BinanceError(resp)
		if UnmarshallErr != nil {
			LOG_ERRORS("[VERBOSE] Error processing error response body:", UnmarshallErr.Error())
			return nil, UnmarshallErr
		}

		return resp, Err
	}

	return resp, nil
}
