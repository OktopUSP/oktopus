package bridge

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"reflect"
	"strings"
	"sync"

	// "reflect"
	"time"

	"github.com/OktopUSP/oktopus/backend/services/mtp/ws-adapter/internal/config"
	"github.com/OktopUSP/oktopus/backend/services/mtp/ws-adapter/internal/usp/usp_record"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

const (
	NATS_WS_SUBJECT_PREFIX         = "ws.usp.v1."
	NATS_WS_ADAPTER_SUBJECT_PREFIX = "ws-adapter.usp.v1.*."
	WS_TOPIC_PREFIX                = "oktopus/usp/"
	WS_CONNECTION_RETRY            = 10 * time.Second
)

const (
	OFFLINE = iota
	ONLINE
)

type deviceStatus struct {
	Eid    string
	Status string
}

type (
	Publisher  func(string, []byte) error
	Subscriber func(string, func(*nats.Msg)) error
)

type Bridge struct {
	Pub            Publisher
	Sub            Subscriber
	Ws             config.Ws
	NewDeviceQueue map[string]string
	NewDevQMutex   *sync.Mutex
	Ctx            context.Context
}

func NewBridge(p Publisher, s Subscriber, ctx context.Context, w config.Ws) *Bridge {
	return &Bridge{
		Pub: p,
		Sub: s,
		Ws:  w,
		Ctx: ctx,
	}
}

func (b *Bridge) StartBridge() {

	url := b.urlBuild()
	dialer := b.newDialer()
	go func(dialer websocket.Dialer) {
		for {
			wc, _, err := dialer.Dial(url, nil)
			if err != nil {
				log.Printf("Error to connect to %s, err: %s", url, err)
				time.Sleep(WS_CONNECTION_RETRY)
				continue
			}
			log.Println("Connected to WS endpoint--> ", url)
			go b.subscribe(wc)
			go func(wc *websocket.Conn) {
				for {
					msgType, wsMsg, err := wc.ReadMessage()
					if err != nil {
						if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
							log.Printf("websocket error: %v", err)
							b.StartBridge()
							return
						}
						log.Println("websocket unexpected error:", err)
						return
					}
					if msgType == websocket.TextMessage {
						b.statusMsgHandler(wsMsg)
						continue
					}

					var record usp_record.Record
					err = proto.Unmarshal(wsMsg, &record)
					if err != nil {
						log.Println(err)
					}
					device := record.FromId

					noSessionRecord := &usp_record.Record_NoSessionContext{
						NoSessionContext: &usp_record.NoSessionContextRecord{},
					}
					if reflect.TypeOf(record.RecordType) == reflect.TypeOf(noSessionRecord) {
						if _, ok := b.NewDeviceQueue[device]; ok {
							b.newDeviceMsgHandler(wc, device, wsMsg)
							continue
						}
					}

					// log.Println("Handle api request")
					// var msg usp_msg.Msg
					// err = proto.Unmarshal(record.GetNoSessionContext().Payload, &msg)
					// if err != nil {
					// 	log.Println(err)
					// 	continue
					// }
					// if _, ok := w.MsgQueue[msg.Header.MsgId]; ok {
					// 	//m.QMutex.Lock()
					// 	w.MsgQueue[msg.Header.MsgId] <- msg
					// 	//m.QMutex.Unlock()
					// } else {
					// 	log.Printf("Message answer to request %s arrived too late", msg.Header.MsgId)
					// }

				}
			}(wc)
			break
		}
	}(dialer)
}

func (b *Bridge) subscribe(wc *websocket.Conn) {

	b.NewDeviceQueue = make(map[string]string)
	b.NewDevQMutex = &sync.Mutex{}

	b.Sub(NATS_WS_ADAPTER_SUBJECT_PREFIX+"info", func(msg *nats.Msg) {

		log.Printf("Received message on info subject")

		subj := strings.Split(msg.Subject, ".")
		device := subj[len(subj)-2]

		b.NewDevQMutex.Lock()
		b.NewDeviceQueue[device] = ""
		b.NewDevQMutex.Unlock()

		err := wc.WriteMessage(websocket.BinaryMessage, msg.Data)
		if err != nil {
			log.Printf("send websocket msg error: %q", err)
			return
		}
	})
}

func (b *Bridge) newDeviceMsgHandler(wc *websocket.Conn, device string, msg []byte) {
	log.Printf("New device %s response", device)
	b.Pub(NATS_WS_SUBJECT_PREFIX+device+".info", msg)

	b.NewDevQMutex.Lock()
	delete(b.NewDeviceQueue, device)
	b.NewDevQMutex.Unlock()
}

func (b *Bridge) statusMsgHandler(wsMsg []byte) {
	var deviceStatus deviceStatus
	err := json.Unmarshal(wsMsg, &deviceStatus)
	if err != nil {
		log.Println("Websockets Text Message is not about devices status")
		return
	}
	b.Pub(NATS_WS_SUBJECT_PREFIX+deviceStatus.Eid+".status", []byte(deviceStatus.Status))
}

func (b *Bridge) urlBuild() string {
	prefix := "ws://"
	if b.Ws.TlsEnable {
		prefix = "wss://"
	}

	wsUrl := prefix + b.Ws.Addr + b.Ws.Port + b.Ws.Route

	if b.Ws.AuthEnable {
		wsUrl = wsUrl + "?token=" + b.Ws.Token
	}

	return wsUrl
}

func (b *Bridge) newDialer() websocket.Dialer {
	return websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: b.Ws.SkipTlsVerify,
		},
	}
}
