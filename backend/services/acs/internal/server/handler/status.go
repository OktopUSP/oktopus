package handler

import (
	"log"
	"time"
)

// TODO: make these consts dynamic via config
const (
	CHECK_STATUS_INTERVAL = 5 * time.Second
	KEEP_ALIVE_INTERVAL   = 10 * time.Second
)

func (h *Handler) handleCpeStatus(cpe string) {
	for {
		if time.Since(h.cpes[cpe].LastConnection) > KEEP_ALIVE_INTERVAL {
			delete(h.cpes, cpe)
			break
		}
		time.Sleep(CHECK_STATUS_INTERVAL)
	}
	log.Println("CPE", cpe, "is offline")
	h.pub("cwmp.v1."+cpe+".status", []byte("0"))
}
