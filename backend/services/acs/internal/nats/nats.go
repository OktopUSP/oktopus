package nats

import (
	"context"
	"errors"
	"log"
	"oktopUSP/backend/services/acs/internal/config"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	CWMP_STREAM_NAME = "cwmp"
)

type NatsActions struct {
	Publish   func(string, []byte) error
	Subscribe func(string, func(*nats.Msg)) error
}

func StartNatsClient(c config.Nats) NatsActions {

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

	streams := defineStreams()
	err = createStreams(c.Ctx, js, streams)
	if err != nil {
		log.Fatalf("Failed to create Consumer: %v", err)
	}

	consumers := defineConsumers()
	err = createConsumers(c.Ctx, js, consumers)
	if err != nil {
		log.Fatalf("Failed to create Consumer: %v", err)
	}

	return NatsActions{
		Publish:   publisher(js),
		Subscribe: subscriber(nc),
	}
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

func createStreams(ctx context.Context, js jetstream.JetStream, streams []string) error {
	for _, stream := range streams {
		_, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
			Name:        stream,
			Description: "Stream for " + stream + " messages",
			Subjects:    []string{stream + ".>"},
			Retention:   jetstream.InterestPolicy,
		})
		if err != nil {
			return errors.New(err.Error() + " | consumer:" + stream)
		}
	}
	return nil
}

func createConsumers(ctx context.Context, js jetstream.JetStream, consumers []string) error {
	for _, consumer := range consumers {
		_, err := js.CreateOrUpdateConsumer(ctx, consumer, jetstream.ConsumerConfig{
			Name:        consumer,
			Description: "Consumer for " + consumer + " messages",
			AckPolicy:   jetstream.AckExplicitPolicy,
			Durable:     consumer,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func defineStreams() []string {
	return []string{
		CWMP_STREAM_NAME,
	}
}

func defineConsumers() []string {
	return []string{
		CWMP_STREAM_NAME,
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
