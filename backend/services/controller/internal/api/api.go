package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/leandrofars/oktopus/internal/api/cors"
	"github.com/leandrofars/oktopus/internal/api/middleware"
	"github.com/leandrofars/oktopus/internal/bridge"
	"github.com/leandrofars/oktopus/internal/config"
	"github.com/leandrofars/oktopus/internal/db"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Api struct {
	port   string
	js     jetstream.JetStream
	nc     *nats.Conn
	bridge bridge.Bridge
	db     db.Database
	ctx    context.Context
}

const REQUEST_TIMEOUT = time.Second * 30

const (
	NormalUser = iota
	AdminUser
)

func NewApi(c config.RestApi, js jetstream.JetStream, nc *nats.Conn, bridge bridge.Bridge, d db.Database) Api {
	return Api{
		port:   c.Port,
		js:     js,
		nc:     nc,
		ctx:    c.Ctx,
		bridge: bridge,
		db:     d,
	}
}

func (a *Api) StartApi() {
	r := mux.NewRouter()
	authentication := r.PathPrefix("/api/auth").Subrouter()
	authentication.HandleFunc("/login", a.generateToken).Methods("PUT")
	authentication.HandleFunc("/register", a.registerUser).Methods("POST")
	authentication.HandleFunc("/admin/register", a.registerAdminUser).Methods("POST")
	authentication.HandleFunc("/admin/exists", a.adminUserExists).Methods("GET")
	iot := r.PathPrefix("/api/device").Subrouter()
	// iot.HandleFunc("", a.retrieveDevices).Methods("GET")
	// iot.HandleFunc("/{id}", a.retrieveDevices).Methods("GET")
	iot.HandleFunc("/{sn}/{mtp}/get", a.deviceGetMsg).Methods("PUT")
	iot.HandleFunc("/{sn}/{mtp}/add", a.deviceCreateMsg).Methods("PUT")
	iot.HandleFunc("/{sn}/{mtp}/del", a.deviceDeleteMsg).Methods("PUT")
	iot.HandleFunc("/{sn}/{mtp}/set", a.deviceUpdateMsg).Methods("PUT")
	iot.HandleFunc("/{sn}/{mtp}/parameters", a.deviceGetSupportedParametersMsg).Methods("PUT")
	iot.HandleFunc("/{sn}/{mtp}/instances", a.deviceGetParameterInstances).Methods("PUT")
	iot.HandleFunc("/{sn}/{mtp}/operate", a.deviceOperateMsg).Methods("PUT")
	// iot.HandleFunc("/{sn}/{mtp}/fw_update", a.deviceFwUpdate).Methods("PUT")
	// iot.HandleFunc("/{sn}/{mtp}/wifi", a.deviceWifi).Methods("PUT", "GET")
	// mtp := r.PathPrefix("/api/mtp").Subrouter()
	// mtp.HandleFunc("", a.mtpInfo).Methods("GET")
	// dash := r.PathPrefix("/api/info").Subrouter()
	// dash.HandleFunc("/vendors", a.vendorsInfo).Methods("GET")
	// dash.HandleFunc("/status", a.statusInfo).Methods("GET")
	// dash.HandleFunc("/device_class", a.productClassInfo).Methods("GET")
	// dash.HandleFunc("/general", a.generalInfo).Methods("GET")
	users := r.PathPrefix("/api/users").Subrouter()
	users.HandleFunc("", a.retrieveUsers).Methods("GET")

	/* ----- Middleware for requests which requires user to be authenticated ---- */
	iot.Use(func(handler http.Handler) http.Handler {
		return middleware.Middleware(handler)
	})

	// mtp.Use(func(handler http.Handler) http.Handler {
	// 	return middleware.Middleware(handler)
	// })

	// dash.Use(func(handler http.Handler) http.Handler {
	// 	return middleware.Middleware(handler)
	// })

	// users.Use(func(handler http.Handler) http.Handler {
	// 	return middleware.Middleware(handler)
	// })
	/* -------------------------------------------------------------------------- */

	corsOpts := cors.GetCorsConfig()

	srv := &http.Server{
		Addr:         "0.0.0.0:" + a.port,
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 60,
		IdleTimeout:  time.Second * 60,
		Handler:      corsOpts.Handler(r),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	log.Println("Running REST API at port", a.port)
}

// func (a *Api) uspCall(msg usp_msg.Msg, sn string, w http.ResponseWriter, device db.Device) {

// 	encodedMsg, err := proto.Marshal(&msg)
// 	if err != nil {
// 		log.Println(err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	record := utils.NewUspRecord(encodedMsg, sn)
// 	tr369Message, err := proto.Marshal(&record)
// 	if err != nil {
// 		log.Fatalln("Failed to encode tr369 record:", err)
// 	}

// 	a.QMutex.Lock()
// 	a.MsgQueue[msg.Header.MsgId] = make(chan usp_msg.Msg)
// 	a.QMutex.Unlock()
// 	log.Println("Sending Msg:", msg.Header.MsgId)

// 	if device.Mqtt == db.Online {
// 		a.Mqtt.Publish(tr369Message, "oktopus/v1/agent/"+sn, "oktopus/v1/api/"+sn, false)
// 	} else if device.Websockets == db.Online {
// 		a.Websockets.Publish(tr369Message, "", "", false)
// 	} else if device.Stomp == db.Online {
// 		//TODO: send stomp message
// 	}

// 	select {
// 	case msg := <-a.MsgQueue[msg.Header.MsgId]:
// 		log.Printf("Received Msg: %s", msg.Header.MsgId)
// 		a.QMutex.Lock()
// 		delete(a.MsgQueue, msg.Header.MsgId)
// 		a.QMutex.Unlock()
// 		log.Println("requests queue:", a.MsgQueue)
// 		body := msg.Body.GetResponse()
// 		switch body.RespType.(type) {
// 		case *usp_msg.Response_GetResp:
// 			json.NewEncoder(w).Encode(body.GetGetResp())
// 		case *usp_msg.Response_DeleteResp:
// 			json.NewEncoder(w).Encode(body.GetDeleteResp())
// 		case *usp_msg.Response_AddResp:
// 			json.NewEncoder(w).Encode(body.GetAddResp())
// 		case *usp_msg.Response_SetResp:
// 			json.NewEncoder(w).Encode(body.GetSetResp())
// 		case *usp_msg.Response_GetInstancesResp:
// 			json.NewEncoder(w).Encode(body.GetGetInstancesResp())
// 		case *usp_msg.Response_GetSupportedDmResp:
// 			json.NewEncoder(w).Encode(body.GetGetSupportedDmResp())
// 		case *usp_msg.Response_GetSupportedProtocolResp:
// 			json.NewEncoder(w).Encode(body.GetGetSupportedProtocolResp())
// 		case *usp_msg.Response_NotifyResp:
// 			json.NewEncoder(w).Encode(body.GetNotifyResp())
// 		case *usp_msg.Response_OperateResp:
// 			json.NewEncoder(w).Encode(body.GetOperateResp())
// 		default:
// 			json.NewEncoder(w).Encode("Unknown message answer")
// 		}
// 		return
// 	case <-time.After(REQUEST_TIMEOUT):
// 		log.Printf("Request %s Timed Out", msg.Header.MsgId)
// 		w.WriteHeader(http.StatusGatewayTimeout)
// 		a.QMutex.Lock()
// 		delete(a.MsgQueue, msg.Header.MsgId)
// 		a.QMutex.Unlock()
// 		log.Println("requests queue:", a.MsgQueue)
// 		json.NewEncoder(w).Encode("Request Timed Out")
// 		return
// 	}
// }
