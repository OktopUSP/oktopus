package handler

import (
	"log"
	"strconv"

	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/db"
	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/usp"
	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/usp/usp_msg"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) HandleDeviceStatus(device, subject string, data []byte) {
	payload, err := strconv.Atoi(string(data))
	if err != nil {
		log.Printf("Status subject payload message error %q", err)
	}

	switch payload {
	case ONLINE:
		h.deviceOnline(device)
	case OFFLINE:
		h.deviceOffline(device)
	default:
		ignoreMsg(subject, "status", data)
	}
}

func (h *Handler) deviceOnline(device string) {

	log.Printf("Device %s is online", device)

	msg := usp.NewGetMsg(usp_msg.Get{
		ParamPaths: []string{
			"Device.DeviceInfo.Manufacturer",
			"Device.DeviceInfo.ModelName",
			"Device.DeviceInfo.SoftwareVersion",
			"Device.DeviceInfo.SerialNumber",
			"Device.DeviceInfo.ProductClass",
		},
		MaxDepth: 1,
	})

	payload, _ := proto.Marshal(&msg)
	record := usp.NewUspRecord(payload, device, h.cid)

	tr369Message, err := proto.Marshal(&record)
	if err != nil {
		log.Fatalln("Failed to encode tr369 record:", err)
	}

	err = h.nc.Publish(NATS_SUBJ_PREFIX+device+".info", tr369Message)
	if err != nil {
		log.Printf("Failed to publish online device message: %v", err)
	}
}

func (h *Handler) deviceOffline(device string) {
	log.Printf("Device %s is offline", device)

	err := h.db.UpdateStatus(device, db.Offline, db.MQTT)
	if err != nil {
		log.Fatal(err)
	}
}

func ignoreMsg(subject, ctx string, data []byte) {
	log.Printf("Unknown message of %s received, subject: %s, payload: %s. Ignored...", ctx, subject, string(data))
}
