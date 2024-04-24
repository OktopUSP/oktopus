package main

import (
	"oktopUSP/backend/services/acs/internal/bridge"
	"oktopUSP/backend/services/acs/internal/config"
	"oktopUSP/backend/services/acs/internal/nats"
	"oktopUSP/backend/services/acs/internal/server"
	"oktopUSP/backend/services/acs/internal/server/handler"
)

func main() {

	c := config.NewConfig()

	natsActions := nats.StartNatsClient(c.Nats)

	h := handler.NewHandler(natsActions.Publish, natsActions.Subscribe, c.Acs)

	b := bridge.NewBridge(
		natsActions.Publish,
		natsActions.Subscribe,
		h,
	)
	b.StartBridge()

	server.Run(c.Acs, natsActions, h)
}
