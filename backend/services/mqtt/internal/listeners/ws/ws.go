package ws

import (
	"log"

	"github.com/google/uuid"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/listeners"
)

type Ws struct {
	WsPort string
}

func (w *Ws) Start(server *mqtt.Server) {
	ws := listeners.NewWebsocket(uuid.NewString(), w.WsPort, nil)
	err := server.AddListener(ws)
	if err != nil {
		log.Fatal(err)
	}
}
