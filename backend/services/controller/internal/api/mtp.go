package api

import (
	"encoding/json"
	"net"
	"net/http"
	"time"

	"golang.org/x/sys/unix"
)

type mqttInfo struct {
	MqttRtt time.Duration
}

func (a *Api) mtpInfo(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(mqttInfo{
		MqttRtt: rtt / 1000,
	})
}

func tcpInfo(conn *net.TCPConn) (*unix.TCPInfo, error) {
	raw, err := conn.SyscallConn()
	if err != nil {
		return nil, err
	}

	var info *unix.TCPInfo
	ctrlErr := raw.Control(func(fd uintptr) {
		info, err = unix.GetsockoptTCPInfo(int(fd), unix.IPPROTO_TCP, unix.TCP_INFO)
	})
	switch {
	case ctrlErr != nil:
		return nil, ctrlErr
	case err != nil:
		return nil, err
	}
	return info, nil
}
