package server

import (
	"log"
	"net/http"
	"oktopUSP/backend/services/acs/internal/config"
	"oktopUSP/backend/services/acs/internal/nats"
	"oktopUSP/backend/services/acs/internal/server/handler"
	"os"
)

func Run(c config.Acs, natsActions nats.NatsActions, h *handler.Handler) {

	http.HandleFunc(c.Route, h.CwmpHandler)

	log.Printf("ACS running at %s%s", c.Port, c.Route)

	err := http.ListenAndServe(c.Port, nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
