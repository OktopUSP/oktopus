package nats

import (
	"log"
	"time"

	"github.com/OktopUSP/oktopus/backend/services/mtp/ws-adapter/internal/config"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	STREAM_NAME        = "ws"
	BUCKET_NAME        = "devices-auth"
	BUCKET_DESCRIPTION = "Devices authentication"
)

func StartNatsClient(c config.Nats) (
	jetstream.KeyValue,
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

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatalf("Failed to create JetStream client: %v", err)
	}

	kv, err := js.CreateOrUpdateKeyValue(c.Ctx, jetstream.KeyValueConfig{
		Bucket:      BUCKET_NAME,
		Description: BUCKET_DESCRIPTION,
	})
	if err != nil {
		log.Fatalf("Failed to create KeyValue store: %v", err)
	}

	return kv, publisher(js), subscriber(nc)
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

func publisher(js jetstream.JetStream) func(string, []byte) error {
	return func(subject string, payload []byte) error {
		_, err := js.PublishAsync(subject, payload)
		if err != nil {
			log.Printf("error to send jetstream message: %q", err)
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
