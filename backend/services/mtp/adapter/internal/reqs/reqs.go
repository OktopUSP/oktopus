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
	"go.mongodb.org/mongo-driver/bson"
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
			respondMsg(msg.Respond, 200, deviceInfo)
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

	nc.Subscribe(local.ADAPTER_SUBJECT+"devices.count", func(msg *nats.Msg) {
		count, err := db.RetrieveDevicesCount(bson.M{})
		if err != nil {
			respondMsg(msg.Respond, 500, err.Error())
		}
		respondMsg(msg.Respond, 200, count)
	})

	nc.Subscribe(local.ADAPTER_SUBJECT+"devices.retrieve", func(msg *nats.Msg) {

		var filter bson.A

		err := json.Unmarshal(msg.Data, &filter)
		if err != nil {
			respondMsg(msg.Respond, 500, err.Error())
		}

		devicesList, err := db.RetrieveDevices(filter)
		if err != nil {
			respondMsg(msg.Respond, 500, err.Error())
		}
		respondMsg(msg.Respond, 200, devicesList)
	})

	nc.Subscribe(local.ADAPTER_SUBJECT+"devices.class", func(msg *nats.Msg) {
		productClassCount, err := db.RetrieveProductsClassInfo()
		if err != nil {
			respondMsg(msg.Respond, 500, err.Error())
		}
		respondMsg(msg.Respond, 200, productClassCount)
	})

	nc.Subscribe(local.ADAPTER_SUBJECT+"devices.vendors", func(msg *nats.Msg) {
		productClassCount, err := db.RetrieveVendorsInfo()
		if err != nil {
			respondMsg(msg.Respond, 500, err.Error())
		}
		respondMsg(msg.Respond, 200, productClassCount)
	})

	nc.Subscribe(local.ADAPTER_SUBJECT+"devices.status", func(msg *nats.Msg) {
		productClassCount, err := db.RetrieveStatusInfo()
		if err != nil {
			respondMsg(msg.Respond, 500, err.Error())
		}
		respondMsg(msg.Respond, 200, productClassCount)
	})

	nc.Subscribe(local.ADAPTER_SUBJECT+"*.device.alias", func(msg *nats.Msg) {
		subject := strings.Split(msg.Subject, ".")
		device := subject[len(subject)-3]

		err := db.SetDeviceAlias(device, string(msg.Data))
		if err != nil {
			if err == mongo.ErrNoDocuments {
				respondMsg(msg.Respond, 404, "Device not found")
			} else {
				respondMsg(msg.Respond, 500, err.Error())
			}
			return
		}
		respondMsg(msg.Respond, 200, "Alias updated")
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
