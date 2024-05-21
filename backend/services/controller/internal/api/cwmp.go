package api

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"

	"github.com/leandrofars/oktopus/internal/bridge"
	"github.com/leandrofars/oktopus/internal/cwmp"
	"github.com/leandrofars/oktopus/internal/nats"
	"github.com/leandrofars/oktopus/internal/utils"
)

func (a *Api) cwmpGetParameterNamesMsg(w http.ResponseWriter, r *http.Request) {
	sn := getSerialNumberFromRequest(r)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.Marshall(err.Error()))
		return
	}

	data, err := bridge.NatsCwmpInteraction(
		nats.NATS_CWMP_ADAPTER_SUBJECT_PREFIX+sn+".api",
		payload,
		w,
		a.nc,
	)
	if err != nil {
		return
	}

	var response cwmp.GetParameterNamesResponse
	err = xml.Unmarshal(data, &response)
	if err != nil {
		err = json.Unmarshal(data, &response)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(utils.Marshall(err))
			return
		}
		return
	}

	w.Write(data)
}

func (a *Api) cwmpGetParameterValuesMsg(w http.ResponseWriter, r *http.Request) {
	sn := getSerialNumberFromRequest(r)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.Marshall(err.Error()))
		return
	}

	data, err := bridge.NatsCwmpInteraction(
		nats.NATS_CWMP_ADAPTER_SUBJECT_PREFIX+sn+".api",
		payload,
		w,
		a.nc,
	)
	if err != nil {
		return
	}

	var response cwmp.GetParameterValuesResponse
	err = xml.Unmarshal(data, &response)
	if err != nil {
		err = json.Unmarshal(data, &response)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(utils.Marshall(err))
			return
		}
		return
	}

	w.Write(data)
}

func (a *Api) cwmpSetParameterValuesMsg(w http.ResponseWriter, r *http.Request) {
	sn := getSerialNumberFromRequest(r)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.Marshall(err.Error()))
		return
	}

	data, err := bridge.NatsCwmpInteraction(
		nats.NATS_CWMP_ADAPTER_SUBJECT_PREFIX+sn+".api",
		payload,
		w,
		a.nc,
	)
	if err != nil {
		return
	}

	var response cwmp.SetParameterValuesResponse
	err = xml.Unmarshal(data, &response)
	if err != nil {
		err = json.Unmarshal(data, &response)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(utils.Marshall(err))
			return
		}
		return
	}

	w.Write(data)
}
