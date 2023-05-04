package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/leandrofars/oktopus/internal/api/auth"
	"github.com/leandrofars/oktopus/internal/api/middleware"
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
	authentication := r.PathPrefix("/auth").Subrouter()
	authentication.HandleFunc("/login", a.generateToken).Methods("PUT")
	//authentication.HandleFunc("/register", a.registerUser).Methods("POST")
	iot := r.PathPrefix("/device").Subrouter()
	iot.HandleFunc("/", a.retrieveDevices).Methods("GET")
	iot.HandleFunc("/{sn}/get", a.deviceGetMsg).Methods("PUT")
	iot.HandleFunc("/{sn}/add", a.deviceCreateMsg).Methods("PUT")
	iot.HandleFunc("/{sn}/del", a.deviceDeleteMsg).Methods("PUT")
	iot.HandleFunc("/{sn}/set", a.deviceUpdateMsg).Methods("PUT")
	//TODO: Create operation action handler
	iot.HandleFunc("/device/{sn}/act", a.deviceUpdateMsg).Methods("PUT")

	iot.Use(func(handler http.Handler) http.Handler {
		return middleware.Middleware(handler)
	})

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
		log.Printf("Received Msg: %s", msg.Header.MsgId)
		json.NewEncoder(w).Encode(msg.Body.GetResponse().GetAddResp())
		return
	case <-time.After(time.Second * 5):
		log.Printf("Request %s Timed Out", msg.Header.MsgId)
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
		log.Printf("Received Msg: %s", msg.Header.MsgId)
		json.NewEncoder(w).Encode(msg.Body.GetResponse().GetGetResp())
		return
	case <-time.After(time.Second * 5):
		log.Printf("Request %s Timed Out", msg.Header.MsgId)
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
		log.Printf("Received Msg: %s", msg.Header.MsgId)
		json.NewEncoder(w).Encode(msg.Body.GetResponse().GetDeleteResp())
		return
	case <-time.After(time.Second * 5):
		log.Printf("Request %s Timed Out", msg.Header.MsgId)
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
		log.Printf("Received Msg: %s", msg.Header.MsgId)
		json.NewEncoder(w).Encode(msg.Body.GetResponse().GetSetResp())
		return
	case <-time.After(time.Second * 5):
		log.Printf("Request %s Timed Out", msg.Header.MsgId)
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

func (a *Api) registerUser(w http.ResponseWriter, r *http.Request) {
	var user db.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := user.HashPassword(user.Password); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := a.Db.RegisterUser(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *Api) generateToken(w http.ResponseWriter, r *http.Request) {
	var tokenReq TokenRequest

	err := json.NewDecoder(r.Body).Decode(&tokenReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := a.Db.FindUser(tokenReq.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Invalid Credentials")
		return
	}

	credentialError := user.CheckPassword(tokenReq.Password)
	if credentialError != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Invalid Credentials")
		return
	}

	token, err := auth.GenerateJWT(user.Email, user.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
	return
}
