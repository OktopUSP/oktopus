package ws

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/leandrofars/oktopus/internal/db"
	"github.com/leandrofars/oktopus/internal/mtp/handler"
	"github.com/leandrofars/oktopus/internal/usp_record"

	"google.golang.org/protobuf/proto"
)

type Ws struct {
	Addr               string
	Port               string
	Token              string
	Route              string
	Auth               bool
	TLS                bool
	InsecureSkipVerify bool
	Ctx                context.Context
	NewDeviceQueue     map[string]string
	NewDevQMutex       *sync.Mutex
	DB                 db.Database
}

const (
	WS_CONNECTION_RETRY = 10 * time.Second
)

const (
	OFFLINE = "0"
	ONLINE  = "1"
)

type deviceStatus struct {
	Eid    string
	Status string
}

// Global Websocket connection used in this package
var wsConn *websocket.Conn

func (w *Ws) Connect() {
	log.Println("Connecting to WS endpoint...")

	prefix := "ws://"
	if w.TLS {
		prefix = "wss://"
	}

	wsUrl := prefix + w.Addr + ":" + w.Port + w.Route

	if w.Auth {
		log.Println("WS token:", w.Token)
		// e.g. ws://localhost:8080/ws/controller?token=123456
		wsUrl = wsUrl + "?token=" + w.Token
	}

	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: w.InsecureSkipVerify,
		},
	}

	// Keeps trying to connect to the WS endpoint until it succeeds or receives a stop signal
	go func(dialer websocket.Dialer) {
		for {
			c, _, err := dialer.Dial(wsUrl, nil)
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
	}(dialer)
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

	var m sync.Mutex
	w.NewDevQMutex = &m
	w.NewDeviceQueue = make(map[string]string)

	for {
		msgType, wsMsg, err := wsConn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket error: %v", err)
				w.Connect()
				return
			}
			log.Println("websocket unexpected error:", err)
			return
		}

		//TODO: deal with message in new go routine
		if msgType == websocket.TextMessage {
			var deviceStatus deviceStatus
			err = json.Unmarshal(wsMsg, &deviceStatus)
			if err != nil {
				log.Println("Websockets Text Message is not about devices status")
			}

			log.Println("Received device status message")
			var status db.Status
			switch deviceStatus.Status {
			case ONLINE:
				status = db.Online
			case OFFLINE:
				status = db.Offline
			default:
				log.Println("Invalid device status")
				return
			}

			w.DB.UpdateStatus(deviceStatus.Eid, status, db.WEBSOCKETS)

			//TODO: return error 1003 to device
			//TODO: get status messages
			continue
		}

		//log.Printf("binary data: %s", string(wsMsg))

		//TODO: if error at processing message return error 1003 to devicec
		//TODO: deal with received messages in parallel

		var record usp_record.Record
		//var message usp_msg.Msg

		err = proto.Unmarshal(wsMsg, &record)
		if err != nil {
			log.Println(err)
		}

		connRecord := &usp_record.Record_WebsocketConnect{
			WebsocketConnect: &usp_record.WebSocketConnectRecord{},
		}

		noSessionRecord := &usp_record.Record_NoSessionContext{
			NoSessionContext: &usp_record.NoSessionContextRecord{},
		}

		//log.Printf("Record Type: %++v", record.RecordType)
		deviceId := record.FromId

		// New Device Handler
		if reflect.TypeOf(record.RecordType) == reflect.TypeOf(connRecord) {
			log.Println("Websocket new device:", deviceId)
			tr369Message := handler.HandleNewDevice(deviceId)
			w.NewDevQMutex.Lock()
			w.NewDeviceQueue[deviceId] = ""
			w.NewDevQMutex.Unlock()
			w.Publish(tr369Message, "", "", false)
			continue
		}

		//TODO: see what type of message was received
		if reflect.TypeOf(record.RecordType) == reflect.TypeOf(noSessionRecord) {

			//log.Printf("Websocket device %s message", record.FromId)
			// New device answer
			if _, ok := w.NewDeviceQueue[deviceId]; ok {
				log.Printf("New device %s response", deviceId)
				device := handler.HandleNewDevicesResponse(wsMsg, deviceId, db.WEBSOCKETS)
				w.NewDevQMutex.Lock()
				delete(w.NewDeviceQueue, deviceId)
				w.NewDevQMutex.Unlock()
				w.DB.CreateDevice(device)
				if err != nil {
					log.Fatal(err)
				}
				continue
			}

			//TODO: send message to Api Msg Queue

		}

		//log.Printf("recv: %++v", record)
	}
}

func (w *Ws) Publish(msg []byte, topic, respTopic string, retain bool) {
	err := wsConn.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		log.Println("write:", err)
		return
	}
}

/* -------------------------------------------------------------------------- */
