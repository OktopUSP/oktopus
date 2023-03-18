// Made by Leandro Antônio Farias Machado

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/leandrofars/oktopus/internal/mqtt"
	"github.com/leandrofars/oktopus/internal/mtp"
)

const VERSION = "0.0.1"

func main() {
	done := make(chan os.Signal, 1)

	// Locks app running until it receives a stop command as Ctrl+C.
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Starting Oktopus Project TR-369 Controller Version:", VERSION)
	flBroker := flag.Bool("m", false, "Defines if mosquitto container must run or not")
	// fl_endpointId := flag.String("endpoint_id", "proto::oktopus-controller", "Defines the enpoint id the Agent must trust on.")
	flSubTopic := flag.String("s", "oktopus/+/+/agent", "That's the topic agent must publish to, and the controller keeps on listening.")
	// fl_pub_topic := flag.String("pub_topic", "oktopus/v1/controller", "That's the topic controller must publish to, and the agent keeps on listening.")
	flBrokerAddr := flag.String("a", "localhost", "Mqtt broker adrress")
	flBrokerPort := flag.String("p", "1883", "Mqtt broker port")
	flTlsCert := flag.String("ca", "", "TLS ca certificate")
	flBrokerUsername := flag.String("u", "", "Mqtt broker username")
	flBrokerPassword := flag.String("P", "", "Mqtt broker password")
	flBrokerClientId := flag.String("i", "", "A clientid for the Mqtt connection")
	flBrokerQos := flag.Int("q", 2, "Quality of service of mqtt messages delivery")
	flHelp := flag.Bool("help", false, "Help")

	flag.Parse()

	if *flHelp {
		flag.Usage()
		os.Exit(0)
	}
	if *flBroker {
		log.Println("Starting Mqtt Broker")
		mqtt.StartMqttBroker()
	}
	/*
		This context suppress our needs, but we can use a more sofisticate
		approach with cancel and timeout options passing it through paho mqtt functions.
	*/
	ctx := context.Background()

	/*
	 If you want to use another message protocol just make it implement Broker interface.
	*/
	mqttClient := mqtt.Mqtt{
		Addr:     *flBrokerAddr,
		Port:     *flBrokerPort,
		Id:       *flBrokerClientId,
		User:     *flBrokerUsername,
		Passwd:   *flBrokerPassword,
		Ctx:      ctx,
		QoS:      *flBrokerQos,
		SubTopic: *flSubTopic,
		CA:       *flTlsCert,
	}

	log.Println()

	mtp.MtpService(&mqttClient, done)

	<-done

	log.Println("(⌐■_■) Oktopus is out!")

}
