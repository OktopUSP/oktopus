package api

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/leandrofars/oktopus/internal/db"
	"github.com/leandrofars/oktopus/internal/utils"
)

type StatusCount struct {
	Online  int
	Offline int
}

type GeneralInfo struct {
	MqttRtt           time.Duration
	ProductClassCount []db.ProductClassCount
	StatusCount       StatusCount
	VendorsCount      []db.VendorsCount
}

func (a *Api) generalInfo(w http.ResponseWriter, r *http.Request) {

	var result GeneralInfo

	productclasscount, err := a.Db.RetrieveProductsClassInfo()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	vendorcount, err := a.Db.RetrieveVendorsInfo()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	statuscount, err := a.Db.RetrieveStatusInfo()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, v := range statuscount {
		switch v.Status {
		case utils.Online:
			result.StatusCount.Online = v.Count
		case utils.Offline:
			result.StatusCount.Offline = v.Count
		}
	}

	result.VendorsCount = vendorcount
	result.ProductClassCount = productclasscount

	/* ------------ TODO: [mqtt rtt] create common function for this ------------ */
	//TODO: address with value from env or something like that
	conn, err := net.Dial("tcp", "127.0.0.1:1883")
	if err != nil {
		json.NewEncoder(w).Encode("Error to connect to broker")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	info, err := tcpInfo(conn.(*net.TCPConn))
	if err != nil {
		json.NewEncoder(w).Encode("Error to get TCP socket info")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rtt := time.Duration(info.Rtt) * time.Microsecond
	/* -------------------------------------------------------------------------- */

	result.MqttRtt = rtt / 1000

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println(err)
	}

	return
}

func (a *Api) vendorsInfo(w http.ResponseWriter, r *http.Request) {
	vendors, err := a.Db.RetrieveVendorsInfo()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(vendors)
	if err != nil {
		log.Println(err)
	}

	return
}

func (a *Api) productClassInfo(w http.ResponseWriter, r *http.Request) {
	vendors, err := a.Db.RetrieveProductsClassInfo()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(vendors)
	if err != nil {
		log.Println(err)
	}

	return
}

func (a *Api) statusInfo(w http.ResponseWriter, r *http.Request) {
	vendors, err := a.Db.RetrieveStatusInfo()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var status StatusCount
	for _, v := range vendors {
		switch v.Status {
		case utils.Online:
			status.Online = v.Count
		case utils.Offline:
			status.Offline = v.Count
		}
	}

	err = json.NewEncoder(w).Encode(status)
	if err != nil {
		log.Println(err)
	}

	return
}
