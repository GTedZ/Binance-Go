package Binance

import (
	"fmt"
	"strings"
	"time"

	ws "github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

type Websocket struct {
	Conn *ws.Conn
	// Host server's URL
	BaseURL string
	Streams []string

	// This is to show the current state of the stream
	// false -> it's a raw stream
	// true -> it's a combined stream
	IsCombined bool

	// This happens when a response for a request has been received
	// This will get called first, even before the requesting function receives a response
	OnPrivateMessage func(msg []byte)
	OnMessage        func(messageType int, msg []byte)
	OnPing           func(appData string)
	OnPong           func(appData string)
	OnDisconnect     func(code int, text string)
	OnReconnecting   func()
	OnReconnect      func()
	OnClose          func(code int, text string)

	Creation_Timestamp       int64
	Last_Heartbeat_Timestamp int64

	isReconnecting bool

	reconnect bool
	closed    bool

	privateMessageValidator func(msg []byte) (isPrivate bool, Id string)
	pendingRequests         map[string]chan []byte

	subsocket_id int64
}

type CombinedStream_MSG struct {
	Stream string              `json:"stream"`
	Data   jsoniter.RawMessage `json:"data"`
}

func CreateSocket(baseURL string, streams []string, isCombined bool) (*Websocket, *Error) {
	if !isCombined && len(streams) > 1 {
		isCombined = true
	}

	queryStr := CreateQueryStringWS(streams, isCombined)

	fullStreamStr := baseURL + queryStr

	LOG_WS_VERBOSE_FULL("Connecting Websocket to:", fullStreamStr)

	conn, _, err := ws.DefaultDialer.Dial(fullStreamStr, nil)
	if err != nil {
		LOG_WS_VERBOSE("There was an error creating websocket:", err)

		return nil, LocalError(WS_OPEN_ERR, err.Error())
	}

	LOG_WS_VERBOSE("Socket connected:", queryStr)

	currentTime := time.Now().Unix()

	websocket := &Websocket{
		Conn:                     conn,
		BaseURL:                  baseURL,
		Streams:                  streams,
		IsCombined:               isCombined,
		Creation_Timestamp:       currentTime,
		Last_Heartbeat_Timestamp: currentTime,
		reconnect:                true,
		closed:                   false,
		pendingRequests:          make(map[string]chan []byte),
	}

	setUpSocket(websocket, conn)

	return websocket, nil
}

func setUpSocket(websocket *Websocket, conn *ws.Conn) {
	conn.SetCloseHandler(websocket.CloseHandler)
	conn.SetPingHandler(websocket.PingHandler)
	conn.SetPongHandler(websocket.PongHandler)

	// Handle system interrupts to close the connection gracefully
	// interrupt := make(chan os.Signal, 1)
	// signal.Notify(interrupt, os.Interrupt)
	// Not sure how to best do this, so will leave it empty for now
	// TODO

	current_subSocketID := websocket.subsocket_id

	// Goroutine to read messages
	go func() {
		for {
			if websocket.closed {
				return
			}
			if websocket.isReconnecting {
				continue
			}
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				LOG_WS_VERBOSE("Error reading message:", err)

				if current_subSocketID != websocket.subsocket_id {
					return
				}
				time.Sleep(500 * time.Millisecond)
				continue
			}
			if PRINT_WS_MESSAGES {
				fmt.Printf("Type: %d, message: %s\n", msgType, string(msg))
			}

			websocket.RecordLastHeartbeat()

			if websocket.privateMessageValidator != nil {
				isPrivate, Id := websocket.privateMessageValidator(msg)
				if isPrivate {
					LOG_WS_VERBOSE("[VERBOSE] Private Message detected:", string(msg))

					websocket.pendingRequests[Id] <- msg
					if websocket.OnPrivateMessage != nil {
						websocket.OnPrivateMessage(msg)
					}
					delete(websocket.pendingRequests, Id)
					continue
				}
			}

			if websocket.IsCombined {
				var tempData CombinedStream_MSG
				err := json.Unmarshal(msg, &tempData)
				if err != nil {
					fmt.Println("ERR PARSING COMBINED MESSAGE!!!")
					panic(fmt.Sprintln(LocalError(PARSING_ERROR, err.Error()).Error()))
				}
				msg = tempData.Data
			}

			if websocket.OnMessage != nil {
				websocket.OnMessage(msgType, msg)
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Duration(Constants.Websocket.HEARTBEAT_CHECK_INTERVAL_SEC) * time.Second)
		defer ticker.Stop()

		for {

			<-ticker.C // Wait the appropriate amount of time
			if websocket.closed {
				LOG_WS_VERBOSE("[HEARTBEAT] Websocket is closed, skipping check.")

				return
			}

			if websocket.isReconnecting {
				continue
			}

			currentTime := time.Now().Unix()

			elapsed := currentTime - websocket.Last_Heartbeat_Timestamp

			// Check if the last heartbeat is older than the close interval
			if elapsed >= Constants.Websocket.HEARTBEAT_CLOSE_ON_NO_HEARTBEAT_SEC {
				LOG_WS_VERBOSE("[HEARTBEAT] Reconnecting to websocket...")

				websocket.Reconnect()
				return
			}

			// Check if the last heartbeat is older than the heartbeat check interval
			if elapsed >= Constants.Websocket.HEARTBEAT_CHECK_INTERVAL_SEC {
				err := websocket.Conn.WriteMessage(ws.PingMessage, nil)
				if err != nil {
					LOG_WS_VERBOSE("[HEARTBEAT] Error sending ping:", err)

				} else {
					LOG_WS_VERBOSE_FULL("[HEARTBEAT] Ping sent.")
				}
			}
		}
	}()
}

func (websocket *Websocket) SendRequest_sync(req map[string]interface{}, timeout_sec ...int) (data []byte, hasTimedOut bool, WS_send_err *Error) {
	respChan := make(chan []byte, 1)

	LOG_WS_VERBOSE("[VERBOSE] Sending request:", req)

	websocket.pendingRequests[req["id"].(string)] = respChan

	// Send the message including the ID
	err := websocket.Conn.WriteJSON(req)
	if err != nil {
		return nil, false, LocalError(WS_SEND_MESSAGE_ERR, err.Error())
	}

	timeout := 4
	if len(timeout_sec) > 0 {
		timeout = timeout_sec[0]
	}

	// Wait for response or timeout
	var timer <-chan time.Time
	if timeout > 0 {
		timer = time.After(time.Duration(timeout) * time.Second)
	}

	select {
	case resp := <-respChan:
		return resp, false, nil
	case <-timer:
		delete(websocket.pendingRequests, req["id"].(string))
		return nil, true, LocalError(REQUEST_TIMEOUT_ERR, fmt.Sprintf("The request has timed out after %d seconds...", timeout))
	}
}

func CreateQueryStringWS(streams []string, isCombined bool) string {
	streamsStr := ""
	if isCombined {
		streamsStr += "/stream?streams=" + strings.Join(streams, "/")
	} else {
		streamsStr += "/ws/" + streams[0]
	}

	return streamsStr
}

func (websocket *Websocket) Reconnect() {
	websocket.isReconnecting = true
	websocket.subsocket_id++

	fmt.Println("[*Websocket.RECONNECT()] Reconnecting socket, closed:", websocket.closed)
	if !websocket.closed {
		err := websocket.silentCloseSocket()
		if err != nil {
			fmt.Println("[*Websocket.RECONNECT()] [SilentClose] There was an error silent closing the socket.")
		}
	}

	if websocket.OnReconnecting != nil {
		websocket.OnReconnecting()
	}

	for {
		queryStr := CreateQueryStringWS(websocket.Streams, websocket.IsCombined)
		conn, _, err := ws.DefaultDialer.Dial(websocket.BaseURL+queryStr, nil)
		if err != nil {
			LOG_WS_VERBOSE("There was an error reconnecting socket:", err)

			time.Sleep(500 * time.Millisecond)
			continue // Retry until successful
		}

		// Assign the newly created socket and break the loop
		websocket.Conn = conn
		websocket.closed = false
		websocket.RecordLastHeartbeat()
		break
	}

	setUpSocket(websocket, websocket.Conn)
	websocket.isReconnecting = false

	if websocket.OnReconnect != nil {
		websocket.OnReconnect()
	}

	LOG_WS_VERBOSE_FULL("[*Websocket.RECONNECT()] Successfully reconnected the socket!!!")
}

// This terminates the socket indefinitely
func (websocket *Websocket) Close() error {
	websocket.reconnect = false

	LOG_WS_VERBOSE_FULL("[*Websocket.CLOSE()] Forcefully closing socket indefinitely, closed:", websocket.closed)

	if !websocket.closed {
		err := websocket.silentCloseSocket()
		if err != nil {
			return err
		}
		websocket.CloseHandler(-1, "")
	}

	return fmt.Errorf("[LIB] Socket was already closed before closing")
}

func (websocket *Websocket) silentCloseSocket() error {
	fmt.Println("[*Websocket.silentClose()] Forcefully closing socket, closed:", websocket.closed)
	err := websocket.Conn.Close()
	if err != nil {
		fmt.Println("[*Websocket.silentClose()] There was an error closing the socket:", err)
		return err
	}

	return nil
}

func (websocket *Websocket) CloseHandler(code int, text string) error {
	fmt.Println("[*Websocket.CloseHandler()] code", code, "text", text, "is closed", websocket.closed)
	websocket.closed = true

	if websocket.reconnect {
		if websocket.OnDisconnect != nil {
			websocket.OnDisconnect(code, text)
		}
		websocket.Reconnect()
	} else {
		if websocket.OnClose != nil {
			websocket.OnClose(code, text)
		}
	}

	return nil
}

func (websocket *Websocket) PingHandler(appData string) error {
	LOG_WS_VERBOSE_FULL("Received a ping:", appData)

	err := websocket.Conn.WriteMessage(ws.PongMessage, []byte(appData))
	if err != nil {
		LOG_WS_VERBOSE("Error sending Pong:", err)

		return err
	}
	if websocket.OnPing != nil {
		websocket.OnPing(appData)
	}

	websocket.RecordLastHeartbeat()

	return nil
}

func (websocket *Websocket) PongHandler(appData string) error {

	LOG_WS_VERBOSE_FULL("Received a pong:", appData)

	if websocket.OnPong != nil {
		websocket.OnPong(appData)
	}

	websocket.RecordLastHeartbeat()

	return nil
}

func (websocket *Websocket) RecordLastHeartbeat() {
	websocket.Last_Heartbeat_Timestamp = time.Now().Unix()
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
