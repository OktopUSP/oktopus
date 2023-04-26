package mqtt

//
//import (
//	usp_msg "github.com/leandrofars/oktopus/internal/usp_message"
//	"github.com/leandrofars/oktopus/internal/usp_record"
//	"google.golang.org/protobuf/proto"
//	"log"
//)
//
//func SendGetMsg(sn string) {
//	payload := usp_msg.Msg{
//		Header: &usp_msg.Header{
//			MsgId:   "uniqueIdentifierForThismessage",
//			MsgType: usp_msg.Header_GET,
//		},
//		Body: &usp_msg.Body{
//			MsgBody: &usp_msg.Body_Request{
//				Request: &usp_msg.Request{
//					ReqType: &usp_msg.Request_Get{
//						Get: &usp_msg.Get{
//							ParamPaths: []string{
//								"Device.DeviceInfo.Manufacturer",
//								"Device.DeviceInfo.ModelName",
//								"Device.DeviceInfo.SoftwareVersion",
//							},
//							MaxDepth: 1,
//						},
//					},
//				},
//			},
//		},
//	}
//	teste, _ := proto.Marshal(&payload)
//	record := usp_record.Record{
//		Version:         "0.1",
//		ToId:            sn,
//		FromId:          "leleco",
//		PayloadSecurity: usp_record.Record_PLAINTEXT,
//		RecordType: &usp_record.Record_NoSessionContext{
//			NoSessionContext: &usp_record.NoSessionContextRecord{
//				Payload: teste,
//			},
//		},
//	}
//
//	tr369Message, err := proto.Marshal(&record)
//	if err != nil {
//		log.Fatalln("Failed to encode address book:", err)
//	}
//	m.Publish(tr369Message, "oktopus/v1/agent/"+deviceMac, "oktopus/v1/controller/"+deviceMac)
//}
