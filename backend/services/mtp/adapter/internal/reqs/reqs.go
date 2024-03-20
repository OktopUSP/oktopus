/*
Provide answers to nats request-reply messages, executing queries to the database
*/
package reqs

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/db"
	local "github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/nats"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
)

type msgAnswer struct {
	Code int
	Msg  any
}

func StartRequestsListener(ctx context.Context, nc *nats.Conn, db db.Database) {
	log.Println("Listening for nats requests")
	nc.Subscribe(local.ADAPTER_SUBJECT+"*.device", func(msg *nats.Msg) {
		subject := strings.Split(msg.Subject, ".")
		device := subject[len(subject)-2]

		deviceInfo, err := db.RetrieveDevice(device)
		if deviceInfo.SN != "" {
			body, _ := json.Marshal(deviceInfo)
			respondMsg(msg.Respond, 200, body)
		} else {
			if err != nil {
				if err == mongo.ErrNoDocuments {
					respondMsg(msg.Respond, 404, "Device not found")
				} else {
					respondMsg(msg.Respond, 500, err.Error())
				}
			}
		}
	})
}

func respondMsg(respond func(data []byte) error, code int, msgData any) {

	msg, err := json.Marshal(msgAnswer{
		Code: code,
		Msg:  msgData,
	})
	if err != nil {
		log.Printf("Failed to marshal message: %q", err)
		respond([]byte(err.Error()))
		return
	}

	respond([]byte(msg))
}
