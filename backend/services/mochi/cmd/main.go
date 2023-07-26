package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	rv8 "github.com/go-redis/redis/v8"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/hooks/storage/redis"
	"github.com/mochi-co/mqtt/v2/packets"
	"github.com/rs/zerolog"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mochi-co/mqtt/v2/hooks/auth"
	"github.com/mochi-co/mqtt/v2/listeners"
)

var (
	//TODO: create custom mqtt server options
	server = mqtt.New(&mqtt.Options{
		//Capabilities: &mqtt.Capabilities{
		//	ServerKeepAlive:              10000,
		//	ReceiveMaximum:               math.MaxUint16,
		//	MaximumMessageExpiryInterval: math.MaxUint32,
		//	MaximumSessionExpiryInterval: math.MaxUint32, // maximum number of seconds to keep disconnected sessions
		//	MaximumClientWritesPending:   65536,
		//	MaximumPacketSize:            0,
		//	MaximumQos:                   2,
		//},
	})
)

func main() {
	tcpAddr := flag.String("tcp", ":1883", "network address for TCP listener")
	redisAddr := flag.String("redis", "172.17.0.2:6379", "host address of redis db")
	redisPassword := flag.String("redis_passwd", "", "redis db password")
	wsAddr := flag.String("ws", "", "network address for Websocket listener")
	infoAddr := flag.String("info", ":8080", "network address for web info dashboard listener")
	path := flag.String("path", "", "path to data auth file")
	fullchain := flag.String("full_chain_path", "", "path to fullchain.pem certificate")
	privkey := flag.String("private_key_path", "", "path to privkey.pem certificate")
	logLevel := flag.Int("logLevel", 1, "log level, default is INFO, 0 value is DEBUG")

	flag.Parse()

	if *logLevel > 2 || *logLevel < 0 {
		log.Println("Log level not valid, choose a number between 0 and 7")
		log.Println("For more info access zeroLog documentation: https://github.com/rs/zerolog")
		os.Exit(1)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	serverForTLS := mqtt.New(&mqtt.Options{})

	lTls := serverForTLS.Log.Level(zerolog.Level(*logLevel))
	serverForTLS.Log = &lTls

	l := server.Log.Level(zerolog.Level(*logLevel))
	server.Log = &l

	if *path != "" {
		data, err := os.ReadFile(*path)
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

	if *fullchain != "" && *privkey != "" {
		chain, err := ioutil.ReadFile(*fullchain)
		if err != nil {
			log.Fatal(err)
		}

		pv, err := ioutil.ReadFile(*privkey)
		if err != nil {
			log.Fatal(err)
		}

		cert, err := tls.X509KeyPair(chain, pv)
		if err != nil {
			log.Fatal(err)
		}

		//Basic TLS Config
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		if *tcpAddr != "" {
			tcp := listeners.NewTCP("t1", ":8883", &listeners.Config{
				TLSConfig: tlsConfig,
			})
			err := serverForTLS.AddListener(tcp)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = serverForTLS.AddHook(new(MyHook), map[string]any{})
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Mqtt Broker is running with TLS at port 8883")
	}
	if *tcpAddr != "" {
		//tcp := listeners.NewTCP("t1", *tcpAddr, &listeners.Config{
		//	TLSConfig: tlsConfig,
		//})
		tcp := listeners.NewTCP("t1", ":1883", nil)
		err := server.AddListener(tcp)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("Mqtt Broker is running without TLS at port 1883")

	if *wsAddr != "" {
		ws := listeners.NewWebsocket("ws1", *wsAddr, nil)
		err := server.AddListener(ws)
		if err != nil {
			log.Fatal(err)
		}
	}

	if *infoAddr != "" {
		stats := listeners.NewHTTPStats("stats", *infoAddr, nil, server.Info)
		err := server.AddListener(stats)
		if err != nil {
			log.Fatal(err)
		}
	}

	err := server.AddHook(new(MyHook), map[string]any{})
	if err != nil {
		log.Fatal(err)
	}

	if *redisAddr != "" {
		err = server.AddHook(new(redis.Hook), &redis.Options{
			Options: &rv8.Options{
				Addr:     *redisAddr,     // default redis address
				Password: *redisPassword, // your password
				DB:       0,              // your redis db
			},
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
		err = serverForTLS.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-done
	server.Log.Warn().Msg("caught signal, stopping...")
	server.Close()
	server.Log.Info().Msg("main.go finished")

}

type MyHook struct {
	mqtt.HookBase
}

func (h *MyHook) ID() string {
	return "events-controller"
}

func (h *MyHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnSubscribed,
		mqtt.OnDisconnect,
		mqtt.OnClientExpired,
		mqtt.OnPacketEncode,
	}, []byte{b})
}

func (h *MyHook) Init(config any) error {
	h.Log.Info().Msg("initialised")
	return nil
}

func (h *MyHook) OnClientExpired(cl *mqtt.Client) {
	log.Printf("Client id %s expired", cl.ID)
}

func (h *MyHook) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
	var clUser string
	if len(cl.Properties.Props.User) > 0 {
		clUser = cl.Properties.Props.User[0].Val
	}

	if clUser != "" {
		err := server.Publish("oktopus/v1/status/"+clUser, []byte("1"), false, 1)
		if err != nil {
			log.Println("server publish error: ", err)
		}
	}
}

func (h *MyHook) OnSubscribed(cl *mqtt.Client, pk packets.Packet, reasonCodes []byte) {
	// Verifies if it's a device who is subscribed
	if strings.Contains(pk.Filters[0].Filter, "oktopus/v1/agent") {
		var clUser string

		if len(cl.Properties.Props.User) > 0 {
			clUser = cl.Properties.Props.User[0].Val
		}

		if clUser != "" {
			cl.Properties.Will = mqtt.Will{
				Qos:       1,
				TopicName: "oktopus/v1/status/" + clUser,
				Payload:   []byte("1"),
				Retain:    false,
			}
			log.Println("new device:", clUser)
			err := server.Publish("oktopus/v1/status/"+clUser, []byte("0"), false, 1)
			if err != nil {
				log.Println("server publish error: ", err)
			}
		}

	}
}

func (h *MyHook) OnPacketEncode(cl *mqtt.Client, pk packets.Packet) packets.Packet {
	var clUser string
	if len(cl.Properties.Props.User) > 0 {
		clUser = cl.Properties.Props.User[0].Val
	}
	if pk.FixedHeader.Type == packets.Connack {
		pk.Properties.User = []packets.UserProperty{{Key: "subscribe-topic", Val: "oktopus/v1/agent/" + clUser}}
	}

	return pk
}
