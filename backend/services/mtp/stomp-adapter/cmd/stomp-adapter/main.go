package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/oktopUSP/oktopus/backend/services/mtp/stomp-adapter/internal/bridge"
	"github.com/oktopUSP/oktopus/backend/services/mtp/stomp-adapter/internal/config"
	"github.com/oktopUSP/oktopus/backend/services/mtp/stomp-adapter/internal/nats"
)

func main() {

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	c := config.NewConfig()

	_, publisher, subscriber := nats.StartNatsClient(c.Nats)

	bridge := bridge.NewBridge(publisher, subscriber, c.Nats.Ctx, c.Stomp)
	bridge.StartBridge()

	<-done

	log.Println("stomp adapter is shutting down...")
}
