// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 mochi-co
// SPDX-FileContributor: mochi-co

package main

import (
	"bytes"
	"flag"
	"github.com/mochi-co/mqtt/v2/packets"
	"github.com/rs/zerolog"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	rv8 "github.com/go-redis/redis/v8"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/hooks/auth"
	"github.com/mochi-co/mqtt/v2/hooks/storage/redis"
	"github.com/mochi-co/mqtt/v2/listeners"
)

var server = mqtt.New(&mqtt.Options{
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

func main() {
	tcpAddr := flag.String("tcp", ":1883", "network address for TCP listener")
	redisAddr := flag.String("redis", "172.17.0.2:6379", "host address of redis db")
	wsAddr := flag.String("ws", "", "network address for Websocket listener")
	infoAddr := flag.String("info", ":8080", "network address for web info dashboard listener")
	path := flag.String("path", "", "path to data auth file")
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	l := server.Log.Level(zerolog.DebugLevel)
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

	if *tcpAddr != "" {
		tcp := listeners.NewTCP("t1", *tcpAddr, nil)
		err := server.AddListener(tcp)
		if err != nil {
			log.Fatal(err)
		}
	}

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

	err = server.AddHook(new(redis.Hook), &redis.Options{
		Options: &rv8.Options{
			Addr:     *redisAddr, // default redis address
			Password: "",         // your password
			DB:       0,          // your redis db
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := server.Serve()
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
	}, []byte{b})
}

func (h *MyHook) Init(config any) error {
	h.Log.Info().Msg("initialised")
	return nil
}

func (h *MyHook) Red(config any) error {
	h.Log.Info().Msg("initialised")
	return nil
}

func (h *MyHook) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
	var clUser string
	if len(cl.Properties.Props.User) > 0 {
		clUser = cl.Properties.Props.User[0].Val
	}

	if clUser != "" {
		err := server.Publish("oktopus/disconnect", []byte(clUser), false, 1)
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
			log.Println("new device:", clUser)
			err := server.Publish("oktopus/devices", []byte(clUser), false, 1)
			if err != nil {
				log.Println("server publish error: ", err)
			}
		}

	}
}
