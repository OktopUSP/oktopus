package stomp

import (
	"log"
	"os"
	"time"

	"github.com/go-stomp/stomp"
)

type Stomp struct {
	Addr      string
	Conn      *stomp.Conn
	StopConn  os.Signal
	Connected bool
	Username  string
	Password  string
}

func (s *Stomp) Connect() {

	log.Println("STOMP username:", s.Username)
	log.Println("STOMP password:", s.Password)

	var options []func(*stomp.Conn) error = []func(*stomp.Conn) error{
		stomp.ConnOpt.Login(s.Username, s.Password),
		stomp.ConnOpt.Host("/"),
	}

	const MAX_TRIES = 3

	for i := 0; i < MAX_TRIES; i++ {
		log.Println("Starting new STOMP client")
		stompConn, err := stomp.Dial("tcp", s.Addr, options...)
		if err != nil {
			log.Println("Error connecting to STOMP server:", err.Error())
			if i == MAX_TRIES-1 {
				log.Printf("Reached max tries count: %d, stop trying to connect", MAX_TRIES)
				return
			}
			time.Sleep(1 * time.Second)
			continue
		}
		s.Conn = stompConn
		s.Connected = true
		break
	}

	log.Println("Connected to STOMP broker-->", s.Addr)
}

func (s *Stomp) Disconnect() {
	if s.Connected {
		s.Conn.Disconnect()
	}
	return
}

func (s *Stomp) Publish(msg []byte, topic, respTopic string, retain bool) {
	//s.Conn.Send()
}

func (s *Stomp) Subscribe() {
	//s.Conn.Subscribe()
}
