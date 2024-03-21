package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/leandrofars/oktopus/internal/api"
	"github.com/leandrofars/oktopus/internal/bridge"
	"github.com/leandrofars/oktopus/internal/config"
	"github.com/leandrofars/oktopus/internal/db"
	"github.com/leandrofars/oktopus/internal/nats"
)

func main() {
	done := make(chan os.Signal, 1)

	signal.Notify(done, syscall.SIGINT)

	c := config.NewConfig()

	js, nc := nats.StartNatsClient(c.Nats)

	bridge := bridge.NewBridge(js, nc)

	db := db.NewDatabase(c.Mongo.Ctx, c.Mongo.Uri)

	api := api.NewApi(c.RestApi, js, nc, bridge, db)
	api.StartApi()

	<-done

	log.Println("rest api is shutting down...")
}
