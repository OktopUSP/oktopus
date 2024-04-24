package handler

import (
	"log"
	"time"
)

func (h *Handler) handleCpeStatus(cpe string) {
	for {
		if time.Since(h.Cpes[cpe].LastConnection) > h.acsConfig.KeepAliveInterval {
			delete(h.Cpes, cpe)
			break
		}
		time.Sleep(h.acsConfig.KeepAliveInterval)
	}
	log.Println("CPE", cpe, "is offline")
	h.pub("cwmp.v1."+cpe+".status", []byte("0"))
}
