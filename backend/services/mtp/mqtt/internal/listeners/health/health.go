package health

import (
	"log"

	"github.com/google/uuid"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/listeners"
)

type HttpHealth struct {
	HttpPort string
}

func (h *HttpHealth) Start(server *mqtt.Server) {
	healthCheckEndpoint := listeners.NewHTTPHealthCheck(uuid.NewString(), h.HttpPort, nil)
	err := server.AddListener(healthCheckEndpoint)
	if err != nil {
		log.Fatal(err)
	}
}
