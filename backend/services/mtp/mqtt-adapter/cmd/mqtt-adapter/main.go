package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/OktopUSP/oktopus/backend/services/mqtt-adapter/internal/bridge"
	"github.com/OktopUSP/oktopus/backend/services/mqtt-adapter/internal/config"
	"github.com/OktopUSP/oktopus/backend/services/mqtt-adapter/internal/nats"
)

func main() {

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	c := config.NewConfig()

	_, publisher, subscriber := nats.StartNatsClient(c.Nats)

	bridge := bridge.NewBridge(publisher, subscriber, c.Mqtt.Ctx, c.Mqtt)
	bridge.StartBridge()

	<-done

	log.Println("mqtt adapter is shutting down...")
}
