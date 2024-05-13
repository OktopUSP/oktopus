package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/oktopUSP/backend/services/bulkdata/internal/api"
	"github.com/oktopUSP/backend/services/bulkdata/internal/config"
	"github.com/oktopUSP/backend/services/bulkdata/internal/nats"
)

func main() {
	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGINT)

	c := config.NewConfig()

	js, nc, kv := nats.StartNatsClient(c.Nats)

	server := api.NewApi(c.RestApi, js, nc, kv)

	server.StartApi()

	<-done
	log.Println("bulk data collector is saying adios ...")
}
