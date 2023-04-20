package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/leandrofars/oktopus/internal/db"
	"log"
	"net/http"
	"time"
)

type Api struct {
	Port string
	Db   db.Database
}

func NewApi(port string, db db.Database) Api {
	return Api{
		Port: port,
		Db:   db,
	}
}

func StartApi(a Api) {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		return
	})
	r.HandleFunc("/devices", a.retrieveDevices)
	//r.HandleFunc("/devices/{sn}", a.devicesMessaging)

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
