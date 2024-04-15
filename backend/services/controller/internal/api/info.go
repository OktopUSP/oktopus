package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/leandrofars/oktopus/internal/bridge"
	"github.com/leandrofars/oktopus/internal/entity"
	local "github.com/leandrofars/oktopus/internal/nats"
	"github.com/leandrofars/oktopus/internal/utils"
)

type StatusCount struct {
	Online  int
	Offline int
}

type GeneralInfo struct {
	MqttRtt           string
	WebsocketsRtt     string
	StompRtt          string
	ProductClassCount []entity.ProductClassCount
	StatusCount       StatusCount
	VendorsCount      []entity.VendorsCount
}

func (a *Api) generalInfo(w http.ResponseWriter, r *http.Request) {

	var result GeneralInfo

	productclasscount, err := bridge.NatsReq[[]entity.ProductClassCount](
		local.NATS_ADAPTER_SUBJECT+"devices.class",
		[]byte(""),
		w,
		a.nc,
	)
	if err != nil {
		return
	}

	vendorcount, err := bridge.NatsReq[[]entity.VendorsCount](
		local.NATS_ADAPTER_SUBJECT+"devices.vendors",
		[]byte(""),
		w,
		a.nc,
	)
	if err != nil {
		return
	}

	statusCount, err := bridge.NatsReq[[]entity.StatusCount](
		local.NATS_ADAPTER_SUBJECT+"devices.status",
		[]byte(""),
		w,
		a.nc,
	)
	if err != nil {
		return
	}

	for _, v := range statusCount.Msg {
		switch entity.Status(v.Status) {
		case entity.Online:
			result.StatusCount.Online = v.Count
		case entity.Offline:
			result.StatusCount.Offline = v.Count
		}
	}

	result.VendorsCount = vendorcount.Msg
	result.ProductClassCount = productclasscount.Msg

	now := time.Now()
	_, err = bridge.NatsReqWithoutHttpSet[time.Duration](
		local.NATS_WS_ADAPTER_SUBJECT_PREFIX+"rtt",
		[]byte(""),
		a.nc,
	)
	if err == nil {
		result.WebsocketsRtt = time.Until(now).String()
	}

	now = time.Now()
	_, err = bridge.NatsReqWithoutHttpSet[time.Duration](
		local.NATS_STOMP_ADAPTER_SUBJECT_PREFIX+"rtt",
		[]byte(""),
		a.nc,
	)
	if err == nil {
		result.StompRtt = time.Until(now).String()
	}

	now = time.Now()
	_, err = bridge.NatsReqWithoutHttpSet[time.Duration](
		local.NATS_MQTT_ADAPTER_SUBJECT_PREFIX+"rtt",
		[]byte(""),
		a.nc,
	)
	if err == nil {
		result.MqttRtt = time.Until(now).String()
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println(err)
	}
}

func (a *Api) vendorsInfo(w http.ResponseWriter, r *http.Request) {
	vendors, err := bridge.NatsReq[[]entity.VendorsCount](
		local.NATS_ADAPTER_SUBJECT+"devices.vendors",
		[]byte(""),
		w,
		a.nc,
	)
	if err != nil {
		return
	}
	utils.MarshallEncoder(vendors.Msg, w)
}

func (a *Api) productClassInfo(w http.ResponseWriter, r *http.Request) {
	vendors, err := bridge.NatsReq[[]entity.ProductClassCount](
		local.NATS_ADAPTER_SUBJECT+"devices.class",
		[]byte(""),
		w,
		a.nc,
	)
	if err != nil {
		return
	}
	utils.MarshallEncoder(vendors.Msg, w)
}

func (a *Api) statusInfo(w http.ResponseWriter, r *http.Request) {
	vendors, err := bridge.NatsReq[[]entity.StatusCount](
		local.NATS_ADAPTER_SUBJECT+"devices.status",
		[]byte(""),
		w,
		a.nc,
	)
	if err != nil {
		return
	}

	var status StatusCount
	for _, v := range vendors.Msg {
		switch entity.Status(v.Status) {
		case entity.Online:
			status.Online = v.Count
		case entity.Offline:
			status.Offline = v.Count
		}
	}

	utils.MarshallEncoder(status, w)
}
