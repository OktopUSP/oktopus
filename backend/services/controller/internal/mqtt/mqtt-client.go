package mqtt

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	usp_msg "github.com/leandrofars/oktopus/internal/usp_message"
	"github.com/leandrofars/oktopus/internal/usp_record"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/eclipse/paho.golang/paho"
	"github.com/leandrofars/oktopus/internal/utils"
)

type Mqtt struct {
	Addr         string
	Port         string
	Id           string
	User         string
	Passwd       string
	Ctx          context.Context
	QoS          int
	SubTopic     string
	DevicesTopic string
	CA           string
}

var c *paho.Client

/* ------------------- Implementations of broker interface ------------------ */

func (m *Mqtt) Connect() {
	devices := make(chan *paho.Publish)
	controller := make(chan *paho.Publish)
	go m.messageHandler(devices, controller)
	clientConfig := m.startClient(devices, controller)
	connParameters := startConnection(m.Id, m.User, m.Passwd)

	conn, err := clientConfig.Connect(m.Ctx, &connParameters)
	if err != nil {
		log.Println(err)
	}
	if conn.ReasonCode != 0 {
		log.Fatalf("Failed to connect to %s : %d - %s", m.Addr+m.Port, conn.ReasonCode, conn.Properties.ReasonString)
	}

	// Sets global client to be used by other mqtt functions
	c = clientConfig

	if conn.ReasonCode != 0 {
		log.Fatalf("Failed to connect to %s : %d - %s", m.Addr, conn.ReasonCode, conn.Properties.ReasonString)
	}

	log.Printf("Connected to broker--> %s:%s", m.Addr, m.Port)
}

func (m *Mqtt) Disconnect() {
	d := &paho.Disconnect{ReasonCode: 0}
	err := c.Disconnect(d)
	if err != nil {
		log.Fatalf("failed to send Disconnect: %s", err)
	}
}

func (m *Mqtt) Subscribe() {
	if _, err := c.Subscribe(m.Ctx, &paho.Subscribe{
		Subscriptions: map[string]paho.SubscribeOptions{
			m.SubTopic:     {QoS: byte(m.QoS), NoLocal: true},
			m.DevicesTopic: {QoS: byte(m.QoS), NoLocal: true},
		},
	}); err != nil {
		log.Fatalln(err)
	}

	log.Printf("Subscribed to %s", m.SubTopic)
	log.Printf("Subscribed to %s", m.DevicesTopic)
}

func (m *Mqtt) Publish(msg []byte, topic, respTopic string) {
	if _, err := c.Publish(context.Background(), &paho.Publish{
		Topic:   topic,
		QoS:     byte(m.QoS),
		Retain:  false,
		Payload: msg,
		Properties: &paho.PublishProperties{
			ResponseTopic: respTopic,
		},
	}); err != nil {
		log.Println("error sending message:", err)
	}

	log.Printf("Published to %s", topic)
}

/* -------------------------------------------------------------------------- */

func (m *Mqtt) startClient(devices, controller chan *paho.Publish) *paho.Client {
	singleHandler := paho.NewSingleHandlerRouter(func(p *paho.Publish) {
		if p.Topic == m.DevicesTopic {
			devices <- p
		} else if strings.Contains(p.Topic, "controller") {
			controller <- p
		} else {
			log.Println("No handler for topic: ", p.Topic)
		}
	})

	if m.CA != "" {
		conn := connWithTls(m.CA, m.Addr+":"+m.Port, m.Ctx)
		clientConfig := paho.ClientConfig{
			Conn:   conn,
			Router: singleHandler,
			OnServerDisconnect: func(disconnect *paho.Disconnect) {
				log.Println("disconnected from mqtt server, reason code: ", disconnect.ReasonCode)
			},
			OnClientError: func(err error) {
				log.Println(err)
			},
		}
		return paho.NewClient(clientConfig)
	}

	conn, err := net.Dial("tcp", m.Addr+":"+m.Port)
	if err != nil {
		log.Println(err)
	}

	clientConfig := paho.ClientConfig{
		Conn:   conn,
		Router: singleHandler,
	}

	return paho.NewClient(clientConfig)
}

func connWithTls(tlsCa, address string, ctx context.Context) net.Conn {
	ca, err := ioutil.ReadFile(tlsCa)
	if err != nil {
		log.Fatal(err)
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(ca)
	if !ok {
		panic("failed to parse root certificate")
	}

	config := &tls.Config{
		// After going to cloud, certificates must match names, and we must take this option below
		InsecureSkipVerify: true,
		RootCAs:            roots,
	}

	d := tls.Dialer{
		Config: config,
	}

	conn, err := d.DialContext(ctx, "tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	conn = newThreadSafeConnection(conn)

	return conn
}

// Custom net.Conn with thread safety
func newThreadSafeConnection(c net.Conn) net.Conn {
	type threadSafeConn struct {
		net.Conn
		sync.Locker
	}

	return &threadSafeConn{
		Conn:   c,
		Locker: &sync.Mutex{},
	}
}

func startConnection(id, user, pass string) paho.Connect {

	connParameters := paho.Connect{
		KeepAlive:  30,
		ClientID:   id,
		CleanStart: true,
	}

	if id != "" {
		connParameters.ClientID = id
	} else {
		mac, err := utils.GetMacAddr()
		if err != nil {
			log.Fatal(err)
		}
		connParameters.ClientID = mac[0]
	}

	if user != "" {
		connParameters.Username = user
		connParameters.UsernameFlag = true
	}
	if pass != "" {
		connParameters.Password = []byte(pass)
		connParameters.PasswordFlag = true
	}

	return connParameters
}

func (m *Mqtt) messageHandler(devices, controller chan *paho.Publish) {
	for {
		select {
		case d := <-devices:
			payload := string(d.Payload)
			log.Println("New device: ", payload)
			m.handleNewDevice(payload)
		case c := <-controller:
			m.handleDevicesResponse(c.Payload)
		}
	}
}

func (m *Mqtt) handleNewDevice(deviceMac string) {
	payload := usp_msg.Msg{
		Header: &usp_msg.Header{
			MsgId:   "uniqueIdentifierForThismessage",
			MsgType: usp_msg.Header_GET,
		},
		Body: &usp_msg.Body{
			MsgBody: &usp_msg.Body_Request{
				Request: &usp_msg.Request{
					ReqType: &usp_msg.Request_Get{
						Get: &usp_msg.Get{
							ParamPaths: []string{"Device.DeviceInfo."},
							MaxDepth:   1,
						},
					},
				},
			},
		},
	}
	teste, _ := proto.Marshal(&payload)
	record := usp_record.Record{
		Version:         "0.1",
		ToId:            deviceMac,
		FromId:          "leleco",
		PayloadSecurity: usp_record.Record_PLAINTEXT,
		RecordType: &usp_record.Record_NoSessionContext{
			NoSessionContext: &usp_record.NoSessionContextRecord{
				Payload: teste,
			},
		},
	}

	tr369Message, err := proto.Marshal(&record)
	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	}
	time.Sleep(5 * time.Second)
	m.Publish(tr369Message, "oktopus/v1/agent/"+deviceMac, "oktopus/v1/controller/"+deviceMac)
}

func (m *Mqtt) handleDevicesResponse(p []byte) {
	var record usp_record.Record
	var message usp_msg.Msg

	err := proto.Unmarshal(p, &record)
	if err != nil {
		log.Fatal(err)
	}
	err = proto.Unmarshal(record.GetNoSessionContext().Payload, &message)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Received a usp_message: %s\n", message.String())
}
