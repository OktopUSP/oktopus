package mqtt

import (
	"broker/internal/config"
	"broker/internal/nats"
	"bytes"
	"context"
	"log"
	"strings"

	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/hooks/auth"
	"github.com/mochi-co/mqtt/v2/packets"
	"github.com/nats-io/nats.go/jetstream"
)

type MyHook struct {
	mqtt.HookBase
}

type NatsAuthHook struct {
	mqtt.HookBase
	kv jetstream.KeyValue
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
		err := server.Publish("oktopus/usp/v1/status/"+clUser, []byte("0"), false, 1)
		if err != nil {
			log.Println("server publish error: ", err)
		}
	}
}

func (h *MyHook) OnSubscribed(cl *mqtt.Client, pk packets.Packet, reasonCodes []byte) {
	//Verifies if it's a device who is subscribed
	if strings.Contains(pk.Filters[0].Filter, "oktopus/usp/v1/agent") {
		var clUser string

		if len(cl.Properties.Props.User) > 0 {
			clUser = cl.Properties.Props.User[0].Val
		}

		if clUser != "" {
			cl.Properties.Will = mqtt.Will{
				Qos:       1,
				TopicName: "oktopus/usp/v1/status/" + clUser,
				Payload:   []byte("0"),
				Retain:    false,
			}
			log.Println("new device:", clUser)
			err := server.Publish("oktopus/usp/v1/status/"+clUser, []byte("1"), false, 1)
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
		pk.Properties.User = []packets.UserProperty{{Key: "subscribe-topic", Val: "oktopus/usp/v1/agent/" + clUser}}
	}

	return pk
}

func (h *NatsAuthHook) ID() string {
	return "device-auth"
}

func (h *NatsAuthHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnectAuthenticate,
		mqtt.OnACLCheck,
	}, []byte{b})
}

func (h *NatsAuthHook) Init(c any) error {

	_, _, kv := nats.StartNatsClient(c.(config.Nats))
	h.kv = kv

	h.Log.Info().Msg("initialised device auth nats hook")
	return nil
}

func (h *NatsAuthHook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {

	username := string(pk.Connect.Username)

	entry, err := h.kv.Get(context.TODO(), username)
	if err != nil {
		if err == jetstream.ErrKeyNotFound {
			log.Println("user access not found, blocked user:", username)
			return false
		}
		log.Println("error getting key value: ", err)
		return false
	}

	if bytes.Equal(entry.Value(), pk.Connect.Password) {
		return true
	}

	return false
}

func (h *NatsAuthHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {

	username := string(cl.Properties.Username)

	_, err := h.kv.Get(context.TODO(), username)
	if err != nil {
		if err == jetstream.ErrKeyNotFound {
			log.Println("user access not found, blocked user:", username)
			return false
		}
		log.Println("error getting key value: ", err)
		return false
	}

	if username == "oktopusController" {
		return true
	}

	if !write {
		deviceAllowedTopics := []string{
			"oktopus/usp/v1/agent/" + username,
		}
		for _, allowedTopic := range deviceAllowedTopics {
			_, ok := auth.MatchTopic(allowedTopic, topic)
			if ok {
				return true
			}
		}
		return false
	}

	if write {
		deviceAllowedTopics := []string{
			"oktopus/usp/v1/controller",
			"oktopus/usp/v1/api/" + username,
		}
		for _, allowedTopic := range deviceAllowedTopics {
			_, ok := auth.MatchTopic(allowedTopic, topic)
			if ok {
				return true
			}
		}
	}

	return false
}
