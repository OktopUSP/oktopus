package cwmp_handler

import (
	"log"
	"strconv"

	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/db"
)

func (h *Handler) HandleDeviceStatus(device, subject string, data []byte, ack func()) {
	defer ack()
	payload, err := strconv.Atoi(string(data))
	if err != nil {
		log.Printf("Status subject payload message error %q", err)
	}

	switch payload {
	case OFFLINE:
		h.deviceOffline(device)
	default:
		ignoreMsg(subject, "status", data)
	}
}

func (h *Handler) deviceOffline(device string) {
	log.Printf("Device %s is offline", device)

	err := h.db.UpdateStatus(device, db.Offline, db.CWMP)
	if err != nil {
		log.Fatal(err)
	}
}

func ignoreMsg(subject, ctx string, data []byte) {
	log.Printf("Unknown message of %s received, subject: %s, payload: %s. Ignored...", ctx, subject, string(data))
}
