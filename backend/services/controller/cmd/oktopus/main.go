// Made by Leandro Antônio Farias Machado

package main

import (
	"context"
	"flag"
	"github.com/leandrofars/oktopus/internal/api"
	"github.com/leandrofars/oktopus/internal/db"
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
	// fl_endpointId := flag.String("endpoint_id", "proto::oktopus-controller", "Defines the enpoint id the Agent must trust on.")
	flDevicesTopic := flag.String("d", "oktopus/devices", "That's the topic mqtt broker end new devices info.")
	flDisconTopic := flag.String("dis", "oktopus/disconnect", "It's where disconnected IoTs are known.")
	flSubTopic := flag.String("sub", "oktopus/+/controller/+", "That's the topic agent must publish to, and the controller keeps on listening.")
	flBrokerAddr := flag.String("a", "localhost", "Mqtt broker adrress")
	flBrokerPort := flag.String("p", "1883", "Mqtt broker port")
	flTlsCert := flag.String("ca", "", "TLS ca certificate")
	flBrokerUsername := flag.String("u", "", "Mqtt broker username")
	flBrokerPassword := flag.String("P", "", "Mqtt broker password")
	flBrokerClientId := flag.String("i", "", "A clientid for the Mqtt connection")
	flBrokerQos := flag.Int("q", 2, "Quality of service of mqtt messages delivery")
	flAddrDB := flag.String("mongo", "mongodb://localhost:27017/", "MongoDB URI")
	flApiPort := flag.String("ap", "8000", "Rest api port")
	flHelp := flag.Bool("help", false, "Help")

	flag.Parse()

	if *flHelp {
		flag.Usage()
		os.Exit(0)
	}
	/*
		This context suppress our needs, but we can use a more sofisticate
		approach with cancel and timeout options passing it through paho mqtt functions.
	*/
	ctx, cancel := context.WithCancel(context.Background())
	database := db.NewDatabase(ctx, *flAddrDB)
	/*
	 If you want to use another message protocol just make it implement Broker interface.
	*/
	mqttClient := mqtt.Mqtt{
		Addr:            *flBrokerAddr,
		Port:            *flBrokerPort,
		Id:              *flBrokerClientId,
		User:            *flBrokerUsername,
		Passwd:          *flBrokerPassword,
		Ctx:             ctx,
		QoS:             *flBrokerQos,
		SubTopic:        *flSubTopic,
		DevicesTopic:    *flDevicesTopic,
		DisconnectTopic: *flDisconTopic,
		CA:              *flTlsCert,
		DB:              database,
	}

	mtp.MtpService(&mqttClient, done)
	a := api.NewApi(*flApiPort, database)
	api.StartApi(a)

	<-done
	cancel()

	log.Println("(⌐■_■) Oktopus is out!")

}
