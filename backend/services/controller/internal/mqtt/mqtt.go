package mqtt

import (
	"context"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/leandrofars/oktopus/internal/db"
	usp_msg "github.com/leandrofars/oktopus/internal/usp_message"
	"github.com/leandrofars/oktopus/internal/usp_record"
	"github.com/leandrofars/oktopus/internal/utils"
	"google.golang.org/protobuf/proto"
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
	TLS          bool
	DB           db.Database
	MsgQueue     map[string](chan usp_msg.Msg)
	QMutex       *sync.Mutex
}

const (
	ONLINE = iota
	OFFLINE
)

var c *autopaho.ConnectionManager

/* ------------------- Implementations of broker interface ------------------ */

func (m *Mqtt) Connect() {

	broker, _ := url.Parse("tcp://" + m.Addr + ":" + m.Port)

	status := make(chan *paho.Publish)
	controller := make(chan *paho.Publish)
	apiMsg := make(chan *paho.Publish)

	go m.messageHandler(status, controller, apiMsg)
	pahoClientConfig := m.buildClientConfig(status, controller, apiMsg)

	autopahoClientConfig := autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{broker},
		KeepAlive:         30,
		ConnectRetryDelay: 5 * time.Second,
		ConnectTimeout:    5 * time.Second,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			log.Printf("Connected to MQTT broker--> %s:%s", m.Addr, m.Port)
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
			m.SubTopic:        {QoS: byte(m.QoS)},
			m.DevicesTopic:    {QoS: byte(m.QoS)},
			"oktopus/+/api/+": {QoS: byte(m.QoS)},
		},
	}); err != nil {
		log.Fatalln(err)
	}

	log.Printf("Subscribed to %s", m.SubTopic)
	log.Printf("Subscribed to %s", m.DevicesTopic)
	log.Printf("Subscribed to %s", "oktopus/+/api/+")
}

func (m *Mqtt) Publish(msg []byte, topic, respTopic string, retain bool) {
	if _, err := c.Publish(context.Background(), &paho.Publish{
		Topic:   topic,
		QoS:     byte(m.QoS),
		Retain:  retain,
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

func (m *Mqtt) buildClientConfig(status, controller, apiMsg chan *paho.Publish) *paho.ClientConfig {
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
		//Conn:   conn,
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

func (m *Mqtt) messageHandler(status, controller, apiMsg chan *paho.Publish) {
	for {
		select {
		case d := <-status:
			paths := strings.Split(d.Topic, "/")
			device := paths[len(paths)-1]
			payload, err := strconv.Atoi(string(d.Payload))
			if err != nil {
				log.Println("Status topic payload message type error")
				log.Fatalln(err)
			}
			if payload == ONLINE {
				log.Println("Device connected:", device)
				m.handleNewDevice(device)
				//m.deleteRetainedMessage(d, device)
			} else if payload == OFFLINE {
				log.Println("Device disconnected:1", device)
				m.handleDevicesDisconnect(device)
				//m.deleteRetainedMessage(d, device)
			} else {
				log.Println("Status topic payload message type error")
			}
		case c := <-controller:
			topic := c.Topic
			sn := strings.Split(topic, "/")
			m.handleNewDevicesResponse(c.Payload, sn[3])
		case api := <-apiMsg:
			log.Println("Handle api request")
			m.handleApiRequest(api.Payload)
		}
	}
}

//TODO: handle device status at mochi redis
//func (m *Mqtt) deleteRetainedMessage(message *paho.Publish, deviceMac string) {
//	m.Publish([]byte(""), "oktopus/v1/status/"+deviceMac, "", true)
//	log.Println("Message contains the retain flag, deleting it, as it's already received")
//}

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
	m.Publish(tr369Message, "oktopus/v1/agent/"+deviceMac, "oktopus/v1/controller/"+deviceMac, false)
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
