package nats

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/OktopUSP/oktopus/backend/services/mtp/adapter/internal/config"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	MQTT_STREAM_NAME  = "mqtt"
	WS_STREAM_NAME    = "ws"
	STOMP_STREAM_NAME = "stomp"
	LORA_STREAM_NAME  = "lora"
	OPC_STREAM_NAME   = "opc"
)

func StartNatsClient(c config.Nats) (jetstream.JetStream, *nats.Conn) {

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

	return js, nc
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
		MQTT_STREAM_NAME,
		WS_STREAM_NAME,
		STOMP_STREAM_NAME,
		LORA_STREAM_NAME,
		OPC_STREAM_NAME,
	}
}

func defineConsumers() []string {
	return []string{
		MQTT_STREAM_NAME,
		WS_STREAM_NAME,
		STOMP_STREAM_NAME,
		LORA_STREAM_NAME,
		OPC_STREAM_NAME,
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
