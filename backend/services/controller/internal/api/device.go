package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leandrofars/oktopus/internal/db"
	usp_msg "github.com/leandrofars/oktopus/internal/usp_message"
	"github.com/leandrofars/oktopus/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func (a *Api) deviceGetSupportedParametersMsg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sn := vars["sn"]
	device := a.deviceExists(sn, w)

	var receiver usp_msg.GetSupportedDM

	err := json.NewDecoder(r.Body).Decode(&receiver)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := utils.NewGetSupportedParametersMsg(receiver)
	a.uspCall(msg, sn, w, device)
}

func (a *Api) retrieveDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := a.Db.RetrieveDevices()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(devices)
	if err != nil {
		log.Println(err)
	}

	return
}

func (a *Api) deviceCreateMsg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sn := vars["sn"]
	device := a.deviceExists(sn, w)

	var receiver usp_msg.Add

	err := json.NewDecoder(r.Body).Decode(&receiver)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := utils.NewCreateMsg(receiver)
	a.uspCall(msg, sn, w, device)
}

func (a *Api) deviceGetMsg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sn := vars["sn"]

	device := a.deviceExists(sn, w)

	var receiver usp_msg.Get

	err := json.NewDecoder(r.Body).Decode(&receiver)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := utils.NewGetMsg(receiver)
	a.uspCall(msg, sn, w, device)
}

func (a *Api) deviceOperateMsg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sn := vars["sn"]

	device := a.deviceExists(sn, w)

	var receiver usp_msg.Operate

	err := json.NewDecoder(r.Body).Decode(&receiver)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := utils.NewOperateMsg(receiver)
	a.uspCall(msg, sn, w, device)
}

func (a *Api) deviceDeleteMsg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sn := vars["sn"]
	device := a.deviceExists(sn, w)

	var receiver usp_msg.Delete

	err := json.NewDecoder(r.Body).Decode(&receiver)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := utils.NewDelMsg(receiver)
	a.uspCall(msg, sn, w, device)

	//a.Broker.Request(tr369Message, usp_msg.Header_GET, "oktopus/v1/agent/"+sn, "oktopus/v1/get/"+sn)

}

func (a *Api) deviceUpdateMsg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sn := vars["sn"]
	device := a.deviceExists(sn, w)

	var receiver usp_msg.Set

	err := json.NewDecoder(r.Body).Decode(&receiver)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := utils.NewSetMsg(receiver)
	a.uspCall(msg, sn, w, device)
}

func (a *Api) deviceExists(sn string, w http.ResponseWriter) db.Device {
	device, err := a.Db.RetrieveDevice(sn)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("No device with serial number " + sn + " was found")
		}
		w.WriteHeader(http.StatusInternalServerError)
		return device
	}
	return device
}

func (a *Api) deviceGetParameterInstances(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sn := vars["sn"]
	device := a.deviceExists(sn, w)

	var receiver usp_msg.GetInstances

	err := json.NewDecoder(r.Body).Decode(&receiver)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := utils.NewGetParametersInstancesMsg(receiver)
	a.uspCall(msg, sn, w, device)
}
