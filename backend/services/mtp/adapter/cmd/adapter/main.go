package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/config"
	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/db"
	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/events"
	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/events/cwmp_handler"
	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/events/usp_handler"
	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/nats"
	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/reqs"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	c := config.NewConfig()

	js, nc := nats.StartNatsClient(c.Nats, c.Controller)

	db := db.NewDatabase(c.Mongo.Ctx, c.Mongo.Uri)

	usp_handler := usp_handler.NewHandler(nc, js, db, c.Controller.ControllerId)
	cwmp_handler := cwmp_handler.NewHandler(nc, js, db, c.Controller.ControllerId)

	events.StartEventsListener(c.Nats.Ctx, js, usp_handler, cwmp_handler)

	reqs.StartRequestsListener(c.Nats.Ctx, nc, db)

	<-done

	log.Println("mtp adapter is shutting down...")
}
