package bridge

import (
	"encoding/json"
	"log"
	"net/http"
	"oktopUSP/backend/services/acs/internal/server/handler"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type Bridge struct {
	pub  func(string, []byte) error
	sub  func(string, func(*nats.Msg)) error
	cpes map[string]handler.CPE
	h    *handler.Handler
}

type msgAnswer struct {
	Code int
	Msg  any
}

const DEVICE_ANSWER_TIMEOUT = 5 * time.Second

func NewBridge(
	pub func(string, []byte) error,
	sub func(string, func(*nats.Msg)) error,
	h *handler.Handler,
) *Bridge {
	return &Bridge{
		pub:  pub,
		sub:  sub,
		cpes: h.Cpes,
		h:    h,
	}
}

func (b *Bridge) StartBridge() {
	b.sub(handler.NATS_CWMP_ADAPTER_SUBJECT_PREFIX+"*.api", func(msg *nats.Msg) {
		log.Printf("Received message: %s", string(msg.Data))
		log.Printf("Subject: %s", msg.Subject)
		log.Printf("Reply: %s", msg.Reply)

		device := getDeviceFromSubject(msg.Subject)
		cpe, ok := b.cpes[device]
		if !ok {
			log.Printf("Device %s not found", device)
			respondMsg(msg.Respond, http.StatusNotFound, "Device not found")
			return
		}
		if cpe.Queue.Size() > 0 {
			log.Printf("Device %s is busy", device)
			respondMsg(msg.Respond, http.StatusConflict, "Device is busy")
			return
		}

		deviceAnswer := make(chan []byte)
		defer close(deviceAnswer)

		cpe.Queue.Enqueue(handler.Request{ //TODO: pass user and password too
			Id:       uuid.NewString(),
			CwmpMsg:  msg.Data,
			Callback: deviceAnswer,
		})

		err := b.h.ConnectionRequest(cpe)
		if err != nil {
			log.Println("Failed to do connection request", err)
			cpe.Queue.Dequeue()
			respondMsg(msg.Respond, http.StatusBadRequest, err.Error())
			return
		}

		//req := cpe.Queue.Dequeue().(handler.Request)
		//cpe.Waiting = &req

		select {
		case response := <-deviceAnswer:
			log.Println("Received response from device: ", string(response))
			respondMsg(msg.Respond, http.StatusOK, response)
		case <-time.After(DEVICE_ANSWER_TIMEOUT):
			log.Println("Device response timed out")
			respondMsg(msg.Respond, http.StatusRequestTimeout, "Request timeout")
		}
	})
}

func respondMsg(respond func(data []byte) error, code int, msgData any) {

	msg, err := json.Marshal(msgAnswer{
		Code: code,
		Msg:  msgData,
	})
	if err != nil {
		log.Printf("Failed to marshal message: %q", err)
		respond([]byte(err.Error()))
		return
	}

	respond(msg)
	//log.Println("Responded with message: ", string(msg))
}

func getDeviceFromSubject(subject string) string {
	paths := strings.Split(subject, ".")
	device := paths[len(paths)-2]
	return device
}
