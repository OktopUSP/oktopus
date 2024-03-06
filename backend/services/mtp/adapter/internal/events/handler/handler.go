package handler

import (
	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/db"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	ONLINE = iota
	OFFLINE
)

const NATS_SUBJ_PREFIX = "mqtt-adapter.usp.v1."

type Handler struct {
	nc  *nats.Conn
	js  jetstream.JetStream
	db  db.Database
	cid string
}

func NewHandler(nc *nats.Conn, js jetstream.JetStream, d db.Database, cid string) Handler {
	return Handler{
		nc:  nc,
		js:  js,
		db:  d,
		cid: cid,
	}
}
