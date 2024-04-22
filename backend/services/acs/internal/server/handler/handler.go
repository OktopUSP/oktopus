package handler

import (
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/oleiade/lane"
)

const Version = "1.0.0"

type Request struct {
	Id       string
	User     string
	Password string
	CwmpMsg  []byte
	Callback chan []byte
}

type CPE struct {
	SerialNumber         string
	Manufacturer         string
	OUI                  string
	ConnectionRequestURL string
	SoftwareVersion      string
	ExternalIPAddress    string
	Queue                *lane.Queue
	Waiting              *Request
	HardwareVersion      string
	LastConnection       time.Time
	DataModel            string
	Username             string
	Password             string
}

type Message struct {
	SerialNumber string
	Message      string
}

type WsMessage struct {
	Cmd string
}

type NatsSendMessage struct {
	MsgType string
	Data    json.RawMessage
}

type MsgCPEs struct {
	CPES map[string]CPE
}

type Handler struct {
	pub  func(string, []byte) error
	sub  func(string, func(*nats.Msg)) error
	Cpes map[string]CPE
}

const (
	NATS_CWMP_SUBJECT_PREFIX         = "cwmp.v1."
	NATS_CWMP_ADAPTER_SUBJECT_PREFIX = "cwmp-adapter.v1."
	NATS_ADAPTER_SUBJECT_PREFIX      = "adapter.v1."
)

func NewHandler(pub func(string, []byte) error, sub func(string, func(*nats.Msg)) error) *Handler {
	return &Handler{
		pub:  pub,
		sub:  sub,
		Cpes: make(map[string]CPE),
	}
}
