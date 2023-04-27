package utils

import (
	"github.com/google/uuid"
	usp_msg "github.com/leandrofars/oktopus/internal/usp_message"
	"github.com/leandrofars/oktopus/internal/usp_record"
	"net"
)

//Status are saved at database as numbers
const (
	Online = iota
	Associating
	Offline
)

// Get interfaces MACs, and the first interface MAC is gonna be used as mqtt clientId
func GetMacAddr() ([]string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
}

func NewUspRecord(p []byte, toId string) usp_record.Record {
	return usp_record.Record{
		Version:         "0.1",
		ToId:            toId,
		FromId:          "leleco",
		PayloadSecurity: usp_record.Record_PLAINTEXT,
		RecordType: &usp_record.Record_NoSessionContext{
			NoSessionContext: &usp_record.NoSessionContextRecord{
				Payload: p,
			},
		},
	}
}

func NewCreateMsg(createStuff usp_msg.Add) usp_msg.Msg {
	return usp_msg.Msg{
		Header: &usp_msg.Header{
			MsgId:   uuid.NewString(),
			MsgType: usp_msg.Header_ADD,
		},
		Body: &usp_msg.Body{
			MsgBody: &usp_msg.Body_Request{
				Request: &usp_msg.Request{
					ReqType: &usp_msg.Request_Add{
						Add: &createStuff,
					},
				},
			},
		},
	}
}

func NewGetMsg(getStuff usp_msg.Get) usp_msg.Msg {
	return usp_msg.Msg{
		Header: &usp_msg.Header{
			MsgId:   uuid.NewString(),
			MsgType: usp_msg.Header_GET,
		},
		Body: &usp_msg.Body{
			MsgBody: &usp_msg.Body_Request{
				Request: &usp_msg.Request{
					ReqType: &usp_msg.Request_Get{
						Get: &getStuff,
					},
				},
			},
		},
	}
}

func NewDelMsg(getStuff usp_msg.Delete) usp_msg.Msg {
	return usp_msg.Msg{
		Header: &usp_msg.Header{
			MsgId:   uuid.NewString(),
			MsgType: usp_msg.Header_DELETE,
		},
		Body: &usp_msg.Body{
			MsgBody: &usp_msg.Body_Request{
				Request: &usp_msg.Request{
					ReqType: &usp_msg.Request_Delete{
						Delete: &getStuff,
					},
				},
			},
		},
	}
}
