package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/leandrofars/oktopus/internal/bridge"
	local "github.com/leandrofars/oktopus/internal/nats"
	"github.com/leandrofars/oktopus/internal/usp/usp_msg"
	"github.com/leandrofars/oktopus/internal/usp/usp_record"
	"github.com/leandrofars/oktopus/internal/usp/usp_utils"
	"github.com/leandrofars/oktopus/internal/utils"
	"google.golang.org/protobuf/proto"
)

type fwUpdate struct {
	Url string
}

func (a *Api) deviceFwUpdate(w http.ResponseWriter, r *http.Request) {
	sn := getSerialNumberFromRequest(r)
	mtp, err := getMtpFromRequest(r, w)
	if err != nil {
		return
	}

	if mtp == "" {
		var ok bool
		mtp, ok = deviceStateOK(w, a.nc, sn)
		if !ok {
			return
		}
	}

	var payload fwUpdate

	utils.MarshallDecoder(&payload, r.Body)

	msg := usp_utils.NewGetMsg(usp_msg.Get{
		ParamPaths: []string{"Device.DeviceInfo.FirmwareImage.*.Status"},
		MaxDepth:   1,
	})

	protoMsg, err := proto.Marshal(&msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.Marshall(err.Error()))
		return
	}

	record := usp_utils.NewUspRecord(protoMsg, sn)
	protoRecord, err := proto.Marshal(&record)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.Marshall(err.Error()))
		return
	}

	data, err := bridge.NatsUspInteraction(
		local.DEVICE_SUBJECT_PREFIX+sn+".api",
		mtp+"-adapter.usp.v1."+sn+".api",
		protoRecord,
		w,
		a.nc,
	)
	if err != nil {
		return
	}

	var receivedRecord usp_record.Record
	err = proto.Unmarshal(data, &receivedRecord)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.Marshall(err.Error()))
		return
	}
	var receivedMsg usp_msg.Msg
	err = proto.Unmarshal(receivedRecord.GetNoSessionContext().Payload, &receivedMsg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.Marshall(err.Error()))
		return
	}

	getMsgAnswer := receivedMsg.Body.GetResponse().GetGetResp()

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

	msg = usp_utils.NewOperateMsg(receiver)
	err = sendUspMsg(msg, sn, w, a.nc, mtp)
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
