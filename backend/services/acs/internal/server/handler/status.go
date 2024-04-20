package handler

import (
	"log"
	"time"
)

// TODO: make these consts dynamic via config
const (
	CHECK_STATUS_INTERVAL = 10 * time.Second
	KEEP_ALIVE_INTERVAL   = 600 * time.Second
)

func (h *Handler) handleCpeStatus(cpe string) {
	for {
		if time.Since(h.Cpes[cpe].LastConnection) > KEEP_ALIVE_INTERVAL {
			delete(h.Cpes, cpe)
			break
		}
		time.Sleep(CHECK_STATUS_INTERVAL)
	}
	log.Println("CPE", cpe, "is offline")
	h.pub("cwmp.v1."+cpe+".status", []byte("0"))
}
