package main

import (
	"broker/internal/config"
	"broker/internal/listeners"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	done := make(chan os.Signal, 1)

	conf := config.NewConfig()

	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go listeners.StartServers(conf)

	if conf.WsEnable {
		log.Printf("websocket is running at port %s", conf.WsPort)
	}

	if conf.HttpEnable {
		log.Printf("http is running at port %s", conf.HttpPort)
	}

	log.Printf("mqtt is running at port %s", conf.MqttPort)

	<-done

	log.Println("server is shutting down...")
}
