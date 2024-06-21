package nats

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/oktopUSP/backend/services/bulkdata/internal/config"
)

func StartNatsClient(c config.Nats) (
	func(string, []byte) error,
	func(string, func(*nats.Msg)) error,
) {

	var (
		nc  *nats.Conn
		err error
	)

	opts := defineOptions(c)

	log.Printf("Connecting to NATS server %s", c.Url)

	for {
		nc, err = nats.Connect(c.Url, opts...)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	log.Printf("Successfully connected to NATS server %s", c.Url)

	return publisher(nc), subscriber(nc)
}

func subscriber(nc *nats.Conn) func(string, func(*nats.Msg)) error {
	return func(subject string, handler func(*nats.Msg)) error {
		_, err := nc.Subscribe(subject, handler)
		if err != nil {
			log.Printf("error to subscribe to subject %s error: %q", subject, err)
		}
		return err
	}
}

func publisher(nc *nats.Conn) func(string, []byte) error {
	return func(subject string, payload []byte) error {
		err := nc.Publish(subject, payload)
		if err != nil {
			log.Printf("error to send nats core message: %q", err)
		}
		return err
	}
}

func defineOptions(c config.Nats) []nats.Option {
	var opts []nats.Option

	opts = append(opts, nats.Name(c.Name))
	opts = append(opts, nats.MaxReconnects(-1))
	opts = append(opts, nats.ReconnectWait(5*time.Second))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Got disconnected! Reason: %q\n", err)
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Got reconnected to %v!\n", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Printf("Connection closed. Reason: %q\n", nc.LastError())
	}))
	if c.VerifyCertificates {
		opts = append(opts, nats.RootCAs())
	}

	return opts
}
