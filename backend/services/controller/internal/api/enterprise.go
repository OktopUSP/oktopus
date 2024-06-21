package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leandrofars/oktopus/internal/bridge"
	"github.com/leandrofars/oktopus/internal/entity"
	"github.com/leandrofars/oktopus/internal/utils"
)

func (a *Api) getEnterpriseResource(
	resource string,
	action string,
	device *entity.Device,
	sn string,
	w http.ResponseWriter,
	body []byte,
	protocol, datamodel string,
) error {
	model, err := cwmpGetDeviceModel(device, w)
	if err != nil {
		return err
	}

	err = bridge.NatsEnterpriseInteraction("enterprise.v1."+protocol+"."+datamodel+"."+model+"."+sn+"."+resource+"."+action, body, w, a.nc)
	return err
}

func (a *Api) deviceSiteSurvey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sn := vars["sn"]

	device, err := getDeviceInfo(w, sn, a.nc)
	if err != nil {
		return
	}

	if r.Method == http.MethodGet {

		if device.Cwmp == entity.Online {
			a.getEnterpriseResource("sitesurvey", "get", device, sn, w, []byte{}, "cwmp", "098")
			return
		}

		if device.Mqtt == entity.Online || device.Stomp == entity.Online || device.Websockets == entity.Online {
			w.WriteHeader(http.StatusNotImplemented)
			w.Write(utils.Marshall("This feature is only working with CWMP devices"))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.Marshall("Device is Offline"))
	}
}

func (a *Api) deviceConnectedDevices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sn := vars["sn"]

	device, err := getDeviceInfo(w, sn, a.nc)
	if err != nil {
		return
	}

	if r.Method == http.MethodGet {

		if device.Cwmp == entity.Online {
			a.getEnterpriseResource("connecteddevices", "get", device, sn, w, []byte{}, "cwmp", "098")
			return
		}

		if device.Mqtt == entity.Online || device.Stomp == entity.Online || device.Websockets == entity.Online {
			w.WriteHeader(http.StatusNotImplemented)
			w.Write(utils.Marshall("This feature is only working with CWMP devices"))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.Marshall("Device is Offline"))
	}
}
