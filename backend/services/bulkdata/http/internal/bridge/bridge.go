package bridge

import (
	"github.com/nats-io/nats.go"
)

type (
	Publisher  func(string, []byte) error
	Subscriber func(string, func(*nats.Msg)) error
)

type Bridge struct {
	pub Publisher
	sub Subscriber
}

const BULK_DATA_SUBJECT = "bulk"

func NewBridge(p Publisher, s Subscriber) *Bridge {
	return &Bridge{
		pub: p,
		sub: s,
	}
}

func (b *Bridge) SendDeviceData(deviceId string, payload []byte) error {
	return b.pub("oi", payload)
}
