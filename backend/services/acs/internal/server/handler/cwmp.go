package handler

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"oktopUSP/backend/services/acs/internal/auth"
	"oktopUSP/backend/services/acs/internal/cwmp"
	"time"

	"github.com/oleiade/lane"
)

func (h *Handler) CwmpHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("--> Connection from %s", r.RemoteAddr)

	defer r.Body.Close()
	defer log.Printf("<-- Connection from %s closed", r.RemoteAddr)

	tmp, _ := ioutil.ReadAll(r.Body)
	body := string(tmp)

	if h.acsConfig.DebugMode {
		log.Println("Received message: ", body)
	}

	var envelope cwmp.SoapEnvelope
	xml.Unmarshal(tmp, &envelope)

	messageType := envelope.Body.CWMPMessage.XMLName.Local
	log.Println("messageType: ", messageType)

	var cpe CPE
	var exists bool

	w.Header().Set("Server", "Oktopus "+Version)

	if messageType != "Inform" {
		if cookie, err := r.Cookie("oktopus"); err == nil {
			if cpe, exists = h.Cpes[cookie.Value]; !exists {
				log.Printf("CPE with serial number %s not found", cookie.Value)
			}
			log.Printf("CPE with serial number %s found", cookie.Value)
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

		if _, exists := h.Cpes[sn]; !exists {
			log.Println("New device: " + sn)
			h.Cpes[sn] = CPE{
				SerialNumber:         sn,
				LastConnection:       time.Now().UTC(),
				SoftwareVersion:      Inform.GetSoftwareVersion(),
				HardwareVersion:      Inform.GetHardwareVersion(),
				ExternalIPAddress:    addr,
				ConnectionRequestURL: Inform.GetConnectionRequest(),
				OUI:                  Inform.DeviceId.OUI,
				Queue:                lane.NewQueue(),
				DataModel:            Inform.GetDataModelType(),
			}
			go h.handleCpeStatus(sn)
			h.pub(NATS_CWMP_SUBJECT_PREFIX+sn+".info", tmp)
		}
		obj := h.Cpes[sn]
		cpe := &obj
		cpe.LastConnection = time.Now().UTC()

		log.Printf("Received an Inform from  device %s withEventCodes %s", addr, Inform.GetEvents())

		expiration := time.Now().AddDate(0, 0, 1)

		cookie := http.Cookie{Name: "oktopus", Value: sn, Expires: expiration}
		http.SetCookie(w, &cookie)
		data, _ := xml.Marshal(cwmp.InformResponse(envelope.Header.Id))
		w.Write(data)
	} else if messageType == "TransferComplete" {

	} else if messageType == "GetRPC" {

	} else {

		if len(body) == 0 {
			log.Println("Got Empty Post")
		}

		if cpe.Waiting != nil {
			log.Println("CPE was waiting for a response, now received something")
			var e cwmp.SoapEnvelope
			xml.Unmarshal([]byte(body), &e)
			log.Println("Kind of envelope: ", e.KindOf())

			if e.KindOf() == "GetParameterNamesResponse" {
				// var envelope cwmp.GetParameterNamesResponse
				// xml.Unmarshal([]byte(body), &envelope)

				// msg := new(NatsSendMessage)
				// msg.MsgType = "GetParameterNamesResponse"
				// msg.Data, _ = json.Marshal(envelope)
				log.Println("Receive GetParameterNamesResponse from CPE:", cpe.SerialNumber)
				cpe.Waiting.Callback <- tmp

			} else if e.KindOf() == "GetParameterValuesResponse" {
				var envelope cwmp.GetParameterValuesResponse
				xml.Unmarshal([]byte(body), &envelope)

				msg := new(NatsSendMessage)
				msg.MsgType = "GetParameterValuesResponse"
				msg.Data, _ = json.Marshal(envelope)

				cpe.Waiting.Callback <- tmp

			} else {
				log.Println("Unknown message type")
				cpe.Waiting.Callback <- tmp
			}
			cpe.Waiting = nil
		} else {
			log.Println("CPE was not waiting for a response")
		}

		log.Printf("CPE %s Queue size: %d", cpe.SerialNumber, cpe.Queue.Size())

		if cpe.Queue.Size() > 0 {
			req := cpe.Queue.Dequeue().(Request)
			cpe.Waiting = &req
			log.Println("Sending request to CPE:", req.Id)
			w.Header().Set("Connection", "keep-alive")
			w.Write(req.CwmpMsg)
		} else {
			w.Header().Set("Connection", "close")
			w.WriteHeader(204)
		}
	}
	h.Cpes[cpe.SerialNumber] = cpe
}

func (h *Handler) ConnectionRequest(cpe CPE) error {
	log.Println("--> ConnectionRequest, CPE: ", cpe.SerialNumber)
	// log.Println("ConnectionRequestURL: ", cpe.ConnectionRequestURL)
	// log.Println("ConnectionRequestUsername: ", cpe.Username)
	// log.Println("ConnectionRequestPassword: ", cpe.Password)

	ok, err := auth.Auth("", "", cpe.ConnectionRequestURL)
	if !ok {
		log.Println("Error while authenticating to CPE: ", err)
	}
	return err
}
