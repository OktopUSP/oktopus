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

type Ws struct {
	Token         string
	AuthEnable    bool
	Addr          string
	Port          string
	Route         string
	TlsEnable     bool
	SkipTlsVerify bool
	Ctx           context.Context
}

type Config struct {
	Nats Nats
	Ws   Ws
}

func NewConfig() *Config {
	loadEnvVariables()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	natsUrl := flag.String("nats_url", lookupEnvOrString("NATS_URL", "nats://localhost:4222"), "url for nats server")
	natsName := flag.String("nats_name", lookupEnvOrString("NATS_NAME", "ws-adapter"), "name for nats client")
	natsVerifyCertificates := flag.Bool("nats_verify_certificates", lookupEnvOrBool("NATS_VERIFY_CERTIFICATES", false), "verify validity of certificates from nats server")
	wsToken := flag.String("ws_token", lookupEnvOrString("WS_TOKEN", ""), "websocket server auth token (if authentication is enabled)")
	wsAuthEnable := flag.Bool("ws_auth_enable", lookupEnvOrBool("WS_AUTH_ENABLE", false), "enable authentication for websocket server")
	wsAddr := flag.String("ws_addr", lookupEnvOrString("WS_ADDR", "localhost"), "websocket server address (domain or ip)")
	wsPort := flag.String("ws_port", lookupEnvOrString("WS_PORT", ":8080"), "websocket server port")
	wsRoute := flag.String("ws_route", lookupEnvOrString("WS_ROUTE", "/ws/controller"), "websocket server route")
	wsTlsEnable := flag.Bool("ws_tls_enable", lookupEnvOrBool("WS_TLS_ENABLE", false), "access websocket via tls protocol (wss)")
	wsSkipTlsVerify := flag.Bool("ws_skip_tls_verify", lookupEnvOrBool("WS_SKIP_TLS_VERIFY", false), "skip tls verification for websocket server")
	flHelp := flag.Bool("help", false, "Help")

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
		Ws: Ws{
			Token:         *wsToken,
			AuthEnable:    *wsAuthEnable,
			Addr:          *wsAddr,
			Port:          *wsPort,
			Route:         *wsRoute,
			TlsEnable:     *wsTlsEnable,
			SkipTlsVerify: *wsSkipTlsVerify,
			Ctx:           ctx,
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
