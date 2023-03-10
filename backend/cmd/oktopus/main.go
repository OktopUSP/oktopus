// Made by Leandro Antônio Farias Machado (leandrofars@gmail.com)

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/eclipse/paho.golang/paho"
	"github.com/leandrofars/oktopus/internal/mqtt"
)

const VERSION = "0.0.1"

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Starting Oktopus Project TR-369 Controller Version:", VERSION)
	fl_broker := flag.Bool("mosquitto", false, "Defines if mosquitto container must run or not")
	// fl_endpointId := flag.String("endpoint_id", "proto::oktopus-controller", "Defines the enpoint id the Agent must trust on.")
	// fl_sub_topic := flag.String("sub_topic", "oktopus/v1/agent", "That's the topic agent must publish to, and the controller keeps on listening.")
	// fl_pub_topic := flag.String("pub_topic", "oktopus/v1/controller", "That's the topic controller must publish to, and the agent keeps on listening.")
	fl_broker_addr := flag.String("broker_addr", "localhost", "Mqtt broker adrress")
	fl_broker_port := flag.String("broker_port", "1883", "Mqtt broker port")
	fl_broker_username := flag.String("broker_user", "", "Mqtt broker username")
	fl_broker_password := flag.String("password", "", "Mqtt broker password")
	fl_broker_clientid := flag.String("clientid", "", "A clientid for the Mqtt connection")
	fl_help := flag.Bool("help", false, "Help")

	flag.Parse()

	if *fl_help {
		flag.Usage()
		os.Exit(0)
	}
	if *fl_broker {
		log.Println("Starting Mqtt Broker")
		mqtt.StartMqttBroker()
	}

	newClient := mqtt.StartMqttClient(fl_broker_addr, fl_broker_port)

	newConnection := mqtt.StartNewConnection(*fl_broker_clientid, *fl_broker_username, *fl_broker_password)

	mqtt.ConnectMqttBroker(newClient, newConnection, fl_broker_addr)

	<-done

	log.Println("Disconnecting broker")
	if newClient != nil {
		d := &paho.Disconnect{ReasonCode: 0}
		err := newClient.Disconnect(d)
		if err != nil {
			log.Fatalf("failed to send Disconnect: %s", err)
		}
	}

	log.Println("Oktopus is out!")

}
