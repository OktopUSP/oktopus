// Loads environemnt variables and returns a config struct
package config

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string // server port: e.g. ":8080"
	Auth          bool   // server auth enable/disable
	Token         string // controller auth token
	ControllerEID string // controller endpoint id
	Tls           bool   // enable/diable websockets server tls
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
	flPort := flag.String("port", lookupEnvOrString("SERVER_PORT", ":8080"), "Server port")
	flToken := flag.String("token", lookupEnvOrString("SERVER_AUTH_TOKEN", ""), "Controller auth token")
	flAuth := flag.Bool("auth", lookupEnvOrBool("SERVER_AUTH_ENABLE", false), "Server auth enable/disable")
	flControllerEid := flag.String("controller-eid", lookupEnvOrString("CONTROLLER_EID", "oktopusController"), "Controller eid")
	flTls := flag.Bool("tls", lookupEnvOrBool("SERVER_TLS_ENABLE", false), "Enable/diable websockets server tls")
	flHelp := flag.Bool("help", false, "Help")
	flag.Parse()
	/* -------------------------------------------------------------------------- */

	if *flHelp {
		flag.Usage()
		os.Exit(0)
	}

	return Config{
		Port:          *flPort,
		Token:         *flToken,
		Auth:          *flAuth,
		ControllerEID: *flControllerEid,
		Tls:           *flTls,
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
