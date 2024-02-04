package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func StartNewServer() {
	r := mux.NewRouter()
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		header := http.Header{
			"Sec-Websocket-Protocol": {"v1.usp"},
			"Sec-Websocket-Version":  {"13"},
		}

		conn, err := upgrader.Upgrade(w, r, header)
		if err != nil {
			log.Println(err)
		}
		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error to read message:", err)
				return
			}
			log.Println("Message", string(p))
		}
	})

	log.Println("Websockets server running")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
