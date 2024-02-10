package handler

import (
	"log"

	"github.com/leandrofars/oktopus/internal/db"
	usp_msg "github.com/leandrofars/oktopus/internal/usp_message"
	"github.com/leandrofars/oktopus/internal/usp_record"
	"github.com/leandrofars/oktopus/internal/utils"
	"google.golang.org/protobuf/proto"
)

func HandleNewDevice(deviceMac string) []byte {

	payload := utils.NewGetMsg(usp_msg.Get{
		ParamPaths: []string{
			"Device.DeviceInfo.Manufacturer",
			"Device.DeviceInfo.ModelName",
			"Device.DeviceInfo.SoftwareVersion",
			"Device.DeviceInfo.SerialNumber",
			"Device.DeviceInfo.ProductClass",
		},
		MaxDepth: 1,
	})

	teste, _ := proto.Marshal(&payload)
	record := utils.NewUspRecord(teste, deviceMac)

	tr369Message, err := proto.Marshal(&record)
	if err != nil {
		log.Fatalln("Failed to encode tr369 record:", err)
	}

	return tr369Message
}

func HandleNewDevicesResponse(p []byte, sn string, mtp db.MTP) db.Device {
	var record usp_record.Record
	var message usp_msg.Msg

	err := proto.Unmarshal(p, &record)
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
