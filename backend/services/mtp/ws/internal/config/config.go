// Loads environemnt variables and returns a config struct
package config

import (
	"context"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string // server port: e.g. ":8080"
	Auth          bool   // server auth enable/disable
	ControllerEID string // controller endpoint id
	Tls           bool   // enable/diable websockets server tls
	TlsPort       string
	NoTls         bool
	FullChain     string
	PrivateKey    string
	Nats          Nats
}

type Nats struct {
	Url       string
	Name      string
	EnableTls bool
	Cert      Tls
	Ctx       context.Context
}

type Tls struct {
	CertFile string
	KeyFile  string
	CaFile   string
}

func NewConfig() Config {

	//Defines log format
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	loadEnv()

	/*
		App variables priority:
		1º - Flag through command line.
		2º - Env variables.
		3º - Default flag value.
	*/

	/* ------------------------------ define flags ------------------------------ */
	natsUrl := flag.String("nats_url", lookupEnvOrString("NATS_URL", "nats://localhost:4222"), "url for nats server")
	natsName := flag.String("nats_name", lookupEnvOrString("NATS_NAME", "ws-adapter"), "name for nats client")
	natsEnableTls := flag.Bool("nats_enable_tls", lookupEnvOrBool("NATS_ENABLE_TLS", false), "enbale TLS to nats server")
	clientCrt := flag.String("client_crt", lookupEnvOrString("CLIENT_CRT", "cert.pem"), "client certificate file to TLS connection")
	clientKey := flag.String("client_key", lookupEnvOrString("CLIENT_KEY", "key.pem"), "client key file to TLS connection")
	serverCA := flag.String("server_ca", lookupEnvOrString("SERVER_CA", "rootCA.pem"), "server CA file to TLS connection")
	flPort := flag.String("port", lookupEnvOrString("SERVER_PORT", ":8080"), "Server port")
	flAuth := flag.Bool("auth", lookupEnvOrBool("SERVER_AUTH_ENABLE", false), "Server auth enable/disable")
	flControllerEid := flag.String("controller-eid", lookupEnvOrString("CONTROLLER_EID", "oktopusController"), "Controller eid")
	flTls := flag.Bool("tls", lookupEnvOrBool("SERVER_TLS_ENABLE", false), "Enable/disable websockets server tls")
	flTlsPort := flag.String("tls_port", lookupEnvOrString("SERVER_TLS_PORT", ":8081"), "Server Port to use if TLS is enabled")
	flNoTls := flag.Bool("no_tls", lookupEnvOrBool("SERVER_NO_TLS", false), "Disable/enable websockets serevr without tls")
	flFullchain := flag.String("fullchain_path", lookupEnvOrString("FULL_CHAIN_PATH", "cert.pem"), "Fullchain file path")
	flPrivKey := flag.String("privkey_path", lookupEnvOrString("PRIVATE_KEY_PATH", "key.pem"), "Private key file path")
	flHelp := flag.Bool("help", false, "Help")
	flag.Parse()
	/* -------------------------------------------------------------------------- */

	if *flHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *flNoTls && !*flTls {
		log.Fatalf("You must at least choose one between tls and no_tls configs")
	}

	ctx := context.TODO()

	return Config{
		Port:          *flPort,
		TlsPort:       *flTlsPort,
		NoTls:         *flNoTls,
		Auth:          *flAuth,
		ControllerEID: *flControllerEid,
		Tls:           *flTls,
		FullChain:     *flFullchain,
		PrivateKey:    *flPrivKey,
		Nats: Nats{
			Url:       *natsUrl,
			Name:      *natsName,
			EnableTls: *natsEnableTls,
			Ctx:       ctx,
			Cert: Tls{
				CertFile: *clientCrt,
				KeyFile:  *clientKey,
				CaFile:   *serverCA,
			},
		},
	}
}

// Load environment variables from .env or .env.local file
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error to load environment variables:", err)
	}

	localEnv := ".env.local"
	if _, err := os.Stat(localEnv); err == nil {
		_ = godotenv.Overload(localEnv)
		log.Println("Loaded variables from '.env.local'")
	} else {
		log.Println("Loaded variables from '.env'")
	}
}

/* ---------------------------- helper functions ---------------------------- */
/*
	They are used to lookup if a environment variable is set with a value
	different of "" and return it.
	In case the var doesn't exist, it returns the default value.
	Also, they're useful to convert the string value of vars to the desired type.
*/

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

/* -------------------------------------------------------------------------- */
