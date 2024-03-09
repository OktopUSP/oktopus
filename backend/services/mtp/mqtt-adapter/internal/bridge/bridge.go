package bridge

import (
	"context"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/OktopUSP/oktopus/backend/services/mqtt-adapter/internal/config"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

const (
	OFFLINE = iota
	ONLINE
)

const NATS_MQTT_SUBJECT_PREFIX = "mqtt.usp.v1."
const NATS_MQTT_ADAPTER_SUBJECT_PREFIX = "mqtt-adapter.usp.v1.*."
const MQTT_TOPIC_PREFIX = "oktopus/usp/"

type (
	Publisher  func(string, []byte) error
	Subscriber func(string, func(*nats.Msg)) error
)

type Bridge struct {
	Pub  Publisher
	Sub  Subscriber
	Mqtt config.Mqtt
	Ctx  context.Context
}

func NewBridge(p Publisher, s Subscriber, ctx context.Context, m config.Mqtt) *Bridge {
	return &Bridge{
		Pub:  p,
		Sub:  s,
		Mqtt: m,
		Ctx:  ctx,
	}
}

func (b *Bridge) StartBridge() {

	broker, _ := url.Parse(b.Mqtt.Url)

	status := make(chan *paho.Publish)
	controller := make(chan *paho.Publish)
	apiMsg := make(chan *paho.Publish)

	go b.mqttMessageHandler(status, controller, apiMsg)

	pahoClientConfig := buildClientConfig(status, controller, apiMsg, b.Mqtt.ClientId)

	autopahoClientConfig := autopaho.ClientConfig{
		BrokerUrls: []*url.URL{
			broker,
		},
		KeepAlive:         30,
		ConnectRetryDelay: 5 * time.Second,
		ConnectTimeout:    5 * time.Second,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			log.Printf("Connected to MQTT broker--> %s", b.Mqtt.Url)
			subscribe(b.Mqtt.Ctx, b.Mqtt.Qos, cm)
		},
		OnConnectError: func(err error) {
			log.Printf("Error while attempting connection: %s\n", err)
		},
		ClientConfig: *pahoClientConfig,
	}

	if b.Mqtt.Username != "" && b.Mqtt.Password != "" {
		autopahoClientConfig.SetUsernamePassword(b.Mqtt.Username, []byte(b.Mqtt.Password))
	}

	log.Println("MQTT client id:", pahoClientConfig.ClientID)
	log.Println("MQTT username:", b.Mqtt.Username)
	log.Println("MQTT password:", b.Mqtt.Password)

	cm, err := autopaho.NewConnection(b.Ctx, autopahoClientConfig)
	if err != nil {
		log.Fatalln(err)
	}

	b.natsMessageHandler(cm)
}

func (b *Bridge) natsMessageHandler(cm *autopaho.ConnectionManager) {
	b.Sub(NATS_MQTT_ADAPTER_SUBJECT_PREFIX+"info", func(m *nats.Msg) {

		log.Printf("Received message on info subject")
		cm.Publish(b.Ctx, &paho.Publish{
			QoS:     byte(b.Mqtt.Qos),
			Topic:   MQTT_TOPIC_PREFIX + "v1/agent/" + getDeviceFromSubject(m.Subject),
			Payload: m.Data,
			Properties: &paho.PublishProperties{
				ResponseTopic: "oktopus/usp/v1/controller/" + getDeviceFromSubject(m.Subject),
			},
		})
	})
}

func getDeviceFromSubject(subject string) string {
	paths := strings.Split(subject, ".")
	device := paths[len(paths)-2]
	return device
}

func (b *Bridge) mqttMessageHandler(status, controller, apiMsg chan *paho.Publish) {
	for {
		select {
		case d := <-status:
			b.Pub(NATS_MQTT_SUBJECT_PREFIX+getDeviceFromTopic(d.Topic)+".status", d.Payload)
		case c := <-controller:
			b.Pub(NATS_MQTT_SUBJECT_PREFIX+getDeviceFromTopic(c.Topic)+".info", c.Payload)
		case a := <-apiMsg:
			b.Pub(NATS_MQTT_SUBJECT_PREFIX+getDeviceFromTopic(a.Topic)+".api", a.Payload)
		}
	}
}

func getDeviceFromTopic(topic string) string {
	paths := strings.Split(topic, "/")
	device := paths[len(paths)-1]
	return device
}

func subscribe(ctx context.Context, qos int, c *autopaho.ConnectionManager) {
	if _, err := c.Subscribe(ctx, &paho.Subscribe{
		Subscriptions: []paho.SubscribeOptions{
			{
				Topic: MQTT_TOPIC_PREFIX + "+/api/+",
				QoS:   byte(qos),
			},
			{
				Topic: MQTT_TOPIC_PREFIX + "+/controller/+",
				QoS:   byte(qos),
			},
			{
				Topic: MQTT_TOPIC_PREFIX + "+/status/+",
				QoS:   byte(qos),
			},
		},
	}); err != nil {
		log.Fatalln(err)
	}

	log.Printf("Subscribed to %s", MQTT_TOPIC_PREFIX+"+/controller/+")
	log.Printf("Subscribed to %s", MQTT_TOPIC_PREFIX+"+/status/+")
	log.Printf("Subscribed to %s", MQTT_TOPIC_PREFIX+"+/api/+")
}

func buildClientConfig(status, controller, apiMsg chan *paho.Publish, id string) *paho.ClientConfig {
	log.Println("Starting new MQTT client")
	singleHandler := paho.NewSingleHandlerRouter(func(p *paho.Publish) {

		if strings.Contains(p.Topic, "status") {
			status <- p
		} else if strings.Contains(p.Topic, "controller") {
			controller <- p
		} else if strings.Contains(p.Topic, "api") {
			apiMsg <- p
		} else {
			log.Println("No handler for topic: ", p.Topic)
		}

	})

	clientConfig := paho.ClientConfig{}

	clientConfig = paho.ClientConfig{
		Router: singleHandler,
		OnServerDisconnect: func(d *paho.Disconnect) {
			if d.Properties != nil {
				log.Printf("Requested disconnect: %s\n , properties reason: %s\n", clientConfig.ClientID, d.Properties.ReasonString)
			} else {
				log.Printf("Requested disconnect; %s reason code: %d\n", clientConfig.ClientID, d.ReasonCode)
			}
		},
		OnClientError: func(err error) {
			log.Println(err)
		},
	}

	if id != "" {
		clientConfig.ClientID = id
	} else {
		clientConfig.ClientID = uuid.NewString()
	}

	return &clientConfig
}
