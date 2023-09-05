package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	usp_msg "github.com/leandrofars/oktopus/internal/usp_message"
	"github.com/leandrofars/oktopus/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/proto"
)

func (a *Api) deviceFwUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sn := vars["sn"]
	a.deviceExists(sn, w)

	msg := utils.NewGetMsg(usp_msg.Get{
		ParamPaths: []string{"Device.DeviceInfo.FirmwareImage.*.Status"},
		MaxDepth:   1,
	})
	encodedMsg, err := proto.Marshal(&msg)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	record := utils.NewUspRecord(encodedMsg, sn)
	tr369Message, err := proto.Marshal(&record)
	if err != nil {
		log.Fatalln("Failed to encode tr369 record:", err)
	}

	a.QMutex.Lock()
	a.MsgQueue[msg.Header.MsgId] = make(chan usp_msg.Msg)
	a.QMutex.Unlock()
	log.Println("Sending Msg:", msg.Header.MsgId)
	a.Broker.Publish(tr369Message, "oktopus/v1/agent/"+sn, "oktopus/v1/api/"+sn, false)

	var getMsgAnswer *usp_msg.GetResp

	select {
	case msg := <-a.MsgQueue[msg.Header.MsgId]:
		log.Printf("Received Msg: %s", msg.Header.MsgId)
		a.QMutex.Lock()
		delete(a.MsgQueue, msg.Header.MsgId)
		a.QMutex.Unlock()
		log.Println("requests queue:", a.MsgQueue)
		getMsgAnswer = msg.Body.GetResponse().GetGetResp()
	case <-time.After(time.Second * 55):
		log.Printf("Request %s Timed Out", msg.Header.MsgId)
		w.WriteHeader(http.StatusGatewayTimeout)
		a.QMutex.Lock()
		delete(a.MsgQueue, msg.Header.MsgId)
		a.QMutex.Unlock()
		log.Println("requests queue:", a.MsgQueue)
		json.NewEncoder(w).Encode("Request Timed Out")
		return
	}

	// Check which fw image is activated
	partition := checkAvaiableFwPartition(getMsgAnswer.ReqPathResults)
	if partition < 0 {
		log.Println("Error to get device available firmware partition, probably it has only one partition")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Server don't have the hability to update device with only one partition")
		return
		//TODO: update device with only one partition
	}

	var receiver = usp_msg.Operate{
		Command:    "Device.DeviceInfo.FirmwareImage.1.Download()",
		CommandKey: "Download()",
		SendResp:   true,
		InputArgs: map[string]string{
			"URL":          "http://cronos.intelbras.com.br/download/PON/121AC/beta/121AC-2.3-230620-77753201df4f1e2c607a7236746c8491.tar", //TODO: use dynamic url
			"AutoActivate": "true",
			//"Username": "",
			//"Password": "",
			"FileSize": "0", //TODO: send firmware length
			//"CheckSumAlgorithm": "",
			//"CheckSum":          "",
		},
	}

	msg = utils.NewOperateMsg(receiver)
	encodedMsg, err = proto.Marshal(&msg)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	record = utils.NewUspRecord(encodedMsg, sn)
	tr369Message, err = proto.Marshal(&record)
	if err != nil {
		log.Fatalln("Failed to encode tr369 record:", err)
	}

	a.QMutex.Lock()
	a.MsgQueue[msg.Header.MsgId] = make(chan usp_msg.Msg)
	a.QMutex.Unlock()
	log.Println("Sending Msg:", msg.Header.MsgId)
	a.Broker.Publish(tr369Message, "oktopus/v1/agent/"+sn, "oktopus/v1/api/"+sn, false)

	select {
	case msg := <-a.MsgQueue[msg.Header.MsgId]:
		log.Printf("Received Msg: %s", msg.Header.MsgId)
		a.QMutex.Lock()
		delete(a.MsgQueue, msg.Header.MsgId)
		a.QMutex.Unlock()
		log.Println("requests queue:", a.MsgQueue)
		json.NewEncoder(w).Encode(msg.Body.GetResponse().GetSetResp())
		return
	case <-time.After(time.Second * 55):
		log.Printf("Request %s Timed Out", msg.Header.MsgId)
		w.WriteHeader(http.StatusGatewayTimeout)
		a.QMutex.Lock()
		delete(a.MsgQueue, msg.Header.MsgId)
		a.QMutex.Unlock()
		log.Println("requests queue:", a.MsgQueue)
		json.NewEncoder(w).Encode("Request Timed Out")
		return
	}
}

// Check which fw image is activated
func checkAvaiableFwPartition(reqPathResult []*usp_msg.GetResp_RequestedPathResult) int {
	for _, x := range reqPathResult {
		partitionsNumber := len(x.ResolvedPathResults)
		if partitionsNumber > 1 {
			log.Printf("Device has %d firmware partitions", partitionsNumber)
		}
		for i, y := range x.ResolvedPathResults {
			if y.ResultParams["Status"] == "Available" {
				log.Printf("Partition %d is avaiable", i)
				return i
			}
		}
	}
	return -1
}

func (a *Api) deviceGetSupportedParametersMsg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sn := vars["sn"]
	a.deviceExists(sn, w)

	var receiver usp_msg.GetSupportedDM

	err := json.NewDecoder(r.Body).Decode(&receiver)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := utils.NewGetSupportedParametersMsg(receiver)
	encodedMsg, err := proto.Marshal(&msg)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	record := utils.NewUspRecord(encodedMsg, sn)
	tr369Message, err := proto.Marshal(&record)
	if err != nil {
		log.Fatalln("Failed to encode tr369 record:", err)
	}

	//a.Broker.Request(tr369Message, usp_msg.Header_GET, "oktopus/v1/agent/"+sn, "oktopus/v1/get/"+sn)
	a.QMutex.Lock()
	a.MsgQueue[msg.Header.MsgId] = make(chan usp_msg.Msg)
	a.QMutex.Unlock()
	log.Println("Sending Msg:", msg.Header.MsgId)
	a.Broker.Publish(tr369Message, "oktopus/v1/agent/"+sn, "oktopus/v1/api/"+sn, false)

	select {
	case msg := <-a.MsgQueue[msg.Header.MsgId]:
		log.Printf("Received Msg: %s", msg.Header.MsgId)
		a.QMutex.Lock()
		delete(a.MsgQueue, msg.Header.MsgId)
		a.QMutex.Unlock()
		log.Println("requests queue:", a.MsgQueue)
		json.NewEncoder(w).Encode(msg.Body.GetResponse().GetGetSupportedDmResp())
		return
	case <-time.After(time.Second * 55):
		log.Printf("Request %s Timed Out", msg.Header.MsgId)
		w.WriteHeader(http.StatusGatewayTimeout)
		a.QMutex.Lock()
		delete(a.MsgQueue, msg.Header.MsgId)
		a.QMutex.Unlock()
		log.Println("requests queue:", a.MsgQueue)
		json.NewEncoder(w).Encode("Request Timed Out")
		return
	}
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
	a.deviceExists(sn, w)

	var receiver usp_msg.Add

	err := json.NewDecoder(r.Body).Decode(&receiver)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := utils.NewCreateMsg(receiver)
	encodedMsg, err := proto.Marshal(&msg)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	record := utils.NewUspRecord(encodedMsg, sn)
	tr369Message, err := proto.Marshal(&record)
	if err != nil {
		log.Fatalln("Failed to encode tr369 record:", err)
	}

	//a.Broker.Request(tr369Message, usp_msg.Header_GET, "oktopus/v1/agent/"+sn, "oktopus/v1/get/"+sn)
	a.QMutex.Lock()
	a.MsgQueue[msg.Header.MsgId] = make(chan usp_msg.Msg)
	a.QMutex.Unlock()
	log.Println("Sending Msg:", msg.Header.MsgId)
	a.Broker.Publish(tr369Message, "oktopus/v1/agent/"+sn, "oktopus/v1/api/"+sn, false)

	select {
	case msg := <-a.MsgQueue[msg.Header.MsgId]:
		log.Printf("Received Msg: %s", msg.Header.MsgId)
		a.QMutex.Lock()
		delete(a.MsgQueue, msg.Header.MsgId)
		a.QMutex.Unlock()
		log.Println("requests queue:", a.MsgQueue)
		json.NewEncoder(w).Encode(msg.Body.GetResponse().GetAddResp())
		return
	case <-time.After(time.Second * 55):
		log.Printf("Request %s Timed Out", msg.Header.MsgId)
		w.WriteHeader(http.StatusGatewayTimeout)
		a.QMutex.Lock()
		delete(a.MsgQueue, msg.Header.MsgId)
		a.QMutex.Unlock()
		log.Println("requests queue:", a.MsgQueue)
		json.NewEncoder(w).Encode("Request Timed Out")
		return
	}
}

func (a *Api) deviceGetMsg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sn := vars["sn"]

	a.deviceExists(sn, w)

	var receiver usp_msg.Get

	err := json.NewDecoder(r.Body).Decode(&receiver)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := utils.NewGetMsg(receiver)
	encodedMsg, err := proto.Marshal(&msg)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	record := utils.NewUspRecord(encodedMsg, sn)
	tr369Message, err := proto.Marshal(&record)
	if err != nil {
		log.Fatalln("Failed to encode tr369 record:", err)
	}

	a.QMutex.Lock()
	a.MsgQueue[msg.Header.MsgId] = make(chan usp_msg.Msg)
	a.QMutex.Unlock()

	log.Println("Sending Msg:", msg.Header.MsgId)
	a.Broker.Publish(tr369Message, "oktopus/v1/agent/"+sn, "oktopus/v1/api/"+sn, false)

	select {
	case msg := <-a.MsgQueue[msg.Header.MsgId]:
		log.Printf("Received Msg: %s", msg.Header.MsgId)
		a.QMutex.Lock()
		delete(a.MsgQueue, msg.Header.MsgId)
		a.QMutex.Unlock()
		log.Println("requests queue:", a.MsgQueue)
		json.NewEncoder(w).Encode(msg.Body.GetResponse().GetGetResp())
		return
	case <-time.After(time.Second * 55):
		log.Printf("Request %s Timed Out", msg.Header.MsgId)
		w.WriteHeader(http.StatusGatewayTimeout)
		a.QMutex.Lock()
		delete(a.MsgQueue, msg.Header.MsgId)
		a.QMutex.Unlock()
		log.Println("requests queue:", a.MsgQueue)
		json.NewEncoder(w).Encode("Request Timed Out")
		return
	}
}

func (a *Api) deviceDeleteMsg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sn := vars["sn"]
	a.deviceExists(sn, w)

	var receiver usp_msg.Delete

	err := json.NewDecoder(r.Body).Decode(&receiver)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := utils.NewDelMsg(receiver)
	encodedMsg, err := proto.Marshal(&msg)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	record := utils.NewUspRecord(encodedMsg, sn)
	tr369Message, err := proto.Marshal(&record)
	if err != nil {
		log.Fatalln("Failed to encode tr369 record:", err)
	}

	//a.Broker.Request(tr369Message, usp_msg.Header_GET, "oktopus/v1/agent/"+sn, "oktopus/v1/get/"+sn)
	a.QMutex.Lock()
	a.MsgQueue[msg.Header.MsgId] = make(chan usp_msg.Msg)
	a.QMutex.Unlock()
	log.Println("Sending Msg:", msg.Header.MsgId)
	a.Broker.Publish(tr369Message, "oktopus/v1/agent/"+sn, "oktopus/v1/api/"+sn, false)

	select {
	case msg := <-a.MsgQueue[msg.Header.MsgId]:
		log.Printf("Received Msg: %s", msg.Header.MsgId)
		a.QMutex.Lock()
		delete(a.MsgQueue, msg.Header.MsgId)
		a.QMutex.Unlock()
		log.Println("requests queue:", a.MsgQueue)
		json.NewEncoder(w).Encode(msg.Body.GetResponse().GetDeleteResp())
		return
	case <-time.After(time.Second * 55):
		log.Printf("Request %s Timed Out", msg.Header.MsgId)
		w.WriteHeader(http.StatusGatewayTimeout)
		a.QMutex.Lock()
		delete(a.MsgQueue, msg.Header.MsgId)
		a.QMutex.Unlock()
		log.Println("requests queue:", a.MsgQueue)
		json.NewEncoder(w).Encode("Request Timed Out")
		return
	}
}

func (a *Api) deviceUpdateMsg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sn := vars["sn"]
	a.deviceExists(sn, w)

	var receiver usp_msg.Set

	err := json.NewDecoder(r.Body).Decode(&receiver)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := utils.NewSetMsg(receiver)
	encodedMsg, err := proto.Marshal(&msg)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	record := utils.NewUspRecord(encodedMsg, sn)
	tr369Message, err := proto.Marshal(&record)
	if err != nil {
		log.Fatalln("Failed to encode tr369 record:", err)
	}

	//a.Broker.Request(tr369Message, usp_msg.Header_GET, "oktopus/v1/agent/"+sn, "oktopus/v1/get/"+sn)
	a.QMutex.Lock()
	a.MsgQueue[msg.Header.MsgId] = make(chan usp_msg.Msg)
	a.QMutex.Unlock()
	log.Println("Sending Msg:", msg.Header.MsgId)
	a.Broker.Publish(tr369Message, "oktopus/v1/agent/"+sn, "oktopus/v1/api/"+sn, false)

	select {
	case msg := <-a.MsgQueue[msg.Header.MsgId]:
		log.Printf("Received Msg: %s", msg.Header.MsgId)
		a.QMutex.Lock()
		delete(a.MsgQueue, msg.Header.MsgId)
		a.QMutex.Unlock()
		log.Println("requests queue:", a.MsgQueue)
		json.NewEncoder(w).Encode(msg.Body.GetResponse().GetSetResp())
		return
	case <-time.After(time.Second * 55):
		log.Printf("Request %s Timed Out", msg.Header.MsgId)
		w.WriteHeader(http.StatusGatewayTimeout)
		a.QMutex.Lock()
		delete(a.MsgQueue, msg.Header.MsgId)
		a.QMutex.Unlock()
		log.Println("requests queue:", a.MsgQueue)
		json.NewEncoder(w).Encode("Request Timed Out")
		return
	}
}

func (a *Api) deviceExists(sn string, w http.ResponseWriter) {
	_, err := a.Db.RetrieveDevice(sn)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("No device with serial number " + sn + " was found")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *Api) deviceGetParameterInstances(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sn := vars["sn"]
	a.deviceExists(sn, w)

	var receiver usp_msg.GetInstances

	err := json.NewDecoder(r.Body).Decode(&receiver)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := utils.NewGetParametersInstancesMsg(receiver)
	encodedMsg, err := proto.Marshal(&msg)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	record := utils.NewUspRecord(encodedMsg, sn)
	tr369Message, err := proto.Marshal(&record)
	if err != nil {
		log.Fatalln("Failed to encode tr369 record:", err)
	}

	a.QMutex.Lock()
	a.MsgQueue[msg.Header.MsgId] = make(chan usp_msg.Msg)
	a.QMutex.Unlock()
	log.Println("Sending Msg:", msg.Header.MsgId)
	a.Broker.Publish(tr369Message, "oktopus/v1/agent/"+sn, "oktopus/v1/api/"+sn, false)

	select {
	case msg := <-a.MsgQueue[msg.Header.MsgId]:
		log.Printf("Received Msg: %s", msg.Header.MsgId)
		a.QMutex.Lock()
		delete(a.MsgQueue, msg.Header.MsgId)
		a.QMutex.Unlock()
		log.Println("requests queue:", a.MsgQueue)
		json.NewEncoder(w).Encode(msg.Body.GetResponse().GetGetInstancesResp())
		return
	case <-time.After(time.Second * 55):
		log.Printf("Request %s Timed Out", msg.Header.MsgId)
		w.WriteHeader(http.StatusGatewayTimeout)
		a.QMutex.Lock()
		delete(a.MsgQueue, msg.Header.MsgId)
		a.QMutex.Unlock()
		log.Println("requests queue:", a.MsgQueue)
		json.NewEncoder(w).Encode("Request Timed Out")
		return
	}
}
