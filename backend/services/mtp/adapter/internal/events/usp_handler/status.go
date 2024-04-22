package usp_handler

import (
	"log"
	"strconv"

	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/db"
)

func (h *Handler) HandleDeviceStatus(device, subject string, data []byte, mtp string, ack func()) {
	defer ack()
	payload, err := strconv.Atoi(string(data))
	if err != nil {
		log.Printf("Status subject payload message error %q", err)
	}

	switch payload {
	case OFFLINE:
		h.deviceOffline(device, mtp)
	default:
		ignoreMsg(subject, "status", data)
	}
}

func (h *Handler) deviceOffline(device, mtp string) {
	log.Printf("Device %s is offline", device)

	mtpLayer := getMtp(mtp)

	err := h.db.UpdateStatus(device, db.Offline, mtpLayer)
	if err != nil {
		log.Fatal(err)
	}
}

func ignoreMsg(subject, ctx string, data []byte) {
	log.Printf("Unknown message of %s received, subject: %s, payload: %s. Ignored...", ctx, subject, string(data))
}
