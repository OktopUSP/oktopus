package config

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const LOCAL_ENV = ".env.local"

type Config struct {
	MqttPort      string
	Tls           bool
	Fullchain     string
	Privkey       string
	AuthFile      string
	RedisEnable   bool
	RedisAddr     string
	RedisPassword string
	WsEnable      bool
	WsPort        string
	HttpEnable    bool
	HttpPort      string
	LogLevel      int
}

//TODO: debug websocket and http listeners

func NewConfig() Config {

	loadEnvVariables()

	/*
		App variables priority:
		1º - Flag through command line.
		2º - Env variables.
		3º - Default flag value.
	*/

	mqttPort := flag.String("mqtt_port", lookupEnvOrString("MQTT_PORT", ":1883"), "port for MQTT listener")
	tls := flag.Bool("mqtt_tls", lookupEnvOrBool("MQTT_TLS", false), "enable/disable TLS")
	fullchain := flag.String("full_chain_path", lookupEnvOrString("FULL_CHAIN_PATH", ""), "path to fullchain.pem certificate")
	privkey := flag.String("private_key_path", lookupEnvOrString("PRIVATE_KEY_PATH", ""), "path to privkey.pem certificate")
	authFile := flag.String("auth_file_path", lookupEnvOrString("AUTH_FILE_PATH", ""), "path to MQTT RBAC auth file")
	redisEnable := flag.Bool("redis_enable", lookupEnvOrBool("REDIS_ENABLE", true), "enable/disable Redis db")
	redisAddr := flag.String("redis_addr", lookupEnvOrString("REDIS_ADDR", "localhost:6379"), "address of redis db")
	redisPassword := flag.String("redis_passwd", lookupEnvOrString("REDIS_PASSWD", ""), "redis db password")
	wsEnable := flag.Bool("ws_enable", lookupEnvOrBool("WS_ENABLE", false), "enable/disable Websocket listener")
	wsPort := flag.String("ws_port", lookupEnvOrString("WS_PORT", ":80"), "port for Websocket listener")
	httpEnable := flag.Bool("http_enable", lookupEnvOrBool("HTTP_ENABLE", false), "enable/disable HTTP listener of mqtt metrics")
	httpPort := flag.String("http_port", lookupEnvOrString("HTTP_PORT", ":8080"), "port for HTTP listener of mqtt metrics")
	logLevel := flag.Int("log_level", lookupEnvOrInt("LOG_LEVEL", 1), "0=DEBUG, 1=INFO, 2=WARNING, 3=ERROR")

	flag.Parse()
	flHelp := flag.Bool("help", false, "Help")

	if *flHelp {
		flag.Usage()
		os.Exit(0)
	}

	conf := Config{
		MqttPort:      *mqttPort,
		Tls:           *tls,
		Fullchain:     *fullchain,
		Privkey:       *privkey,
		AuthFile:      *authFile,
		RedisEnable:   *redisEnable,
		RedisAddr:     *redisAddr,
		RedisPassword: *redisPassword,
		WsEnable:      *wsEnable,
		WsPort:        *wsPort,
		HttpEnable:    *httpEnable,
		HttpPort:      *httpPort,
		LogLevel:      *logLevel,
	}

	conf.validate()

	return conf
}

func (c *Config) validate() {

	valid := true

	if c.Tls && (c.Fullchain == "" || c.Privkey == "") {
		log.Println("TLS is enabled, but fullchain and privkey are not set")
		valid = false
	}
	if c.LogLevel > 3 || c.LogLevel < 0 {
		log.Println("Log level not valid, choose a number between 0 and 3")
		valid = false
	}

	if !valid {
		log.Println("For more info execute --help")
		os.Exit(1)
	}

}

func loadEnvVariables() {
	err := godotenv.Load()

	if _, err := os.Stat(LOCAL_ENV); err == nil {
		_ = godotenv.Overload(LOCAL_ENV)
		log.Printf("Loaded variables from '%s'", LOCAL_ENV)
	} else {
		log.Println("Loaded variables from '.env'")
	}
	if err != nil {
		log.Println("Error to load environment variables:", err)
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
