package bridge

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/leandrofars/oktopus/internal/entity"
	local "github.com/leandrofars/oktopus/internal/nats"
	"github.com/leandrofars/oktopus/internal/utils"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

var errNatsMsgReceivedWithErrorData = errors.New("Nats message received with error data")
var errNatsRequestTimeout = errors.New("Nats message response timeout")

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

func NatsUspInteraction(
	subSubj, pubSubj string,
	body []byte,
	w http.ResponseWriter,
	nc *nats.Conn,
) ([]byte, error) {

	var answer []byte

	log.Println("Sending usp message")
	log.Println("subSubj: ", subSubj)
	log.Println("pubSubj: ", pubSubj)

	ch := make(chan *nats.Msg, 64)
	done := make(chan error)
	_, err := nc.ChanSubscribe(subSubj, ch)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.Marshall("Error to communicate with nats: " + err.Error()))
		return []byte{}, err
	}

	go func() {
		select {
		case msg := <-ch:
			log.Println("Received an usp message response")
			answer = msg.Data
			done <- nil
		case <-time.After(local.NATS_REQUEST_TIMEOUT):
			log.Println("usp message response timeout")
			w.WriteHeader(http.StatusGatewayTimeout)
			w.Write(utils.Marshall("usp message response timeout"))
			done <- errNatsRequestTimeout
		}
	}()

	err = nc.Publish(pubSubj, body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.Marshall("Error to communicate with nats: " + err.Error()))
		return nil, err
	}

	err = <-done

	return answer, err
}

func NatsCustomReq[T entity.DataType](
	subSubj, pubSubj string,
	body []byte,
	w http.ResponseWriter,
	nc *nats.Conn,
) (interface{}, error) {

	var answer T

	ch := make(chan *nats.Msg, 64)
	done := make(chan string)
	_, err := nc.ChanSubscribe(subSubj, ch)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.Marshall("Error to communicate with nats: " + err.Error()))
		return nil, err
	}

	select {
	case msg := <-ch:
		log.Println("Received an api message response")
		err = json.Unmarshal(msg.Data, &answer)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(msg.Data)
			return nil, err
		}
		done <- "done"
	case <-time.After(local.NATS_REQUEST_TIMEOUT):
		log.Println("Api message response timeout")
		done <- "timeout"
	}

	err = nc.Publish(pubSubj, body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.Marshall("Error to communicate with nats: " + err.Error()))
		return nil, err
	}

	<-done

	return nil, nil
}

/*
- makes a request to nats topic

- handle nats communication

- verify if received data is of error type
*/
func NatsReq[T entity.DataType](
	subj string,
	body []byte,
	w http.ResponseWriter,
	nc *nats.Conn,
) (*entity.MsgAnswer[T], error) {

	var answer *entity.MsgAnswer[T]

	msg, err := nc.Request(subj, body, local.NATS_REQUEST_TIMEOUT)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.Marshall("Error to communicate with nats: " + err.Error()))
		return nil, err
	}

	err = json.Unmarshal(msg.Data, &answer)
	if err != nil {
		var errMsg *entity.MsgAnswer[*string]
		err = json.Unmarshal(msg.Data, &errMsg)

		if err != nil {
			log.Println("Bad answer message formatting: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(msg.Data)
			return nil, err
		}

		log.Printf("message received, msg: %s, code: %d", *errMsg.Msg, errMsg.Code)
		w.WriteHeader(errMsg.Code)
		w.Write(utils.Marshall(*errMsg.Msg))
		return nil, errNatsMsgReceivedWithErrorData
	}

	return answer, nil
}

func NatsReqWithoutHttpSet[T entity.DataType](
	subj string,
	body []byte,
	nc *nats.Conn,
) (*entity.MsgAnswer[T], error) {

	var answer *entity.MsgAnswer[T]

	msg, err := nc.Request(subj, body, local.NATS_REQUEST_TIMEOUT)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = json.Unmarshal(msg.Data, &answer)
	if err != nil {

		var errMsg *entity.MsgAnswer[*string]
		err = json.Unmarshal(msg.Data, &errMsg)

		if err != nil {
			log.Println("Bad answer message formatting: ", err.Error())
			return nil, err
		}

		log.Printf("Error message received, msg: %s, code: %d", *errMsg.Msg, errMsg.Code)
		return nil, errNatsMsgReceivedWithErrorData
	}

	return answer, nil
}

func NatsCwmpInteraction(
	subj string,
	body []byte,
	w http.ResponseWriter,
	nc *nats.Conn,
) ([]byte, error) {

	log.Println("Sending cwmp message")
	log.Println("Subject: ", subj)

	var answer entity.MsgAnswer[[]byte]

	msg, err := nc.Request(subj, body, local.NATS_REQUEST_TIMEOUT)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.Marshall("Error to communicate with nats: " + err.Error()))
		return nil, err
	}

	err = json.Unmarshal(msg.Data, &answer)
	if err != nil {

		var errMsg *entity.MsgAnswer[*string]
		err = json.Unmarshal(msg.Data, &errMsg)

		if err != nil {
			log.Println("Bad answer message formatting: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(msg.Data)
			return nil, err
		}

		log.Printf("Error message received, msg: %s, code: %d", *errMsg.Msg, errMsg.Code)
		w.WriteHeader(errMsg.Code)
		w.Write(utils.Marshall(*errMsg.Msg))
		return nil, errNatsMsgReceivedWithErrorData
	}

	return answer.Msg, nil
}

func NatsEnterpriseInteraction(
	subj string,
	body []byte,
	w http.ResponseWriter,
	nc *nats.Conn,
) error {

	log.Println("Sending enterprise message")
	log.Println("Subject: ", subj)

	var answer entity.MsgAnswer[[]byte]

	msg, err := nc.Request(subj, body, local.NATS_REQUEST_TIMEOUT+20*time.Second)
	if err != nil {
		if err == nats.ErrNoResponders {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(utils.Marshall("You have no enterprise license, to get one contact: sales@oktopus.app.br"))
			return err
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.Marshall("Error to communicate with nats:" + err.Error()))
		return err
	}

	err = json.Unmarshal(msg.Data, &answer)
	if err != nil {

		var errMsg *entity.MsgAnswer[*string]
		err = json.Unmarshal(msg.Data, &errMsg)

		if err != nil {
			log.Println("Bad answer message formatting: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(msg.Data)
			return err
		}

		log.Printf("Error message received, msg: %s, code: %d", *errMsg.Msg, errMsg.Code)
		w.WriteHeader(errMsg.Code)
		w.Write(utils.Marshall(*errMsg.Msg))
		return errNatsMsgReceivedWithErrorData
	}

	w.Write(answer.Msg)
	return nil
}
