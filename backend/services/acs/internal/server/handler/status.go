package handler

import (
	"log"
	"time"
)

func (h *Handler) HandleCpeStatus() {
	for {
		for cpe := range h.Cpes {
			if cpe == "" {
				continue
			}
			log.Println("Checking CPE " + cpe + " status")
			if time.Since(h.Cpes[cpe].LastConnection) > h.acsConfig.KeepAliveInterval {
				log.Printf("LastConnection: %s, KeepAliveInterval: %s", h.Cpes[cpe].LastConnection, h.acsConfig.KeepAliveInterval)
				log.Println("CPE", cpe, "is offline")
				h.pub("cwmp.v1."+cpe+".status", []byte("0"))
				delete(h.Cpes, cpe)
				break
			}
		}
		time.Sleep(10 * time.Second)
	}
}
