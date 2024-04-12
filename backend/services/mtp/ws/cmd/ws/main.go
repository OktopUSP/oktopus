package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/OktopUSP/oktopus/ws/internal/config"
	"github.com/OktopUSP/oktopus/ws/internal/nats"
	"github.com/OktopUSP/oktopus/ws/internal/ws"
)

func main() {

	done := make(chan os.Signal, 1)

	conf := config.NewConfig()

	// Locks app running until it receives a stop command as Ctrl+C.
	signal.Notify(done, syscall.SIGINT)

	_, kv := nats.StartNatsClient(conf.Nats)
	
	ws.StartNewServer(conf, kv)

	<-done

	log.Println("(⌐■_■) Websockets server is out!")
}
