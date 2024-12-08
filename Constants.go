package Binance

const SECOND = 1000
const MINUTE = 60 * SECOND
const HOUR = 60 * MINUTE
const DAY = 24 * HOUR

var Constants = struct {
	Methods   Methods
	Websocket WebsocketConstants
}{
	Methods: Methods{
		GET:    "GET",
		POST:   "POST",
		PUT:    "PUT",
		PATCH:  "PATCH",
		DELETE: "DELETE",
	},
	Websocket: WebsocketConstants{
		MAX_STREAMS_PER_SOCKET:              1024,
		MAX_OUTGOING_MESSAGES_PER_SECOND:    5,
		CONNECTION_ATTEMPTS_PER_5MINS:       300,
		RESPONSE_TIMEOUT_SECONDS:            20,
		HEARTBEAT_CHECK_INTERVAL_SEC:        5,
		HEARTBEAT_CLOSE_ON_NO_HEARTBEAT_SEC: 20,
		EXPECTED_DISCONNECTION_TIME_SEC:     (DAY - 5*MINUTE) / 1000,
	},
}

type Methods struct {
	GET    string
	POST   string
	PUT    string
	PATCH  string
	DELETE string
}

type WebsocketConstants struct {
	MAX_STREAMS_PER_SOCKET              uint64
	MAX_OUTGOING_MESSAGES_PER_SECOND    uint64
	CONNECTION_ATTEMPTS_PER_5MINS       uint64
	RESPONSE_TIMEOUT_SECONDS            int64
	HEARTBEAT_CHECK_INTERVAL_SEC        int64
	HEARTBEAT_CLOSE_ON_NO_HEARTBEAT_SEC int64
	EXPECTED_DISCONNECTION_TIME_SEC     int64
}
