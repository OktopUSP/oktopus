// Made by Leandro Antônio Farias Machado

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/leandrofars/oktopus/internal/api"
	"github.com/leandrofars/oktopus/internal/db"
	"github.com/leandrofars/oktopus/internal/mqtt"
	"github.com/leandrofars/oktopus/internal/mtp"
	"github.com/leandrofars/oktopus/internal/stomp"
	usp_msg "github.com/leandrofars/oktopus/internal/usp_message"
	"github.com/leandrofars/oktopus/internal/ws"
)

// TODO: refact where this version number comes from
const VERSION = "0.0.1"

func main() {
	done := make(chan os.Signal, 1)

	err := godotenv.Load()

	localEnv := ".env.local"
	if _, err := os.Stat(localEnv); err == nil {
		_ = godotenv.Overload(localEnv)
		log.Println("Loaded variables from '.env.local'")
	} else {
		log.Println("Loaded variables from '.env'")
	}

	if err != nil {
		log.Println("Error to load environment variables:", err)
	}

	// Locks app running until it receives a stop command as Ctrl+C.
	signal.Notify(done, syscall.SIGINT)

	//TODO: refact app confiurations and env loading to another package
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	/*
		App variables priority:
		1º - Flag through command line.
		2º - Env variables.
		3º - Default flag value.
	*/

	log.Println("Starting Oktopus Project TR-369 Controller Version:", VERSION)
	// fl_endpointId := flag.String("endpoint_id", "proto::oktopus-controller", "Defines the enpoint id the Agent must trust on.")
	flDevicesTopic := flag.String("d", lookupEnvOrString("DEVICES_STATUS_TOPIC", "oktopus/+/status/+"), "That's the topic mqtt broker send new devices info.")
	flSubTopic := flag.String("sub", lookupEnvOrString("DEVICE_PUB_TOPIC", "oktopus/+/controller/+"), "That's the topic agent must publish to")
	flBrokerAddr := flag.String("a", lookupEnvOrString("BROKER_ADDR", "localhost"), "Mqtt broker adrress")
	flBrokerPort := flag.String("p", lookupEnvOrString("BROKER_PORT", "1883"), "Mqtt broker port")
	flTlsCert := flag.Bool("tls", lookupEnvOrBool("BROKER_TLS", false), "Connect to broker over TLS")
	flBrokerUsername := flag.String("u", lookupEnvOrString("BROKER_USERNAME", ""), "Mqtt broker username")
	flBrokerPassword := flag.String("P", lookupEnvOrString("BROKER_PASSWORD", ""), "Mqtt broker password")
	flBrokerClientId := flag.String("i", lookupEnvOrString("BROKER_CLIENTID", ""), "A clientid for the Mqtt connection")
	flBrokerQos := flag.Int("q", lookupEnvOrInt("BROKER_QOS", 0), "Quality of service of mqtt messages delivery")
	flAddrDB := flag.String("mongo", lookupEnvOrString("MONGO_URI", "mongodb://localhost:27017"), "MongoDB URI")
	flApiPort := flag.String("ap", lookupEnvOrString("REST_API_PORT", "8000"), "Rest api port")
	flStompAddr := flag.String("stomp", lookupEnvOrString("STOMP_ADDR", "127.0.0.1:61613"), "Stomp broker address")
	flStompUser := flag.String("stomp_user", lookupEnvOrString("STOMP_USERNAME", ""), "Stomp broker username")
	flStompPasswd := flag.String("stomp_passwd", lookupEnvOrString("STOMP_PASSWORD", ""), "Stomp broker password")
	flWsToken := flag.String("ws_token", lookupEnvOrString("WS_TOKEN", ""), "Websocket token")
	flWsAuth := flag.Bool("ws_auth", lookupEnvOrBool("WS_AUTH", true), "Websocket auth enable or not")
	flWsAddr := flag.String("ws_addr", lookupEnvOrString("WS_ADDR", "localhost"), "Websocket server address")
	flWsPort := flag.String("ws_port", lookupEnvOrString("WS_PORT", "8080"), "Websocket server port")
	flWsRoute := flag.String("ws_route", lookupEnvOrString("WS_ROUTE", "/ws/controller"), "Websocket server route")
	flWsTls := flag.Bool("ws_tls", lookupEnvOrBool("WS_TLS", false), "Websocket server tls")
	flWsSkipVerify := flag.Bool("ws_skip_verify", lookupEnvOrBool("WS_SKIP_VERIFY", false), "Websocket skip tls certificate verify")
	flDisableWs := flag.Bool("ws_disable", lookupEnvOrBool("WS_DISABLE", false), "Disable WS MTP")
	flDisableStomp := flag.Bool("stomp_disable", lookupEnvOrBool("STOMP_DISABLE", false), "Disable STOMP MTP")
	flDisableMqtt := flag.Bool("mqtt_disable", lookupEnvOrBool("MQTT_DISABLE", false), "Disable MQTT MTP")
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
	apiMsgQueue := make(map[string](chan usp_msg.Msg))
	var m sync.Mutex

	//TODO: refact mtps initialization through main.go
	/*
	 If you want to use another message protocol just make it implement Broker interface.
	*/
	log.Println("Start MTP protocols: MQTT | Websockets | STOMP")

	if *flDisableMqtt && *flDisableStomp && *flDisableWs {
		log.Println("ERROR: you have to enable at least one MTP")
		os.Exit(0)
	}

	wg := new(sync.WaitGroup)
	wg.Add(3) // Three wait groups (mqtt, stomp, ws)

	/* ------------------------------ MTPs clients ------------------------------ */
	var stompClient stomp.Stomp
	var mqttClient mqtt.Mqtt
	var wsClient ws.Ws
	/* -------------------------------------------------------------------------- */

	/* ------------------------ MTPs disconnect channels ------------------------ */
	var mqttDone chan os.Signal
	var wsDone chan os.Signal
	var stompDone chan os.Signal
	/* -------------------------------------------------------------------------- */

	go func() {
		mqttClient = mqtt.Mqtt{
			Addr:         *flBrokerAddr,
			Port:         *flBrokerPort,
			Id:           *flBrokerClientId,
			User:         *flBrokerUsername,
			Passwd:       *flBrokerPassword,
			Ctx:          ctx,
			QoS:          *flBrokerQos,
			SubTopic:     *flSubTopic,
			DevicesTopic: *flDevicesTopic,
			TLS:          *flTlsCert,
			DB:           database,
			MsgQueue:     apiMsgQueue,
			QMutex:       &m,
		}

		mqttDone = make(chan os.Signal, 1)

		if !*flDisableMqtt {
			// MQTT will try connect to broker forever
			go mtp.MtpService(&mqttClient, mqttDone, wg)
		} else {
			wg.Done()
		}
	}()

	go func() {
		stompClient = stomp.Stomp{
			Addr:     *flStompAddr,
			Username: *flStompUser,
			Password: *flStompPasswd,
		}

		stompDone = make(chan os.Signal, 1)

		if !*flDisableStomp {
			// STOMP will try to connect for a bunch of times and then exit
			go mtp.MtpService(&stompClient, stompDone, wg)
		} else {
			wg.Done()
		}
	}()

	go func() {
		wsClient = ws.Ws{
			Addr:               *flWsAddr,
			Port:               *flWsPort,
			Token:              *flWsToken,
			Route:              *flWsRoute,
			Auth:               *flWsAuth,
			TLS:                *flWsTls,
			InsecureSkipVerify: *flWsSkipVerify,
			DB:                 database,
			Ctx:                ctx,
		}

		wsDone = make(chan os.Signal, 1)

		if !*flDisableWs {
			go mtp.MtpService(&wsClient, wsDone, wg)
		} else {
			wg.Done()
		}
	}()

	wg.Wait()

	a := api.NewApi(*flApiPort, database, &mqttClient, apiMsgQueue, &m)
	api.StartApi(a)

	<-done
	cancel()
	// send done signal to all MTPs
	wsDone <- os.Interrupt
	mqttDone <- os.Interrupt
	stompDone <- os.Interrupt

	log.Println("(⌐■_■) Oktopus is out!")

}

//TODO: refact functions below to another package

func lookupEnvOrString(key string, defaultVal string) string {
	if val, _ := os.LookupEnv(key); val != "" {
		return val
	}
	return defaultVal
}

func lookupEnvOrInt(key string, defaultVal int) int {
	if val, _ := os.LookupEnv(key); val != "" {
		v, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("LookupEnvOrInt[%s]: %v", key, err)
		}
		return v
	}
	return defaultVal
}

func lookupEnvOrBool(key string, defaultVal bool) bool {
	if val, _ := os.LookupEnv(key); val != "" {
		v, err := strconv.ParseBool(val)
		if err != nil {
			log.Fatalf("LookupEnvOrBool[%s]: %v", key, err)
		}
		return v
	}
	return defaultVal
}
