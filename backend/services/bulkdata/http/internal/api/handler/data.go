package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handler) Data(w http.ResponseWriter, r *http.Request) {

	oui := r.URL.Query().Get("oui")
	pc := r.URL.Query().Get("pc")
	sn := r.URL.Query().Get("sn")
	eid := r.URL.Query().Get("eid")

	log.Println("oui: ", oui)
	log.Println("pc: ", pc)
	log.Println("sn: ", sn)
	log.Println("eid: ", eid)

	var body map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("Body: ", body)
}
