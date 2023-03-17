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

	log.Println("Starting Oktopus Project TR-369 Controller Version:", VERSION)
	fl_broker := flag.Bool("m", false, "Defines if mosquitto container must run or not")
	// fl_endpointId := flag.String("endpoint_id", "proto::oktopus-controller", "Defines the enpoint id the Agent must trust on.")
	fl_sub_topic := flag.String("s", "oktopus/v1/agent", "That's the topic agent must publish to, and the controller keeps on listening.")
	// fl_pub_topic := flag.String("pub_topic", "oktopus/v1/controller", "That's the topic controller must publish to, and the agent keeps on listening.")
	fl_broker_addr := flag.String("a", "localhost", "Mqtt broker adrress")
	fl_broker_port := flag.String("p", "1883", "Mqtt broker port")
	fl_tls_cert := flag.String("ca", "", "TLS ca certificate")
	fl_broker_username := flag.String("u", "", "Mqtt broker username")
	fl_broker_password := flag.String("P", "", "Mqtt broker password")
	fl_broker_clientid := flag.String("i", "", "A clientid for the Mqtt connection")
	fl_broker_qos := flag.Int("q", 2, "Quality of service of mqtt messages delivery")
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
	/*
		This context suppress our needs, but we can use a more sofisticade
		approach with cancel and timeout options passing it through paho mqtt functions.
	*/
	ctx := context.Background()

	/*
	 If you want to use another message protocol just make it implement Broker interface.
	*/
	mqttClient := mqtt.Mqtt{
		Addr:     *fl_broker_addr,
		Port:     *fl_broker_port,
		Id:       *fl_broker_clientid,
		User:     *fl_broker_username,
		Passwd:   *fl_broker_password,
		Ctx:      ctx,
		QoS:      *fl_broker_qos,
		SubTopic: *fl_sub_topic,
		CA:       *fl_tls_cert,
	}

	mtp.MtpService(&mqttClient, done)

	<-done

	log.Println("(⌐■_■) Oktopus is out!")

}
