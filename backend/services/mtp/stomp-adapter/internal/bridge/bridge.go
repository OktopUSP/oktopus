package bridge

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"strings"
	"time"

	"github.com/go-stomp/stomp/v3"
	"github.com/go-stomp/stomp/v3/frame"
	"github.com/nats-io/nats.go"
	"github.com/oktopUSP/oktopus/backend/services/mtp/stomp-adapter/internal/config"
	"golang.org/x/sys/unix"
)

const (
	OFFLINE = iota
	ONLINE
)

const STOMP_CONNECTION_RETRY = 5 * time.Second

type msgAnswer struct {
	Code int
	Msg  any
}

const (
	NATS_STOMP_SUBJECT_PREFIX         = "stomp.usp.v1."
	NATS_STOMP_ADAPTER_SUBJECT_PREFIX = "stomp-adapter.usp.v1."
	DEVICE_SUBJECT_PREFIX             = "device.usp.v1."
	STOMP_QUEUE_PREFIX                = "oktopus/usp/v1/"
	STOMP_STATUS_QUEUE                = STOMP_QUEUE_PREFIX + "status"
	DEVICE_TIMEOUT_RESPONSE           = 5 * time.Second
	USP_CONTENT_TYPE                  = "application/vnd.bbf.usp.msg"
)

type (
	Publisher  func(string, []byte) error
	Subscriber func(string, func(*nats.Msg)) error
)

type Bridge struct {
	Pub   Publisher
	Sub   Subscriber
	Stomp config.Stomp
	Ctx   context.Context
}

func NewBridge(p Publisher, s Subscriber, ctx context.Context, stomp config.Stomp) *Bridge {
	return &Bridge{
		Pub:   p,
		Sub:   s,
		Stomp: stomp,
		Ctx:   ctx,
	}
}

func (b *Bridge) StartBridge() {

	options := []func(*stomp.Conn) error{
		stomp.ConnOpt.Login(b.Stomp.User, b.Stomp.Password),
		stomp.ConnOpt.Host("/"),
	}

	var conn *stomp.Conn
	var err error

	go func() {
		for {
			conn, err = connectToServer(b.Stomp.Url, options)
			if err != nil {
				continue
			}
			b.subscribe(conn)

			sub, err := conn.Subscribe(STOMP_STATUS_QUEUE, stomp.AckAuto)
			if err != nil {
				log.Println("cannot subscribe to", STOMP_STATUS_QUEUE, err.Error())
				return
			}
			log.Println("Subscribed to", STOMP_STATUS_QUEUE)

			for {
				if !sub.Active() {
					log.Println("Subscription is no longer active")
					break
				}
				msg := <-sub.C
				body := msg.Header.Get("message")
				if body != "connection closed" {
					log.Println("Received message", body)
					fmtBody := strings.Split(body, "|")
					if len(fmtBody) == 2 {
						deviceQueue := strings.Split(fmtBody[0], "/")
						device := deviceQueue[len(deviceQueue)-1]
						status := fmtBody[1]
						log.Println("Device:", device, "Status:", status)
						b.Pub(NATS_STOMP_SUBJECT_PREFIX+device+".status", []byte(status))
					} else {
						log.Println("Invalid status message", body)
					}
				}
			}
		}
	}()

}

func connectToServer(url string, options []func(*stomp.Conn) error) (*stomp.Conn, error) {

	conn, err := stomp.Dial("tcp", url, options...)

	if err != nil {
		log.Printf("Error to connect to %s, err: %s", url, err)
		time.Sleep(STOMP_CONNECTION_RETRY)
	} else {
		log.Println("Connected to STOMP server", url)
	}

	return conn, err
}

func (b *Bridge) subscribe(st *stomp.Conn) {

	b.Sub(NATS_STOMP_ADAPTER_SUBJECT_PREFIX+"*.info", func(msg *nats.Msg) {

		log.Printf("Received message on info subject")

		subj := strings.Split(msg.Subject, ".")
		device := subj[len(subj)-2]

		deviceInfoQueue := STOMP_QUEUE_PREFIX + "controller/" + device + "/info"

		sub, err := st.Subscribe(deviceInfoQueue, stomp.AckAuto)
		if err != nil {
			log.Println("cannot subscribe to", deviceInfoQueue, err.Error())
			return
		}
		log.Println("Subscribed to", deviceInfoQueue)

		err = st.Send(STOMP_QUEUE_PREFIX+"agent/"+device, "application/vnd.bbf.usp.msg", msg.Data, func(f *frame.Frame) error {
			f.Header.Set("reply-to-dest", deviceInfoQueue)
			return nil
		})

		if err != nil {
			log.Printf("send stomp msg error: %q", err)
			return
		}

		select {
		case data := <-sub.C:
			body := data.Body
			log.Println("Received message answer")
			err = b.Pub(NATS_STOMP_SUBJECT_PREFIX+device+".info", body)
			if err != nil {
				log.Printf("send nats msg error: %q", err)
			}
		case <-time.After(DEVICE_TIMEOUT_RESPONSE):
			log.Println("Timeout waiting for device info response")
		}
		sub.Unsubscribe()
	})

	b.Sub(NATS_STOMP_ADAPTER_SUBJECT_PREFIX+"*.api", func(msg *nats.Msg) {

		log.Printf("Received message on api subject")

		subj := strings.Split(msg.Subject, ".")
		device := subj[len(subj)-2]

		deviceApiQueue := STOMP_QUEUE_PREFIX + "controller/" + device + "/api"

		sub, err := st.Subscribe(deviceApiQueue, stomp.AckAuto)
		if err != nil {
			log.Println("cannot subscribe to", STOMP_STATUS_QUEUE, err.Error())
			return
		}
		log.Println("Subscribed to", deviceApiQueue)

		err = st.Send(STOMP_QUEUE_PREFIX+"agent/"+device, "application/vnd.bbf.usp.msg", msg.Data, func(f *frame.Frame) error {
			f.Header.Set("reply-to-dest", deviceApiQueue)
			return nil
		})

		if err != nil {
			log.Printf("send stomp msg error: %q", err)
			return
		}

		select {
		case data := <-sub.C:
			body := data.Body
			err = b.Pub(DEVICE_SUBJECT_PREFIX+device+".api", body)
			if err != nil {
				log.Printf("send nats msg error: %q", err)
			}
		case <-time.After(DEVICE_TIMEOUT_RESPONSE):
			log.Println("Timeout waiting for device info response")
		}
		sub.Unsubscribe()
	})

	b.Sub(NATS_STOMP_ADAPTER_SUBJECT_PREFIX+"rtt", func(msg *nats.Msg) {

		log.Printf("Received message on rtt subject")

		conn, err := net.Dial("tcp", b.Stomp.Url)
		if err != nil {
			respondMsg(msg.Respond, 500, err.Error())
			return
		}
		defer conn.Close()

		info, err := tcpInfo(conn.(*net.TCPConn))
		if err != nil {
			respondMsg(msg.Respond, 500, err.Error())
			return
		}
		rtt := time.Duration(info.Rtt) * time.Microsecond

		respondMsg(msg.Respond, 200, rtt/1000)

	})
}

func respondMsg(respond func(data []byte) error, code int, msgData any) {

	msg, err := json.Marshal(msgAnswer{
		Code: code,
		Msg:  msgData,
	})
	if err != nil {
		log.Printf("Failed to marshal message: %q", err)
		respond([]byte(err.Error()))
		return
	}

	respond([]byte(msg))
}

func tcpInfo(conn *net.TCPConn) (*unix.TCPInfo, error) {
	raw, err := conn.SyscallConn()
	if err != nil {
		return nil, err
	}

	var info *unix.TCPInfo
	ctrlErr := raw.Control(func(fd uintptr) {
		info, err = unix.GetsockoptTCPInfo(int(fd), unix.IPPROTO_TCP, unix.TCP_INFO)
	})
	switch {
	case ctrlErr != nil:
		return nil, ctrlErr
	case err != nil:
		return nil, err
	}
	return info, nil
}
