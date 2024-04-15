package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/OktopUSP/oktopus/backend/services/mtp/ws-adapter/internal/bridge"
	"github.com/OktopUSP/oktopus/backend/services/mtp/ws-adapter/internal/config"
	"github.com/OktopUSP/oktopus/backend/services/mtp/ws-adapter/internal/nats"
)

func main() {

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	c := config.NewConfig()

	kv, publisher, subscriber := nats.StartNatsClient(c.Nats)

	bridge := bridge.NewBridge(publisher, subscriber, c.Ws.Ctx, c.Ws, kv)
	bridge.StartBridge()

	<-done

	log.Println("websockets adapter is shutting down...")

}
