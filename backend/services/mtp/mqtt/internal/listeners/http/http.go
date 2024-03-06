package http

import (
	"log"

	"github.com/google/uuid"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/listeners"
)

type Http struct {
	HttpPort string
}

func (h *Http) Start(server *mqtt.Server) {
	stats := listeners.NewHTTPStats(uuid.NewString(), h.HttpPort, nil, server.Info)
	err := server.AddListener(stats)
	if err != nil {
		log.Fatal(err)
	}
}
