package ws

import (
	"context"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Ws struct {
	Addr  string
	Port  string
	Token string
	Route string
	Auth  bool
	TLS   bool
	Ctx   context.Context
}

const (
	WS_CONNECTION_RETRY = 10 * time.Second
)

// Global Websocket connection used in this package
var wsConn *websocket.Conn

func (w *Ws) Connect() {

	var wsUrl string
	if w.Auth {
		log.Println("WS token: ", w.Token)
		// e.g. ws://localhost:8080/ws/controller?token=123456
		wsUrl = "ws://" + w.Addr + ":" + w.Port + w.Route + "?token=" + w.Token
	} else {
		// e.g. ws://localhost:8080/ws/controller
		wsUrl = "ws://" + w.Addr + ":" + w.Port + w.Route
	}

	// Keeps trying to connect to the WS endpoint until it succeeds or receives a stop signal
	go func() {
		for {
			c, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
			if err != nil {
				log.Printf("Error to connect to %s, err: %s", wsUrl, err)
				time.Sleep(WS_CONNECTION_RETRY)
				continue
			}
			// instantiate global ws connection
			wsConn = c
			log.Println("Connected to WS endpoint--> ", wsUrl)
			go w.Subscribe()
			break
		}
	}()
}

func (w *Ws) Disconnect() {
	log.Println("Disconnecting from WS endpoint...")

	if wsConn != nil {
		err := wsConn.Close()
		if err != nil {
			log.Println("Error while disconnecting from WS endpoint:", err.Error())
		}
		log.Println("Succesfully disconnected from WS endpoint")
	} else {
		log.Println("No WS connection to close")
	}
}

// Websockets doesn't follow pub/sub architecture, but we use these
// naming here to implement the Broker interface and abstract the MTP layer.
/* -------------------------------------------------------------------------- */

func (w *Ws) Subscribe() {
	for {
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
	}
}

func (w *Ws) Publish(msg []byte, topic, respTopic string, retain bool) {
	err := wsConn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Println("write:", err)
		return
	}
}

/* -------------------------------------------------------------------------- */
