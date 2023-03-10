package mqtt

import (
	"context"
	"log"
	"net"

	"github.com/eclipse/paho.golang/paho"
	"github.com/leandrofars/oktopus/internal/utils"
)

func StartMqttClient(addr, port *string) *paho.Client {

	conn, err := net.Dial("tcp", *addr+":"+*port)
	if err != nil {
		log.Fatal(err)
	}

	clientConfig := paho.ClientConfig{
		Conn: conn,
	}

	return paho.NewClient(clientConfig)
}

func StartNewConnection(id, user, pass string) paho.Connect {

	connParameters := paho.Connect{
		KeepAlive:  30,
		ClientID:   id,
		CleanStart: true,
		Username:   user,
		Password:   []byte(pass),
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
		connParameters.UsernameFlag = true
	}
	if pass != "" {
		connParameters.PasswordFlag = true
	}

	return connParameters

}

func ConnectMqttBroker(c *paho.Client, cp paho.Connect, addr *string) {
	conn, err := c.Connect(context.Background(), &cp)
	if err != nil {
		log.Fatal(err)
	}

	if conn.ReasonCode != 0 {
		log.Fatalf("Failed to connect to %s : %d - %s", *addr, conn.ReasonCode, conn.Properties.ReasonString)
	}
}
