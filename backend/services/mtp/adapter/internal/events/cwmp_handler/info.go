package cwmp_handler

import (
	"encoding/xml"
	"log"

	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/cwmp"
	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/db"
)

func (h *Handler) HandleDeviceInfo(device string, data []byte, ack func()) {
	defer ack()
	log.Printf("Device %s info", device)
	deviceInfo := parseDeviceInfoMsg(data)
	err := h.db.CreateDevice(deviceInfo)
	if err != nil {
		log.Printf("Failed to create device: %v", err)
	}
}

func parseDeviceInfoMsg(data []byte) db.Device {

	var inform cwmp.CWMPInform
	err := xml.Unmarshal(data, &inform)
	if err != nil {
		log.Println("Error unmarshalling xml:", err)
	}

	var device db.Device

	device.Vendor = inform.DeviceId.Manufacturer
	device.Model = ""
	device.Version = inform.GetSoftwareVersion()
	device.ProductClass = inform.DeviceId.ProductClass
	device.SN = inform.DeviceId.SerialNumber
	device.Cwmp = db.Online
	device.Status = db.Online

	return device
}
