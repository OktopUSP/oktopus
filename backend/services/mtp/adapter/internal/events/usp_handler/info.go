package usp_handler

import (
	"log"

	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/db"
	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/nats"
	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/usp/usp_msg"
	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/usp/usp_record"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) HandleDeviceInfo(device, subject string, data []byte, mtp string, ack func()) {
	defer ack()
	log.Printf("Device %s info, mtp: %s", device, mtp)
	deviceInfo := parseDeviceInfoMsg(device, subject, data, getMtp(mtp))
	err := h.db.CreateDevice(deviceInfo)
	if err != nil {
		log.Printf("Failed to create device: %v", err)
	}
}

func getMtp(mtp string) db.MTP {
	switch mtp {
	case nats.MQTT_STREAM_NAME:
		return db.MQTT
	case nats.WS_STREAM_NAME:
		return db.WEBSOCKETS
	case nats.STOMP_STREAM_NAME:
		return db.STOMP
	default:
		return db.UNDEFINED
	}
}

func parseDeviceInfoMsg(sn, subject string, data []byte, mtp db.MTP) db.Device {
	var record usp_record.Record
	var message usp_msg.Msg

	err := proto.Unmarshal(data, &record)
	if err != nil {
		log.Fatal(err)
	}
	err = proto.Unmarshal(record.GetNoSessionContext().Payload, &message)
	if err != nil {
		log.Fatal(err)
	}

	var device db.Device
	msg := message.Body.MsgBody.(*usp_msg.Body_Response).Response.GetGetResp()

	device.Vendor = msg.ReqPathResults[0].ResolvedPathResults[0].ResultParams["Manufacturer"]
	device.Model = msg.ReqPathResults[1].ResolvedPathResults[0].ResultParams["ModelName"]
	device.Version = msg.ReqPathResults[2].ResolvedPathResults[0].ResultParams["SoftwareVersion"]
	device.ProductClass = msg.ReqPathResults[4].ResolvedPathResults[0].ResultParams["ProductClass"]
	device.SN = sn
	switch db.MTP(mtp) {
	case db.MQTT:
		device.Mqtt = db.Online
	case db.WEBSOCKETS:
		device.Websockets = db.Online
	case db.STOMP:
		device.Stomp = db.Online
	}

	device.Status = db.Online

	return device
}
