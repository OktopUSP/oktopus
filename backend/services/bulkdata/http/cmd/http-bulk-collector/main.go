package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/oktopUSP/backend/services/bulkdata/internal/api"
	"github.com/oktopUSP/backend/services/bulkdata/internal/bridge"
	"github.com/oktopUSP/backend/services/bulkdata/internal/config"
	"github.com/oktopUSP/backend/services/bulkdata/internal/nats"
)

func main() {
	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGINT)

	c := config.NewConfig()

	pub, sub := nats.StartNatsClient(c.Nats)

	b := bridge.NewBridge(pub, sub)

	server := api.NewApi(c.RestApi, b)

	server.StartApi()

	<-done
	log.Println("bulk data collector is saying adios ...")
}
