package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	rv8 "github.com/go-redis/redis/v8"
	"github.com/mochi-co/mqtt/v2/hooks/storage/redis"
	"github.com/mochi-co/mqtt/v2/packets"
	"github.com/rs/zerolog"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/hooks/auth"
	"github.com/mochi-co/mqtt/v2/listeners"
)

var (
	//	testCertificate = []byte(`-----BEGIN CERTIFICATE-----
	//MIIB/zCCAWgCCQDm3jV+lSF1AzANBgkqhkiG9w0BAQsFADBEMQswCQYDVQQGEwJB
	//VTETMBEGA1UECAwKU29tZS1TdGF0ZTERMA8GA1UECgwITW9jaGkgQ28xDTALBgNV
	//BAsMBE1RVFQwHhcNMjAwMTA0MjAzMzQyWhcNMjEwMTAzMjAzMzQyWjBEMQswCQYD
	//VQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTERMA8GA1UECgwITW9jaGkgQ28x
	//DTALBgNVBAsMBE1RVFQwgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAKz2bUz3
	//AOymssVLuvSOEbQ/sF8C/Ill8nRTd7sX9WBIxHJZf+gVn8lQ4BTQ0NchLDRIlpbi
	//OuZgktpd6ba8sIfVM4jbVprctky5tGsyHRFwL/GAycCtKwvuXkvcwSwLvB8b29EI
	//MLQ/3vNnYuC3eZ4qqxlODJgRsfQ7mUNB8zkLAgMBAAEwDQYJKoZIhvcNAQELBQAD
	//gYEAiMoKnQaD0F/J332arGvcmtbHmF2XZp/rGy3dooPug8+OPUSAJY9vTfxJwOsQ
	//qN1EcI+kIgrGxzA3VRfVYV8gr7IX+fUYfVCaPGcDCfPvo/Ihu757afJRVvpafWgy
	//zSpDZYu6C62h3KSzMJxffDjy7/2t8oYbTzkLSamsHJJjLZw=
	//-----END CERTIFICATE-----`)
	//
	//	testPrivateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
	//MIICXAIBAAKBgQCs9m1M9wDsprLFS7r0jhG0P7BfAvyJZfJ0U3e7F/VgSMRyWX/o
	//FZ/JUOAU0NDXISw0SJaW4jrmYJLaXem2vLCH1TOI21aa3LZMubRrMh0RcC/xgMnA
	//rSsL7l5L3MEsC7wfG9vRCDC0P97zZ2Lgt3meKqsZTgyYEbH0O5lDQfM5CwIDAQAB
	//AoGBAKlmVVirFqmw/qhDaqD4wBg0xI3Zw/Lh+Vu7ICoK5hVeT6DbTW3GOBAY+M8K
	//UXBSGhQ+/9ZZTmyyK0JZ9nw2RAG3lONU6wS41pZhB7F4siatZfP/JJfU6p+ohe8m
	//n22hTw4brY/8E/tjuki9T5e2GeiUPBhjbdECkkVXMYBPKDZhAkEA5h/b/HBcsIZZ
	//mL2d3dyWkXR/IxngQa4NH3124M8MfBqCYXPLgD7RDI+3oT/uVe+N0vu6+7CSMVx6
	//INM67CuE0QJBAMBpKW54cfMsMya3CM1BfdPEBzDT5kTMqxJ7ez164PHv9CJCnL0Z
	//AuWgM/p2WNbAF1yHNxw1eEfNbUWwVX2yhxsCQEtnMQvcPWLSAtWbe/jQaL2scGQt
	///F9JCp/A2oz7Cto3TXVlHc8dxh3ZkY/ShOO/pLb3KOODjcOCy7mpvOrZr6ECQH32
	//WoFPqImhrfryaHi3H0C7XFnC30S7GGOJIy0kfI7mn9St9x50eUkKj/yv7YjpSGHy
	//w0lcV9npyleNEOqxLXECQBL3VRGCfZfhfFpL8z+5+HPKXw6FxWr+p5h8o3CZ6Yi3
	//OJVN3Mfo6mbz34wswrEdMXn25MzAwbhFQvCVpPZrFwc=
	//-----END RSA PRIVATE KEY-----`)

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
		err := server.Publish("oktopus/v1/disconnect/"+clUser, []byte(""), false, 1)
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
				TopicName: "oktopus/v1/disconnect/" + clUser,
			}
			log.Println("new device:", clUser)
			err := server.Publish("oktopus/v1/devices/"+clUser, []byte(""), false, 1)
			if err != nil {
				log.Println("server publish error: ", err)
			}
		}

	}
}
