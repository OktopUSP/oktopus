package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/leandrofars/oktopus/internal/db"
	"github.com/leandrofars/oktopus/internal/mtp"
	usp_msg "github.com/leandrofars/oktopus/internal/usp_message"
	"github.com/leandrofars/oktopus/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
	"sync"
	"time"
)

type Api struct {
	Port     string
	Db       db.Database
	Broker   mtp.Broker
	MsgQueue map[string](chan usp_msg.Msg)
	QMutex   *sync.Mutex
}

func NewApi(port string, db db.Database, b mtp.Broker, msgQueue map[string](chan usp_msg.Msg), m *sync.Mutex) Api {
	return Api{
		Port:     port,
		Db:       db,
		Broker:   b,
		MsgQueue: msgQueue,
		QMutex:   m,
	}
}

func StartApi(a Api) {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		return
	})
	r.HandleFunc("/devices", a.retrieveDevices).Methods("GET")
	r.HandleFunc("/device/{sn}/get", a.deviceGetMsg).Methods("PUT")
	r.HandleFunc("/device/{sn}/add", a.deviceCreateMsg).Methods("PUT")
	r.HandleFunc("/device/{sn}/del", a.deviceDeleteMsg).Methods("PUT")
	r.HandleFunc("/device/{sn}/set", a.deviceUpdateMsg).Methods("PUT")

	srv := &http.Server{
		Addr: "0.0.0.0:" + a.Port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
}

func (a *Api) retrieveDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := a.Db.RetrieveDevices()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(devices)
	if err != nil {
		log.Println(err)
	}

	return
}

func (a *Api) deviceCreateMsg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sn := vars["sn"]
	a.deviceExists(sn, w)

	var receiver usp_msg.Add

	err := json.NewDecoder(r.Body).Decode(receiver)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := utils.NewCreateMsg(receiver)
	encodedMsg, err := proto.Marshal(&msg)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	record := utils.NewUspRecord(encodedMsg, sn)
	tr369Message, err := proto.Marshal(&record)
	if err != nil {
		log.Fatalln("Failed to encode tr369 record:", err)
	}

	//a.Broker.Request(tr369Message, usp_msg.Header_GET, "oktopus/v1/agent/"+sn, "oktopus/v1/get/"+sn)
	a.MsgQueue[msg.Header.MsgId] = make(chan usp_msg.Msg)
	a.Broker.Publish(tr369Message, "oktopus/v1/agent/"+sn, "oktopus/v1/api/"+sn)

	select {
	case msg := <-a.MsgQueue[msg.Header.MsgId]:
		log.Printf("Received Msg")
		json.NewEncoder(w).Encode(msg.Body.GetResponse().GetAddResp())
		return
	case <-time.After(time.Second * 5):
		log.Printf("Request Timed Out")
		w.WriteHeader(http.StatusGatewayTimeout)
		json.NewEncoder(w).Encode("Request Timed Out")
		return
	}
}

func (a *Api) deviceGetMsg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sn := vars["sn"]

	a.deviceExists(sn, w)

	var receiver usp_msg.Get

	err := json.NewDecoder(r.Body).Decode(&receiver)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := utils.NewGetMsg(receiver)
	encodedMsg, err := proto.Marshal(&msg)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	record := utils.NewUspRecord(encodedMsg, sn)
	tr369Message, err := proto.Marshal(&record)
	if err != nil {
		log.Fatalln("Failed to encode tr369 record:", err)
	}

	a.MsgQueue[msg.Header.MsgId] = make(chan usp_msg.Msg)
	a.Broker.Publish(tr369Message, "oktopus/v1/agent/"+sn, "oktopus/v1/api/"+sn)

	select {
	case msg := <-a.MsgQueue[msg.Header.MsgId]:
		log.Printf("Received Msg")
		json.NewEncoder(w).Encode(msg.Body.GetResponse().GetGetResp())
		return
	case <-time.After(time.Second * 5):
		log.Printf("Request Timed Out")
		w.WriteHeader(http.StatusGatewayTimeout)
		json.NewEncoder(w).Encode("Request Timed Out")
		return
	}
}

func (a *Api) deviceDeleteMsg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sn := vars["sn"]
	a.deviceExists(sn, w)

	var receiver usp_msg.Delete

	err := json.NewDecoder(r.Body).Decode(&receiver)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := utils.NewDelMsg(receiver)
	encodedMsg, err := proto.Marshal(&msg)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	record := utils.NewUspRecord(encodedMsg, sn)
	tr369Message, err := proto.Marshal(&record)
	if err != nil {
		log.Fatalln("Failed to encode tr369 record:", err)
	}

	//a.Broker.Request(tr369Message, usp_msg.Header_GET, "oktopus/v1/agent/"+sn, "oktopus/v1/get/"+sn)
	a.MsgQueue[msg.Header.MsgId] = make(chan usp_msg.Msg)
	a.Broker.Publish(tr369Message, "oktopus/v1/agent/"+sn, "oktopus/v1/api/"+sn)

	select {
	case msg := <-a.MsgQueue[msg.Header.MsgId]:
		log.Printf("Received Msg")
		json.NewEncoder(w).Encode(msg.Body.GetResponse().GetDeleteResp())
		return
	case <-time.After(time.Second * 5):
		log.Printf("Request Timed Out")
		w.WriteHeader(http.StatusGatewayTimeout)
		json.NewEncoder(w).Encode("Request Timed Out")
		return
	}
}

func (a *Api) deviceUpdateMsg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sn := vars["sn"]
	a.deviceExists(sn, w)

	var receiver usp_msg.Set

	err := json.NewDecoder(r.Body).Decode(&receiver)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := utils.NewSetMsg(receiver)
	encodedMsg, err := proto.Marshal(&msg)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	record := utils.NewUspRecord(encodedMsg, sn)
	tr369Message, err := proto.Marshal(&record)
	if err != nil {
		log.Fatalln("Failed to encode tr369 record:", err)
	}

	//a.Broker.Request(tr369Message, usp_msg.Header_GET, "oktopus/v1/agent/"+sn, "oktopus/v1/get/"+sn)
	a.MsgQueue[msg.Header.MsgId] = make(chan usp_msg.Msg)
	a.Broker.Publish(tr369Message, "oktopus/v1/agent/"+sn, "oktopus/v1/api/"+sn)

	select {
	case msg := <-a.MsgQueue[msg.Header.MsgId]:
		log.Printf("Received Msg")
		json.NewEncoder(w).Encode(msg.Body.GetResponse().GetSetResp())
		return
	case <-time.After(time.Second * 5):
		log.Printf("Request Timed Out")
		w.WriteHeader(http.StatusGatewayTimeout)
		json.NewEncoder(w).Encode("Request Timed Out")
		return
	}
}

func (a *Api) deviceExists(sn string, w http.ResponseWriter) {
	_, err := a.Db.RetrieveDevice(sn)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("No device with serial number " + sn + " was found")
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
