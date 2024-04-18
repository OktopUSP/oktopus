package handler

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"oktopUSP/backend/services/acs/internal/cwmp"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/oleiade/lane"
	"golang.org/x/net/websocket"
)

const Version = "1.0.0"

type Request struct {
	Id          string
	Websocket   *websocket.Conn
	CwmpMessage string
	Callback    func(msg *WsSendMessage) error
}

type CPE struct {
	SerialNumber         string
	Manufacturer         string
	OUI                  string
	ConnectionRequestURL string
	XmppId               string
	XmppUsername         string
	XmppPassword         string
	SoftwareVersion      string
	ExternalIPAddress    string
	State                string
	Queue                *lane.Queue
	Waiting              *Request
	HardwareVersion      string
	LastConnection       time.Time
	DataModel            string
	KeepConnectionOpen   bool
}

type Message struct {
	SerialNumber string
	Message      string
}

type WsMessage struct {
	Cmd string
}

type WsSendMessage struct {
	MsgType string
	Data    json.RawMessage
}

type MsgCPEs struct {
	CPES map[string]CPE
}

type Handler struct {
	pub  func(string, []byte) error
	sub  func(string, func(*nats.Msg)) error
	cpes map[string]CPE
}

func NewHandler(pub func(string, []byte) error, sub func(string, func(*nats.Msg)) error) *Handler {
	return &Handler{
		pub:  pub,
		sub:  sub,
		cpes: make(map[string]CPE),
	}
}

func (h *Handler) CwmpHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("--> Connection from %s", r.RemoteAddr)

	defer r.Body.Close()
	defer log.Printf("<-- Connection from %s closed", r.RemoteAddr)

	tmp, _ := ioutil.ReadAll(r.Body)
	body := string(tmp)
	len := len(body)

	log.Printf("body:\n %v", body)

	var envelope cwmp.SoapEnvelope
	xml.Unmarshal(tmp, &envelope)

	messageType := envelope.Body.CWMPMessage.XMLName.Local

	var cpe *CPE

	w.Header().Set("Server", "Oktopus "+Version)

	if messageType != "Inform" {
		if cookie, err := r.Cookie("oktopus"); err == nil {
			if _, exists := h.cpes[cookie.Value]; !exists {
				log.Printf("CPE with serial number %s not found", cookie.Value)
			}
		} else {
			fmt.Println("cookie 'oktopus' missing")
			w.WriteHeader(401)
			return
		}
	}

	if messageType == "Inform" {
		var Inform cwmp.CWMPInform
		xml.Unmarshal(tmp, &Inform)

		var addr string
		if r.Header.Get("X-Real-Ip") != "" {
			addr = r.Header.Get("X-Real-Ip")
		} else {
			addr = r.RemoteAddr
		}

		sn := Inform.DeviceId.SerialNumber

		if _, exists := h.cpes[sn]; !exists {
			fmt.Println("New device: " + sn)
			h.cpes[sn] = CPE{
				SerialNumber:         sn,
				LastConnection:       time.Now().UTC(),
				SoftwareVersion:      Inform.GetSoftwareVersion(),
				HardwareVersion:      Inform.GetHardwareVersion(),
				ExternalIPAddress:    addr,
				ConnectionRequestURL: Inform.GetConnectionRequest(),
				OUI:                  Inform.DeviceId.OUI,
				Queue:                lane.NewQueue(),
				DataModel:            Inform.GetDataModelType(),
				KeepConnectionOpen:   false}
		}
		obj := h.cpes[sn]
		cpe := &obj
		cpe.LastConnection = time.Now().UTC()

		log.Printf("Received an Inform from %s (%d bytes) with SerialNumber %s and EventCodes %s", addr, len, sn, Inform.GetEvents())

		expiration := time.Now().AddDate(0, 0, 1)

		cookie := http.Cookie{Name: "oktopus", Value: sn, Expires: expiration}
		http.SetCookie(w, &cookie)
	} else if messageType == "TransferComplete" {

	} else if messageType == "GetRPC" {

	} else {
		if len == 0 {
			log.Printf("Got Empty Post")
		}

		if cpe.Waiting != nil {
			var e cwmp.SoapEnvelope
			xml.Unmarshal([]byte(body), &e)

			if e.KindOf() == "GetParameterNamesResponse" {
				var envelope cwmp.GetParameterNamesResponse
				xml.Unmarshal([]byte(body), &envelope)

				msg := new(WsSendMessage)
				msg.MsgType = "GetParameterNamesResponse"
				msg.Data, _ = json.Marshal(envelope)

				cpe.Waiting.Callback(msg)
				//				if err := websocket.JSON.Send(cpe.Waiting.Websocket, msg); err != nil {
				//					fmt.Println("error while sending back answer:", err)
				//				}

			} else if e.KindOf() == "GetParameterValuesResponse" {
				var envelope cwmp.GetParameterValuesResponse
				xml.Unmarshal([]byte(body), &envelope)

				msg := new(WsSendMessage)
				msg.MsgType = "GetParameterValuesResponse"
				msg.Data, _ = json.Marshal(envelope)

				cpe.Waiting.Callback(msg)
				//				if err := websocket.JSON.Send(cpe.Waiting.Websocket, msg); err != nil {
				//					fmt.Println("error while sending back answer:", err)
				//				}

			} else {
				msg := new(WsMessage)
				msg.Cmd = body

				if err := websocket.JSON.Send(cpe.Waiting.Websocket, msg); err != nil {
					fmt.Println("error while sending back answer:", err)
				}

			}

			cpe.Waiting = nil
		}

		// Got Empty Post or a Response. Now check for any event to send, otherwise 204
		if cpe.Queue.Size() > 0 {
			req := cpe.Queue.Dequeue().(Request)
			// fmt.Println("sending "+req.CwmpMessage)
			fmt.Fprintf(w, req.CwmpMessage)
			cpe.Waiting = &req
		} else {
			if cpe.KeepConnectionOpen {
				fmt.Println("I'm keeping connection open")
			} else {
				w.WriteHeader(204)
			}
		}
	}
}
