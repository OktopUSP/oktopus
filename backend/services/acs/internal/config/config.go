package config

import (
	"context"
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const LOCAL_ENV = ".env.local"

type Nats struct {
	Url                string
	Name               string
	VerifyCertificates bool
	Ctx                context.Context
}

type Acs struct {
	Port              string
	Tls               bool
	TlsPort           bool
	NoTls             bool
	KeepAliveInterval time.Duration
	Username          string
	Password          string
	Route             string
	DebugMode         bool
}

type Config struct {
	Acs  Acs
	Nats Nats
}

func NewConfig() *Config {

	loadEnvVariables()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	natsUrl := flag.String("nats_url", lookupEnvOrString("NATS_URL", "nats://localhost:4222"), "url for nats server")
	natsName := flag.String("nats_name", lookupEnvOrString("NATS_NAME", "adapter"), "name for nats client")
	natsVerifyCertificates := flag.Bool("nats_verify_certificates", lookupEnvOrBool("NATS_VERIFY_CERTIFICATES", false), "verify validity of certificates from nats server")
	acsPort := flag.String("acs_port", lookupEnvOrString("ACS_PORT", ":9292"), "port for acs server")
	acsRoute := flag.String("acs_route", lookupEnvOrString("ACS_ROUTE", "/acs"), "route for acs server")
	acsKeepAliveInterval := flag.Int("acs_keep_alive_interval", lookupEnvOrInt("KEEP_ALIVE_INTERVAL", 300), "keep alive interval in seconds for acs server")
	cwmpDebugMode := flag.Bool("debug_mode", lookupEnvOrBool("CWMP_DEBUG", false), "enable or disable cwmp logs in debug mode")
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
		Acs: Acs{
			Port:              *acsPort,
			Route:             *acsRoute,
			KeepAliveInterval: time.Duration(*acsKeepAliveInterval),
			DebugMode:         *cwmpDebugMode,
		},
	}
}

func loadEnvVariables() {
	err := godotenv.Load()

	if _, err := os.Stat(LOCAL_ENV); err == nil {
		_ = godotenv.Overload(LOCAL_ENV)
		log.Printf("Loaded variables from '%s'", LOCAL_ENV)
		return
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
