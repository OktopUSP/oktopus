package ws

import (
	"fmt"
	"log"
	"net/http"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/gorilla/mux"
	"github.com/leandrofars/oktopus/internal/api/cors"
)

/* ----------- [Deprecated code] migrated to Socketio with NodeJs ----------- */
func Ws() {
	server := socketio.NewServer(&engineio.Options{
		PingTimeout:  5 * time.Second,
		PingInterval: 10 * time.Second,
		Transports: []transport.Transport{
			&polling.Transport{
				Client: &http.Client{
					Timeout: time.Minute,
				},
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
			&websocket.Transport{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		},
	})

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "offer", func(s socketio.Conn, msg string) string {
		log.Printf("offer: %s", msg)

		return "test"
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()

	const wsPort = "5000"

	r := mux.NewRouter()
	r.Handle("/socket.io/", server)

	//r.Use(func(handler http.Handler) http.Handler {
	//	return middleware.Middleware(handler)
	//})

	corsOpts := cors.GetCorsConfig()

	srv := &http.Server{
		Addr:    "0.0.0.0:" + wsPort,
		Handler: corsOpts.Handler(r), // Pass our instance of gorilla/mux in.
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	log.Printf("Running websocket at port %s", wsPort)
}
