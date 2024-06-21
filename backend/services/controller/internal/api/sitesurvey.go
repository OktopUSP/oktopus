package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leandrofars/oktopus/internal/entity"
	"github.com/leandrofars/oktopus/internal/utils"
)

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
