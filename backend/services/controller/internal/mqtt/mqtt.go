package mqtt

import (
	"context"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/leandrofars/oktopus/internal/db"
	usp_msg "github.com/leandrofars/oktopus/internal/usp_message"
	"github.com/leandrofars/oktopus/internal/usp_record"
	"github.com/leandrofars/oktopus/internal/utils"
	"google.golang.org/protobuf/proto"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Mqtt struct {
	Addr            string
	Port            string
	Id              string
	User            string
	Passwd          string
	Ctx             context.Context
	QoS             int
	SubTopic        string
	DevicesTopic    string
	DisconnectTopic string
	TLS             bool
	DB              db.Database
	MsgQueue        map[string](chan usp_msg.Msg)
	QMutex          *sync.Mutex
}

var c *autopaho.ConnectionManager

/* ------------------- Implementations of broker interface ------------------ */

func (m *Mqtt) Connect() {

	broker, _ := url.Parse("tcp://" + m.Addr + ":" + m.Port)

	devices := make(chan *paho.Publish)
	controller := make(chan *paho.Publish)
	disconnect := make(chan *paho.Publish)
	apiMsg := make(chan *paho.Publish)

	go m.messageHandler(devices, controller, disconnect, apiMsg)
	pahoClientConfig := m.buildClientConfig(devices, controller, disconnect, apiMsg)

	autopahoClientConfig := autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{broker},
		KeepAlive:         30,
		ConnectRetryDelay: 5 * time.Second,
		ConnectTimeout:    5 * time.Second,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			log.Printf("Connected to broker--> %s:%s", m.Addr, m.Port)
			m.Subscribe()
		},
		OnConnectError: func(err error) {
			log.Printf("Error while attempting connection: %s\n", err)
		},
		ClientConfig: *pahoClientConfig,
	}

	if m.User != "" && m.Passwd != "" {
		autopahoClientConfig.SetUsernamePassword(m.User, []byte(m.Passwd))
	}

	log.Println("MQTT client id:", pahoClientConfig.ClientID)
	log.Println("MQTT username:", m.User)
	log.Println("MQTT password:", m.Passwd)

	cm, err := autopaho.NewConnection(m.Ctx, autopahoClientConfig)
	if err != nil {
		log.Fatalln(err)
	}

	c = cm

	//devices := make(chan *paho.Publish)
	//controller := make(chan *paho.Publish)
	//disconnect := make(chan *paho.Publish)
	//apiMsg := make(chan *paho.Publish)
	//go m.messageHandler(devices, controller, disconnect, apiMsg)
	//clientConfig := m.startClient(devices, controller, disconnect, apiMsg)
	//connParameters := startConnection(m.Id, m.User, m.Passwd)
	//
	//conn, err := clientConfig.Connect(m.Ctx, &connParameters)
	//if err != nil {
	//	log.Println(err)
	//}
	//if conn.ReasonCode != 0 {
	//	log.Fatalf("Failed to connect to %s : %d - %s", m.Addr+m.Port, conn.ReasonCode, conn.Properties.ReasonString)
	//}
	//
	//// Sets global client to be used by other mqtt functions
	//c = clientConfig
	//
	//log.Printf("Connected to broker--> %s:%s", m.Addr, m.Port)
}

func (m *Mqtt) Disconnect() {
	err := c.Disconnect(m.Ctx)
	if err != nil {
		log.Fatalf("failed to send Disconnect: %s", err)
	}
}

func (m *Mqtt) Subscribe() {
	if _, err := c.Subscribe(m.Ctx, &paho.Subscribe{
		Subscriptions: map[string]paho.SubscribeOptions{
			m.SubTopic:        {QoS: byte(m.QoS), NoLocal: true},
			m.DevicesTopic:    {QoS: byte(m.QoS), NoLocal: true},
			m.DisconnectTopic: {QoS: byte(m.QoS), NoLocal: true},
			"oktopus/+/api/+": {QoS: byte(m.QoS), NoLocal: true},
		},
	}); err != nil {
		log.Fatalln(err)
	}

	log.Printf("Subscribed to %s", m.SubTopic)
	log.Printf("Subscribed to %s", m.DevicesTopic)
	log.Printf("Subscribed to %s", m.DisconnectTopic)
	log.Println("Subscribed to %s", "oktopus/+/api/+")
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

func (m *Mqtt) buildClientConfig(devices, controller, disconnect, apiMsg chan *paho.Publish) *paho.ClientConfig {
	log.Println("Starting new mqtt client")
	singleHandler := paho.NewSingleHandlerRouter(func(p *paho.Publish) {
		if p.Topic == m.DevicesTopic {
			devices <- p
		} else if strings.Contains(p.Topic, "controller") {
			controller <- p
		} else if p.Topic == m.DisconnectTopic {
			disconnect <- p
		} else if strings.Contains(p.Topic, "api") {
			apiMsg <- p
		} else {
			log.Println("No handler for topic: ", p.Topic)
		}
	})

	//if m.TLS {
	//	conn := connWithTls(m.Addr+":"+m.Port, m.Ctx)
	//	clientConfig := paho.ClientConfig{
	//		Conn:   conn,
	//		Router: singleHandler,
	//		OnServerDisconnect: func(discon *paho.Disconnect) {
	//			log.Println("disconnected from mqtt server, reason code: ", discon.ReasonCode)
	//		},
	//		OnClientError: func(err error) {
	//			log.Println(err)
	//		},
	//	}
	//	return &clientConfig
	//}

	//conn, _ := connWithoutTLS(m.Addr, m.Port)

	clientConfig := paho.ClientConfig{}

	clientConfig = paho.ClientConfig{
		//Conn:   conn,
		Router: singleHandler,
		OnServerDisconnect: func(d *paho.Disconnect) {
			if d.Properties != nil {
				log.Printf("Requested disconnect: %s\n", clientConfig.ClientID, d.Properties.ReasonString)
			} else {
				log.Printf("Requested disconnect; reason code: %d\n", clientConfig.ClientID, d.ReasonCode)
			}
		},
		OnClientError: func(err error) {
			log.Println(err)
		},
	}

	if m.Id != "" {
		clientConfig.ClientID = m.Id
	} else {
		mac, err := utils.GetMacAddr()
		if err != nil {
			log.Fatal(err)
		}
		clientConfig.ClientID = mac[0]
	}

	return &clientConfig
}

//Not used by autopaho package
//func connWithoutTLS(addr, port string) (net.Conn, error) {
//	x := 1
//	for {
//		log.Println("Trying connection to broker...")
//		conn, err := net.Dial("tcp", addr+":"+port)
//		if err == nil {
//			log.Println("Successfully connected to broker!")
//			return conn, err
//		} else {
//			if x == 5 {
//				log.Println("Couldn't connect to broker")
//				os.Exit(1)
//			}
//			x = x + 1
//			time.Sleep(5 * time.Second)
//		}
//	}
//}

//Not used by autopaho package
//func connWithTls(address string, ctx context.Context) net.Conn {
//	config := &tls.Config{
//		// After going to cloud, certificates must match names, and we must take this option below
//		InsecureSkipVerify: true,
//	}
//
//	d := tls.Dialer{
//		Config: config,
//	}
//
//	conn, err := d.DialContext(ctx, "tcp", address)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	conn = newThreadSafeConnection(conn)
//
//	return conn
//}

// Custom net.Conn with thread safety
//func newThreadSafeConnection(c net.Conn) net.Conn {
//	type threadSafeConn struct {
//		net.Conn
//		sync.Locker
//	}
//
//	return &threadSafeConn{
//		Conn:   c,
//		Locker: &sync.Mutex{},
//	}
//}

//func startConnection(id, user, pass string) paho.Connect {
//
//	connParameters := paho.Connect{
//		KeepAlive:  30,
//		ClientID:   id,
//		CleanStart: true,
//	}
//
//	if id != "" {
//		connParameters.ClientID = id
//	} else {
//		mac, err := utils.GetMacAddr()
//		if err != nil {
//			log.Fatal(err)
//		}
//		connParameters.ClientID = mac[0]
//		log.Println("MQTT client id:", connParameters.ClientID)
//	}
//
//	if user != "" {
//		connParameters.Username = user
//		connParameters.UsernameFlag = true
//		log.Println("MQTT username:", connParameters.Username)
//	}
//	if pass != "" {
//		connParameters.Password = []byte(pass)
//		connParameters.PasswordFlag = true
//		log.Println("MQTT password:", pass)
//	}
//
//	return connParameters
//}

func (m *Mqtt) messageHandler(devices, controller, disconnect, apiMsg chan *paho.Publish) {
	for {
		select {
		case d := <-devices:
			payload := string(d.Payload)
			log.Println("New device: ", payload)
			m.handleNewDevice(payload)
		case c := <-controller:
			topic := c.Topic
			sn := strings.Split(topic, "/")
			m.handleNewDevicesResponse(c.Payload, sn[3])
		case dis := <-disconnect:
			payload := string(dis.Payload)
			log.Println("Device disconnected: ", payload)
			m.handleDevicesDisconnect(payload)
		case api := <-apiMsg:
			log.Println("Handle api request")
			m.handleApiRequest(api.Payload)
		}
	}
}

func (m *Mqtt) handleApiRequest(api []byte) {
	var record usp_record.Record
	err := proto.Unmarshal(api, &record)
	if err != nil {
		log.Println(err)
	}

	var msg usp_msg.Msg
	err = proto.Unmarshal(record.GetNoSessionContext().Payload, &msg)
	if err != nil {
		log.Println(err)
	}

	if _, ok := m.MsgQueue[msg.Header.MsgId]; ok {
		//m.QMutex.Lock()
		m.MsgQueue[msg.Header.MsgId] <- msg
		//m.QMutex.Unlock()
	} else {
		log.Printf("Message answer to request %s arrived too late", msg.Header.MsgId)
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
							ParamPaths: []string{
								"Device.DeviceInfo.Manufacturer",
								"Device.DeviceInfo.ModelName",
								"Device.DeviceInfo.SoftwareVersion",
								"Device.DeviceInfo.SerialNumber",
							},
							MaxDepth: 1,
						},
					},
				},
			},
		},
	}
	teste, _ := proto.Marshal(&payload)
	record := utils.NewUspRecord(teste, deviceMac)

	tr369Message, err := proto.Marshal(&record)
	if err != nil {
		log.Fatalln("Failed to encode tr369 record:", err)
	}
	m.Publish(tr369Message, "oktopus/v1/agent/"+deviceMac, "oktopus/v1/controller/"+deviceMac)
}

func (m *Mqtt) handleNewDevicesResponse(p []byte, sn string) {
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

	var device db.Device
	msg := message.Body.MsgBody.(*usp_msg.Body_Response).Response.GetGetResp()

	device.Vendor = msg.ReqPathResults[0].ResolvedPathResults[0].ResultParams["Manufacturer"]
	device.Model = msg.ReqPathResults[1].ResolvedPathResults[0].ResultParams["ModelName"]
	device.Version = msg.ReqPathResults[2].ResolvedPathResults[0].ResultParams["SoftwareVersion"]
	device.SN = sn
	device.Status = utils.Online

	err = m.DB.CreateDevice(device)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *Mqtt) handleDevicesDisconnect(p string) {
	// Update status of device at database
	err := m.DB.UpdateStatus(p, utils.Offline)
	if err != nil {
		log.Fatal(err)
	}
}

/*
func (m *Mqtt) Request(msg []byte, msgType usp_msg.Header_MsgType, pubTopic string, respTopic string) {
	m.Publish(msg, pubTopic, respTopic)
}*/
