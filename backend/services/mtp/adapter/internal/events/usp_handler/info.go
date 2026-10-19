package usp_handler

import (
	"encoding/json"
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
	if deviceExists, _ := h.db.DeviceExists(deviceInfo.SN); !deviceExists {
		fmtDeviceInfo, _ := json.Marshal(deviceInfo)
		h.nc.Publish("device.v1.new", fmtDeviceInfo)
	}
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
        log.Println("Error unmarshaling USP Record:", err)
        return db.Device{}
    }

    noSessionCtx := record.GetNoSessionContext()
    if noSessionCtx == nil || noSessionCtx.Payload == nil {
        log.Printf("Record has no NoSessionContext payload (record type: %T) - registering device with basic info", record.RecordType)
        var device db.Device
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

    err = proto.Unmarshal(noSessionCtx.Payload, &message)
    if err != nil {
        log.Println("Error unmarshaling USP Message:", err)
        return db.Device{}
    }

    var device db.Device

    respBody, isResponse := message.Body.MsgBody.(*usp_msg.Body_Response)

    if !isResponse {
        log.Printf("Ignored message for DeviceInfo: Expected Body_Response but got %T", message.Body.MsgBody)
        return device
    }

    msg := respBody.Response.GetGetResp()
    if msg == nil {
        log.Println("Ignored message: Response does not contain GetResp")
        return device
    }

    if len(msg.ReqPathResults) < 5 {
        log.Printf("Error: Expected 5 params in GetResp, got %d", len(msg.ReqPathResults))
        return device
    }

    if len(msg.ReqPathResults[0].ResolvedPathResults) > 0 {
         device.Vendor = msg.ReqPathResults[0].ResolvedPathResults[0].ResultParams["Manufacturer"]
    }
    if len(msg.ReqPathResults[1].ResolvedPathResults) > 0 {
         device.Model = msg.ReqPathResults[1].ResolvedPathResults[0].ResultParams["ModelName"]
    }
    if len(msg.ReqPathResults[2].ResolvedPathResults) > 0 {
         device.Version = msg.ReqPathResults[2].ResolvedPathResults[0].ResultParams["SoftwareVersion"]
    }
    if len(msg.ReqPathResults[4].ResolvedPathResults) > 0 {
         device.ProductClass = msg.ReqPathResults[4].ResolvedPathResults[0].ResultParams["ProductClass"]
    }

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
