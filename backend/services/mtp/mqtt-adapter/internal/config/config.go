package config

import (
	"context"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const LOCAL_ENV = ".env.local"

type Nats struct {
	Url                string
	Name               string
	VerifyCertificates bool
	Ctx                context.Context
}

type Mqtt struct {
	Url        string
	UrlForTls  string
	SkipVerify bool
	ClientId   string
	Username   string
	Password   string
	Qos        int
	Ctx        context.Context
}

type Config struct {
	Nats Nats
	Mqtt Mqtt
}

func NewConfig() *Config {

	loadEnvVariables()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	natsUrl := flag.String("nats_url", lookupEnvOrString("NATS_URL", "nats://localhost:4222"), "url for nats server")
	natsName := flag.String("nats_name", lookupEnvOrString("NATS_NAME", "mqtt-adapter"), "name for nats client")
	natsVerifyCertificates := flag.Bool("nats_verify_certificates", lookupEnvOrBool("NATS_VERIFY_CERTIFICATES", false), "verify validity of certificates from nats server")
	mqttUrl := flag.String("mqtt_url", lookupEnvOrString("MQTT_URL", "tcp://localhost:1883"), "url for mqtt server")
	mqttsUrl := flag.String("mqtts_url", lookupEnvOrString("MQTTS_URL", ""), "url for mqtts server")
	mqttsSkipVerify := flag.Bool("mqtts_skip_verify", lookupEnvOrBool("MQTTS_SKIP_VERIFY", false), "skip verification of server certificate for mqtts")
	mqttClientId := flag.String("mqtt_client_id", lookupEnvOrString("MQTT_CLIENT_ID", "mqtt-adapter"), "client id for mqtt")
	mqttUsername := flag.String("mqtt_username", lookupEnvOrString("MQTT_USERNAME", "oktopusController"), "username for mqtt")
	mqttQos := flag.Int("mqtt_qos", lookupEnvOrInt("MQTT_QOS", 1), "quality of service for mqtt")
	flHelp := flag.Bool("help", false, "Help")

	/*
		App variables priority:
		1º - Flag through command line.
		2º - Env variables.
		3º - Default flag value.
	*/

	flag.Parse()

	if *flHelp {
		flag.Usage()
		os.Exit(0)
	}

	ctx := context.TODO()

	return &Config{
		Nats: Nats{
			Url:                *natsUrl,
			Name:               *natsName,
			VerifyCertificates: *natsVerifyCertificates,
			Ctx:                ctx,
		},
		Mqtt: Mqtt{
			Url:        *mqttUrl,
			UrlForTls:  *mqttsUrl,
			SkipVerify: *mqttsSkipVerify,
			ClientId:   *mqttClientId,
			Username:   *mqttUsername,
			Ctx:        ctx,
			Qos:        *mqttQos,
		},
	}
}

func loadEnvVariables() {
	err := godotenv.Load()

	if _, err := os.Stat(LOCAL_ENV); err == nil {
		_ = godotenv.Overload(LOCAL_ENV)
		log.Printf("Loaded variables from '%s'", LOCAL_ENV)
	}

	if err != nil {
		log.Println("Error to load environment variables:", err)
	} else {
		log.Println("Loaded variables from '.env'")
	}
}

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
