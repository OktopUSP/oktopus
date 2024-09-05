package api

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leandrofars/oktopus/internal/bridge"
	"github.com/leandrofars/oktopus/internal/entity"
	local "github.com/leandrofars/oktopus/internal/nats"
	"github.com/leandrofars/oktopus/internal/utils"
	"github.com/nats-io/nats.go"
)

var errInvalidMtp = errors.New("Invalid MTP, valid options are: mqtt, ws, stomp")

func deviceStateOK(w http.ResponseWriter, nc *nats.Conn, sn string) (string, bool) {

	device, err := getDeviceInfo(w, sn, nc)
	if err != nil {
		return "", false
	}

	if !isDeviceOnline(w, device.Status) {
		return "", false
	}

	if device.Mqtt == entity.Online {
		return entity.Mqtt, true
	}

	if device.Websockets == entity.Online {
		return entity.Websockets, true
	}

	if device.Stomp == entity.Online {
		return entity.Stomp, true
	}

	return "", false
}

func getSerialNumberFromRequest(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["sn"]
}

func getMtpFromRequest(r *http.Request, w http.ResponseWriter) (string, error) {
	vars := mux.Vars(r)
	switch vars["mtp"] {
	case entity.Mqtt:
		return entity.Mqtt, nil
	case entity.Websockets:
		return entity.Websockets, nil
	case entity.Stomp:
		return entity.Stomp, nil
	case "any":
		return "", nil
	case ":mtp":
		return "", nil
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.Marshall("Invalid MTP, valid options are: " + entity.Mqtt + ", " + entity.Websockets + ", " + entity.Stomp))
		return "", errInvalidMtp
	}
}

func isDeviceOnline(w http.ResponseWriter, deviceStatus entity.Status) bool {
	if deviceStatus != entity.Online {
		w.WriteHeader(http.StatusServiceUnavailable)
		switch deviceStatus {
		case entity.Offline:
			w.Write(utils.Marshall("Device is offline"))
		case entity.Associating:
			w.Write(utils.Marshall("Device status is associating"))
		default:
			w.Write(utils.Marshall("Unknown device status"))
		}
		return false
	}
	return true
}

func getDeviceInfo(w http.ResponseWriter, sn string, nc *nats.Conn) (device *entity.Device, err error) {
	msg, err := bridge.NatsReq[entity.Device](
		local.NATS_ADAPTER_SUBJECT+sn+".device",
		[]byte(""),
		w,
		nc,
	)
	if msg != nil {
		return &msg.Msg, err
	}
	return nil, err
}

func getDeviceCount(w http.ResponseWriter, nc *nats.Conn) (int64, error) {
	msg, err := bridge.NatsReq[int64](
		local.NATS_ADAPTER_SUBJECT+"devices.count",
		[]byte(""),
		w,
		nc,
	)
	return msg.Msg, err
}

func getDevices(w http.ResponseWriter, filter map[string]interface{}, nc *nats.Conn) (*entity.DevicesList, error) {
	msg, err := bridge.NatsReq[entity.DevicesList](
		local.NATS_ADAPTER_SUBJECT+"devices.retrieve",
		utils.Marshall(filter),
		w,
		nc,
	)
	if msg != nil {
		return &msg.Msg, err
	}
	return nil, err
}
