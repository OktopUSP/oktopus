package ws

// Websockets server implementation inspired by https://github.com/gorilla/websocket/tree/main/examples/chat

import (
	"log"
	"net/http"

	"github.com/OktopUSP/oktopus/ws/internal/config"
	"github.com/OktopUSP/oktopus/ws/internal/ws/handler"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go/jetstream"
)

// Starts New Websockets Server
func StartNewServer(c config.Config, kv jetstream.KeyValue) {
	// Initialize handlers of websockets events
	go handler.InitHandlers(c.ControllerEID)

	r := mux.NewRouter()
	r.HandleFunc("/ws/agent/{passwd}", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeAgent(w, r, c.ControllerEID, kv, c.Auth)
	})
	r.HandleFunc("/ws/agent", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeAgent(w, r, c.ControllerEID, kv, c.Auth)
	})
	r.HandleFunc("/ws/controller", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeController(w, r, c.ControllerEID, c.Auth, kv)
	})

	go func() {
		if c.Tls {
			log.Println("Websockets server running with TLS")
			err := http.ListenAndServeTLS(c.Port, "cert.pem", "key.pem", r)
			if err != nil {
				log.Fatal("ListenAndServeTLS: ", err)
			}
		} else {
			log.Println("Websockets server running at port", c.Port)
			err := http.ListenAndServe(c.Port, r)
			if err != nil {
				log.Fatal("ListenAndServe: ", err)
			}
		}
	}()
}
