package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	usp_msg "github.com/leandrofars/oktopus/internal/usp_message"
	"github.com/leandrofars/oktopus/internal/utils"
	"google.golang.org/protobuf/proto"
)

type FwUpdate struct {
	Url string
}

func (a *Api) deviceFwUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sn := vars["sn"]
	a.deviceExists(sn, w)

	var payload FwUpdate

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Bad body, err: " + err.Error())
		return
	}

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
	case <-time.After(REQUEST_TIMEOUT):
		log.Printf("Request %s Timed Out", msg.Header.MsgId)
		w.WriteHeader(http.StatusGatewayTimeout)
		a.QMutex.Lock()
		delete(a.MsgQueue, msg.Header.MsgId)
		a.QMutex.Unlock()
		log.Println("requests queue:", a.MsgQueue)
		json.NewEncoder(w).Encode("Request Timed Out")
		return
	}

	partition := checkAvaiableFwPartition(getMsgAnswer.ReqPathResults)
	if partition == "" {
		log.Println("Error to get device available firmware partition, probably it has only one partition")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Server don't have the hability to update device with only one partition")
		return
		//TODO: update device with only one partition
	}

	log.Println("URL to download firmware:", payload.Url)

	receiver := usp_msg.Operate{
		Command:    "Device.DeviceInfo.FirmwareImage." + partition + "Download()",
		CommandKey: "Download()",
		SendResp:   true,
		InputArgs: map[string]string{
			"URL":          payload.Url,
			"AutoActivate": "true",
			//"Username": "",
			//"Password": "",
			"FileSize": "0", //TODO: send firmware length
			//"CheckSumAlgorithm": "",
			//"CheckSum":          "", //TODO: send firmware with checksum
		},
	}

	msg = utils.NewOperateMsg(receiver)
	a.uspCall(msg, sn, w)
}

// Check which fw image is activated
func checkAvaiableFwPartition(reqPathResult []*usp_msg.GetResp_RequestedPathResult) string {
	for _, x := range reqPathResult {
		partitionsNumber := len(x.ResolvedPathResults)
		if partitionsNumber > 1 {
			log.Printf("Device has %d firmware partitions", partitionsNumber)
			for _, y := range x.ResolvedPathResults {
				//TODO: verify if validation failed is trustable
				if y.ResultParams["Status"] == "Available" || y.ResultParams["Status"] == "ValidationFailed" {
					partition := y.ResolvedPath[len(y.ResolvedPath)-2:]
					log.Printf("Partition %s is avaiable", partition)
					return partition
				}
			}
		} else {
			return ""
		}
	}
	return ""
}
