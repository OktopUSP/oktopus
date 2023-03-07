// Made by Leandro Antônio Farias Machado (leandrofars@gmail.com)

package main

import (
	//"flag"
	//"fmt"
	//"github.com/leandrofars/oktopus/internal/usp_record"
	//"github.com/leandrofars/oktopus/internal/usp_message"
	//"github.com/golang/protobuf/proto"
	"github.com/leandrofars/oktopus/internal/mqtt"
	"log"
	//"os/exec"
)

func main() {
	done := make(chan bool)
	log.Println("Starting Oktopus Project TR-369 Controller...")
	log.Println("Starting Mosquitto Broker")
	go mqtt.StartMqttBroker()

	//TODO: Create more options to set using flags
	//TODO: Read user inputs

	// usp_record.Record{
	// 	Version:         "1.0",
	// 	ToId:            "os::4851CF-000000000002",
	// 	FromId:          "leleco",
	// 	PayloadSecurity: usp_record.Record_PLAINTEXT,
	// 	RecordType: &usp_record.Record_NoSessionContext{
	// 		NoSessionContext: &usp_record.NoSessionContextRecord{
	// 			Payload: []byte("payload"),
	// 		},
	// 	},
	// }

	<-done
}
