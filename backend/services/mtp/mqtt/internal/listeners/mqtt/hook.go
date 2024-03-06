package mqtt

import (
	"bytes"
	"log"
	"strings"

	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/packets"
)

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
		err := server.Publish("oktopus/usp/v1/status/"+clUser, []byte("1"), false, 1)
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
				Payload:   []byte("1"),
				Retain:    false,
			}
			log.Println("new device:", clUser)
			err := server.Publish("oktopus/usp/v1/status/"+clUser, []byte("0"), false, 1)
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
