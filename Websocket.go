package Binance

import (
	"fmt"
	"strings"
	"time"

	ws "github.com/gorilla/websocket"
)

type Websocket struct {
	Conn *ws.Conn
	// Host server's URL
	BaseURL    string
	StreamName string
	Streams    []string

	// This is to show the current state of the stream
	// false -> it's a raw stream
	// true -> it's a combined stream
	IsCombined bool

	OnMessage    func(messageType int, msg []byte, err error)
	OnReconnect  func()
	OnPing       func(appData string)
	OnPong       func(appData string)
	OnDisconnect func(code int, text string)
	OnClose      func(code int, text string)

	Creation_Timestamp       int64
	Last_Heartbeat_Timestamp int64

	reconnect bool
	closed    bool
}

func CreateSocket(baseURL string, streams []string, isCombined bool) (*Websocket, error) {
	queryStr := CreateQueryStringWS(streams, isCombined)

	fullStreamStr := baseURL + queryStr

	if SHOWQUERIES {
		fmt.Println("Websocket:", fullStreamStr)
	}

	conn, _, err := ws.DefaultDialer.Dial(fullStreamStr, nil)
	if err != nil {
		if VERBOSE {
			fmt.Println("There was an error creating websocket:", err)
		}
		return nil, err
	}

	currentTime := time.Now().Unix()

	websocket := &Websocket{
		Conn:                     conn,
		BaseURL:                  baseURL,
		Streams:                  streams,
		StreamName:               queryStr,
		IsCombined:               isCombined,
		Creation_Timestamp:       currentTime,
		Last_Heartbeat_Timestamp: currentTime,
		reconnect:                true,
		closed:                   false,
	}

	conn.SetCloseHandler(websocket.CloseHandler)
	conn.SetPingHandler(websocket.PingHandler)
	conn.SetPongHandler(websocket.PongHandler)

	// Handle system interrupts to close the connection gracefully
	// interrupt := make(chan os.Signal, 1)
	// signal.Notify(interrupt, os.Interrupt)
	// Not sure how to best do this, so will leave it empty for now
	// TODO

	// Goroutine to read messages
	go func() {
		for {
			if websocket.closed {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				if VERBOSE {
					fmt.Println("Error reading message:", err)
				}
				return
			}
			if PRINT_WS_MESSAGES {
				fmt.Printf("Type: %d, message: %s\n", msgType, string(msg))
			}

			websocket.RecordLastHeartbeat()

			if websocket.OnMessage != nil {
				websocket.OnMessage(msgType, msg, err)
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Duration(Constants.Websocket.HEARTBEAT_CHECK_INTERVAL_SEC) * time.Second)
		defer ticker.Stop()

		for {

			<-ticker.C // Wait the appropriate amount of time
			if websocket.closed {
				if VERBOSE {
					fmt.Println("[HEARTBEAT] Websocket is closed, skipping check.")
				}
				continue
			}

			currentTime := time.Now().Unix()

			elapsed := currentTime - websocket.Last_Heartbeat_Timestamp

			// Check if the last heartbeat is older than the close interval
			if elapsed >= Constants.Websocket.HEARTBEAT_CLOSE_ON_NO_HEARTBEAT_SEC {
				if VERBOSE {
					fmt.Println("[HEARTBEAT] Reconnecting to websocket...")
				}
				websocket.Reconnect()
				continue
			}

			// Check if the last heartbeat is older than the heartbeat check interval
			if elapsed >= Constants.Websocket.HEARTBEAT_CHECK_INTERVAL_SEC {
				err := websocket.Conn.WriteMessage(ws.PingMessage, nil)
				if err != nil {
					if VERBOSE {
						fmt.Println("[HEARTBEAT] Error sending ping:", err)
					}
				} else if VERBOSE {
					fmt.Println("[HEARTBEAT] Ping sent.")
				}
			}
		}
	}()

	return websocket, nil
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
	if !websocket.closed {
		fmt.Println("[TO REMOVE] [0] Forcefully disconnecting the socket before reconnecting...")
		err := websocket.Disconnect()
		if err != nil {
			if VERBOSE {
				fmt.Println("There was an error disconnecting the socket...", err)
			}
		}
	}

	for {
		queryStr := CreateQueryStringWS(websocket.Streams, websocket.IsCombined)
		conn, _, err := ws.DefaultDialer.Dial(websocket.BaseURL+queryStr, nil)
		if err != nil {
			if VERBOSE {
				fmt.Println("There was an error reconnecting socket:", err)
			}
			continue // Retry until successful
		}

		// Assign the newly created socket and break the loop
		websocket.Conn = conn
		websocket.closed = false
		websocket.RecordLastHeartbeat()
		break
	}

	fmt.Println("Successfully reconnected the socket!!!")
}

func (websocket *Websocket) Disconnect() error {
	if !websocket.closed {
		fmt.Println("[TO REMOVE] [1] Forecfully closing the connection before restarting...")
		err := websocket.Conn.Close()
		if err != nil {
			fmt.Println("[DISCONNECT] There was an error disconnecting the socket:", err)
			return err
		}
		websocket.CloseHandler(-1, "")
	}

	return fmt.Errorf("socket was already closed before disconnecting")
}

func (websocket *Websocket) CloseHandler(code int, text string) error {
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

// func (websocket *Websocket) CloseHandler(code int, text string) error {
// 	if VERBOSE {
// 		fmt.Println("Received a close message, code:", code, "text:", text)
// 	}

// 	err := websocket.Conn.Close()
// 	if err != nil {
// 		if VERBOSE {
// 			fmt.Println("There was an error when handling Websocket Close:", err, "Is closed already?", websocket.closed)
// 		}
// 		return err
// 	}

// 	websocket.closed = true

// 	if websocket.OnClose != nil {
// 		websocket.OnClose(code, text)
// 	}
// 	return nil
// }

func (websocket *Websocket) PingHandler(appData string) error {
	err := websocket.Conn.WriteMessage(ws.PongMessage, []byte(appData))
	if err != nil {
		if VERBOSE {
			fmt.Println("Error sending Pong:", err)
		}
		return err
	}
	if websocket.OnPing != nil {
		websocket.OnPing(appData)
	}

	if VERBOSE {
		fmt.Println("Received a ping:", appData)
	}

	websocket.RecordLastHeartbeat()

	return nil
}

func (websocket *Websocket) PongHandler(appData string) error {
	if VERBOSE {
		fmt.Println("Received a pong:", appData)
	}
	if websocket.OnPong != nil {
		websocket.OnPong(appData)
	}

	websocket.RecordLastHeartbeat()

	return nil
}

func (websocket *Websocket) RecordLastHeartbeat() {
	websocket.Last_Heartbeat_Timestamp = time.Now().Unix()
}
