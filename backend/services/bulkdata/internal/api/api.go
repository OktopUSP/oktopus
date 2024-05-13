package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/oktopUSP/backend/services/bulkdata/internal/api/cors"
	"github.com/oktopUSP/backend/services/bulkdata/internal/api/handler"
	"github.com/oktopUSP/backend/services/bulkdata/internal/api/middleware"
	"github.com/oktopUSP/backend/services/bulkdata/internal/config"
)

type Api struct {
	port    string
	handler handler.Handler
}

const REQUEST_TIMEOUT = time.Second * 30

func NewApi(c config.RestApi, js jetstream.JetStream, nc *nats.Conn, kv jetstream.KeyValue) Api {
	return Api{
		port:    c.Port,
		handler: handler.NewHandler(js, nc, kv),
	}
}

func (a *Api) StartApi() {
	r := mux.NewRouter()
	r.HandleFunc("/healthcheck", a.handler.Healthcheck).Methods("GET")
	r.HandleFunc("/", a.handler.Data).Methods("POST")

	/* ----- Middleware for requests which requires user to be authenticated ---- */
	r.Use(func(handler http.Handler) http.Handler {
		return middleware.Middleware(handler)
	})
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
	log.Println("Running Bulk Data Collector HTTP Server at port", a.port)
}
