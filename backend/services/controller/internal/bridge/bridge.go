package bridge

import (
	"encoding/json"
	"log"
	"net/http"

	local "github.com/leandrofars/oktopus/internal/nats"
	"github.com/leandrofars/oktopus/internal/utils"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type DataType interface {
	[]map[string]interface{}
}

type Bridge struct {
	js jetstream.JetStream
	nc *nats.Conn
}

func NewBridge(js jetstream.JetStream, nc *nats.Conn) Bridge {
	return Bridge{
		js: js,
		nc: nc,
	}
}

func NatsReq[T DataType](
	subj string,
	body []byte,
	w http.ResponseWriter,
	nc *nats.Conn,
) (T, error) {

	var answer T

	msg, err := nc.Request(subj, body, local.NATS_REQUEST_TIMEOUT)
	if err != nil {
		log.Println(err)
		w.Write(utils.Marshall("Error to communicate with nats: " + err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	err = json.Unmarshal(msg.Data, &answer)
	if err != nil {
		log.Println(err)
		w.Write(msg.Data)
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	}

	return answer, nil
}
