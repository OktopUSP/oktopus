package bridge

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"oktopUSP/backend/services/acs/internal/config"
	"oktopUSP/backend/services/acs/internal/server/handler"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"golang.org/x/sys/unix"
)

type Bridge struct {
	pub  func(string, []byte) error
	sub  func(string, func(*nats.Msg)) error
	cpes map[string]handler.CPE
	h    *handler.Handler
	conf *config.Acs
}

type msgAnswer struct {
	Code int
	Msg  any
}

func NewBridge(
	pub func(string, []byte) error,
	sub func(string, func(*nats.Msg)) error,
	h *handler.Handler,
	c *config.Acs,
) *Bridge {
	return &Bridge{
		pub:  pub,
		sub:  sub,
		cpes: h.Cpes,
		h:    h,
		conf: c,
	}
}

func (b *Bridge) StartBridge() {

	b.sub(handler.NATS_CWMP_ADAPTER_SUBJECT_PREFIX+"*.api", func(msg *nats.Msg) {
		if b.conf.DebugMode {
			log.Printf("Received message: %s", string(msg.Data))
			log.Printf("Subject: %s", msg.Subject)
			log.Printf("Reply: %s", msg.Reply)
		}

		device := getDeviceFromSubject(msg.Subject)
		cpe, ok := b.cpes[device]
		if !ok {
			log.Printf("Device %s not found", device)
			respondMsg(msg.Respond, http.StatusNotFound, "Device not found")
			return
		}
		if cpe.Queue.Size() > 0 {
			log.Println("Queue size: ", cpe.Queue.Size())
			log.Println("Queue data: ", cpe.Queue)
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
			Time:     time.Now(),
		})

		err := b.h.ConnectionRequest(cpe)
		if err != nil {
			log.Println("Failed to do connection request", err)
			cpe.Queue.Dequeue()
			respondMsg(msg.Respond, http.StatusBadRequest, err.Error())
			return
		}

		defer cpe.Queue.Dequeue()

		select {
		case response := <-deviceAnswer:
			if b.conf.DebugMode {
				log.Printf("Received response from cpe: %s payload: %s ", cpe.SerialNumber, string(response))
			}
			respondMsg(msg.Respond, http.StatusOK, response)
		case <-time.After(b.conf.DeviceAnswerTimeout):
			log.Println("Device response timed out")
			respondMsg(msg.Respond, http.StatusRequestTimeout, "Request timeout")
		}

	})

	b.sub(handler.NATS_CWMP_ADAPTER_SUBJECT_PREFIX+"rtt", func(msg *nats.Msg) {
		log.Printf("Received message on rtt subject")
		url := "127.0.0.1" + b.conf.Port
		conn, err := net.Dial("tcp", url)
		if err != nil {
			respondMsg(msg.Respond, 500, err.Error())
			return
		}
		defer conn.Close()

		info, err := tcpInfo(conn.(*net.TCPConn))
		if err != nil {
			respondMsg(msg.Respond, 500, err.Error())
			return
		}
		rtt := time.Duration(info.Rtt) * time.Microsecond

		respondMsg(msg.Respond, 200, rtt/1000)
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
}

func getDeviceFromSubject(subject string) string {
	paths := strings.Split(subject, ".")
	device := paths[len(paths)-2]
	return device
}

func tcpInfo(conn *net.TCPConn) (*unix.TCPInfo, error) {
	raw, err := conn.SyscallConn()
	if err != nil {
		return nil, err
	}

	var info *unix.TCPInfo
	ctrlErr := raw.Control(func(fd uintptr) {
		info, err = unix.GetsockoptTCPInfo(int(fd), unix.IPPROTO_TCP, unix.TCP_INFO)
	})
	switch {
	case ctrlErr != nil:
		return nil, ctrlErr
	case err != nil:
		return nil, err
	}
	return info, nil
}
