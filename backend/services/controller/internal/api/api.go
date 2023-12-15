package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/leandrofars/oktopus/internal/api/cors"
	"github.com/leandrofars/oktopus/internal/api/middleware"
	"github.com/leandrofars/oktopus/internal/db"
	"github.com/leandrofars/oktopus/internal/mqtt"
	"github.com/leandrofars/oktopus/internal/mtp"
	usp_msg "github.com/leandrofars/oktopus/internal/usp_message"
	"github.com/leandrofars/oktopus/internal/utils"
	"google.golang.org/protobuf/proto"
)

type Api struct {
	Port     string
	Db       db.Database
	Broker   mtp.Broker
	MsgQueue map[string](chan usp_msg.Msg)
	QMutex   *sync.Mutex
	Mqtt     mqtt.Mqtt
}

const REQUEST_TIMEOUT = time.Second * 30

const (
	NormalUser = iota
	AdminUser
)

func NewApi(port string, db db.Database, mqtt *mqtt.Mqtt, msgQueue map[string](chan usp_msg.Msg), m *sync.Mutex) Api {
	return Api{
		Port:     port,
		Db:       db,
		Broker:   mqtt,
		MsgQueue: msgQueue,
		QMutex:   m,
		Mqtt:     *mqtt,
	}
}

func StartApi(a Api) {
	r := mux.NewRouter()
	authentication := r.PathPrefix("/api/auth").Subrouter()
	authentication.HandleFunc("/login", a.generateToken).Methods("PUT")
	authentication.HandleFunc("/register", a.registerUser).Methods("POST")
	authentication.HandleFunc("/admin/register", a.registerAdminUser).Methods("POST")
	authentication.HandleFunc("/admin/exists", a.adminUserExists).Methods("GET")
	iot := r.PathPrefix("/api/device").Subrouter()
	//TODO: create query for devices
	iot.HandleFunc("", a.retrieveDevices).Methods("GET")
	iot.HandleFunc("/{id}", a.retrieveDevices).Methods("GET")
	iot.HandleFunc("/{sn}/get", a.deviceGetMsg).Methods("PUT")
	iot.HandleFunc("/{sn}/add", a.deviceCreateMsg).Methods("PUT")
	iot.HandleFunc("/{sn}/del", a.deviceDeleteMsg).Methods("PUT")
	iot.HandleFunc("/{sn}/set", a.deviceUpdateMsg).Methods("PUT")
	iot.HandleFunc("/{sn}/parameters", a.deviceGetSupportedParametersMsg).Methods("PUT")
	iot.HandleFunc("/{sn}/instances", a.deviceGetParameterInstances).Methods("PUT")
	iot.HandleFunc("/{sn}/operate", a.deviceOperateMsg).Methods("PUT")
	iot.HandleFunc("/{sn}/fw_update", a.deviceFwUpdate).Methods("PUT")
	iot.HandleFunc("/{sn}/wifi", a.deviceWifi).Methods("PUT", "GET")
	mtp := r.PathPrefix("/api/mtp").Subrouter()
	mtp.HandleFunc("", a.mtpInfo).Methods("GET")
	dash := r.PathPrefix("/api/info").Subrouter()
	dash.HandleFunc("/vendors", a.vendorsInfo).Methods("GET")
	dash.HandleFunc("/status", a.statusInfo).Methods("GET")
	dash.HandleFunc("/device_class", a.productClassInfo).Methods("GET")
	dash.HandleFunc("/general", a.generalInfo).Methods("GET")
	users := r.PathPrefix("/api/users").Subrouter()
	users.HandleFunc("", a.retrieveUsers).Methods("GET")

	/* ----- Middleware for requests which requires user to be authenticated ---- */
	iot.Use(func(handler http.Handler) http.Handler {
		return middleware.Middleware(handler)
	})

	mtp.Use(func(handler http.Handler) http.Handler {
		return middleware.Middleware(handler)
	})

	dash.Use(func(handler http.Handler) http.Handler {
		return middleware.Middleware(handler)
	})

	users.Use(func(handler http.Handler) http.Handler {
		return middleware.Middleware(handler)
	})
	/* -------------------------------------------------------------------------- */

	// Verifies CORS configs for requests
	corsOpts := cors.GetCorsConfig()

	srv := &http.Server{
		Addr: "0.0.0.0:" + a.Port,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 60,
		IdleTimeout:  time.Second * 60,
		Handler:      corsOpts.Handler(r), // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	log.Println("Running Api at port", a.Port)
}

func (a *Api) uspCall(msg usp_msg.Msg, sn string, w http.ResponseWriter, device db.Device) {

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

	a.QMutex.Lock()
	a.MsgQueue[msg.Header.MsgId] = make(chan usp_msg.Msg)
	a.QMutex.Unlock()
	log.Println("Sending Msg:", msg.Header.MsgId)
	//TODO: Check what MTP the device is connected to
	a.Broker.Publish(tr369Message, "oktopus/v1/agent/"+sn, "oktopus/v1/api/"+sn, false)

	select {
	case msg := <-a.MsgQueue[msg.Header.MsgId]:
		log.Printf("Received Msg: %s", msg.Header.MsgId)
		a.QMutex.Lock()
		delete(a.MsgQueue, msg.Header.MsgId)
		a.QMutex.Unlock()
		log.Println("requests queue:", a.MsgQueue)
		body := msg.Body.GetResponse()
		switch body.RespType.(type) {
		case *usp_msg.Response_GetResp:
			json.NewEncoder(w).Encode(body.GetGetResp())
		case *usp_msg.Response_DeleteResp:
			json.NewEncoder(w).Encode(body.GetDeleteResp())
		case *usp_msg.Response_AddResp:
			json.NewEncoder(w).Encode(body.GetAddResp())
		case *usp_msg.Response_SetResp:
			json.NewEncoder(w).Encode(body.GetSetResp())
		case *usp_msg.Response_GetInstancesResp:
			json.NewEncoder(w).Encode(body.GetGetInstancesResp())
		case *usp_msg.Response_GetSupportedDmResp:
			json.NewEncoder(w).Encode(body.GetGetSupportedDmResp())
		case *usp_msg.Response_GetSupportedProtocolResp:
			json.NewEncoder(w).Encode(body.GetGetSupportedProtocolResp())
		case *usp_msg.Response_NotifyResp:
			json.NewEncoder(w).Encode(body.GetNotifyResp())
		case *usp_msg.Response_OperateResp:
			json.NewEncoder(w).Encode(body.GetOperateResp())
		default:
			json.NewEncoder(w).Encode("Unknown message answer")
		}
		return
	case <-time.After(REQUEST_TIMEOUT):
		log.Printf("Request %s Timed Out", msg.Header.MsgId)
		w.WriteHeader(http.StatusGatewayTimeout)
		a.QMutex.Lock()
		delete(a.MsgQueue, msg.Header.MsgId)
		a.QMutex.Unlock()
		log.Println("requests queue:", a.MsgQueue)
		json.NewEncoder(w).Encode("Request Timed Out")
		return
	}
}
