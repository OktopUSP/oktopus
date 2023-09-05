package api

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/leandrofars/oktopus/internal/api/cors"
	"github.com/leandrofars/oktopus/internal/api/middleware"
	"github.com/leandrofars/oktopus/internal/db"
	"github.com/leandrofars/oktopus/internal/mtp"
	usp_msg "github.com/leandrofars/oktopus/internal/usp_message"
)

type Api struct {
	Port     string
	Db       db.Database
	Broker   mtp.Broker
	MsgQueue map[string](chan usp_msg.Msg)
	QMutex   *sync.Mutex
}

type WiFi struct {
	SSID                 string   `json:"ssid"`
	Password             string   `json:"password"`
	Security             string   `json:"security"`
	SecurityCapabilities []string `json:"securityCapabilities"`
	AutoChannelEnable    bool     `json:"autoChannelEnable"`
	Channel              int      `json:"channel"`
	ChannelBandwidth     string   `json:"channelBandwidth"`
	FrequencyBand        string   `json:"frequencyBand"`
	//PossibleChannels     		[]int    `json:"PossibleChannels"`
	SupportedChannelBandwidths []string `json:"supportedChannelBandwidths"`
}

const (
	NormalUser = iota
	AdminUser
)

func NewApi(port string, db db.Database, b mtp.Broker, msgQueue map[string](chan usp_msg.Msg), m *sync.Mutex) Api {
	return Api{
		Port:     port,
		Db:       db,
		Broker:   b,
		MsgQueue: msgQueue,
		QMutex:   m,
	}
}

//TODO: restructure http api calls for mqtt, to use golang generics and avoid code repetition
//TODO: standardize timeouts through code
//TODO: fix api methods

func StartApi(a Api) {
	r := mux.NewRouter()
	authentication := r.PathPrefix("/api/auth").Subrouter()
	authentication.HandleFunc("/login", a.generateToken).Methods("PUT")
	authentication.HandleFunc("/register", a.registerUser).Methods("POST")
	authentication.HandleFunc("/admin/register", a.registerAdminUser).Methods("POST")
	authentication.HandleFunc("/admin/exists", a.adminUserExists).Methods("GET")
	iot := r.PathPrefix("/api/device").Subrouter()
	iot.HandleFunc("", a.retrieveDevices).Methods("GET")
	iot.HandleFunc("/{sn}/get", a.deviceGetMsg).Methods("PUT")
	iot.HandleFunc("/{sn}/add", a.deviceCreateMsg).Methods("PUT")
	iot.HandleFunc("/{sn}/del", a.deviceDeleteMsg).Methods("PUT")
	iot.HandleFunc("/{sn}/set", a.deviceUpdateMsg).Methods("PUT")
	iot.HandleFunc("/{sn}/parameters", a.deviceGetSupportedParametersMsg).Methods("PUT")
	iot.HandleFunc("/{sn}/instances", a.deviceGetParameterInstances).Methods("PUT")
	iot.HandleFunc("/{sn}/update", a.deviceFwUpdate).Methods("PUT")
	iot.HandleFunc("/{sn}/wifi", a.deviceWifi).Methods("PUT", "GET")

	// Middleware for requests which requires user to be authenticated
	iot.Use(func(handler http.Handler) http.Handler {
		return middleware.Middleware(handler)
	})

	users := r.PathPrefix("/api/users").Subrouter()
	users.HandleFunc("", a.retrieveUsers).Methods("GET")

	users.Use(func(handler http.Handler) http.Handler {
		return middleware.Middleware(handler)
	})

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
