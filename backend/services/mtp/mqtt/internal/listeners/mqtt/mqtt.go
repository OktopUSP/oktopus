package mqtt

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"os"

	rv8 "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/hooks/auth"
	"github.com/mochi-co/mqtt/v2/hooks/storage/redis"
	"github.com/mochi-co/mqtt/v2/listeners"
	"github.com/rs/zerolog"
)

var (
	server *mqtt.Server
)

type Mqtt struct {
	Port      string
	Tls       bool
	Fullchain string
	Privkey   string
	AuthFile  string
	Redis     Redis
	LogLevel  int
}

type Redis struct {
	RedisEnable   bool
	RedisAddr     string
	RedisPassword string
}

func (m *Mqtt) Start(mqttServer *mqtt.Server) {

	defineSeverLog(mqttServer, m.LogLevel)
	defineServerAuth(mqttServer, m.AuthFile)

	server = mqttServer

	var tlsConfig *listeners.Config
	if m.Tls {
		tlsConfig = defineServerTls(m.Fullchain, m.Privkey)
	}

	createListener(mqttServer, m.Port, tlsConfig)
	addHooks(mqttServer, m.Redis)
}

func addHooks(server *mqtt.Server, redisConf Redis) {

	err := server.AddHook(new(MyHook), map[string]any{})
	if err != nil {
		log.Fatal(err)
	}

	if redisConf.RedisEnable {
		if redisConf.RedisAddr != "" {
			err = server.AddHook(new(redis.Hook), &redis.Options{
				Options: &rv8.Options{
					Addr:     redisConf.RedisAddr,
					Password: redisConf.RedisPassword,
					DB:       0,
				},
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func createListener(server *mqtt.Server, port string, listenersConf *listeners.Config) {
	tcp := listeners.NewTCP(uuid.NewString(), port, listenersConf)

	err := server.AddListener(tcp)
	if err != nil {
		log.Fatal(err)
	}
}

func defineServerTls(fullchain, privkey string) *listeners.Config {
	if fullchain != "" && privkey != "" {
		chain, err := ioutil.ReadFile(fullchain)
		if err != nil {
			log.Fatal(err)
		}

		pv, err := ioutil.ReadFile(privkey)
		if err != nil {
			log.Fatal(err)
		}

		cert, err := tls.X509KeyPair(chain, pv)
		if err != nil {
			log.Fatal(err)
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		return &listeners.Config{
			TLSConfig: tlsConfig,
		}

	}
	return nil
}

func defineServerAuth(server *mqtt.Server, authFile string) {
	if authFile != "" {
		data, err := os.ReadFile(authFile)
		if err != nil {
			log.Fatal(err)
		}
		err = server.AddHook(new(auth.Hook), &auth.Options{
			Data: data,
		})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := server.AddHook(new(auth.AllowHook), nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func defineSeverLog(server *mqtt.Server, logLevel int) {
	l := server.Log.Level(zerolog.Level(logLevel))
	server.Log = &l
}
