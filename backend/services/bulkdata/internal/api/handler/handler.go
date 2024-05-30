package handler

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Handler struct {
	js  jetstream.JetStream
	nc  *nats.Conn
	kv  jetstream.KeyValue
	ctx context.Context
}

func NewHandler(js jetstream.JetStream, nc *nats.Conn, kv jetstream.KeyValue) Handler {
	return Handler{
		js: js,
		nc: nc,
		kv: kv,
	}
}
