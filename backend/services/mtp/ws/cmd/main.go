package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/OktopUSP/oktopus/ws/internal/config"
	"github.com/OktopUSP/oktopus/ws/internal/ws"
)

// TODO: refact from where this version is get
const VERSION = "0.0.1"

func main() {

	done := make(chan os.Signal, 1)

	conf := config.NewConfig()

	// Locks app running until it receives a stop command as Ctrl+C.
	signal.Notify(done, syscall.SIGINT)

	log.Println("Starting Oktopus Websockets Version:", VERSION)
	ws.StartNewServer(conf)

	<-done

	log.Println("(⌐■_■) Websockets server is out!")
}
