package ws

// Websockets server implementation inspired by https://github.com/gorilla/websocket/tree/main/examples/chat

import (
	"log"
	"net/http"

	"github.com/OktopUSP/oktopus/ws/internal/ws/handler"
	"github.com/gorilla/mux"
)

// Starts New Websockets Server
func StartNewServer() {
	// Initialize handlers of websockets events
	go handler.InitHandlers()

	r := mux.NewRouter()
	r.HandleFunc("/ws/agent", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeAgent(w, r)
	})
	r.HandleFunc("/ws/controller", func(w http.ResponseWriter, r *http.Request) {
		//TODO: Implement controller handler
	})

	log.Println("Websockets server running")

	// Blocks application running until it receives a KILL signal
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
