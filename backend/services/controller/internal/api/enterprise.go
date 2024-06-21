package api

import (
	"net/http"

	"github.com/leandrofars/oktopus/internal/bridge"
	"github.com/leandrofars/oktopus/internal/entity"
)

func (a *Api) getEnterpriseResource(
	resource string,
	action string,
	device *entity.Device,
	sn string,
	w http.ResponseWriter,
	body []byte,
) error {
	model, err := cwmpGetDeviceModel(device, w)
	if err != nil {
		return err
	}
	err = bridge.NatsEnterpriseInteraction("enterprise.v1.cwmp."+model+"."+sn+"."+resource+"."+action, body, w, a.nc)
	return err
}
