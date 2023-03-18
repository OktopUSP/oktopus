package mqtt

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"sync"

	"github.com/eclipse/paho.golang/paho"
	"github.com/leandrofars/oktopus/internal/utils"
)

type Mqtt struct {
	Addr     string
	Port     string
	Id       string
	User     string
	Passwd   string
	Ctx      context.Context
	QoS      int
	SubTopic string
	CA       string
}

var c *paho.Client

/* ------------------- Implementations of broker interface ------------------ */

func (m *Mqtt) Connect() {
	msgChan := make(chan *paho.Publish)
	go messageHandler(msgChan)
	clientConfig := startClient(m.Addr, m.Port, m.CA, m.Ctx, msgChan)
	connParameters := startConnection(m.Id, m.User, m.Passwd)

	conn, err := clientConfig.Connect(m.Ctx, &connParameters)
	if err != nil {
		log.Println(err)
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
			m.SubTopic: {QoS: byte(m.QoS), NoLocal: true},
		},
	}); err != nil {
		log.Fatalln(err)
	}

	log.Printf("Subscribed to %s", m.SubTopic)
}

/* -------------------------------------------------------------------------- */

func startClient(addr string, port string, tlsCa string, ctx context.Context, msgChan chan *paho.Publish) *paho.Client {
	singleHandler := paho.NewSingleHandlerRouter(func(m *paho.Publish) {
		msgChan <- m
	})

	if tlsCa != "" {
		conn := connWithTls(tlsCa, addr+":"+port, ctx)
		clientConfig := paho.ClientConfig{
			Conn:   conn,
			Router: singleHandler,
		}
		return paho.NewClient(clientConfig)
	}

	conn, err := net.Dial("tcp", addr+":"+port)
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

func messageHandler(msg chan *paho.Publish) {
	for m := range msg {
		log.Println("Received message:", string(m.Payload))
	}
}
